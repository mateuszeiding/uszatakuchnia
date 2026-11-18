// Przywraca jeden przepis z kopii zapasowej do bazy docelowej (zwykle produkcji).
//
// Kopię najpierw odtwarzasz do lokalnej bazy (docs/backup-bazy.md), potem:
//
//	go run ./cmd/restore-recipe -from "<DSN lokalnej kopii>" -id 12           # podgląd, nic nie zapisuje
//	go run ./cmd/restore-recipe -from "<DSN lokalnej kopii>" -id 12 -apply    # zapis
//
// Baza docelowa: flaga -to, a bez niej DB_DATABASE_URL.
//
// Co robi (w jednej transakcji): wiersz w `recipes` wstawia albo nadpisuje (to samo id),
// kroki, składniki, zamienniki, tagi, sprzęt i link do zdjęcia przepisu kasuje i wstawia z kopii,
// a brakujące zdjęcie (`photos`) dokłada. Inne przepisy zostają nietknięte. Kolumny dopasowuje
// po nazwie, więc kopia sprzed zmiany schematu też zadziała (nowe kolumny dostają NULL).
// Składniki, jednostki, tagi i sprzęt muszą istnieć w bazie docelowej — jeśli któregoś brakuje,
// transakcja się wycofuje z błędem klucza obcego.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Kolejność wstawiania; kasowanie idzie od końca (zamienniki przed składnikami).
var childTables = []string{
	"recipe_steps",
	"recipe_ingredients",
	"recipe_ingredient_substitutes",
	"recipe_taxonomy",
	"recipe_equipment",
	"recipe_photos",
}

var photoURL = regexp.MustCompile(`/api/photos/(\d+)$`)

func main() {
	from := flag.String("from", "", "DSN bazy z odtworzoną kopią (źródło)")
	to := flag.String("to", os.Getenv("DB_DATABASE_URL"), "DSN bazy docelowej (domyślnie DB_DATABASE_URL)")
	id := flag.Uint("id", 0, "id przepisu")
	apply := flag.Bool("apply", false, "zapisz zmiany (bez tej flagi: podgląd i ROLLBACK)")
	flag.Parse()

	if *from == "" || *to == "" || *id == 0 {
		flag.Usage()
		os.Exit(2)
	}
	if *from == *to {
		log.Fatal("-from i -to wskazują tę samą bazę")
	}

	src := open(*from)
	dst := open(*to)
	fmt.Printf("źródło (kopia): %s\ncel:            %s\n\n", hostOf(*from), hostOf(*to))

	var name string
	if err := src.Raw(`SELECT name FROM recipes WHERE id = ?`, *id).Scan(&name).Error; err != nil || name == "" {
		log.Fatalf("w kopii nie ma przepisu id=%d", *id)
	}
	fmt.Printf("przepis w kopii:  #%d %q\n", *id, name)

	tx := dst.Begin()
	if tx.Error != nil {
		log.Fatalf("transakcja: %v", tx.Error)
	}
	defer tx.Rollback() // po Commit nic nie robi

	var current string
	tx.Raw(`SELECT name FROM recipes WHERE id = ?`, *id).Scan(&current)
	if current == "" {
		fmt.Println("przepis w celu:   brak (zostanie wstawiony)")
	} else {
		fmt.Printf("przepis w celu:   %q (zostanie nadpisany)\n", current)
	}
	fmt.Println()

	must(upsertRecipe(src, tx, *id))

	for i := len(childTables) - 1; i >= 0; i-- {
		t := childTables[i]
		if !tableExists(tx, t) {
			continue
		}
		must(tx.Exec(fmt.Sprintf(`DELETE FROM public.%s WHERE recipe_id = ?`, quote(t)), *id).Error)
	}

	for _, t := range childTables {
		if !tableExists(src, t) {
			fmt.Printf("  %-30s brak w kopii, pomijam\n", t)
			continue
		}
		if !tableExists(tx, t) {
			fmt.Printf("  %-30s brak w celu, pomijam\n", t)
			continue
		}
		rows := sourceRows(src, t, "recipe_id", *id)
		must(insertRows(src, tx, t, rows))
		fmt.Printf("  %-30s %d\n", t, len(rows))

		if t == "recipe_photos" {
			for _, r := range rows {
				must(copyPhoto(src, tx, r))
			}
		}
	}

	if !*apply {
		fmt.Println("\nPODGLĄD — nic nie zapisano (ROLLBACK). Żeby zapisać, dodaj -apply.")
		return
	}
	if err := tx.Commit().Error; err != nil {
		log.Fatalf("commit: %v", err)
	}
	fmt.Println("\nZapisano.")
}

func open(dsn string) *gorm.DB {
	// Cichy logger: SQL z danymi (np. bajtami zdjęć) nie trafia na ekran.
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Fatalf("połączenie z %s: %v", hostOf(dsn), err)
	}
	return db
}

