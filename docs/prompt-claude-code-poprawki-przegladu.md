# Prompt dla Claude Code — poprawki z przeglądu 18.11 (bez designu)

## Kontekst

18.11 ktoś przejrzał żywą stronę oczami znajomego, który dostał link do przepisu: `docs/przeglad-uzycia-2026-11-18.md`. Ten prompt obejmuje znaleziska, które **nie wymagają nowego designu**: to błędy albo rozjazdy z decyzjami, które już są. Pozostałe idą do designera albo do właściciela (treść) — tych nie ruszasz.

Numery poniżej to numery znalezisk z przeglądu. Przeczytaj w nim opis każdego punktu przed pracą.

## Zadania

### P1 (#1) Tryb gotowania ucina początek długiego kroku — sprawdź
20.11, przy wdrażaniu Uwag (commit `ef3ac72`), `CookMode` dostał poprawkę: wysoki krok przewija się od góry. **Sprawdź, czy to zamyka #1**:
- Bułki krok 3 i 8 przy 375×812,
- Ragù krok 3 i 7 przy 375×640,
- w obu motywach.
Jeśli coś dalej jest ucięte — popraw. Jeśli nie — odhacz w raporcie z zrzutem.

### P2 (#3) Lista zakupów gubi doprecyzowanie składnika
Na liście zakupów pokaż doprecyzowanie („mielona”, „wędzony”, „duże”), w tym samym stylu, w jakim stoi pod nazwą na stronie przepisu (te same klasy i tokeny). Nie projektuj nowego wyglądu.
- **Sumowanie:** jeśli ten sam składnik ma w przepisie różne doprecyzowania, pokaż wszystkie różne, bez powtórzeń.
- **Druk:** ten sam tekst.
- Jeśli układ wiersza listy zakupów nie ma miejsca na drugą linię i trzeba by go przeprojektować — zapytaj.

### P3 (#4) Podgląd linku w komunikatorach (meta)
Cel: link do przepisu i do listy zakupów wysłany na Messenger / WhatsApp pokazuje nazwę dania i krótki opis.
- **`web/index.html`:**
  - `lang="pl"`,
  - ogólne `meta description`, `og:title`, `og:description`, `og:type`, `og:site_name`,
  - tytuł w jednym formacie wszędzie (patrz P9).
- **Strony przepisu** (`/recipes/:id` i `/recipes/:id/zakupy`) dostają meta konkretnego przepisu. Komunikatory nie uruchamiają JS, więc meta muszą przyjść w HTML z serwera.
- **Propozycja:** trasa w routerze API, która oddaje `index.html` z podmienionymi meta. Rewrite w `vercel.json` kieruje tam **tylko boty podglądu** po nagłówku `user-agent`: Vercel `has`, np. facebookexternalhit, WhatsApp, Twitterbot, Slackbot, TelegramBot, Discordbot, LinkedInBot. Ludzie dalej dostają statyczny plik, bez zimnego startu funkcji.
- **Treść meta:**
  - przepis: tytuł = nazwa przepisu, opis = tagline przepisu albo pierwsze ~150 znaków opisu,
  - lista zakupów: „Lista zakupów: <nazwa> · <n> porcji”.
- **Szkic albo nieistniejący przepis:** meta ogólne, bez tytułu szkicu.
- **`og:image`:** zdjęć przepisów jeszcze nie ma, więc na razie bez obrazka. Obraz domyślny dojdzie od designera. Zostaw w kodzie miejsce, żeby go dopiąć jedną linijką.
- Jeśli zobaczysz lepszą drogę niż rewrite po user-agencie, opisz ją i zapytaj, zanim wdrożysz.
- Sprawdzenie: curl z user-agentem bota (lokalnie i po pushu) pokazuje właściwe meta. Zwykły curl dostaje statyczny `index.html`.

### P4 (#7) Przełącznik motywu przy systemowym trybie ciemnym
Bug: strona jest ciemna za systemem, a przycisk mówi, że nie jest. Pierwsze dotknięcie nic nie zmienia.
- Stan przycisku (`aria-pressed`, ikona, etykieta) ma wynikać z **aktualnie widocznego** motywu, także gdy pochodzi z systemu.
- Pierwsze dotknięcie zawsze przełącza na przeciwny.
- **Etykieta:**
  - po polsku i zgodna ze stanem, np. „Włącz tryb jasny” / „Włącz tryb ciemny”,
  - zamiast angielskiego „light” na desktopie: polski tekst albo sama ikona — wybierz to, co mniej zmienia układ, i zgłoś jako „do potwierdzenia”.

### P5 (#8) Wyszukiwarka bez polskich znaków
- **Polskie znaki i akcenty:** „bulki” znajduje „Bułki”, „ragu” znajduje „Ragù”. Normalizuj obie strony: małe litery, bez diakrytyków (ł → l).
- **Zamienniki:** szukaj też po zamiennikach składników.
- **Formy słowa** („bulka” → „Bułki”): zaproponuj proste dopasowanie bez słownika, np. słowo z zapytania ≥ 4 znaki pasuje po wspólnym początku. Wdroż je tylko wtedy, gdy nie daje wyraźnie głupich trafień na obecnych przepisach i seedzie. Opisz w raporcie jako „do potwierdzenia”.
- Sprawdź, gdzie dziś dzieje się wyszukiwanie (front czy API) i popraw tam.

