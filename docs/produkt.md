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
| Poprawki z przeglądu 18.11 | do rozdzielenia: inżynier / designer / treść | `docs/przeglad-uzycia-2026-11-18.md` |
| Wyzwania | design zaakceptowany (artefakt: Handoff → „Wyzwania”, prototyp `Wyzwania.html`); **na produkcji 19.11** (push po kopii bazy z 18:18 UTC; migracja 0009 przeszła, 3 formaty startowe są). Pierwsze wyzwanie dodane przez właściciela 20.11. Zakładki Klasyfikacji: przewijanie (decyzja 19.11). Do potwierdzenia: 7 decyzji inżyniera | `docs/prompt-claude-design-wyzwania.md` |
| Jedno wejście do API (router) | **na produkcji 19.11** (290b90c), sprawdzone curlem: odczyty 200, zapis bez klucza 401, nieznany adres 404. Zmiana celowa: nieznany `/api/…` → 404 JSON zamiast index.html. Znosi limit 12 funkcji Vercela (ryzyko 8), adresy API bez zmian. Bez migracji | `docs/prompt-claude-code-router-api.md` |
| Hover kart bez rozmytego tekstu | design zaakceptowany 20.11 (wariant B: ramka `--rule-strong` + tytuł w akcencie, przycisk akcentowy ciemnieje zamiast podjeżdżać; zasada „hover nie przesuwa tekstu”); wdrożone lokalnie 20.11 (9ce7752), niewypchnięte, bez migracji. Odkrycie: lista składników (`IngredientsView`) nie ma trasy — martwy widok | `docs/prompt-claude-design-hover-kart.md` |
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
