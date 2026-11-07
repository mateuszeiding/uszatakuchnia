# Uszata Kuchnia — instrukcje dla Claude

## Stack

- **Backend:** Go, GORM, Vercel Serverless Functions (`api/**/*.go`)
- **Frontend:** Vue 3 + TypeScript, SCSS (DS v2), Vite, Auth0
- **DB:** PostgreSQL (Neon), connection string w `DB_DATABASE_URL`

## Zasady kodu — Frontend (Vue)

### Formatowanie
Zawsze uruchamiaj ESLint po zmianach w plikach `.vue` / `.ts`:
```
cd web && npx eslint --fix src/...
```
Projekt używa Prettier przez `eslint-plugin-prettier`. Nie commituj kodu bez przejścia przez linter.

### Struktura plików
- Widoki z logiką > ~150 linii rozbijaj na podkomponenty w katalogu `components/` obok widoku.
- Przykład: `RecipesView.vue` → `components/FilterSidebar.vue`, `components/RecipeCard.vue`
- Nie twórz mega-plików z całą logiką w jednym miejscu.

### Style
- **Zero inline `style="..."`** — cała prezentacja idzie do `<style scoped>` lub klas DS.
- Używaj CSS custom properties (`var(--*)`), klas DS (`btn`, `badge`, `card`, `input`, `chip`, itd.).
- Nie używaj starych klas utility (`d-flex`, `bg-*`, `tx-*`, `col-*`, `border-color-*`).
- Patrz: `web/src/scss/config/_components.scss` i `_tokens.scss`.

### Formularze
- Każdy `<input>` / `<textarea>` / `<select>` musi mieć powiązany `<label>`:
  - przez atrybut `for` + `id` na inpucie, **lub**
  - przez owinięcie inputa wewnątrz `<label>`.
- `field-label` to tylko wizualny styl — samo `class="field-label"` na `<div>` nie wystarczy.

### DS (Design System)
- Tokeny: `web/src/scss/config/_tokens.scss`
- Komponenty: `web/src/scss/config/_components.scss`
- Designy źródłowe: `.design/` w roocie repo (handoff, decyzje, prototypy JSX). To jedyne źródło prawdy — `design/` i kopie w `~/.claude/...` są nieaktualne.

## Wdrażanie designu

Design przygotowuje Claude Design w `.design/` (`Handoff - Claude Code.md`, `Decyzje i backlog.md`, prototypy `*.jsx`).
- Każdy punkt **checklisty odbioru** z handoffu = zadanie w planie. Raport końcowy odhacza checklistę: zrobione / pominięte + powód. Niedokończonego zadania nie zostawiaj bez słowa.
- Prototyp to referencja wyglądu i zachowania, nie kod do kopiowania — ale funkcje wskazane w handoffie jako „1:1” przenoś 1:1.
- Jeśli design koliduje z istniejącym zachowaniem albo decyzją z `Decyzje i backlog.md` (np. przyklejony panel, kolejność sekcji) — nie rozwiązuj po cichu. Opisz kolizję i zapytaj.
- Decyzje podjęte samodzielnie (czego design nie rozstrzygnął) wypisz na końcu raportu jako „do potwierdzenia”.
- Przed commitem UI porównaj zrzut ekranu z artboardem: light, dark, 375 px.
- Tokeny: `_tokens.scss` jest źródłem prawdy. Jeśli prototyp używa innej wartości, bierz token i zgłoś rozjazd.

## Zasady kodu — Backend (Go)

- Każda funkcja API to oddzielny handler w `api/**/*.go`.
- Encje DB: `db/entities/`, mappery: `db/mappers/`, DTOs: `dtos/`.
- GORM AutoMigrate — nie piszesz ręcznych migracji SQL.

## Zmiany schematu bazy danych

Serwer odpala migracje (AutoMigrate + `db/migrations`) przy pierwszym połączeniu z bazą (`db/connection.go` → `initDB`), pod advisory lockiem. Lokalnie po zmianie encji i tak uruchom migrację ręcznie, żeby od razu zobaczyć błędy.

### Dodanie kolumn / tabel (bez utraty danych)
```
go run ./cmd/migrate
```
Tylko dodaje brakujące kolumny i tabele. Używaj zawsze gdy zmieniasz encje.

### Zmiany, które przenoszą / przekształcają / usuwają dane
AutoMigrate tego nie zrobi (stara kolumna zostaje, nowa jest pusta). Każda taka zmiana = wersjonowana migracja w `db/migrations/`:
- plik `NNNN_opis.go`, funkcja dopisana na końcu listy `migrations.All`,
- wykonuje się raz (tabela `schema_migrations`), w transakcji, **po** AutoMigrate,
- musi działać też na świeżej bazie (resetdb) — sprawdzaj `HasColumn` / `HasTable`,
- wdrożonej migracji nie edytujesz — poprawka = nowa migracja.

Dopóki w bazie są tylko dane z seeda, zamiast migracji wystarczy `resetdb`. Od pierwszych prawdziwych danych — zawsze migracja.

### Reset lokalnej bazy (niszczy dane)
```
go run ./cmd/resetdb
```
Drop → migrate → seed. Tylko lokalnie, nigdy na produkcji.

Oba skrypty wymagają `DB_DATABASE_URL` z `.env.local`.

## Uruchamianie lokalnie

```
cd web && pnpm dev          # frontend
vercel dev                  # backend (Vercel CLI)
```

Zmienne środowiskowe: `.env.local` w rocie projektu (`DB_DATABASE_URL`, `AUTH0_*`, itd.).
