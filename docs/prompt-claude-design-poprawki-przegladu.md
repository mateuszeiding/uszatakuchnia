# Prompt dla Claude Design — poprawki po przeglądzie żywej strony (18.11)

## Kontekst

Uszata Kuchnia to aplikacja z przepisami właściciela. Do czytania jest publiczna, bez kont użytkowników. Admin jest jeden i pracuje na desktopie. Strony publiczne muszą działać na telefonie. DS v2, tokeny i komponenty bierz z istniejących prototypów.

Dziś aplikacja trafia do ludzi tak: właściciel wysyła znajomemu link do przepisu albo do listy zakupów. 18.11 ktoś przejrzał żywą stronę oczami takiego znajomego (telefon 375 px, light i dark, plus desktop). Opublikowane są dwa przepisy: Ragù alla bolognese i Bułki do burgerów. Żaden nie ma zdjęcia.

Część znalezisk naprawia inżynier bez designu: wyszukiwarka bez polskich znaków, lista zakupów, podgląd linku, przełącznik motywu, zaokrąglanie ilości. Tu jest **6 rzeczy, które wymagają decyzji projektowej**.

Ogólna zasada z 20.11 obowiązuje: najechanie kursorem nie przesuwa i nie skaluje niczego, co zawiera tekst.

## 1. Przepis bez zdjęcia — pierwszy ekran to fioletowy prostokąt
- **Gdzie:** strona przepisu (`Przepis.html`), 375 px, oba motywy.
- **Co się dzieje:**
  - zastępczy obrazek 1:1 zajmuje ~670 z 812 px pierwszego ekranu,
  - nazwa przepisu wypada na samym dole, a składniki zaczynają się około dwóch ekranów niżej,
  - w trybie ciemnym obrazek jest jasnolawendowy, czyli to największa i najjaśniejsza plama na ciemnej stronie.
- **Skala problemu:** dotyczy 100% przepisów. Własne zdjęcia są odłożone, więc tak zostanie na dłużej.
- **Zadanie:** zaproponuj, jak wygląda przepis bez zdjęcia na telefonie i na desktopie. Znajomy ma od razu zobaczyć, co to za danie, i szybko dojść do składników.
- **Do rozstrzygnięcia:**
  - czy karta na liście przepisów też się zmienia,
  - jak to wróci, kiedy zdjęcia już będą.

## 2. „Zacznij od nowa” kasuje postęp jednym dotknięciem
- **Gdzie:** nagłówek „Kroki” na stronie przepisu.
- **Co się dzieje:**
  - przycisk pojawia się po pierwszym odhaczeniu, tuż obok „Tryb gotowania”,
  - jedno dotknięcie, bez pytania, czyści odhaczone kroki i składniki oraz przywraca domyślne porcje,
  - wybrane zamienniki składników zostają, więc nie wiadomo, co właściwie zostało wyczyszczone.
- **Kontekst użycia:** przy kuchence, mokrymi palcami, łatwo trafić obok.
- **Zadanie:** zabezpiecz przed przypadkowym resetem: potwierdzenie, „Cofnij” po resecie albo inne miejsce przycisku. Rozstrzygnij też, co dokładnie resetuje się razem: porcje, zamienniki, odhaczenia.

## 3. Filtry, w których prawie wszystko daje 0 wyników
- **Gdzie:** lista przepisów: arkusz „Filtry” na telefonie, panel po lewej na desktopie (`Przepisy.html`).
- **Co się dzieje:**
  - na 54 opcje filtrów 45 daje „Pokaż 0 przepisów”,
  - filtr czasu kończy się na „do 2 h”, a oba przepisy trwają dłużej.
- **Kontekst:** przepisów będzie przybywać, ale powoli.
- **Zadanie:** co robimy z opcjami, które dają 0 wyników:
  - ukrywamy,
  - przygaszamy z licznikiem,
  - chowamy cały panel, póki przepisów jest mało,
  - albo coś innego.
- **Plus:** czy skala czasu ma zakres powyżej 2 h.

## 4. Przepis, którego nie ma, i nieznany adres
- **Co się dzieje:**
  - `/recipes/999` albo link do przepisu usuniętego lub cofniętego do szkicu pokazuje „Nie udało się wczytać przepisu.” — brzmi jak awaria sieci, więc ludzie odświeżają i czekają,
  - nieznany adres (np. `/przepis-od-kuby`) bez słowa przekierowuje na listę przepisów.
- **Zadanie:** zaprojektuj stronę „nie ma takiego przepisu” i stronę 404, telefon i desktop. Rozdziel też „przepisu nie ma” od „nie udało się wczytać”, czyli prawdziwego błędu sieci z „Spróbuj ponownie”.

## 5. Arkusz „Sprzęt i składniki” w trybie gotowania
- **Gdzie:** tryb gotowania (jeden krok naraz, duże „Wstecz” / „Dalej”) → przycisk „Składniki”, telefon.
- **Co się dzieje:**
  - okno na środku ekranu bez przycisku zamknięcia — zamyka się tylko dotknięciem w tło, czego nic nie sugeruje,
  - nie ma licznika porcji ani informacji „Z zamiennikami: n”, tylko liczba porcji w tytule.
- **Zadanie:** jak ten arkusz wygląda i jak się go zamyka. Do rozważenia: arkusz od dołu jak przy zamiennikach (`Zamienniki.html`), wyraźne „Zamknij”. Rozstrzygnij, czy jest tam licznik porcji.

## 6. Czyja to strona i „Zaloguj”
- **Gdzie:** górny pasek i lista przepisów.
- **Co się dzieje:**
  - na desktopie jedynym przyciskiem w pasku (obok motywu) jest „Zaloguj”,
  - nigdzie nie pada, czyje to przepisy i że do czytania nie trzeba konta,
  - znajomy z linku może uznać, że czegoś mu brakuje, albo że to kolejny serwis z przepisami.
- **Stan obecny:** na telefonie „Zaloguj” jest już schowane (decyzja przy Wyzwaniach). Logować się może tylko właściciel.
- **Zadanie:**
  - gdzie i jak „Zaloguj” przestaje się rzucać w oczy, np. dyskretny link w stopce,
  - jedno zdanie, czyje to przepisy: zaproponuj miejsce i 2–3 warianty tekstu w luźnym tonie, bez marketingu.

## Dodatkowo: obraz podglądu linku
Inżynier dodaje meta dla komunikatorów: nazwa dania i opis przy wysłanym linku. Zdjęć przepisów nie ma, więc potrzebny jest **jeden domyślny obraz podglądu** 1200×630 (logo i nazwa aplikacji), czytelny w miniaturze Messengera i WhatsAppa. Oddaj go jako plik w projekcie.

## Poza zakresem
- Wgrywanie zdjęć (odłożone).
- Kroki zmieniające się razem z zamiennikiem: wróci przy osobnym featurze.
- Rzeczy, które robi inżynier: wyszukiwarka, lista zakupów (doprecyzowania, zamienniki w linku, powrót), meta, przełącznik motywu, zaokrąglanie ilości, podpowiedź o strzałkach na telefonie.

## Do oddania
- **Artboardy:** light, dark i 375 px dla punktów 1, 2, 4, 5 i 6. Desktop dla 1, 3, 4 i 6.
- **Punkt 3:** wariant albo dwa z rekomendacją.
- **Handoff:** checklista odbioru, sekcja „Nie zmieniaj”, otwarte pytania.
- **„Decyzje i backlog.md”:** wpis z decyzjami. Uwaga: „przepis bez zdjęcia” i „strona 404” są tam dziś w „Później”. Przenieś je.
