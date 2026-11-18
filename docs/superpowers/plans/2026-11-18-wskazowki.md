# Wskazówki przy krokach — plan implementacji

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Wskazówki (0–n) przy każdym kroku przepisu: dane + API, strona i tryb gotowania, formularz z przeniesieniem starej `recipes.tip`.

**Architecture:** `recipe_steps.hints` jako jsonb (`StringList`), zapisywane w całości przy upsercie kroków (`ReplaceRecipeChildren`). Front: wspólny `StepHints.vue` dla strony i trybu gotowania; edycja w `StepsEditor.vue`; stara `tip` jako „Zanim zaczniesz” (strona) i „Przenieś do kroku…” (formularz).

**Tech Stack:** Go + GORM (Postgres), Vue 3 + TS, vitest, ESLint/Prettier.

**Spec:** `docs/superpowers/specs/2026-11-18-wskazowki-design.md` + handoff „Wskazówki przy krokach” (`.design/Handoff - Claude Code.md`), prototyp `.design/wk-main.jsx` (artboardy A).

## Global Constraints

- Baza produkcyjna ma prawdziwe dane: tylko AutoMigrate + migracja wersjonowana, **nigdy `resetdb`**, testy na lokalnym Postgresie (embedded, port 5433), nie na `.env.local`.
- Max **300** znaków na wskazówkę; puste (po trim) pomijane po cichu.
- Zero inline `style`, tylko `var(--*)` i `<style scoped>`; każde pole z etykietą (`aria-label` dopuszczalny jak w prototypie).
- Commit na etap: (1) dane+API, (2) strona+tryb gotowania, (3) formularz. Przed commitem UI: zrzut light, dark, 375 px vs artboard A.
- Nie zmieniać: odhaczania kroków, „przygotuj wcześniej”, grup kroków.

## Review Focus

1. Przepis bez żadnych wskazówek (wszystkie istniejące na prod) → strona, tryb gotowania i formularz wyglądają dokładnie jak dziś (brak pustej linii/`border-top`). Test: `TestMapStepHintsNilIsEmpty` + zrzut.
2. Stary wiersz z `hints = NULL` (gdyby AutoMigrate dodał kolumnę bez wartości) → odczyt daje `[]`, nie błąd. Test: `TestStringListScanNil`.
3. Wskazówka z samymi spacjami albo pusty wiersz dodany „+ wskazówka” i zostawiony → nie trafia do bazy. Test: `TestCleanHints`, vitest `toRequest drops blank hints`.
4. Przeniesienie `tip` do kroku, potem „Anuluj” → w bazie nic się nie zmienia (zmiana tylko w stanie formularza). Sprawdzenie ręczne w Task 3.
5. Krok usunięty w formularzu, do którego właśnie przeniesiono `tip` → wskazówka znika razem z krokiem (świadomie, jak każda edycja przed zapisem). Bez testu, zgłoszone w raporcie.

---

### Task 1: Dane + API (`hints`)

**Files:**
- Create: `db/entities/string_list.go`, `db/entities/string_list_test.go`, `db/migrations/0006_recipe_step_hints.go`, `services/recipes_hints_test.go`
- Modify: `db/entities/recipe_step.go`, `db/migrations/migrations.go` (dopisz `recipeStepHints` na końcu `All`), `dtos/recipe_step_dto.go`, `dtos/recipe_dto.go` (`UpsertStepRequest`), `db/mappers/` (mapper kroku), `services/recipes.go` (`ValidateRecipe`, `ReplaceRecipeChildren`)

**Interfaces:**
- Produces: `entities.StringList []string` (`Value() (driver.Value, error)` → JSON, nil → `[]`; `Scan(src any) error`, nil → pusta lista); `RecipeStep.Hints StringList \`gorm:"type:jsonb;not null;default:'[]'"\``; `RecipeStepDto.Hints []string \`json:"hints"\`` (zawsze tablica, nigdy null); `UpsertStepRequest.Hints []string \`json:"hints"\``; `services.CleanHints(in []string) []string`.

