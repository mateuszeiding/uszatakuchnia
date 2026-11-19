# Prompt dla Claude Code — jedno wejście do API (router) zamiast 12 funkcji Vercela

## Kontekst

Uszata Kuchnia: backend w Go na Vercel Serverless (`api/**/*.go`). Na darmowym planie Vercela jedno wdrożenie może mieć **najwyżej 12 funkcji**, a każdy plik w `api/` to osobna funkcja. Wszystkie 12 są już zajęte.

Przez ten limit nowe endpointy są doklejane do istniejących plików przez `?kind=…` w `vercel.json`:
- sprzęt, wyzwania i formaty wyzwań obsługuje `api/taxonomy/units.go`,
- handlery wyzwań leżą w `handlers/challenges.go`.

To działa, ale każdy kolejny feature z nowym endpointem dokłada takie obejście.

Właściciel kiedyś przeniesie API z Vercela (kiedy i dokąd — nie wiadomo). Dlatego router ma być zwykłym `http.Handler`, który da się później uruchomić jako zwykły serwer bez przepisywania.

## Cel

1. Całe API przechodzi przez **jedną funkcję Vercela**, która rozdziela żądania po metodzie i ścieżce.
2. Limit 12 funkcji przestaje nas dotyczyć.
3. Dodanie endpointu to jedna linijka w tablicy tras, a nie wpis w `vercel.json` i `?kind=`.

## Wymagania

### Zachowanie dla świata — bez zmian
- **Te same publiczne adresy, metody, parametry, kody odpowiedzi i ciała.** Frontend (`web/src/data/api/`) i skrypty z kluczem (`curl` z `CLAUDE.md`) nie mogą tego zauważyć. Pełna lista dzisiejszych tras jest w `vercel.json` (rewrites).
- Jeśli coś dziś działa tylko dlatego, że rewrite dokleja parametry, router ma dawać ten sam efekt. Przykłady:
  - `/api/taxonomy/equipment/:id` → `?kind=equipment&id=:id`,
  - `?moveTo=` przy usuwaniu.
- Nieznana ścieżka pod `/api/…` → 404 w formacie JSON, jak dziś `resp.Error`.
- Zła metoda → 405, jak dziś.
- SPA (`/(.*)` → `/index.html`) działa jak dotąd.

### Bezpieczeństwo — nie może się poluzować
- Dziś każda funkcja w `api/` zaczyna się od `if resp.Begin(w, r) { return }` (CORS, preflight, bramka zapisu). Pilnuje tego test `go test ./http` (`TestEveryHandlerStartsWithBegin`).
- Po zmianie `Begin` ma się wykonać **raz, w routerze, przed każdą trasą**, także przed 404 i 405. Żadna trasa nie może go ominąć.
- Pilnujący test przepisz tak, żeby nadal łapał regresję. Na przykład test tabelaryczny po **wszystkich** zarejestrowanych trasach:
  - zapis bez tokenu → 401,
  - zapis ze złym kluczem → 401,
  - `OPTIONS` z dozwolonego originu → preflight OK,
  - nowa trasa automatycznie trafia do testu.
  Zostaw też test, że jedyna funkcja w `api/` zaczyna od `Begin`.
- Handlery, które dziś same wołają `Begin`, nie mogą go wołać drugi raz (podwójne nagłówki CORS).
- Szkice dalej niewidoczne publicznie (404), `resp.IsOwner(r)` działa jak dziś.

### Struktura
- W `api/` zostaje **jeden plik** z jedną funkcją (nazwę wybierz sam). Logika handlerów przechodzi do pakietu poza `api/` (np. `handlers/` — już istnieje).
- Uwaga, specyfika Vercel Go: każdy plik w `api/` z eksportowaną funkcją to funkcja, a pakiety pomocnicze muszą leżeć poza `api/`.
- **Router na standardowej bibliotece.** Go ≥ 1.22 ma `http.ServeMux` ze wzorcami `"GET /api/recipes/{id}"`, więc nowa zależność jest zbędna. Jeśli `go.mod` ma starsze Go, podbij wersję albo uzasadnij inne wyjście.
- `?kind=` znika z `units.go`: sprzęt, wyzwania i formaty dostają własne trasy.
- `id` z adresu: handlery czytają dziś `r.URL.Query().Get("id")`, bo tak przekazuje go rewrite. Ujednolić (`r.PathValue("id")`) albo zostawić adapter — twój wybór, byle spójnie.
- **Sprawdź na żywym wdrożeniu (preview), jaką ścieżkę funkcja Go dostaje w `r.URL.Path` po rewrite** `/api/(.*)` → jedna funkcja. To rozstrzyga, czy router może dopasowywać po oryginalnej ścieżce. Nie zgaduj.
- `vercel.json`: rewrites dla API zwijają się do jednego wpisu (plus SPA). `regions`, `buildCommand` itd. bez zmian.
- Nic nie przestaje działać lokalnie:
  - `go run ./cmd/...` (migracje, restore) korzysta tylko z bazy, nie z API, więc nie powinno go to dotknąć — sprawdź,
  - sprawdź, czy `vercel dev` dalej podnosi API.

