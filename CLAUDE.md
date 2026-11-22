# Uszata Kuchnia — instrukcje dla Claude

## Stack

- **Backend:** Go, GORM, Vercel Serverless — jedna funkcja (`api/index.go`) z routerem na `http.ServeMux`
- **Frontend:** Vue 3 + TypeScript, SCSS (DS v2), Vite, Auth0
- **DB:** PostgreSQL (Neon), connection string w `DB_DATABASE_URL`

## Język

- **Po angielsku:** kod (nazwy, komentarze), commit messages, tytuły i opisy PR, odpowiedzi na review w PR.
- **Po polsku:** dokumenty dla nas — `docs/`, `CLAUDE.md`, briefy, prompty, notatki — oraz teksty w UI i treść przepisów/wyzwań (`content/`).
- Gałąź PR z polskimi commitami → squash merge z angielskim komunikatem. Historii `main` nie przepisujemy.

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
- Designy źródłowe: projekt Claude Design — artefakt https://claude.ai/artifact/YEuK2L22otH24hdvcUC5GR (handoff, decyzje, prototypy JSX). To jedyne źródło prawdy; w repo nie trzymamy kopii (`.design/`, `design/` usunięte 18.11.2026), a eksporty w `~/.claude/...` są nieaktualne.

## Wdrażanie designu

Design przygotowuje Claude Design w artefakcie powyżej, pliki pod `project/`: `Handoff - Claude Code.md`, `Decyzje i backlog.md`, ekrany `*.html` → prototypy `*-main.jsx` (np. `Wskazowki.html` → `wk-main.jsx`).
- Czytaj przez narzędzie Artifact: `action: "read"`, `paths: ["project/…"]` (lista plików: `action: "list"`, `scope: "files"`). Zawsze bieżąca wersja — nie pracuj na zapamiętanych kopiach. Tylko odczyt: nie publikuj do tego artefaktu i nie uruchamiaj `/design-sync`.
- Każdy punkt **checklisty odbioru** z handoffu = zadanie w planie. Raport końcowy odhacza checklistę: zrobione / pominięte + powód. Niedokończonego zadania nie zostawiaj bez słowa.
- Prototyp to referencja wyglądu i zachowania, nie kod do kopiowania — ale funkcje wskazane w handoffie jako „1:1” przenoś 1:1.
- Jeśli design koliduje z istniejącym zachowaniem albo decyzją z `Decyzje i backlog.md` (np. przyklejony panel, kolejność sekcji) — nie rozwiązuj po cichu. Opisz kolizję i zapytaj.
- Decyzje podjęte samodzielnie (czego design nie rozstrzygnął) wypisz na końcu raportu jako „do potwierdzenia”.
- Przed commitem UI porównaj zrzut ekranu z artboardem: light, dark, 375 px.
- Tokeny: `_tokens.scss` jest źródłem prawdy. Jeśli prototyp używa innej wartości, bierz token i zgłoś rozjazd.

## Zasady kodu — Backend (Go)

- Całe API to **jedna funkcja Vercela**: `api/index.go` (`resp.Begin`, potem router). `vercel.json` kieruje do niej całe `/api/*` jednym rewrite'em. Limitu 12 funkcji na wdrożenie już nie dotykamy.
- **Nowy endpoint:** handler (jedna metoda jednej trasy) w `handlers/` + jedna linijka w tabeli `router/routes.go`, np. `{http.MethodPut, "/api/foo/{id}", h.FooUpdate}`. Nie dodawaj plików do `api/` (każdy plik tam to osobna funkcja; test tego pilnuje), nie ruszaj `vercel.json`, nie używaj `?kind=`.
  - `{id}` w ścieżce przepuszcza tylko liczbę ≥ 1 (inaczej 404); handler czyta go przez `pathID(r)` (`r.PathValue("id")`). Parametry typu `?moveTo=` zostają w query.
  - Router sam odpowiada 404 (nieznana ścieżka) i 405 (zła metoda, z `Allow`) w JSON — handler nie sprawdza metody.
  - Handler **nie woła** `resp.Begin` (robi to raz `api/index.go`).
- Wyjątek: `GET /recipes/{id}` i `/recipes/{id}/zakupy` (`handlers/share.go`) to strony dla botów podglądu linków (meta OG). `vercel.json` ma dla nich dwa rewrite'y z `has: user-agent`, a `router_test` pilnuje, żeby poza `/api/` były tylko takie trasy.
- Vercel buduje tylko `main` (`ignoreCommand` w `vercel.json`, fail-closed: build tylko gdy `VERCEL_GIT_COMMIT_REF` = `main`; pusty albo inny ref — także wdrożenie z CLI bez gita — jest pomijany): gałęzie i PR-y nie dostają wdrożeń podglądowych, bo podgląd łączyłby się z produkcyjną bazą i odpalał na niej migracje z niezatwierdzonego kodu. Włączenie podglądów = najpierw osobna baza dla Preview (gałąź Neon w `DB_DATABASE_URL`).
- Router to zwykły `http.Handler` — poza Vercelem wystarczy `http.ListenAndServe(addr, http.HandlerFunc(api.Handler))`.
- Encje DB: `db/entities/`, mappery: `db/mappers/`, DTOs: `dtos/`.
- Połączenie z bazą tylko przez `db.Open` (API i `cmd/*`; test pilnuje): pgx w trybie `exec`, bez cache'u przygotowanych zapytań (pooler Neona). Skutek: `Value()` typu zapisywanego do `jsonb` zwraca `string`, nie `[]byte` (`[]byte` idzie jako `bytea`).
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