### P6 (#12) Skalowanie: „2,7 łodygi”, „2⅔ łyżki”
- **Decyzja istnieje:** jednostki sztukowe zaokrąglamy do całości lub zakresu.
  - Sprawdź, dlaczego „łodyga” zachowuje się jak waga (dane jednostki? brak flagi?) i napraw tak, żeby szła jak sztuka.
  - Jeśli naprawa to zmiana danych na produkcji, zrób ją migracją wersjonowaną i powiedz o tym w raporcie, bo przed pushem trzeba będzie zrobić kopię bazy.
- **Łyżki i łyżeczki:** dziś „2⅔”. Propozycja analityka: zaokrąglaj do ½ (2⅔ → 2½, 1⅓ → 1½).
  - Jeśli w kodzie jest już mechanizm ułamków z listą dozwolonych mianowników, zastosuj go dla łyżek i łyżeczek. Jeśli nie ma — zapytaj.
  - Opisz jako „do potwierdzenia”.
- Gramy i mililitry („665 g”) zostają.

### P7 (#13) Link do listy zakupów nie niesie zamienników
- Wybrane zamienniki trafiają do adresu listy zakupów (obok `?porcje=`), np. `?zamienniki=<id>:<id>`. Nazwę parametru i format wybierz sam, spójnie z istniejącym.
- Odbiorca widzi ten sam wariant co nadawca.
- **Adres ma pierwszeństwo przed pamięcią urządzenia**, ale niczego w pamięci odbiorcy nie nadpisuje na stałe. Jeśli to się gryzie z tym, jak działają porcje w adresie, zrób tak samo jak porcje i opisz.

### P8 (#14) Powrót z listy zakupów
- „‹ Przepis”:
  - jeśli użytkownik przyszedł z przepisu, wraca w historii (jak systemowe „wstecz”) do miejsca, z którego wyszedł,
  - jeśli wszedł z linku prosto na listę, otwiera przepis przewinięty do składników.
- Bez nowych wpisów w historii.

### P9 (#15 częściowo, #18) Drobne rzeczy
- **Tryb gotowania:** podpowiedź „Strzałki ← → zmieniają krok” tylko na urządzeniach z klawiaturą lub kursorem, np. `@media (hover: hover) and (pointer: fine)`. Na telefonie znika. Reszta #15 (zamknięcie arkusza, licznik porcji) idzie do designera — nie ruszaj.
- **Tytuł karty przeglądarki:** jeden format wszędzie, np. „Ragù alla bolognese · Uszata Kuchnia” na przepisie i „Przepisy · Uszata Kuchnia” na liście. Ten sam format w P3.
- **Alergeny zamiennika:** dopisek „zawiera: X” tylko wtedy, gdy zamiennik **dodaje** alergen, którego przepis jeszcze nie ma (zgodnie z decyzją). Bulion warzywny w Ragù nie powinien pokazywać „zawiera: seler”.

## Poza zakresem — nie ruszaj
- #5, #6, #9, #10, #15 (arkusz w trybie gotowania), #16 — idą do designera.
- #2, #17 — treść, poprawia właściciel w panelu.
- #11 (kroki a zamienniki) — wróci przy feature „odwołania w treści kroku”.
- Zdjęcia (odłożone).

## Zasady pracy
- **Kolejność:** zanim napiszesz kod, rozpisz plan punkt po punkcie (P1–P9). Wyślij go do sesji analityka i poczekaj na OK. Analityk przekazuje pytania i odpowiedzi właściciela. Pytania i kolizje kieruj tak samo, nie rozstrzygaj po cichu.
- **Baza:**
  - `.env.local` wskazuje na **produkcję**, więc schemat i dane zmieniasz tylko migracją wersjonowaną. Nigdy `cmd/resetdb`.
  - Testy z zapisem: embedded Postgres z seedem projektu, port 5433, **wyłącznie `127.0.0.1`**.
  - Serwer testowy też na `127.0.0.1`, nigdy na `:3000` ani `0.0.0.0`.
  - Nie używaj `vercel dev` do zapisu.
  - Po pracy zabij wszystkie swoje procesy, także Postgresa.
- **Komputer właściciela:** żadnych okien systemowych Windows ani instalacji globalnych. `python` i `python3` to tu alias Microsoft Store — skrypty pisz w node albo PowerShellu.
- **API:** jeden router. Nowa trasa = linijka w `router/routes.go`, bez nowych plików w `api/` (patrz `CLAUDE.md`).
- **Front:** po zmianach w `.vue` i `.ts` uruchom ESLint. Hover nie może przesuwać ani skalować tekstu (zasada z 20.11).
- **Zrzuty:** po zmianach widocznych w UI — jasny, ciemny i 375 px.
- **Git:**
  - commity etapami na `main`, **bez pusha**; push robi analityk,
  - commituj tylko swoje pliki, bez `git add -A`,
  - nie ruszaj `docs/` ani stasha, nie zmieniaj niezacommitowanych plików w `content/` — to cudza praca.
- **Raport na koniec:**
  - P1–P9 odhaczone: zrobione / pominięte + powód,
  - decyzje „do potwierdzenia”,
  - informacja, czy jest migracja (wtedy przed pushem trzeba zrobić kopię bazy).
