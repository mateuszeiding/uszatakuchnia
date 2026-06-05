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
- Designy źródłowe (JSX): `C:\Users\eidin\.claude\projects\C--Users-eidin-Repos-uszatakuchnia\design-extracted\uszatakuchnia\project\`

## Zasady kodu — Backend (Go)

- Każda funkcja API to oddzielny handler w `api/**/*.go`.
- Encje DB: `db/entities/`, mappery: `db/mappers/`, DTOs: `dtos/`.
- GORM AutoMigrate — nie piszesz ręcznych migracji SQL.

## Zmiany schematu bazy danych

Serwer **nie** odpala migracji automatycznie. Po zmianie encji trzeba uruchomić migrację ręcznie.

### Dodanie kolumn / tabel (bez utraty danych)
```
go run ./cmd/migrate
```
Tylko dodaje brakujące kolumny i tabele. Używaj zawsze gdy zmieniasz encje.

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