**Od 14.11.2026 w bazie są prawdziwe dane — każda zmiana danych/schematu idzie przez migrację. `resetdb` jako skrót odpada.**

### Reset bazy (niszczy dane) — NIE URUCHAMIAĆ
```
go run ./cmd/resetdb
```
Drop → migrate → seed. `.env.local` wskazuje na **produkcyjną** bazę Neon, więc ten skrypt wyczyściłby produkcję. Nie odpalaj go bez osobnej, lokalnej bazy podstawionej w `DB_DATABASE_URL`.

Oba skrypty wymagają `DB_DATABASE_URL` z `.env.local`.

### Kopia zapasowa bazy
- Codziennie (cron ok. 02:17 UTC) workflow `.github/workflows/db-backup.yml` robi `pg_dump` całej bazy, szyfruje gpg (hasło w sekrecie `BACKUP_PASSPHRASE`) i zapisuje jako artefakt na 90 dni: GitHub → Actions → **Kopia bazy** → przebieg → Artifacts.
- **Przed ryzykowną migracją albo operacją na danych produkcji** uruchom kopię ręcznie (Actions → Kopia bazy → Run workflow) albo zrób snapshot / gałąź w Neon Console — i powiedz o tym właścicielowi, zanim ruszysz dane.
- Odtwarzanie (cała baza, jeden przepis przez `go run ./cmd/restore-recipe`, gdzie jest hasło): `docs/backup-bazy.md`. Testuj na lokalnej bazie, nigdy na produkcji.
- Po zmianie wersji głównej Postgresa na Neonie podbij `PG_MAJOR` w workflow.
- `content/` to nie backup — przepisy w stanie z dnia dodania.

## Autoryzacja API

Odczyt opublikowanych treści (lista, przepis, lista zakupów, zdjęcia, słowniki) jest publiczny. **Każdy zapis** (POST / PUT / PATCH / DELETE) i podgląd szkiców (`GET /api/recipes/list?admin=true`, `GET` szczegółów szkicu) wymaga uprawnień właściciela w nagłówku `Authorization: Bearer …`:
- **przeglądarka** — token dostępu Auth0 (JWT, audience `VITE_AUTH0_AUDIENCE`); backend sprawdza podpis (JWKS z Auth0), wystawcę, audience i `sub` = `VITE_AUTH0_OWNER_SUB`. Inne zalogowane konto dostaje 403 i nie widzi UI admina,
- **sesje Claude Code** — klucz `API_ADMIN_KEY` (min. 32 znaki) z `.env.local`; ten sam klucz jest w zmiennych Vercela.

Odpowiedzi: 401 bez tokenu / ze złym kluczem / z wygasłym tokenem, 403 z tokenem cudzego konta, 404 dla szkicu bez uprawnień.

Kod:
- `auth/` — sprawdzenie klucza i JWT (cache kluczy JWKS w pamięci funkcji),
- `http/begin.go` — `resp.Begin(w, r)`: CORS + preflight + bramka zapisu. Woła go **raz jedyna funkcja `api/index.go`, przed routerem** — więc przed każdą trasą, także przed 404 i 405. Nowy endpoint w `router/routes.go` jest chroniony bez dodatkowego kodu. Odczyt, który ma pokazać więcej właścicielowi, sprawdza `resp.IsOwner(r)`. Testy (`go test ./http ./router`):
  - `http/begin_test.go` — w `api/` jest dokładnie jedna funkcja i zaczyna od `if resp.Begin(w, r) { return }`; nikt poza nią nie woła `Begin` (podwójne nagłówki CORS),
  - `router/router_test.go` — po **wszystkich** trasach z tabeli: zapis bez tokenu i ze złym kluczem → 401, preflight → 204 z CORS; Begin przed 404/405; każda trasa trafia do swojego handlera; `vercel.json` ma jeden rewrite API,
- `web/src/data/api/API.ts` — dokleja token do żądań; przy 401 pyta o ponowne logowanie w okienku (`web/src/auth/reauthenticate.ts`) i powtarza żądanie, więc formularz nie traci treści.

CORS: produkcja i podglądy Vercela są z tej samej domeny co API, więc CORS wpuszcza tylko lokalny dev (`localhost:5173`, `:4173`). Dodatkowe domeny: `CORS_ORIGINS="https://a,https://b"`. curl i skrypty CORS nie dotyczy.

