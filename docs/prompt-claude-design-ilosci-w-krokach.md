# Prompt dla Claude Design — ilości składników w treści kroku

## Kontekst

Uszata Kuchnia to aplikacja z przepisami właściciela. Do czytania jest publiczna, bez kont użytkowników. Admin jest jeden i pracuje na desktopie, a strony publiczne muszą działać na telefonie. Tokeny i komponenty DS v2 bierz z istniejących prototypów (`Przepis.html`, `Wskazowki.html`, `Zamienniki.html`, `Uwagi.html`, `Nowy przepis.html`, „Poprawki 18.11”).

Jak dziś działa przepis:
- **Składniki są w grupach.** Na przykład „Zaczyn” i „Ciasto”, a ten sam składnik (mąka) może być w obu grupach z inną ilością.
- **Licznik porcji przelicza ilości na liście składników.** Na liście zakupów też. Ułamki i zakresy wyglądają tak: „1–2 sztuki”, „2½ łyżki”.
- **Zamienniki:** użytkownik może wybrać zamiennik składnika (wieprzowina zamiast wołowiny). Wybór zmienia listę składników, makro i listę zakupów.
- **Kroki to zwykły tekst** z wskazówkami pod spodem. Celowo **nie piszemy w nich ilości**, bo wpisana ręcznie liczba nie zmieniłaby się razem z porcjami.
- **Tryb gotowania:** jeden krok naraz na pełnym ekranie, duże „Wstecz” / „Dalej”, arkusz „Sprzęt i składniki”.

## Problem (feedback od znajomego, który gotował z aplikacji)

> „Dodaj mąkę” — a ja nie pamiętam, ile tej mąki ma być. Tym bardziej że było rozdzielone: ileś na zaczyn, ileś na ciasto. Można by w treści dodać, ile konkretnie sztuk/gramów, nawet jakoś kolorem wyróżnić.

Przy kuchence nikt nie skacze między krokiem a listą składników, a podział składnika na grupy tylko pogarsza sprawę.

## Feature: odwołanie do składnika w treści kroku

Autor wstawia w tekst kroku **odwołanie do konkretnego wiersza składnika** (np. „mąka” z grupy „Zaczyn”). Czytelnik widzi w tym miejscu ilość, wyróżnioną. Ilość:
- **przelicza się z porcjami**, tak samo jak na liście składników (te same zasady zaokrągleń),
- **reaguje na zamiennik**: gdy wybrano zamiennik, odwołanie pokazuje zamiennik i jego ilość.

### Do rozstrzygnięcia przez ciebie

1. **Gramatyka — najtrudniejsze.** Nazwa składnika w bazie jest w mianowniku („Mąka pszenna”), a w zdaniu potrzebna jest inna forma („dodaj 250 g **mąki**”, „wsyp **mąkę**”). Warianty do rozważenia:
   - **(a)** autor pisze słowo sam („dodaj mąkę”), a odwołanie dokleja ilość: „dodaj mąkę **250 g**”,
   - **(b)** odwołanie wstawia ilość i nazwę w stałej formie, np. „mąka pszenna · 250 g”, styl „chip”,
   - **(c)** inny pomysł.

   Weź pod uwagę zamienniki: w (a) słowo autora „wołowinę” nie zmieni się na „wieprzowinę”. Rozstrzygnij, co wtedy widzi czytelnik.
2. **Wygląd odwołania.** Ma być wyraźne („kolorem albo czymś”), ale nie może rozbijać czytania zdania. Musi działać w tekście zwykłym i w dużym tekście trybu gotowania, w jasnym i ciemnym motywie. Hover nie może niczego przesuwać ani skalować (zasada z 20.11).
3. **Część wiersza.** Krok często zużywa część składnika: „dodaj połowę mleka”, „resztę mąki”. Czy odwołanie obsługuje ułamek wiersza (½, ⅓, reszta)? Jeśli tak, jak to wygląda?
4. **Interakcja.** Czy dotknięcie odwołania coś robi, np. pokazuje wiersz na liście składników albo odhacza składnik? Domyślnie nic. Zaproponuj interakcję tylko wtedy, gdy realnie pomaga.
5. **Formularz (panel, desktop).** Jak autor wstawia odwołanie w treści kroku, np. „@” i lista wierszy składników z nazwą grupy? Jak wygląda odwołanie w polu tekstowym? Co się dzieje, gdy wiersz składnika zostanie usunięty albo zmieni grupę?
6. **Słowniczek (później, nie wdrażamy teraz).** Ten sam mechanizm posłuży kiedyś do pojęć z dymkiem („deglasowanie”, „reakcja Maillarda”). Nie projektuj słowniczka. Zadbaj tylko o to, żeby odwołanie do składnika i przyszłe odwołanie do pojęcia dały się od siebie odróżnić.

### Stany brzegowe
- Ten sam składnik w dwóch grupach (mąka na zaczyn i na ciasto) w jednym kroku.
- Składnik opcjonalny w kroku.
- Składnik bez ilości („sól do smaku”).
- Jednostka z odmianą („1 ząbek / 3 ząbki”) i zakres („1–2 sztuki”).
- Wybrany zamiennik i „pół na pół”.
- Krok z 4–5 odwołaniami.
- Krok w trybie gotowania na 375 px.

## Poza zakresem
- Słowniczek, Wiki.
- Minutniki w krokach.
- Automatyczne rozpoznawanie składników w tekście. Odwołania wstawia autor albo Claude przez API.
- Przepisanie istniejących przepisów. Zrobi to Claude po wdrożeniu.

## Do oddania
- **Artboardy:** strona przepisu (krok z odwołaniami) w wersji light, dark i 375 px; tryb gotowania 375 px light i dark; formularz na desktopie (wstawianie i edycja odwołania).
- **Handoff:**
  - checklista odbioru,
  - sekcja „Nie zmieniaj”,
  - otwarte pytania,
  - rozstrzygnięcia punktów 1–6 z uzasadnieniem.
- **„Decyzje i backlog.md”:** wpis. Uwaga: zmienia się zasada „w krokach nie piszemy ilości” — teraz ilość przychodzi z odwołania.