### Dokumentacja
Zaktualizuj `CLAUDE.md`:
- „Każda funkcja API to oddzielny handler w `api/**/*.go`” → gdzie i jak dodaje się endpoint,
- sekcja „Autoryzacja API”: gdzie teraz siedzi `Begin` i co pilnuje testu.

Komentarze o limicie 12 funkcji w kodzie (`units.go`, `handlers/challenges.go`) usuń albo popraw.

## Zasady pracy
- **Baza:** `.env.local` wskazuje na **produkcję**. Ta zmiana nie dotyka schematu ani danych, więc migracji nie ma. `cmd/resetdb` jest zakazany.
- **Testy z zapisem:** lokalny wbudowany Postgres z seedem projektu (Go `fergusstrange/embedded-postgres`, port 5433), bez pytania właściciela.
  - Serwer testowy i baza **tylko na `127.0.0.1`**, nigdy `:3000` ani `0.0.0.0` — inaczej Windows pokazuje okno zapory.
  - Po testach zabij swoje procesy, także Postgresa (`pg_ctl stop -D <data> -m fast`).
  - Nie używaj `vercel dev` do testów zapisu, bo czyta `.env.local`, czyli produkcję.
  - Jeśli lokalny serwer testowy czyta trasy z `vercel.json`, po zmianie ma montować ten sam router, którego używa Vercel. To dodatkowo upraszcza testy.
- **Nic systemowego bez uprzedzenia.** Żadnych okien Windows, instalacji globalnych, procesów zostawionych w tle. Jeśli coś takiego jest konieczne — jedno zdanie do właściciela: co, po co, co wyskoczy.
- **Git:**
  - commity na `main`, bez branchy i PR, **bez pusha** (push tylko na prośbę właściciela),
  - na `main` leżą niewypchnięte commity Wyzwań — pracuj na nich,
  - **nie rób `git stash pop`**: w stashu jest odłożony kod zdjęć, który też dotyka `api/photos/` i `vercel.json`. Po routerze trzeba go będzie dopasować przy wznowieniu — zapisz to jednym zdaniem w raporcie.
- **Kolizje** (coś w wymaganiach się gryzie, wynik testu na preview zmienia plan) — opisz i zapytaj, nie rozwiązuj po cichu.

## Checklista odbioru
1. W `api/` jest dokładnie jedna funkcja; `vercel.json` ma jeden rewrite dla API plus SPA.
2. Tabela „trasa przed → trasa po” w raporcie, z wynikiem sprawdzenia każdej (metoda, adres, kod).
3. `?kind=` usunięte z kodu i `vercel.json`.
4. `Begin` raz, przed każdą trasą, także przed 404 i 405; test tabelaryczny po wszystkich trasach przechodzi; jedyna funkcja w `api/` zaczyna od `Begin`.
5. `go build ./...`, `go vet ./...`, `go test ./...` przechodzą.
6. Frontend na lokalnym routerze (Vite proxy `/api` → `127.0.0.1:3000`, embedded Postgres) — przeklikane:
   - lista i przepis,
   - lista zakupów,
   - zapis przepisu w panelu,
   - Klasyfikacja (osie, jednostki, sprzęt, formaty wyzwań),
   - wyzwania (lista, szczegóły, zapis, usunięcie z `moveTo`).
7. Ścieżka w `r.URL.Path` na Vercelu sprawdzona (preview albo `vercel deploy` bez `--prod`) — **tylko za zgodą właściciela**. Uwaga: preview najpewniej łączy się z tą samą produkcyjną bazą, więc pierwsze wywołanie API odpali na niej migrację `0009_challenges` z niewypchniętych Wyzwań. Przed preview kopia bazy (Actions → Kopia bazy → Run workflow) — powiedz o tym właścicielowi. Jeśli da się sprawdzić ścieżkę na minimalnej funkcji bez migracji (osobny mały projekt, bez bazy), to lepiej.
8. `CLAUDE.md` zaktualizowane.
9. Procesy testowe posprzątane.
10. Raport końcowy:
    - checklista: zrobione / pominięte + powód,
    - decyzje podjęte samodzielnie jako „do potwierdzenia”,
    - lista `curl` do sprawdzenia na produkcji po pushu: odczyty publiczne + zapis bez klucza → 401. Bez zapisów, które zmieniają dane.
