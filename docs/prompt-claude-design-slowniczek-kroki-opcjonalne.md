# Prompt dla Claude Design — słowniczek pojęć i kroki opcjonalne

Dwa osobne zadania. Rób je po kolei, każde z osobnym handoffem i wpisem w „Decyzje i backlog.md”. Ogólne zasady jak dotąd:
- DS v2 i tokeny z istniejących prototypów,
- strony publiczne działają na telefonie (375 px), admin pracuje na desktopie,
- hover nie przesuwa ani nie skaluje tekstu,
- ułamków nie pokazujemy znakami ½ ⅓ ¼ — tylko „1,5” i „1/3”.

---

## 1. Słowniczek pojęć w krokach

### Kontekst
W treści kroku są już **odwołania do składników**: autor pisze „dodaj mąkę”, a za słowem stoi znacznik ilości `[60 g]`. Decyzja i prototyp: „Ilości składników w treści kroku”, `Kroki i ilosci.html`.

Przy tamtym designie zapisałeś na później: **słowniczek** to drugi typ odwołania, „inny wygląd — pojęcie bez tła, kropkowane podkreślenie”. Ten moment właśnie przyszedł.

### Problem
W krokach i wskazówkach padają pojęcia, których część czytelników nie zna: „deglasuj”, „reakcja Maillarda”, „sos ma mrugać”, „zaszklij cebulę”, „al dente”. Dziś autor ma dwa wyjścia: tłumaczy je w każdym przepisie od nowa, co zaśmieca krok, albo nie tłumaczy wcale.

### Feature
- **Słowniczek** to lista pojęć prowadzona przez właściciela. Każde pojęcie ma:
  - nazwę,
  - krótkie wyjaśnienie (1–3 zdania, zwykły tekst),
  - opcjonalnie „więcej”, czyli dłuższy opis. Rozstrzygnij, czy to potrzebne teraz, czy dopiero przy Wiki.
- **W kroku i we wskazówce** autor oznacza słowo jako pojęcie ze słownika. Słowo zostaje w formie, w jakiej je napisał („deglasuj”, „deglasowania”). Odwołanie tylko je oznacza, niczego nie dokleja. Tym różni się od znacznika ilości.
- **Czytelnik** widzi oznaczone słowo (kropkowane podkreślenie albo coś, co zaproponujesz) i po dotknięciu lub kliknięciu dostaje wyjaśnienie.

### Do rozstrzygnięcia
1. **Wygląd pojęcia** w zwykłym tekście kroku, we wskazówce i w dużym tekście trybu gotowania. Ma być odróżnialne od znacznika ilości `[60 g]`.
2. **Wyjaśnienie:**
   - telefon: arkusz od dołu czy dymek,
   - desktop: dymek po kliknięciu czy po najechaniu (z dostępnością z klawiatury),
   - tryb gotowania: jedno dotknięcie, duży tekst, łatwe zamknięcie brudnym palcem.
3. **Panel:**
   - gdzie żyje lista pojęć (zakładka w Klasyfikacji, osobna pozycja w menu?),
   - dodawanie, edycja, usuwanie,
   - co się dzieje z odwołaniami do usuniętego pojęcia (wzorem innych słowników: blokada albo przeniesienie).
4. **Formularz przepisu:** jak autor oznacza słowo. Dla składników jest „@” i lista wierszy. Dla pojęć zaproponuj np. zaznaczenie słowa i „Pojęcie…” albo inny znak. Jest też opcja „dodaj nowe pojęcie” wprost z formularza.
5. **Jedno pojęcie kilka razy w przepisie:** czy oznaczać za każdym razem, czy tylko za pierwszym (zalecenie dla autora).
6. **Wiki (później):** zostaw miejsce, żeby kiedyś z wyjaśnienia prowadził link „Więcej w Wiki”. Strony Wiki nie projektujesz.

