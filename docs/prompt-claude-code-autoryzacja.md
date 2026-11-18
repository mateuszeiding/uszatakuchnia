# Prompt dla Claude Code — autoryzacja zapisów w API

## Kontekst

Uszata Kuchnia działa na produkcji (Vercel + Neon). Od 14.11.2026 w bazie są prawdziwe przepisy właściciela. Frontend loguje admina przez Auth0 (`web/src/auth/createAuth.ts`, `authGuard` w routerze), ale **backend w ogóle nie sprawdza, kto wysyła żądanie**. Login chroni tylko ekrany, nie dane.

Stan produktu: `docs/produkt.md`, sekcja „Ryzyka i dług”, punkt 1.

## Problem

Każdy, kto podejrzy adres API w narzędziach przeglądarki, może:
- tworzyć, zmieniać i usuwać przepisy (`api/recipes/create.go`, `api/recipes/detailed.go` — PUT, DELETE),
- tworzyć, zmieniać i usuwać składniki (`api/ingredients/create.go`, `api/ingredients/detailed.go`),
- wgrywać zdjęcia (`api/photos/upload.go`),
- zmieniać słowniki: osie, wartości, jednostki, sprzęt (`api/taxonomy/values.go`, `api/taxonomy/units.go` i pozostałe zapisy w `api/taxonomy`),
- czytać szkice: `GET /api/recipes?admin=true` zwraca szkice każdemu (`api/recipes/list.go`), a `GET` szczegółów przepisu nie sprawdza statusu (identyfikatory są kolejnymi liczbami, łatwo je zgadnąć).

Do tego CORS (`http/cqrs.go`) wpuszcza każdą domenę i wszystkie metody, więc obca strona może wykonać te żądania z przeglądarki odwiedzającego.

## Kto ma mieć prawo zapisu

1. **Właściciel w przeglądarce** — przez to samo logowanie Auth0 co dziś. Samo zalogowanie nie wystarcza. Zapis tylko dla konta właściciela, nawet gdyby w Auth0 dało się założyć inne konto.
2. **Sesje Claude Code na komputerze właściciela.** Dziś działają tak:
   - analityk / sesja „dodaj przepis” wpisuje przepisy na produkcję przez API aplikacji (`POST /api/ingredients`, `POST` i `PUT /api/recipes`, `GET` szkiców przed edycją),
   - inżynier zapisuje rekordy lokalnie (skrypty `go run`, `vercel dev` na bazie z `.env.local`).

   **Obie ścieżki muszą działać dalej.** Sekret potrzebny sesjom może być tylko w `.env.local` (lokalnie) i w zmiennych środowiskowych Vercela. Nie może trafić do repo ani do kodu frontendu.
3. **Nikt inny.** Odczyt opublikowanych treści (lista, przepis, lista zakupów, zdjęcia, słowniki) zostaje publiczny, bez logowania.

## Wymagania

- Każdy endpoint zapisujący (POST / PUT / PATCH / DELETE) odrzuca żądanie bez ważnych uprawnień (401 bez tokenu, 403 z tokenem, ale bez uprawnień). Dotyczy też endpointów, których nie wymieniłem — przejrzyj całe `api/`.
- Szkice widzi tylko uprawniony: lista z `admin=true` i szczegóły szkicu. Publiczny `GET` szkicu zwraca 404, a nie 403, żeby nie zdradzać, że szkic istnieje.
- Nowy endpoint zapisujący dodany w przyszłości ma być domyślnie chroniony albo mieć ochronę, której nie da się łatwo pominąć. Każdy handler jest osobną funkcją Vercela, więc zapomnieć jest łatwo.
- **Formularz przepisu nie ma autozapisu.** Jeśli sesja logowania wygaśnie w trakcie edycji, zapis nie może zgubić wpisanej treści. Właściciel ma dostać czytelny komunikat i możliwość ponownego zalogowania bez utraty formularza.
- Kto nie jest właścicielem, nie widzi UI admina (badge, link do panelu, akcje), nawet jeśli jest zalogowany.
- CORS: zawęź do domen aplikacji (produkcja, podglądy Vercela, jeśli używane, lokalny dev), o ile nie zepsuje to ścieżek z punktu 2.
- Wdrożenie nie odcina właściciela: backend wymagający tokenu i frontend, który go wysyła, wchodzą razem.

## Decyzje dla ciebie (inżynier)

- Mechanizm dla sesji lokalnych: token maszynowy z Auth0 (client credentials) albo własny klucz API w zmiennych środowiskowych. Wybierz, uzasadnij i opisz koszt utrzymania (rotacja, wygasanie).
- Sposób weryfikacji tokenu w Go: biblioteka, cache kluczy, wspólny helper.
- Jak rozpoznać właściciela: uprawnienie lub rola w Auth0 albo identyfikator użytkownika w konfiguracji.

Jeśli coś z wymagań okaże się drogie albo niewykonalne, opisz to i zapytaj. Nie obchodź tego po cichu.

## Poza zakresem

- Konta użytkowników, role inne niż właściciel, rejestracja.
- Strona logowania (backlog designu: „Później”).
- Zmiany w wyglądzie poza komunikatem o wygasłej sesji. Jeśli komunikat wymaga projektu, opisz potrzebę, a ja przygotuję brief dla designera.

## Testowanie

Nie testuj na produkcji. `.env.local` wskazuje na produkcyjną bazę Neon, więc użyj lokalnej bazy. `resetdb` jest zakazany.

## Do oddania

- Zmiany w kodzie (commity na `main`, bez pushowania bez prośby).
- Aktualizacja `CLAUDE.md`: jak sesja Claude Code uwierzytelnia się przy zapisie przez API, skąd bierze sekret, co zrobić, gdy sekret wygaśnie. Sekcja „Przepisy dodawane przez Claude” musi to opisywać.
- Lista zmiennych środowiskowych do ustawienia w Vercelu i w `.env.local` oraz kroki w panelu Auth0, które właściciel musi wykonać ręcznie.

### Checklista odbioru

- [ ] Każdy POST / PUT / PATCH / DELETE w `api/` bez tokenu → 401 (lista sprawdzonych endpointów w raporcie).
- [ ] Zalogowane konto inne niż właściciel → 403 na zapisie, brak UI admina.
- [ ] `GET /api/recipes?admin=true` bez uprawnień nie zwraca szkiców.
- [ ] `GET` szczegółów szkicu bez uprawnień → 404.
- [ ] Odczyt opublikowanych treści bez logowania działa jak dziś (lista, przepis, lista zakupów, zdjęcia).
- [ ] Właściciel w przeglądarce: tworzenie, edycja, publikacja, usuwanie przepisu; edycja składnika; zmiana w słowniku; wgranie zdjęcia.
- [ ] Wygasła sesja w trakcie edycji formularza: treść nie ginie, da się zalogować ponownie i zapisać.
- [ ] Sesja Claude Code z sekretem z `.env.local`: dodanie składnika i przepisu przez API produkcyjne.
- [ ] Lokalna praca inżyniera (`vercel dev`, skrypty) działa.
- [ ] CORS zawężony (albo opis, dlaczego nie).
- [ ] `CLAUDE.md` zaktualizowany.
