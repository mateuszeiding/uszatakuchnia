# Makroskładniki — plan implementacji

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Orientacyjne makro (zakresy min–max) na składnikach i wyliczane makro przepisu na porcję.

**Architecture:** Składnik trzyma 10 nullable kolumn makro na 100 g i listę przeliczników jednostek. Czysta logika (dziedziczenie z drzewka, przeliczanie, sumowanie zakresów) w `services/nutrition.go`, testowana jednostkowo. Wynik dla przepisu cache'owany w `recipes.nutrition` (jsonb) i przeliczany przy zapisie przepisu oraz każdej zmianie składnika. Front tylko formatuje (zaokrąglenia z prototypu).

**Tech Stack:** Go + GORM (Postgres), Vue 3 + TS, SCSS DS v2.

**Spec:** `docs/superpowers/specs/2026-11-07-makro-skladnikow-design.md` + `.design/Decyzje i backlog.md` („Makroskładniki”) + `.design/Handoff - Claude Code.md` („Makroskładniki — prompt dla Claude Code”). Prototyp UI: `.design/Makro.html` → `.design/mk-main.jsx`.

## Global Constraints

- Pięć wartości, zawsze w tej kolejności i z tymi kluczami JSON: `kcal`, `protein`, `fat`, `carbs`, `fiber`. Etykiety: Kalorie / Białko / Tłuszcz / Węglowodany / Błonnik.
- Wartości składnika na 100 g. `max` null = jedna wartość.
- Format: `20 – 22 g (~21)` albo `21 g`; przecinek dziesiętny; liczby w Geist Mono (`var(--font-mono)` lub odpowiednik z `_tokens.scss`).
- Zaokrąglenia — dane składnika: kcal co 1, g co 0,1. Przepis: kcal co 10 (poniżej 100 co 5), g co 1 (poniżej 10 co 0,5). Gdy po zaokrągleniu min = max → jedna liczba.
- Walidacja: liczby ≥ 0; max bez min = błąd; max < min = „Max nie może być mniejsze niż min”.
- Jednostki bazowe (po nazwie z `units`): `g` = 1, `kg` = 1000, `ml` = 1, `l` = 1000 gramów. Własny/odziedziczony przelicznik nadpisuje `ml`/`l`; dla `g`/`kg` przelicznika nie da się dodać.
- Zero inline `style`, każdy input z `<label>` lub `aria-labelledby`, tylko tokeny DS, light + dark.
- Po zmianach w `.vue`/`.ts`: `cd web && npx eslint --fix <pliki>` i `pnpm type-check`.
- Commity na `main`, bez push.

## Review Focus

