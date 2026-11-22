// Restores a single recipe from a backup into the target database (usually production).
//
// First restore the backup into a local database (docs/backup-bazy.md), then:
//
//	go run ./cmd/restore-recipe -from "<local backup DSN>" -id 12           # dry run, writes nothing
//	go run ./cmd/restore-recipe -from "<local backup DSN>" -id 12 -apply    # write
//
// Target database: the -to flag, or DB_DATABASE_URL without it.
//
// What it does (in one transaction): inserts or overwrites the `recipes` row (same id); deletes
// the recipe's steps, ingredients, substitutes, tags, equipment and photo link and inserts them from
// the backup; adds the photo (`photos`) if it's missing. Other recipes stay untouched. Columns are
// matched by name, so a backup taken before a schema change works too (new columns get NULL).
// Ingredients, units, tags and equipment must exist in the target database; if any is missing,
// the transaction rolls back with a foreign key error.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"uszatakuchnia/db"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Insert order; deletes run in reverse (substitutes before ingredients).
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
	defer tx.Rollback() // no-op after Commit

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

		// a backup from before migration 0011 has no row keys; give them the ones 0011 gives
		// (a target without the column, i.e. before 0011, has nothing to fill)
		if t == "recipe_ingredients" && slices.Contains(columns(tx, t), "key") {
			must(tx.Exec(`UPDATE recipe_ingredients SET "key" = 'i' || sort_order WHERE recipe_id = ? AND "key" = ''`, *id).Error)
		}

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
	// Quiet logger: SQL with data (e.g. photo bytes) doesn't hit the screen.
	conn, err := db.Open(dsn, &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Fatalf("połączenie z %s: %v", hostOf(dsn), err)
	}
	return conn
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

// commonColumns returns the columns present in both databases, in target order.
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

// sourceRows returns rows as JSON (row_to_json) so they can be copied without knowing column types.
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

// copyPhoto copies the photo from the backup if the recipe links to /api/photos/<id> that the target lacks.
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
