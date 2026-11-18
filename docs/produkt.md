# Uszata Kuchnia — stan produktu

Aktualizacja: 18.11.2026 · prowadzi: analityk (sesja Claude Code od featurów)

Żywy dokument: co jest, co się dzieje, co dalej. Szczegóły wyglądu i zachowania są w projekcie Claude Design — [artefakt](https://claude.ai/artifact/YEuK2L22otH24hdvcUC5GR), pliki `project/Decyzje i backlog.md` (decyzje) i `project/Handoff - Claude Code.md` (prompty i checklisty wdrożeń). Ten plik ich nie powtarza, tylko do nich odsyła.

## Po co

Aplikacja z przepisami właściciela, publiczna do czytania, bez kont użytkowników. Jeden autor tworzy całą treść.

Trzy filary:
1. **Przepisy** — wpisywane ręcznie, z naciskiem na gotowanie przy kuchence: skalowanie porcji, postęp, tryb gotowania, lista zakupów.
2. **Wyzwania** — lista zadań kulinarnych w duchu MasterChefa. Jeszcze niezaczęte.
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
| Autoryzacja zapisów w API | brief dla inżyniera | `docs/prompt-claude-code-autoryzacja.md` |

Wdrożone, jeszcze niesprawdzone na żywej stronie:
- linki do przepisów i blok „Przydatne przepisy” (17.11 — żaden opublikowany przepis nie miał jeszcze linku),
- tryb gotowania,
- wskazówki przy krokach, sprzęt, zamienniki (18.11).

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
- Bez autozapisu: `Zapisz szkic` / `Opublikuj`. Walidacja przed publikacją, ostrzeżenie przy wyjściu, przeciąganie wierszy, kadrowanie zdjęcia do 1:1 · 14.10
- Ustawianie linków do przepisów, wskazówek, sprzętu i zamienników · 17–18.11

### Panel `/admin` (admin)
- Przepisy: jedyna lista admina, zakładki statusów, usuwanie z potwierdzeniem · 14.10
- Składniki: alergeny, drzewo (rodzic → odmiany), makro jako zakres min–max, przeliczniki jednostek, dziedziczenie makro z podskładników · 7.11
- Klasyfikacja: osie (typ dania, baza, kuchnia, diety, praktyczne), jednostki z odmianą („1 ząbek / 3 ząbki”), słownik sprzętu · 14.10, sprzęt 18.11

## Następne i pomysły

Kolejność ustalamy w rozmowie. Tu jest pula.

**Mniejsze, odłożone w backlogu jako „później”:**
- Strona 404, przepis bez zdjęcia, strona logowania.
- Czytelne adresy (slugi) i meta OG, czyli podgląd linku po wysłaniu przepisu albo listy zakupów. Rośnie w znaczeniu, kiedy aplikacja ma iść do ludzi.
- „Cofnij” po usunięciu.

**Rozszerzenia istniejących featurów (zapisane, nieplanowane):**
- Filtr po sprzęcie na liście przepisów.
- Zamiennik z inną ilością (drożdże świeże 25 g ↔ suche 7 g).
- Lista zakupów z kilku przepisów naraz.

**Duże obszary:**
- **Wyzwania.** Zarys z planu z czerwca: lista do przeglądania (bez logowania, bez punktów), trudność 1–3, typ wyzwania (konkretne danie, z podanych składników, kuchnia, technika), opcjonalny link do przepisu, „ukończone” zapamiętane na urządzeniu. Brak konceptu, designu i decyzji.
- **Wiki.** Składniki (opis, przechowywanie, odmiany, zamienniki, przepisy z tym składnikiem), techniki, style gotowania; słowa kluczowe w krokach z dymkiem i linkiem. Publiczna karta składnika z makro jest zaprojektowana (`Makro.html`) i odłożona do Wiki.

## Ryzyka i dług

1. **Zapisy w API bez autoryzacji (wysokie).** Handlery Go nie sprawdzają tokenu: każdy może tworzyć, zmieniać i usuwać przepisy, składniki, zdjęcia i słowniki na produkcji. Szkice też da się odczytać (`?admin=true` na liście, szczegóły przepisu bez sprawdzania statusu). → W toku.
2. **`.env.local` wskazuje na produkcję.** Każda lokalna sesja pracuje na prawdziwych danych. `resetdb` jest zakazany.
3. **Zdjęcia w bazie** (tabela `photos`, bytea) zamiast Vercel Blob — rozwiązanie tymczasowe od 10.2026. Baza rośnie razem ze zdjęciami.
4. **Kopia przepisów w `content/` to stan z dnia dodania**, a nie backup bazy. Poprawki zrobione w aplikacji do niej nie wracają.

## Otwarte pytania

- Czy baza ma kopię zapasową poza `content/`? Neon przywraca do punktu w czasie tylko w ramach okna historii, które zależy od planu.
- Kiedy aplikacja ma wyjść do ludzi? Od tego zależy, kiedy robimy slugi, meta OG, 404 i przepis bez zdjęcia.
- `ml = 1 g`, `l = 1000 g` bez przelicznika — rekomendacja z makro: zostawić. Do potwierdzenia.
- Wyzwania: czy wracamy do tego filaru, a jeśli tak, od jakiego typu wyzwania zaczynamy?