- [ ] **Step 1: Testy (fail)**
  - `TestStringListRoundTrip`: `StringList{"a","b"}` → Value → Scan → równe.
  - `TestStringListScanNil`: `Scan(nil)` → `len==0`, bez błędu; `StringList(nil).Value()` → `"[]"`.
  - `TestCleanHints`: `CleanHints([]string{"  x ", "", "   ", "y"})` == `[]string{"x","y"}`; `CleanHints(nil)` → `[]string{}` (nie nil).
  - `TestValidateHintTooLong` (bez bazy, wydziel `validateSteps(req, add)` z `ValidateRecipe`): hint 301 runów → błąd pola `steps` z komunikatem „Krok N: wskazówka może mieć najwyżej 300 znaków”; 300 runów → brak błędu. Liczyć `utf8.RuneCountInString` (polskie znaki).
- [ ] **Step 2:** `go test ./db/entities/ ./services/` → FAIL (brak typów).
- [ ] **Step 3: Implementacja** — typ, pole encji, DTO, mapper (`Hints: []string(e.Hints)` z nil→`[]`), `CleanHints` w `ReplaceRecipeChildren` przy tworzeniu kroku, `validateSteps` w `ValidateRecipe` (na hintach po trim). Migracja `0006_recipe_step_hints`: jeśli `HasColumn("recipe_steps","hints")` → `UPDATE recipe_steps SET hints='[]'::jsonb WHERE hints IS NULL`.
- [ ] **Step 4:** `go test ./...` → PASS; `go vet ./...` czysto.
- [ ] **Step 5: Lokalna baza** — w scratchpadzie serwer Go z embedded-postgres (5433) + `db.Apply` (migrate + seed), handlery z `api/`, trasy z `vercel.json` (wg pamięci `feedback_local_test_db`). Sprawdź: kolumna `hints jsonb not null default '[]'`, migracja 0006 w `schema_migrations`, PUT przepisu z `hints: [" a ", "", "b"]` → GET zwraca `["a","b"]`; hint 301 znaków → 422.
- [ ] **Step 6: Commit** `feat: wskazówki przy krokach — dane i API (recipe_steps.hints)`

### Task 2: Strona przepisu + tryb gotowania

**Files:**
- Create: `web/src/views/recipies/details/components/StepHints.vue`, `web/src/shared/components/icons/BulbIcon.vue` (albo istniejący katalog ikon — sprawdź, gdzie leżą ikony SVG i trzymaj się go)
- Modify: `web/src/data/dtos/recipe/RecipeStepDto.ts` (`hints: string[] = []`), `RecipeDto.ts` (typ żądania kroku: `hints: string[]`), `StepsList.vue`, `CookMode.vue`

**Interfaces:**
- Consumes: `RecipeStepDto.hints: string[]` z API (Task 1).
- Produces: `<StepHints :hints="string[]" :done="boolean" size="md" | "lg" />` — renderuje nic, gdy `hints.length === 0`.

- [ ] **Step 1: `StepHints.vue`** wg handoffu pkt 2 i 4 / `wk-main.jsx` `Hints`: `<ul aria-label="Wskazówki">`, `border-top: 1px solid var(--rule)`; md: margin-top 14 px, padding-top 12 px, gap 8 px, kolumny 16 px + 1fr, gap 10 px, ikona 15 px (`padding-top: 3px`), tekst 15 px / 1.5 `--ink-dim`; lg: margin-top 20 px, padding-top 16 px, gap 12 px, kolumny 20 px + 1fr, gap 12 px, ikona 18 px (`padding-top: 5px`), tekst `clamp(18px, 2.2vw, 22px)`. Ikona `--accent`; `done` → ikona i tekst `--ink-muted`. Ikona żarówki: ścieżka z `BulbI` w `wk-main.jsx`, `aria-hidden`.
- [ ] **Step 2: `StepsList.vue`** — `<StepHints>` pod `.step__text` z `done`; usuń `<aside class="tip-box">` i style `.tip-box*`; nad `.steps` (po `steps-head`) blok „Zanim zaczniesz”, gdy `recipe?.tip?.trim()`: etykieta jak `.steps__group`, wiersz ikona + tekst jak `StepHints` md (bez `border-top`), ramka `1px var(--rule)`, tło `var(--bg-alt)`, `border-radius: var(--r-xl)`, padding 18/20 jak `.step`, margines dolny 28 px (= gap `.steps`). Element: `<aside aria-labelledby>` z etykietą jako `<h3>`.
- [ ] **Step 3: `CookMode.vue`** — `<StepHints :hints="step.hints" size="lg" />` pod `.cook__text`.
- [ ] **Step 4:** `cd web && npx eslint --fix <pliki> && npx vue-tsc --noEmit && pnpm test` → czysto/PASS.
- [ ] **Step 5: Zrzuty** (lokalna baza z Task 1, przepis z 2 krokami ze wskazówkami, 1 bez, `tip` niepusty; playwright-core + systemowy Chrome): strona light 1280, dark 1280, 375 px; krok odhaczony; tryb gotowania light/dark/375. Porównaj z artboardami A (`.design/Wskazowki.html`). Krok bez wskazówek = jak przed zmianą.
- [ ] **Step 6: Commit** `feat(web): wskazówki pod krokami na stronie przepisu i w trybie gotowania`