func hostOf(dsn string) string {
	u, err := url.Parse(dsn)
	if err != nil || u.Host == "" {
		return "(nieczytelny DSN)"
	}
	return u.Host + u.Path
}

func must(err error) {
	if err != nil {
		if strings.Contains(err.Error(), "foreign key") {
			log.Fatalf("%v\n\nW bazie docelowej brakuje czegoś, do czego odwołuje się przepis "+
				"(składnik, jednostka, tag, sprzęt albo podlinkowany przepis). Nic nie zapisano.", err)
		}
		log.Fatalf("%v\n\nNic nie zapisano.", err)
	}
}

func quote(ident string) string { return `"` + strings.ReplaceAll(ident, `"`, `""`) + `"` }

func tableExists(db *gorm.DB, table string) bool {
	var n int64
	db.Raw(`SELECT count(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = ?`, table).Scan(&n)
	return n > 0
}

func columns(db *gorm.DB, table string) []string {
	var cols []string
	db.Raw(`SELECT column_name FROM information_schema.columns
		WHERE table_schema = 'public' AND table_name = ? ORDER BY ordinal_position`, table).Scan(&cols)
	return cols
}

// commonColumns — kolumny obecne w obu bazach, w kolejności z bazy docelowej.
func commonColumns(src, dst *gorm.DB, table string) []string {
	have := map[string]bool{}
	for _, c := range columns(src, table) {
		have[c] = true
	}
	var out []string
	for _, c := range columns(dst, table) {
		if have[c] {
			out = append(out, c)
		}
	}
	return out
}

// sourceRows zwraca wiersze jako JSON (row_to_json), żeby przenieść je bez znajomości typów kolumn.
func sourceRows(src *gorm.DB, table, key string, val uint) []string {
	var rows []string
	err := src.Raw(fmt.Sprintf(`SELECT row_to_json(t)::text FROM public.%s t WHERE %s = ?`, quote(table), quote(key)), val).
		Scan(&rows).Error
	must(err)
	return rows
}

func insertSQL(table string, cols []string, conflict string) string {
	q := make([]string, len(cols))
	for i, c := range cols {
		q[i] = quote(c)
	}
	list := strings.Join(q, ", ")
	return fmt.Sprintf(`INSERT INTO public.%[1]s (%[2]s) SELECT %[2]s FROM json_populate_record(NULL::public.%[1]s, ?::json) %[3]s`,
		quote(table), list, conflict)
}

func insertRows(src, dst *gorm.DB, table string, rows []string) error {
	if len(rows) == 0 {
		return nil
	}
	sql := insertSQL(table, commonColumns(src, dst, table), "")
	for _, r := range rows {
		if err := dst.Exec(sql, r).Error; err != nil {
			return fmt.Errorf("%s: %w", table, err)
		}
	}
	return nil
}

func upsertRecipe(src, dst *gorm.DB, id uint) error {
	rows := sourceRows(src, "recipes", "id", id)
	cols := commonColumns(src, dst, "recipes")
	var set []string
	for _, c := range cols {
		if c != "id" {
			set = append(set, fmt.Sprintf("%s = EXCLUDED.%s", quote(c), quote(c)))
		}
	}
	sql := insertSQL("recipes", cols, "ON CONFLICT (id) DO UPDATE SET "+strings.Join(set, ", "))
	if err := dst.Exec(sql, rows[0]).Error; err != nil {
		return fmt.Errorf("recipes: %w", err)
	}
	fmt.Printf("  %-30s 1\n", "recipes")
	return nil
}

// copyPhoto dokłada zdjęcie z kopii, jeśli link przepisu wskazuje na /api/photos/<id>, którego w celu nie ma.
func copyPhoto(src, dst *gorm.DB, recipePhotoJSON string) error {
	var rp struct {
		ImageURL string `json:"image_url"`
	}
	if err := json.Unmarshal([]byte(recipePhotoJSON), &rp); err != nil {
		return fmt.Errorf("recipe_photos: %w", err)
	}
	m := photoURL.FindStringSubmatch(rp.ImageURL)
	if m == nil {
		return nil
	}
	pid, _ := strconv.ParseUint(m[1], 10, 64)
	var n int64
	dst.Raw(`SELECT count(*) FROM photos WHERE id = ?`, pid).Scan(&n)
	if n > 0 {
		fmt.Printf("  %-30s #%d już jest\n", "photos", pid)
		return nil
	}
	rows := sourceRows(src, "photos", "id", uint(pid))
	if len(rows) == 0 {
		fmt.Printf("  %-30s #%d brak także w kopii\n", "photos", pid)
		return nil
	}
	if err := insertRows(src, dst, "photos", rows); err != nil {
		return err
	}
	fmt.Printf("  %-30s #%d dołożone\n", "photos", pid)
	return nil
}