Zmienne (Vercel: Production + Preview + Development; lokalnie `.env.local`, frontend dodatkowo `web/.env`):
- `API_ADMIN_KEY` — sekret, tylko Vercel i `.env.local`, **nigdy repo ani frontend**,
- `VITE_AUTH0_DOMAIN`, `VITE_AUTH0_AUDIENCE`, `VITE_AUTH0_OWNER_SUB` — czyta je frontend (build) i backend (runtime). Backend przyjmuje też `AUTH0_DOMAIN` / `AUTH0_AUDIENCE` / `AUTH0_OWNER_SUB`, jeśli są ustawione,
- `CORS_ORIGINS` — opcjonalnie.

Brak konfiguracji = zamknięte: bez `API_ADMIN_KEY` klucz nie działa, bez `VITE_AUTH0_OWNER_SUB` żadne konto nie jest właścicielem.

### Zapis przez API z sesji Claude Code

Klucz czytaj z `.env.local` do zmiennej i nie wypisuj go (repo jest publiczne, logi też bywają):
```bash
KEY=$(sed -n 's/^API_ADMIN_KEY="\{0,1\}\([^"]*\)"\{0,1\}$/\1/p' .env.local | tr -d '\r')
curl -s -X POST https://<domena>/api/ingredients \
  -H "Authorization: Bearer $KEY" -H 'Content-Type: application/json' -d @body.json
```
Ten sam nagłówek przy `POST` / `PUT /api/recipes` i `/api/challenges`, `GET /api/recipes/<id>` szkicu przed edycją i `GET /api/recipes/list?admin=true`.

- **401** — w `.env.local` nie ma `API_ADMIN_KEY` albo różni się od klucza w Vercelu. Nie obchodź tego (np. zapisem prosto do bazy) — powiedz właścicielowi; po jego zgodzie `vercel env pull .env.local` ściąga aktualne zmienne (nadpisuje plik).
- Klucz nie wygasa sam. **Rotacja** (np. po wycieku): właściciel generuje nowy (`[Convert]::ToBase64String([Security.Cryptography.RandomNumberGenerator]::GetBytes(32))` w PowerShellu), podmienia `API_ADMIN_KEY` w Vercelu (wszystkie środowiska), robi redeploy i podmienia wartość w `.env.local`. Stary klucz przestaje działać po redeployu.

Lokalnie: `go run ./cmd/...` łączy się z bazą bezpośrednio — autoryzacja API go nie dotyczy. `vercel dev` czyta `.env.local`, więc zapis z lokalnego frontendu wymaga tam `VITE_AUTH0_*` i `VITE_AUTH0_OWNER_SUB`, a curl — `API_ADMIN_KEY`.

## Przepisy i wyzwania dodawane przez Claude (backup)

Przepisy i wyzwania wprowadza się przez API produkcyjne z kluczem `API_ADMIN_KEY` — jak w „Zapis przez API z sesji Claude Code” wyżej.

Przepisy i wyzwania wprowadzane przez API (nie przez formularz) mają kopię w `content/`:
- `content/recipes/<slug>.json` — przepis z nazwami składników, jednostek i wartości taksonomii (nie ID), więc da się go odtworzyć na każdej bazie; `cautions` = lista uwag „czego nie zawalić” (krótkie teksty, maks. 220 znaków, `[]` gdy brak — to nie wskazówki do kroków),
- `content/ingredients/<data>-<temat>.json` — składniki założone przy okazji (makro na 100 g, przeliczniki),
- `content/glossary/<data>-<temat>.json` — pojęcia ze słowniczka założone przy okazji: `[{ "id", "name", "explanation" }]` (id z bazy źródłowej). W treści kroku i wskazówki odwołanie to `{{term:<id>|słowo}}` z id z bazy, na której przepis dodano. **Przed wprowadzeniem przepisu na inną bazę albo po cofnięciu bazy** id trzeba przemapować po nazwie — procedura: `docs/content-na-innej-bazie.md`,
- `content/challenges/<slug>.json` — wyzwanie: dokładnie body dla `POST /api/challenges`, z formatem jako nazwą (`"format": "na temat"` zamiast `formatId`) — API przyjmuje oba. Pola: `title`, `format`, `description` (akapity oddzielone `\n\n`), `difficulty` 1–5, `timeLimitMin` (null = bez limitu), `source` (null = nie pokazujemy), `ingredients` (`[{ "name", "required" }]`, tekst bez powiązania z bazą składników), `cautions` (lista uwag jak w przepisie, `[]` gdy brak). Wyzwanie z API jest od razu publiczne (bez szkicu).

Poprawki zrobione potem w aplikacji nie wracają do tych plików — to backup stanu z dnia dodania.

## Uruchamianie lokalnie

```
cd web && pnpm dev          # frontend
vercel dev                  # backend (Vercel CLI)
```

`vercel dev` podnosi API (jedną funkcję z routerem), ale czyta `.env.local` = produkcja — do testów zapisu go nie używaj.

Zmienne środowiskowe: `.env.local` w rocie projektu (`DB_DATABASE_URL`, `API_ADMIN_KEY`, `VITE_AUTH0_*`, itd. — patrz „Autoryzacja API”). Uwaga: baza z `.env.local` to produkcja.
