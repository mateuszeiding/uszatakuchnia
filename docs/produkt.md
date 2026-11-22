# Uszata Kuchnia — stan produktu

Aktualizacja: 18.11.2026 · prowadzi: analityk (sesja Claude Code od featurów)

Żywy dokument: co jest, co się dzieje, co dalej. Szczegóły wyglądu i zachowania są w projekcie Claude Design — [artefakt](https://claude.ai/artifact/YEuK2L22otH24hdvcUC5GR), pliki `project/Decyzje i backlog.md` (decyzje) i `project/Handoff - Claude Code.md` (prompty i checklisty wdrożeń). Ten plik ich nie powtarza, tylko do nich odsyła.

## Po co

Aplikacja z przepisami właściciela, publiczna do czytania, bez kont użytkowników. Jeden autor tworzy całą treść.

Odbiorcy (18.11.2026): właściciel i kilku znajomych. Później może szersze grono. Główny cel: zbieranie przepisów.

Trzy filary:
1. **Przepisy** — wpisywane ręcznie, z naciskiem na gotowanie przy kuchence: skalowanie porcji, postęp, tryb gotowania, lista zakupów.
2. **Wyzwania** — lista zadań kulinarnych w duchu kulinarnych programów konkursowych. Wdrożone lokalnie 19.11 (patrz „W toku”).
3. **Wiki** — baza wiedzy o składnikach i technikach. Jeszcze niezaczęte.

Zasady:
- Strony publiczne działają na telefonie. Admin pracuje tylko na desktopie.
- Bez funkcji na zapas: elementy „wkrótce” nie trafiają do UI.
- Od 14.11.2026 w bazie są prawdziwe dane. Wszystkie zmiany idą przez migracje (`CLAUDE.md`).

## Jak pracujemy

Pomysł → **brief dla designera** (`docs/prompt-claude-design-<temat>.md`) → **design + handoff** w artefakcie (Claude Design) → **wdrożenie** (Claude Code, commity na `main`) → **przegląd** (designer porównuje wdrożenie z projektem i dopisuje rozjazdy do „Decyzji i backlogu”).

Praca bez UI (backend, dane, infrastruktura) omija designera: brief dla inżyniera w `docs/prompt-claude-code-<temat>.md`.

Przepisy wpisuje Claude przez API, jako szkic. Właściciel je poprawia i publikuje. Kopia przepisów w `content/`.

## W toku

| Feature | Etap | Gdzie |
|---|---|---|
| **Zdjęcia przepisów — ODŁOŻONE (18.11)** | wracamy przy decyzji o przenosinach infrastruktury. Niedokończony kod sesji inżyniera w `git stash` („zdjęcia przepisów — niedokończone (odłożone 18.11)”; zawiera też niechciany import z linku). Porównanie magazynów (18.11): rekomendacja Cloudflare R2 (S3, 10 GB za darmo, transfer bez opłat, przenośne); Vercel Blob prostszy, ale związany z Vercelem (1 GB); bytea odpada (limit Neona). Stan obecny: jedyne obsługiwane źródło obrazów to **Unsplash** — obraz po id zdjęcia z Unsplash, pobierany przez ich API (`db/helpers/unsplash_image.go`, `services/image_provider.go`, encja `images`, wymaga `UNSPLASH_ACCESS_KEY`), używany w seedzie. Ich utrata nic nie szkodzi (właściciel, 18.11) — kandydat do usunięcia. Właściciel nie może dodać własnego zdjęcia. Aplikacja trzyma zdjęcie jako adres lub id z obcej strony, a własnego pliku nie ma gdzie zapisać. Cel: upload własnych zdjęć z pliku | `docs/prompt-claude-code-zdjecia.md` |
| Autoryzacja zapisów w API | na produkcji (c65e2cb, bf16955): zapis bez klucza → 401, klucz z `.env.local` działa. Czeka na sprawdzenie przez właściciela w przeglądarce (logowanie, panel, zapis) | `CLAUDE.md` |
| Codzienna kopia zapasowa bazy | działa: przebiegi nocne i ręczne kończą się sukcesem (sprawdzone 19.11). Odtwarzanie pojedynczej rzeczy tylko dla przepisów (`cmd/restore-recipe`), wyzwania — razem z całą bazą | `docs/backup-bazy.md` |
| Poprawki z przeglądu 18.11 | rozdzielone 20.11. **Inżynier** (bez designu): #1 sprawdzić (pewnie naprawione przy Uwagach), #3, #4 meta OG (rewrite tylko dla botów), #7, #8, #12, #13, #14, #15 (podpowiedź o strzałkach), #18 — oddelegowane 20.11 (sesja „Wdrożenie: poprawki z przeglądu”). **Designer:** #5 przepis bez zdjęcia, #6 reset postępu, #9 filtry z zerem, #10 404 / brak przepisu, #15 arkusz w trybie gotowania, #16 „Zaloguj” i czyja to strona, + obraz podglądu linku — brief gotowy. **Właściciel (treść):** #2 makaron w Ragù, #17. **Odłożone:** #11 (z „odwołaniami w kroku”) | `docs/przeglad-uzycia-2026-11-18.md`, `docs/prompt-claude-code-poprawki-przegladu.md`, `docs/prompt-claude-design-poprawki-przegladu.md` |
| Wyzwania | design zaakceptowany (artefakt: Handoff → „Wyzwania”, prototyp `Wyzwania.html`); **na produkcji 19.11** (push po kopii bazy z 18:18 UTC; migracja 0009 przeszła, 3 formaty startowe są). Pierwsze wyzwanie dodane przez właściciela 20.11. Zakładki Klasyfikacji: przewijanie (decyzja 19.11). Do potwierdzenia: 7 decyzji inżyniera | `docs/prompt-claude-design-wyzwania.md` |
| Jedno wejście do API (router) | **na produkcji 19.11** (290b90c), sprawdzone curlem: odczyty 200, zapis bez klucza 401, nieznany adres 404. Zmiana celowa: nieznany `/api/…` → 404 JSON zamiast index.html. Znosi limit 12 funkcji Vercela (ryzyko 8), adresy API bez zmian. Bez migracji | `docs/prompt-claude-code-router-api.md` |
| Hover kart bez rozmytego tekstu | design zaakceptowany 20.11 (wariant B: ramka `--rule-strong` + tytuł w akcencie, przycisk akcentowy ciemnieje zamiast podjeżdżać; zasada „hover nie przesuwa tekstu”); wdrożone lokalnie 20.11 (9ce7752), niewypchnięte, bez migracji. Odkrycie: lista składników (`IngredientsView`) nie ma trasy — martwy widok | `docs/prompt-claude-design-hover-kart.md` |
| Linki do przepisów w nowej karcie | decyzja właściciela 20.11: ikona przy składniku i karty „Przydatne przepisy” otwierają nową kartę także na stronie przepisu (zmienia decyzję z 17.11 — wcześniej tylko tryb gotowania; designerowi do odnotowania w „Decyzje i backlog”). **Na produkcji 20.11** (8137aa4). Karty „Przydatne przepisy” bez ikony „nowa karta” (decyzja właściciela) | — |
| Sesja Auth0 (błąd, 21.11) | logowanie tylko w pamięci karty → „Podgląd szkicu” w nowej karcie niezalogowany, odświeżenie wylogowuje. Naprawa: `cacheLocation: 'localstorage'` + refresh tokeny z rotacją — **zmergowane 21.11** (PR #6, 799ba88; przy okazji naprawiony test `recipeForm.test.ts`). Do sprawdzenia przez właściciela na żywo. Wygasła sesja (zalogowany bez tokenu → okienko logowania zamiast „Nie ma takiego przepisu”) — zmergowane 21.11 (PR #7, 86aade4), niesprawdzone na prawdziwym Auth0. Auth0 ustawione przez właściciela 21.11: rotacja refresh tokenów (overlap 30 s, reszta domyślnie) i „Allow Offline Access” w API — rotacja jest na darmowym planie | — |
| Słowniczek pojęć + kroki opcjonalne (21.11) | design zaakceptowany 21.11 (`Slowniczek i kroki opcjonalne.html`), **na produkcji 22.11** (PR #10, squash f5d8627, migracja 0012 po kopii z 21.11 20:22 UTC; CodeRabbit: 11 + 3 uwagi poprawione; `/api/glossary` działa, słowniczek pusty — do wypełnienia; procedura przenoszenia pojęć w CLAUDE.md → `docs/content-na-innej-bazie.md`); plan zaakceptowany 21.11: tabela `glossary_terms`, `{{term:id|słowo}}` w krokach i wskazówkach, `/api/glossary`, panel „Słowniczek”, skrót Ctrl+Shift+L; `recipe_steps.optional`, postęp liczy obowiązkowe, tryb gotowania „Pomiń”. Ułamki ujednolicone (decyzja 21.11): połówki i ćwiartki dziesiętnie („0,25”, „1,5”), tylko tercje „1/3” — wszędzie. Przekazane designerowi w briefie kadru zastępczego (21.11). Briefy (jeden plik, dwa zadania): pojęcie oznacza słowo autora w kroku/wskazówce (`term`, kropkowane podkreślenie, wyjaśnienie po dotknięciu), lista pojęć w panelu; krok opcjonalny z flagą, postęp i tryb gotowania do rozstrzygnięcia | `docs/prompt-claude-design-slowniczek-kroki-opcjonalne.md` |
| Kadr zastępczy — redesign (21.11) | właścicielowi obecny (gradient z id + „talerz”) się nie podoba. Kolor ma iść z typu dania (badge, `--c-*`). Brief: 5 szybkich koncepcji (jeden widok, desktop, jasny motyw), wybór, potem pełny design. W briefie też format ułamków dla designera. **Wybrany wariant B** (21.11): ikona typu dania (9 SVG) na tle `color-mix` koloru pierwszego typu; wdrożenie w toku (sesja „kadr: ikona typu dania”, PR); **na produkcji 22.11** (PR #11, 89fd26c; migracja 0013 przeszła, 9 typów ma ikony). Po migracji przez ~1–2 min część zapytań dawała 500 („cached plan must not change result type”, „prepared statement name is already in use” — cache zapytań sterownika vs pooler Neona po zmianie schematu), potem samo przeszło; utwardzone 22.11 (PR #12, d875a95: pgx w trybie exec, wszystkie połączenia przez `db.Open`). Zgłoszony błąd aplikacji Claude (wiadomości „hold” zniknęły po limicie): anthropics/claude-code#111607. Wiązanie ikony: pole `icon` na wartości osi + migracja 0013 (decyzja właściciela 21.11); picker ikony w panelu — później, po designie (do tego czasu ikona przez API) | `docs/prompt-claude-design-placeholder-zdjecia.md` |
| Karta przepisu: kadr zastępczy (21.11) | decyzja właściciela zmienia ustalenie z 20.11: na liście kadr 1:1 zawsze, bez zdjęcia gradient w odcieniu przepisu + „talerz” (h = id·47 mod 360), „wymaga przygotowania” na kadrze; strona przepisu dalej bez kadru. Zmergowane 21.11 (PR #9, c3f686d) | — |
| Ułamki bez znaków ½ ⅓ ¼ (21.11) | decyzja właściciela: glify nieczytelne. Wyświetlanie: połówki „1,5”, tercje/ćwiartki „1/3”, „1 1/4”; parser dalej przyjmuje ½. Zmergowane 21.11 (PR #8, 2d3699a) | — |
| CodeRabbit (próba, 21.11) | podłączony do repo przez właściciela. Pierwszy PR: ilości w krokach + tryb gotowania. Podglądy Vercela dla gałęzi wyłączone (b22dd22, `ignoreCommand`) — inaczej otwarcie podglądu PR-a odpaliłoby migrację na produkcyjnej bazie (jedna `DB_DATABASE_URL` dla wszystkich środowisk). Brak `gh` na komputerze — do decyzji | — |
| Kolejka wdrożeń (20.11) | jeden checkout, więc po kolei: ① poprawki P1–P9 (**na produkcji 20.11**, 42ba50f…56d34bf; meta dla botów sprawdzone curlem z user-agentem Facebooka i WhatsAppa) → ② poprawki z designem (**na produkcji 20.11**, 77a584c…a5dc10f; og-image sprawdzony curlem) → ③ podgląd szkicu (**na produkcji 20.11** przez PR #2, squash f06b413, po 2 rundach review CodeRabbita; → ③½ tłumaczenie i przesiew komentarzy w kodzie (**zmergowane 21.11**, PR #4, 206d788: 1384 → 1056 linii komentarzy, 2 rundy); Vercel Web Analytics na produkcji 21.11 (PR #3, c672342, bez /admin); drobiazg na później: cudza lista zakupów szkicu pokazuje „Nie udało się wczytać” zamiast „Nie ma takiego przepisu”) → ④ ilości w krokach + tryb gotowania bez przewijania (**na produkcji 21.11**, PR #5, squash 491232f, po kopii bazy z 01:32 UTC; migracja 0011 przeszła — wiersze mają klucze i1, i2…; CodeRabbit: 12 uwag, wszystkie poprawione; plan zaakceptowany 21.11: stały klucz wiersza `recipe_ingredients.key`, odwołania `{{ing:key|część}}` w treści kroku, migracja 0011; „Zanim zaczniesz” gdy uwagi lub przygotowanie; desktop bez zmian. 21.11 dopisane odwołania z ilościami do Bułek (12) i Ragù (10) — teksty właściciela co do słowa, bez przeredagowania (decyzja właściciela; za długie kroki Bułek zostają, limit trybu gotowania to tylko ostrzeżenie). Otwarte: makaron w składnikach Ragù; migracja — przed merge'em kopia bazy). Zmiana decyzji „Uwagi”: w trybie gotowania osobny ekran „Zanim zaczniesz” przed krokiem 1 | — |
| Ilości składników w treści kroku | design zaakceptowany 20.11 (wariant (a): słowo autora + znacznik ilości `[350 g]`, zamiennik w znaczniku z ⇄, części ½…¾ i „reszta”, wstawianie przez „@”). Brief: 20.11 (feedback znajomego: „dodaj mąkę — ile?”). Odwołanie do wiersza składnika w kroku, skaluje się z porcjami i zamiennikiem; kluczowa decyzja: gramatyka (słowo autora + ilość vs stała forma). Słowniczek później tym samym mechanizmem. Po wdrożeniu Claude przepisuje kroki istniejących przepisów | `docs/prompt-claude-design-ilosci-w-krokach.md` |
| Tryb gotowania bez przewijania | design zaakceptowany 20.11 (tekst 26→23→20 px, zwijane wskazówki, „Dalej” 3/4 × 72 px, limit kroku 280 znaków). Brief: 20.11 (feedback znajomego: przewijanie brudnymi rękami). Krok mieści się na ekranie, duża strefa „Dalej”, limit długości kroku w formularzu | `docs/prompt-claude-design-tryb-gotowania-bez-przewijania.md` |
| Podgląd szkicu | design zaakceptowany 20.11 (pasek „SZKIC · Widzisz go tylko ty” + Edytuj, bez „Opublikuj”; „Podgląd ↗” w formularzu). Brief: 20.11 (mały): oko przy szkicach w panelu + pasek „Szkic — widzisz go tylko ty”. Backend już zwraca szkic właścicielowi. Potem „Test” (id 7) może wrócić do szkiców | `docs/prompt-claude-design-podglad-szkicu.md` |
| Poprawki z przeglądu — część z designem | design zaakceptowany 20.11 (Handoff → „Poprawki po przeglądzie żywej strony”, prototyp `Poprawki 18.11.html`). Do wdrożenia przez sesję „poprawki z przeglądu” jako 2. etap, po P1–P9 | — |
| Uwagi (przepis i wyzwanie) | design zaakceptowany 20.11 (Handoff → „Uwagi — przepis i wyzwanie”, prototyp `Uwagi.html`); **na produkcji 20.11** (5873b08…8eff6f9, push po kopii bazy z 01:05 UTC; migracja 0010 przeszła — przepisy i wyzwania mają `cautions: []`). 6 drobnych decyzji inżyniera do potwierdzenia (m.in. ramka pola jak reszta formularza zamiast mocniejszej z prototypu). Blok zastępuje baner „Wymaga przygotowania”. Lista 0–n krótkich punktów „czego nie zawalić” (zakupy, bezpieczeństwo, czas) — inna kategoria niż wskazówki przy krokach; brief odnosi się do odrzuconego „Zanim zaczniesz” | `docs/prompt-claude-design-uwagi.md` |

Przegląd żywej strony 18.11: wskazówki, sprzęt i zamienniki działają zgodnie z projektem. Tryb gotowania ucina długie kroki na telefonie (błąd). Nadal niesprawdzone: linki do przepisów i „Przydatne przepisy” (żaden opublikowany przepis nie ma linku), Wake Lock, wydruk listy zakupów.

## Co działa

### Lista przepisów `/recipes`
- Filtry: typ dania, baza, kuchnia, czas, trudność, diety (zawsze wszystkie naraz), praktyczne. Stan filtrów w adresie. Wyszukiwarka szuka też po składnikach. „Pokaż więcej”. Na telefonie filtry w arkuszu od dołu. · 14.10
- `~520 kcal` na karcie · 7.11
- Szkielety w docelowym układzie zamiast spinnerów; panel filtrów nie skacze · 17.11

### Strona przepisu `/recipes/:id`
- Skalowanie porcji. Ilości jako liczby, ułamki i zakresy · 14.10, czytelne ułamki · 14.11
- Alergeny wyliczone ze składników (14 alergenów UE) · 14.10
- Krok „wymaga przygotowania wcześniej” → badge na przepisie i filtr na liście · 14.10
- Postęp gotowania zapamiętany na urządzeniu (składniki, kroki, porcje), „Zacznij od nowa” · 14.10
- Tryb gotowania: jeden krok naraz, ekran nie gaśnie · 14.10
- Wartości odżywcze: orientacyjne zakresy, kcal na porcję, „Nie wliczono: …” · 7.11
- Grupy składników (nagłówki, liczniki, odhaczanie grupy) · 14.11
- Linki ze składnika do przepisu na ten składnik (np. bułki do burgera) + blok „Przydatne przepisy” · 17.11
- Wskazówki pod krokami (zamiast jednej wskazówki na końcu przepisu) · 18.11
- Sprzęt potrzebny do przepisu, karta nad krokami · 18.11
- Zamienniki składnika: wybór zmienia makro, alergeny i listę zakupów; „pół na pół” · 18.11

### Lista zakupów `/recipes/:id/zakupy` · 14.11
- Osobna strona z linkiem do wysłania. Pozycje według działów sklepu, sumowanie tych samych składników, linia „do: …”, druk A4. Odhaczenia niezależne od gotowania.
- Składnik może być pominięty na liście (np. woda) · 17.11

### Formularz przepisu `/recipes/new`, `/recipes/:id/edit` (admin)
- Bez autozapisu: `Zapisz szkic` / `Opublikuj`. Walidacja przed publikacją, ostrzeżenie przy wyjściu, przeciąganie wierszy · 14.10
- Zdjęcie: **nie działa dla właściciela** — patrz „W toku”, zdjęcia przepisów.
- Ustawianie linków do przepisów, wskazówek, sprzętu i zamienników · 17–18.11

### Panel `/admin` (admin)
- Przepisy: jedyna lista admina, zakładki statusów, usuwanie z potwierdzeniem · 14.10
- Składniki: alergeny, drzewo (rodzic → odmiany), makro jako zakres min–max, przeliczniki jednostek, dziedziczenie makro z podskładników · 7.11
- Klasyfikacja: osie (typ dania, baza, kuchnia, diety, praktyczne), jednostki z odmianą („1 ząbek / 3 ząbki”), słownik sprzętu · 14.10, sprzęt 18.11

## Następne i pomysły

Kolejność ustalamy w rozmowie. Tu jest pula.

**Feedback od znajomego (20.11, gotował z aplikacji):**
- W krokach brakuje ilości: „dodaj mąkę”, a nie pamięta ile, tym bardziej gdy mąka jest podzielona na zaczyn i ciasto. Chce ilości w treści kroku, najlepiej wyróżnionych kolorem → potwierdza pomysł „składniki w treści kroku” (niżej), podnosi go na pierwsze miejsce.
- Tryb gotowania: trzeba przewijać w dół, żeby doczytać krok, a ręce w mące brudzą telefon. Krok ma się mieścić na ekranie (krótsze kroki albo mniejszy font), a „Dalej” wcisnąć byle czym.

**NASTĘPNY (22.11, właściciel: „ważny i prosty”, robimy jako kolejny):**
- **Lista pomysłów dla admina.** Właściciel trzyma dziś w telefonie, co chce zrobić albo ugotować. W panelu ma być lista pomysłów (krótki tytuł) z możliwością dopisywania notatek do każdego. Tylko dla admina, bez widoku publicznego. Do rozpisania: telefon (właściciel notuje w biegu — panel jest dziś desktopowy), status (do zrobienia / zrobione / przepis powstał), powiązanie z przepisem albo szkicem, kolejność. Kolejny krok: brief dla designera.

**Pomysły właściciela (22.11, odzyskane po zgubionych wiadomościach):**
- **Rozmiar porcji** — design zaakceptowany 22.11 (`Rozmiar porcji.html`; tryby servings / sized / form, skala z rozmiarów autora albo z pól dna), wdrożenie w toku (sesja „rozmiar porcji”, PR, migracja 0014 → kopia bazy przed merge'em); plan zaakceptowany 22.11: `portion_mode`, tabela `recipe_piece_sizes`, jednostka „placek”, jedna funkcja `portionScale`, URL `?sztuki=&rozmiar=` / `?forma=o26`. PR #13; decyzje z review (22.11): jednostki będącej nazwą sztuki nie da się usunąć bez przeniesienia; ostrzeżenia „za cienkie / podziel na dwie formy” tylko w trybie formy. Kopia bazy 22.11 14:10 UTC. **Na produkcji 22.11** (PR #13, 7ccc347; migracja 0014 przeszła, jednostka „placek” jest, API bez błędów po wdrożeniu). Następny krok: pizza NY (id 8) na tryb „sztuki w rozmiarach”. Po wdrożeniu: pizza NY (id 8) na tryb sized. Przy cieście na pizzę porcja to nie „1 osoba”, tylko placek danego rozmiaru: 30 cm ≈ 250 g ciasta, 40 cm ≈ 420 g. Czytelnik wybiera liczbę i rozmiar („2 placki po 40 cm”), a ilości skalują się z łącznej masy. Ten sam mechanizm pasuje do ciast w formie (średnica formy → skala po powierzchni). Zastąpiłby „ściągawkę rozmiarów” wpisaną w uwagi albo wskazówki.
- **Blok „Informacje”** obok Uwag — ten sam styl i miejsce, ale treść informacyjna zamiast ostrzeżeń. Właściciel nie jest pewien, czy potrzebny; główny przypadek (rozmiary pizzy) lepiej załatwia rozmiar porcji.

**Pomysły właściciela (18.11, w analizie):**
- Składniki w treści kroku, z ilością skalowaną razem z porcjami. Dziś w krokach nie piszemy ilości, bo skaler ich nie zmienia.
- Słowniczek: dymki przy pojęciach w krokach („reakcja Maillarda”, „deglasowanie”). Zaczątek Wiki. Ten sam mechanizm „odwołania w treści kroku” co ilości wyżej.
- Kroki opcjonalne. Dziś opcjonalny może być składnik, a krok nie. Krok opcjonalny nie musi mieć składnika opcjonalnego: może dotyczyć samego przepisu (np. technika, dłuższe odstawienie).
- Nadkategoria jako składnik przepisu („ser twardy” zamiast wymieniania Grana Padano, pecorino…). Drzewo składników już jest (`ParentID`, makro dziedziczone z podskładników) i formularz pozwala wybrać rodzica. Brakuje tego, co widzi czytelnik.

**Mniejsze, odłożone w backlogu jako „później”:**
- Strona 404, przepis bez zdjęcia, strona logowania.
- Czytelne adresy (slugi) i meta OG, czyli podgląd linku po wysłaniu przepisu albo listy zakupów. Rośnie w znaczeniu, kiedy aplikacja ma iść do ludzi.
- „Cofnij” po usunięciu.

**Rozszerzenia istniejących featurów (zapisane, nieplanowane):**
- Filtr po sprzęcie na liście przepisów.
- Zamiennik z inną ilością (drożdże świeże 25 g ↔ suche 7 g).
- Lista zakupów z kilku przepisów naraz.

**Z pierwotnej wizji (czerwiec 2026), nierozpisane:**
- Przepisy wielopoziomowe: kroki w wersji podstawowej / średniej / mistrzowskiej, z coraz większą liczbą szczegółów.
- Kalkulator jednostek imperialnych jako osobna strona, bez zmian w bazie.

**Duże obszary:**
- **Wyzwania** → w toku (brief dla designera, 18.11). Ustalenia z właścicielem: trening i zabawa, zadania w stylu kulinarnych programów konkursowych, ale ogólne (nic pod konkretny program); bez linku do przepisu; trudność 1–5 w atrakcyjnej formie; limit czasu tylko jako informacja, bez zegara; lista + „Losuj”; „zrobione” na urządzeniu; filtry później; bez zadań zespołowych.
- **Wiki.** Składniki (opis, przechowywanie, odmiany, zamienniki, przepisy z tym składnikiem), techniki, style gotowania; słowa kluczowe w krokach z dymkiem i linkiem. Publiczna karta składnika z makro jest zaprojektowana (`Makro.html`) i odłożona do Wiki.

## Ryzyka i dług

1. ~~Zapisy w API bez autoryzacji~~ — zamknięte 18.11 (c65e2cb): zapis tylko z tokenem właściciela albo kluczem `API_ADMIN_KEY`, szkice niewidoczne publicznie. Rotacja klucza: `CLAUDE.md`.
2. **`.env.local` wskazuje na produkcję.** Każda lokalna sesja pracuje na prawdziwych danych. `resetdb` jest zakazany.
3. **Zdjęcia w bazie** (tabela `photos`, bytea) — rozwiązanie tymczasowe od 10.2026. Wgrywanie własnych zdjęć odłożone 18.11 do przenosin infrastruktury (kod w git stash). Codzienny backup obejmuje tylko bazę.
4. **Infrastruktura do przemyślenia.** Baza (darmowy Neon, 0,5 GB), a może i API (limity Vercel serverless) pewnie kiedyś się przeniosą. Nie wiadomo jeszcze kiedy ani dokąd. Do tego czasu nowe rzeczy nie powinny się mocno wiązać z Vercelem.
5. **Brak kopii zapasowej bazy (wysokie).** Neon na darmowym planie przywraca tylko do 6 godzin wstecz i daje 1 ręczny snapshot. Kopia przepisów w `content/` to stan z dnia dodania, a nie backup: poprawki zrobione w aplikacji do niej nie wracają. → W toku.
7. **Test `web` `recipeForm.test.ts` pada od c65e2cb** (autoryzacja): `createAuth.ts` czyta `location` przy imporcie. Nienaprawione.
6. **Repo jest publiczne.** Wszystko, co trafia do repo, logów i artefaktów GitHub Actions, jest widoczne dla świata.
8. ~~Limit 12 funkcji na Vercelu~~ — zamknięte 19.11 (290b90c): całe API to jedna funkcja z routerem, nowy endpoint = linijka w `router/routes.go`. Historia: Każdy plik w `api/` to osobna funkcja. Od Wyzwań (19.11) nowe endpointy doklejamy do istniejących plików (wyzwania siedzą w `api/taxonomy/units.go`). Wyjścia: jeden plik-router na całe API (znosi limit na stałe), plan Pro, albo przenosiny API (pkt 4). Decyzja 19.11: router (prompt `docs/prompt-claude-code-router-api.md`).

## Otwarte pytania

- Kiedy aplikacja ma wyjść do ludzi? Od tego zależy, kiedy robimy slugi, meta OG, 404 i przepis bez zdjęcia.
- `ml = 1 g`, `l = 1000 g` bez przelicznika — rekomendacja z makro: zostawić. Do potwierdzenia.
- Wyzwania: czy wracamy do tego filaru, a jeśli tak, od jakiego typu wyzwania zaczynamy?