1. Składnik z częściowymi danymi (np. kcal jest, błonnika brak) — liczy się do tych pól, które ma; nie trafia do „Nie wliczono”. Test w Task 2.
2. Pozycja przepisu bez jednostki (`unit_id` NULL) z ilością „2” — traktowana jak jednostka `sztuka` (decyzja Claude'a). Test w Task 2.
3. Przepis-szkic z `servings` = 0 — makro liczone jak dla 1 porcji, bez dzielenia przez zero. Test w Task 2.
4. Ustawienie rodzica, który jest potomkiem składnika (cykl w drzewie) — błąd walidacji, nie nieskończona rekurencja. Test w Task 4; `ResolveMacros` dodatkowo chroni się przed cyklem (Task 2).
5. Zmiana makro podskładnika zmienia makro przepisów używających rodzica — przeliczenie obejmuje wszystkie przepisy. Weryfikacja ręczna w Task 4.

---

### Task 1: Schemat — encje, migracja, usunięcie `KcalPerServing`

**Files:**
- Create: `db/entities/nutrition.go`, `db/entities/ingredient_unit_conversion.go`, `db/migrations/0002_drop_recipe_kcal.go`
- Modify: `db/entities/ingredient.go`, `db/entities/recipe.go`, `db/migrate.go`, `db/migrations/migrations.go`, `dtos/recipe_dto.go`, `services/recipes.go:ApplyRecipeRequest`, `db/mappers/recipe_map.go`
- Test: `db/entities/nutrition_test.go`

**Interfaces:**
- Produces:
  - `entities.MacroRange { Min float64; Max *float64 }` (json `min`, `max`) — wartość na 100 g lub w przeliczniku.
  - `entities.SumRange { Min, Max float64 }` (json `min`, `max`) — wartość wyliczona.
  - `entities.RecipeNutrition { PerPortion map[string]SumRange `json:"perPortion"`; Missing []string `json:"missing"` }` z `Scan`/`Value` (jsonb); `*RecipeNutrition` nil ↔ NULL.
  - `entities.NutrientKeys = []string{"kcal","protein","fat","carbs","fiber"}`.
  - `Ingredient`: pola `KcalMin, KcalMax, ProteinMin, ProteinMax, FatMin, FatMax, CarbsMin, CarbsMax, FiberMin, FiberMax *float64` (`type:numeric`), `Conversions []IngredientUnitConversion`.
  - `IngredientUnitConversion { ID; IngredientID uint; UnitID uint; Unit Unit; GramsMin float64 (not null); GramsMax *float64 }`, tabela `ingredient_unit_conversions`, unikalny indeks `(ingredient_id, unit_id)`, FK cascade na składniku, restrict na jednostce.
  - `Recipe.Nutrition *RecipeNutrition` (`type:jsonb`); `RecipeBase.Nutrition *RecipeNutrition` (`->` tylko odczyt).
  - `Ingredient.MacroOf(key string) *MacroRange` — zwraca nil, gdy `*Min` jest nil.

- [ ] **Step 1: Test `TestRecipeNutritionScanValue`** — round-trip `Value()` → `Scan([]byte)` zwraca równą strukturę; `Scan(nil)` daje zero bez błędu. **`TestIngredientMacroOf`** — `KcalMin=nil` → nil; `FatMin=2, FatMax=nil` → `{2, nil}`.
- [ ] **Step 2:** `go test ./db/entities/` → FAIL (brak typów).
- [ ] **Step 3:** Zaimplementuj typy i pola. Usuń `KcalPerServing` z `Recipe`, `RecipeDto`, `UpsertRecipeRequest`, `ApplyRecipeRequest`, mappera. Dopisz `&entities.IngredientUnitConversion{}` do `AutoMigrate`. Migracja `0002_drop_recipe_kcal`: jeśli `HasColumn("recipes","kcal_per_serving")` → `DropColumn`; dopisz na koniec `migrations.All`.
- [ ] **Step 4:** `go test ./... && go vet ./...` → PASS; `go run ./cmd/migrate` → bez błędów.
- [ ] **Step 5:** Commit `feat(db): makro składników i przeliczniki jednostek — schemat`.

### Task 2: Logika makro (czyste funkcje)

**Files:**
- Create: `services/nutrition.go`, `services/nutrition_test.go`

**Interfaces:**
- Consumes: typy z Task 1, `services.Quantity`.
- Produces:
  - `type Macros map[string]entities.MacroRange` — brak klucza = brak danych.
  - `type NutritionIndex struct { Own map[uint]Macros; Children map[uint][]uint; Parent map[uint]*uint; Conversions map[uint]map[uint]entities.MacroRange; Names map[uint]string; UnitIDByName map[string]uint }`
  - `func (ix *NutritionIndex) ResolveMacros(id uint) Macros` — per klucz: własny, inaczej min z `Min` i max z `max(Max ?? Min)` po rozwiązanych dzieciach (rekurencyjnie, z ochroną przed cyklem — odwiedzone węzły pomijane). Wynik zawsze ma `Max` ustawione, gdy dzieci dały różne wartości; gdy min == max → `Max` nil.
  - `func (ix *NutritionIndex) ResolveConversion(id, unitID uint) (entities.MacroRange, bool)` — własny, inaczej u przodków.
  - `func (ix *NutritionIndex) GramsPerUnit(id uint, unitID *uint, unitName string) (entities.MacroRange, bool)` — `unitID` nil → szukaj jednostki `sztuka`; dla `g`/`kg` stała; dla pozostałych przelicznik, a dla `ml`/`l` fallback na stałą.
  - `type NutritionLine struct { IngredientID uint; Qty Quantity; UnitID *uint; UnitName string; NoQty, Optional bool }`
  - `func (ix *NutritionIndex) Compute(lines []NutritionLine, servings int) *entities.RecipeNutrition` — algorytm z handoffu („3. Wyliczenie dla przepisu”); `servings < 1` → 1; zwraca nil, gdy żadna pozycja nic nie wniosła. `Missing` w kolejności pozycji, bez duplikatów.

- [ ] **Step 1: Testy** (dane w budowanym ręcznie `NutritionIndex`):
  - `TestResolveMacrosInherits` — rodzic bez danych, dzieci kcal `{350,360}` i `{330,345}` → `{330,360}`; dziecko `{2,nil}` liczone jako 2–2.
  - `TestResolveMacrosOwnWins` — rodzic z własnym kcal ignoruje dzieci tylko dla kcal; brakujące `fiber` dziedziczy.
  - `TestResolveMacrosCycle` — A→B→A nie zawiesza się.
  - `TestGramsPerUnit` — `g` → 1; `kg` → 1000; `ml` bez przelicznika → 1; `ml` z przelicznikiem 0,92 → 0,92; `łyżka` u rodzica → dziedziczona; `unitID` nil → przelicznik `sztuka`; brak → false.
  - `TestComputeRanges` — 200 g × kcal `{350,360}` + 2 szt jajka (`sztuka` `{50,65}` g, kcal `{140,nil}`), servings 2 → kcal `{(700+140)/2, (720+182)/2}` = `{420, 451}`.
  - `TestComputeSkipsAndMissing` — „do smaku” i opcjonalny: pominięte, nie w `Missing`; pusta ilość: pominięta; składnik z ilością bez makro → `Missing`; bez przelicznika jednostki → `Missing`.
  - `TestComputePartialFields` — składnik z samym kcal wlicza się do kcal, nie trafia do `Missing`, `fiber` z innych pozycji.
  - `TestComputeServingsZero` — servings 0 = całość.
  - `TestComputeNothing` — same braki → nil.
- [ ] **Step 2:** `go test ./services/ -run Nutrition` → FAIL.
- [ ] **Step 3:** Implementacja.
- [ ] **Step 4:** `go test ./services/` → PASS.
- [ ] **Step 5:** Commit `feat(api): wyliczanie makro przepisu z zakresów`.

### Task 3: Przeliczanie i zapis w bazie + seed

**Files:**
- Create: `services/nutrition_store.go`
- Modify: `api/recipes/create.go`, `api/recipes/detailed.go` (put), `db/seed.go`, `db/mappers/recipe_map.go`, `dtos/recipe_dto.go`

**Interfaces:**
- Consumes: Task 2.
- Produces:
  - `func LoadNutritionIndex(tx *gorm.DB) (*NutritionIndex, error)` — wszystkie składniki, przeliczniki, jednostki (tabele są małe).
  - `func RecomputeRecipes(tx *gorm.DB, recipeIDs ...uint) error` — bez argumentów = wszystkie przepisy; zapisuje `recipes.nutrition` (`Select("nutrition").Updates`).
  - `RecipeDto.Nutrition *entities.RecipeNutrition` (json `nutrition`, null gdy brak); `RecipeBaseDto.Kcal *entities.SumRange` (json `kcal`, z `Nutrition.PerPortion["kcal"]`).

- [ ] **Step 1:** Implementuj store; w create/put wywołaj `RecomputeRecipes(tx, id)` w tej samej transakcji po `ReplaceRecipeChildren`. Mappery wypełniają `nutrition` / `kcal`.
- [ ] **Step 2:** Seed: makro na 100 g dla wszystkich seedowanych składników (wartości z tabel żywieniowych, zakres tam, gdzie marki się różnią — np. makarony, śmietana, boczek) + przeliczniki: `łyżka` oliwy ≈ 10 g, `łyżka` masła ≈ 15 g, `łyżeczka` przypraw ≈ 2–3 g, `ząbek` czosnku ≈ 4–6 g, `sztuka` cytryny ≈ 100–130 g, `sztuka` papryczki chili ≈ 10–20 g, `filet` łososia ≈ 120–180 g, `garść` natki ≈ 10–15 g, `pęczek` tymianku ≈ 20–30 g. Na końcu seeda `RecomputeRecipes(tx)`.
- [ ] **Step 3:** Weryfikacja: `go test ./... && go vet ./...`. **Zapytaj użytkownika** przed `go run ./cmd/resetdb` (niszczy dane); po zgodzie uruchom i sprawdź `GET /api/recipes/<id>` (przez `vercel dev`) — `nutrition.perPortion.kcal` ma sensowne wartości, `missing` zawiera tylko składniki bez danych.
- [ ] **Step 4:** Commit `feat(api): makro przepisu w odpowiedzi API, seed z wartościami`.

### Task 4: API składników — makro, przeliczniki, rodzic

**Files:**
- Create: `services/ingredients.go`, `services/ingredients_test.go`
- Modify: `api/ingredients/create.go`, `api/ingredients/detailed.go`, `api/ingredients/list.go`, `dtos/ingredient_dto.go`, `db/mappers/ingredient_map.go`

**Interfaces:**
- Consumes: Task 1–3.
- Produces:
  - Request (create i update): `{ name, typeCode, allergens, parentId: number|null, nutrition: { [key]: {min, max|null} }, conversions: [{unitId, gramsMin, gramsMax|null}] }`. Brak klucza w `nutrition` = brak własnych danych.
  - `func ValidateIngredientRequest(req dtos.UpsertIngredientRequest, ix *NutritionIndex, selfID uint) []dtos.FieldError` — pola `nutrition.<key>.min|max`, `conversions[i].unitId|gramsMin|gramsMax`, `parentId`. Reguły: Global Constraints + duplikat jednostki („Ta jednostka ma już przelicznik”), jednostka `g`/`kg`/„do smaku” odrzucona, rodzic nie może być sobą ani potomkiem („Składnik nie może być własnym podskładnikiem”). Błędy → 422 `{error, fields}` jak w przepisach.
  - `IngredientDto` dodatkowo: `nutrition` (własne, mapa), `resolvedNutrition` (z `ResolveMacros`), `conversions: [{unitId, unitName, gramsMin, gramsMax, inherited}]` (własne + odziedziczone z przodków dla jednostek, których składnik nie ma).
- [ ] **Step 1: Testy `ValidateIngredientRequest`:** max bez min; max < min (dokładny komunikat); ujemne; duplikat jednostki; `g` odrzucone; rodzic = self; rodzic = wnuk; poprawne dane → brak błędów.
- [ ] **Step 2:** `go test ./services/ -run Ingredient` → FAIL.
- [ ] **Step 3:** Implementacja walidacji; handlery zapisują kolumny makro, zastępują przeliczniki (delete + insert w transakcji), ustawiają `ParentID`, potem `RecomputeRecipes(tx)` (wszystkie — zmiana składnika wpływa na przodków i potomków). Lista preloaduje `Conversions.Unit` i liczy `resolvedNutrition` z jednego `LoadNutritionIndex`. Sprawdzanie duplikatu nazwy uwzględnia nowy `parentId`.
- [ ] **Step 4:** `go test ./... && go vet ./...` → PASS. Ręcznie (`vercel dev`): PUT podskładnika z nowym kcal → przepis używający rodzica ma zmienione `nutrition` (Review Focus 5).
- [ ] **Step 5:** Commit `feat(api): edycja makro, przeliczników i rodzica składnika`.

### Task 5: Front — DTO, formatowanie, komponent wartości

**Files:**
- Create: `web/src/shared/utils/nutrition.ts`, `web/src/components/MacroValue.vue`
- Modify: `web/src/data/dtos/recipe/RecipeDto.ts`, `web/src/data/dtos/ingredients/IngredientDto.ts`, `web/src/data/api/ingredients/fetch.ts`, `web/src/views/recipies/auth/recipeForm.ts` (usuń `kcal`), `web/src/views/recipies/details/components/RecipeHero.vue` (usuń stare kcal — nowe w Task 6), formularz przepisu (usuń pole kcal, jeśli jest)

**Interfaces:**
- Produces (`nutrition.ts`):
  - `type NutrientKey = 'kcal'|'protein'|'fat'|'carbs'|'fiber'`; `NUTRIENTS: {key, label, unit}[]` w kolejności z Global Constraints.
  - `type Range = { min: number; max: number | null }`; `type Formatted = { single: true; a: number } | { single: false; a: number; b: number; m: number }`.
  - `rnd(v, key, mode: '100g'|'recipe')`, `fmt(range, key, mode): Formatted|null`, `fmtText(f, unit)`, `spokenText(f, unit)` („od 20 do 22 g, średnio około 21”), `nf` (przecinek dziesiętny) — port 1:1 z `mk-main.jsx`.
  - `energySplit(perPortion): {protein, fat, carbs}` w % (B×4, T×9, W×4 ze średnich).
  - `MacroValue.vue` props `{ value: Range|null; nutrient: NutrientKey; mode: '100g'|'recipe'; size?: 'sm'|'md'|'lg' }` — renderuje jak `Val` w prototypie (sr-only tekst + `aria-hidden` wizualny, średnia przygaszona, „—” gdy brak).
  - `RecipeDto.nutrition: { perPortion: Partial<Record<NutrientKey, {min:number;max:number}>>; missing: string[] } | null`; `RecipeBaseDto.kcal: {min:number;max:number} | null`; `IngredientDto.nutrition`, `resolvedNutrition: Partial<Record<NutrientKey, Range>>`, `conversions: {unitId, unitName, gramsMin, gramsMax, inherited}[]`; `IIngredientPayload` rozszerzony jak request z Task 4.
- [ ] **Step 1:** Implementacja.
- [ ] **Step 2:** `cd web && npx eslint --fix <pliki> && pnpm type-check` → bez błędów.
- [ ] **Step 3:** Commit `feat(web): typy i formatowanie makro`.

### Task 6: Strona przepisu — karta „Wartości odżywcze” + lista

**Files:**
- Create: `web/src/views/recipies/details/components/NutritionCard.vue`
- Modify: `web/src/views/recipies/details/components/RecipeBody.vue`, `.../RecipeHero.vue`, `web/src/views/recipies/list/components/RecipeCard.vue`

- [ ] **Step 1:** `NutritionCard.vue` props `{ nutrition: RecipeDto['nutrition'] & {}; servings: number }` (`servings` = liczba wybrana licznikiem, ta sama, którą dostaje panel składników). Układ wg handoffu „4. UI → Strona przepisu” i `mk-main.jsx`: `id="makro"`, badge „orientacyjnie”, średnia kcal duża + zakres, wiersze B/T/W/błonnik, pasek „Skąd energia” (accent / `--c-mustard` / `--c-teal` + podpisy %), „Całe danie” = na porcję × `servings` (mniejsze), „Nie wliczono: …”. Mobile (≤900px): nagłówek to `button[aria-expanded]` z `~520 kcal/porcja`, domyślnie zwinięty.
- [ ] **Step 2:** `RecipeBody`: karta pod panelem składników (desktop w lewej kolumnie, mobile po składnikach), tylko gdy `recipe.nutrition`. `RecipeHero`: „Kcal / porcja: ~520” jako link `#makro`, ukryte bez danych.
- [ ] **Step 3:** `RecipeCard`: ` · ~520 kcal` (mono) w linii metadanych, gdy `recipe.kcal`.
- [ ] **Step 4:** eslint + type-check. Weryfikacja w przeglądarce (`pnpm dev` + `vercel dev`): przepis z pełnymi danymi, z brakami (dopisek), bez danych (sekcja i metadana ukryte); licznik porcji zmienia tylko „Całe danie”; mobile 375px zwija/rozwija; dark mode.
- [ ] **Step 5:** Commit `feat(web): wartości odżywcze na stronie przepisu i liście`.

### Task 7: Panel — edycja składnika jako strona

**Files:**
- Create: `web/src/views/admin/AdminIngredientEdit.vue`, `web/src/views/admin/components/IngredientBasics.vue`, `.../NutritionTable.vue`, `.../ConversionsTable.vue`
- Modify: `web/src/router/index.ts`, `web/src/views/admin/AdminIngredients.vue`

- [ ] **Step 1:** Trasy `admin/ingredients/new` i `admin/ingredients/:id` (props). `AdminIngredients` — „Dodaj” i „Edytuj” prowadzą do strony zamiast `IngredientModal` (modal zostaje dla szybkiego dodawania w formularzu przepisu). Dane z listy (`fetchIngredients`), po zapisie `invalidateIngredients` i powrót do listy z toastem.
- [ ] **Step 2:** Sekcje wg handoffu „Panel → Składniki”: `01 Podstawowe` (nazwa, rodzic — select z drzewa z wykluczeniem siebie i potomków, typ, alergeny); `02 Wartości odżywcze na 100 g` — tabela 5 × Min / Max (opcj.) / „Na stronie” (podgląd `MacroValue` mode `100g`); dla pustego wiersza z danymi odziedziczonymi: placeholder kursywą z `resolvedNutrition`, badge „odziedziczone”/„własne”, przycisk „przywróć dziedziczenie” czyści wiersz; baner z listą podskładników, gdy dziedziczy; `03 Przeliczniki jednostek` — wiersze jednostka / min g / max g / podgląd „1 szklanka ≈ 90 – 100 g” / usuń, „+ Dodaj przelicznik”, select bez `g`, `kg`, „do smaku”; odziedziczone przeliczniki pokazane jako tylko do odczytu z adnotacją.
- [ ] **Step 3:** Walidacja klienta = te same reguły co serwer; błędy z 422 mapowane po `field` na pola; `role="alert"` + `aria-describedby`; inputy w tabelach z `aria-labelledby` (nagłówek wiersza + kolumny); przecinek i kropka akceptowane. Ostrzeżenie przy wyjściu z niezapisanymi zmianami (jak w formularzu przepisu).
- [ ] **Step 4:** eslint + type-check; w przeglądarce: zapis makro i przelicznika, błąd max < min, cykl rodzica, przywrócenie dziedziczenia, potem strona przepisu pokazuje nowe wartości.
- [ ] **Step 5:** Commit `feat(web): panel — edycja makro i przeliczników składnika`.

### Task 8: Publiczna karta składnika

**Files:**
- Modify: `web/src/views/ingredients/components/IngredientCard.vue`

- [ ] **Step 1:** Makro na 100 g z `resolvedNutrition` (`MacroValue` mode `100g`) + przeliczniki jako `chip` („1 szklanka ≈ 90 – 100 g”); bez danych — nic.
- [ ] **Step 2:** eslint + type-check; przegląd w przeglądarce (light/dark, mobile).
- [ ] **Step 3:** Commit `feat(web): makro na karcie składnika`.