### Task 3: Formularz

**Files:**
- Create: `web/src/views/recipies/auth/components/TipMove.vue`, `web/src/views/recipies/auth/__tests__/recipeForm.test.ts`
- Modify: `recipeForm.ts`, `components/StepsEditor.vue`, `UpsertRecipe.vue` (usuń sekcję `n="04"` „Wskazówka”; `TipMove` w sekcji „Kroki” nad `StepsEditor`)

**Interfaces:**
- Consumes: `hints` w żądaniu/odpowiedzi (Task 1), `BulbIcon` (Task 2).
- Produces: `IStepItem.hints: string[]`; `moveTipToStep(f: IRecipeForm, stepId: number): void` (dopisuje `f.tip.trim()` na koniec `hints` kroku o `id`, ustawia `f.tip = ''`); `blankStep()` z `hints: []`.

- [ ] **Step 1: Testy vitest (fail)**
  - `fromDto maps hints`: krok z `hints: ['a']` → `IStepItem.hints` `['a']`.
  - `toRequest drops blank hints`: `hints: [' a ', '', '  ']` → `['a']`.
  - `moveTipToStep appends and clears tip`: krok z `hints: ['x']`, `tip: ' rada '` → `['x','rada']`, `tip === ''`.
  - `validate flags hint over 300 chars`: hint 301 znaków → `errors.steps` ustawione.
- [ ] **Step 2:** `pnpm test` → FAIL.
- [ ] **Step 3: `recipeForm.ts`** — pole, mapowanie, `moveTipToStep`, walidacja (komunikat „Krok N: wskazówka może mieć najwyżej 300 znaków”, N = numer kroku jak w `toRequest`).
- [ ] **Step 4: `StepsEditor.vue`** wg `wk-main.jsx` `FormInner`: pod textarea kroku wiersze `[ikona 15 px --accent | input .input (13 px) | ✕ 30×30]`, `aria-label="Krok N, wskazówka M"`, `maxlength="300"`, placeholder „np. Jeśli ciasto się klei, odstaw je na 10 minut”, ✕ z `aria-label="Usuń wskazówkę M z kroku N"` bez potwierdzenia; pod nimi tekstowy przycisk „+ wskazówka” (13 px, 600, `--accent`) w linii z istniejącym „Przygotuj wcześniej”, nowa wskazówka dostaje fokus. N = numeracja kroków w edytorze. Jeśli plik przekroczy ~300 linii, wydziel `StepHintsEditor.vue`.
- [ ] **Step 5: `TipMove.vue`** (gdy `form.tip.trim()`): karta nad krokami z etykietą „Stara wskazówka” (styl etykiety pola), treść tip, `<label>` + `<select class="select">` „Przenieś do kroku…” z opcjami „Krok N — początek treści (40 znaków)…”; wybór → `moveTipToStep`, select znika razem z blokiem.
- [ ] **Step 6:** eslint + vue-tsc + `pnpm test` → PASS.
- [ ] **Step 7: Ręcznie na lokalnej bazie:** dodaj/edytuj/usuń wskazówkę → Zapisz → GET ma zmiany; „Przenieś do kroku…” → Zapisz → `tip` null, wskazówka w kroku, strona bez „Zanim zaczniesz”; przenieś → Anuluj → w bazie bez zmian. Zrzuty light/dark (+375 px dla porządku).
- [ ] **Step 8: Commit** `feat(web): formularz — wskazówki przy krokach i przeniesienie starej wskazówki`

### Task 4: Raport
- [ ] Checklista odbioru z handoffu odhaczona (zrobione / pominięte + powód), decyzje „do potwierdzenia” (ze specu + nowe), co niesprawdzone.