### Stany brzegowe
- Pojęcie w środku znacznika składnika, czyli obok siebie „deglasuj winem [100 ml]”.
- Dwa pojęcia w jednym zdaniu.
- Bardzo długie wyjaśnienie.
- Pojęcie usunięte ze słownika, a wciąż użyte w kroku.

### Poza zakresem
- Strona Wiki i publiczna lista pojęć. Zaznacz tylko, gdzie kiedyś byłoby wejście.
- Automatyczne wykrywanie pojęć w tekście.
- Zdjęcia w wyjaśnieniach.

### Do oddania
- **Artboardy:**
  - krok i wskazówka z pojęciem: light, dark, 375 px i desktop,
  - otwarte wyjaśnienie,
  - tryb gotowania na 375 px,
  - panel: lista pojęć i formularz pojęcia,
  - formularz przepisu: oznaczanie pojęcia.
- **Handoff:** checklista odbioru, sekcja „Nie zmieniaj”, otwarte pytania.

---

## 2. Kroki opcjonalne

### Kontekst
Składnik w przepisie może być oznaczony jako opcjonalny, wtedy ma dopisek „(opcjonalnie)”, a na liście zakupów trafia do sekcji „jeśli chcesz”. Krok takiego oznaczenia nie ma.

Kroki mają już:
- sekcje (np. „Zaczyn”, „Ciasto”),
- wskazówki,
- flagę „przygotuj wcześniej”,
- odhaczanie i licznik postępu,
- tryb gotowania (jeden krok naraz, „Dalej” odhacza krok).

### Problem
Część kroków jest do wyboru, na przykład:
- „Jeśli chcesz, odstaw ciasto na noc do lodówki”,
- „Opal wierzch palnikiem”,
- „Przetrzyj sos przez sito”.

Krok opcjonalny **nie musi** dotyczyć składnika opcjonalnego. Często dotyczy samego przepisu: techniki, dłuższego odstawienia, wykończenia. Dziś autor pisze „jeśli chcesz…” w treści, przez co krok wygląda jak obowiązkowy i wlicza się do postępu.

### Feature
Autor może oznaczyć krok jako **opcjonalny**. Czytelnik od razu widzi, że krok można pominąć.

### Do rozstrzygnięcia
1. **Wygląd kroku opcjonalnego** na liście kroków: telefon i desktop, light i dark. Musi być spójny z „(opcjonalnie)” przy składnikach i sprzęcie.
2. **Postęp:**
   - czy krok opcjonalny liczy się do „N z M”,
   - czy przepis jest „ukończony” bez odhaczenia kroków opcjonalnych.
3. **Tryb gotowania:** jak pokazać krok opcjonalny i czy da się go pominąć jednym ruchem („Pomiń” obok „Dalej”?). „Dalej” dziś odhacza krok.
4. **Krok opcjonalny z odwołaniami do składników.** Na przykład krok opcjonalny używa składnika wymaganego. Czy coś z tego wynika dla listy zakupów? Domyślnie nie, lista zakupów wynika tylko ze składników. Potwierdź albo zaproponuj inaczej.
5. **Formularz:** przełącznik „Krok opcjonalny” przy kroku, wzorem „Opcjonalny” przy sprzęcie.

### Stany brzegowe
- Kilka kroków opcjonalnych pod rząd.
- Ostatni krok opcjonalny, czyli kiedy pojawia się „Gotowe” w trybie gotowania.
- Krok opcjonalny z flagą „przygotuj wcześniej”.
- Krok opcjonalny ze wskazówkami.

### Poza zakresem
- Warianty przepisu: wersja podstawowa i mistrzowska. To osobny, przyszły temat.
- Zależności między krokami („jeśli zrobiłeś krok 4, to…”).

### Do oddania
- **Artboardy:**
  - lista kroków z krokiem opcjonalnym: light, dark, 375 px i desktop,
  - tryb gotowania na 375 px,
  - formularz.
- **Handoff:** checklista odbioru, sekcja „Nie zmieniaj”, otwarte pytania.
