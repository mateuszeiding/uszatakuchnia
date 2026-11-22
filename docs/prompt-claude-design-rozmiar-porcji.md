# Prompt dla Claude Design — rozmiar porcji: placki, formy, blachy

Ogólne zasady jak dotąd: DS v2 i tokeny z istniejących prototypów, strony publiczne działają na telefonie (375 px), admin na desktopie, hover nie przesuwa ani nie skaluje tekstu.

**Ułamki:** połówki i ćwiartki dziesiętnie („1,5”, „0,25”), tylko tercje jako „1/3”. Żadnych znaków ½ ¼.

## Kontekst
Strona przepisu ma **licznik porcji** (−/+). Zmiana porcji przelicza:
- ilości na liście składników,
- znaczniki ilości w krokach (`[300 g]`),
- listę zakupów,
- makro na porcję.

Zaokrąglenia i jednostki działają jak dziś. Porcja zawsze znaczy „1 osoba” albo „1 sztuka”.

## Problem
Przy wypiekach porcja to nie osoba, tylko **rozmiar**:
- **Pizza:** przepis jest na 2 placki po 30 cm (ok. 255 g ciasta na placek). Ktoś chce 2 placki po 40 cm (ok. 420 g na placek) albo 3 po 25 cm. Dziś autor dopisuje „ściągawkę rozmiarów” do wskazówek i czytelnik liczy sam.
- **Ciasto w formie:** przepis jest na tortownicę 24 cm, a ktoś ma 26 cm. Ilości trzeba przeliczyć po powierzchni: 26² / 24² ≈ 1,17.
- **Blacha:** przepis jest na blachę 25 × 35 cm, a ktoś ma 30 × 40 albo okrągłą formę.

Właściciel: „przy cieście możesz sobie ustawić, jaką masz blachę, i ci przeliczy”.

## Feature: rozmiar porcji
Autor przepisu wybiera, **czym jest porcja**. Czytelnik ustawia rozmiar (i liczbę), a aplikacja przelicza ilości. Robocze tryby — nazwy i podział do twojego rozstrzygnięcia:
1. **Porcje / sztuki:** jak dziś.
2. **Sztuki w rozmiarach** (np. pizza): autor podaje kilka rozmiarów i masę na sztukę („25 cm — ok. 190 g”, „30 cm — 255 g”…). Czytelnik wybiera rozmiar i liczbę sztuk, a skala = (masa na sztukę × liczba) / masa w przepisie.
3. **Forma lub blacha:** autor podaje formę z przepisu (okrągła Ø albo prostokątna a × b). Czytelnik wpisuje albo wybiera swoją, a skala = stosunek powierzchni. Okrągła ↔ prostokątna też powinna działać.

## Do rozstrzygnięcia przez ciebie
1. **Kontrolka na stronie przepisu** w miejscu licznika porcji: jak wygląda w każdym trybie na telefonie i desktopie. Uwzględnij przyklejony licznik na telefonie i przyklejony panel składników na desktopie.
2. **Komunikat o przeliczeniu:** jak czytelnik widzi, że ilości są przeliczone pod inną formę lub rozmiar (wzorem „Z zamiennikami: n · Przywróć przepis”).
3. **Wpisywanie własnej formy:** wpisywanie wymiarów czy lista typowych rozmiarów (tortownice 18–28 cm, blachy domowe)? Co przy rozmiarze spoza zakresu, np. 2 razy większej formie?
4. **Czas pieczenia i grubość:** przy innej formie zmienia się wysokość albo czas. Czy pokazujemy jedno zdanie ostrzeżenia, np. „Czas pieczenia może się zmienić”? Przeliczania czasu nie robimy.
5. **Makro:** „na porcję” zmienia sens. Na placek? Na kawałek? Na całą formę? Rozstrzygnij, co pokazuje karta „Wartości odżywcze” w trybach 2 i 3.
6. **Lista zakupów, tryb gotowania i link do listy zakupów** (dziś niesie `?porcje=`): co niosą przy rozmiarze.
7. **Formularz przepisu (panel):** jak autor wybiera tryb i wpisuje rozmiary lub formę.
8. **Karta na liście przepisów:** czy coś się zmienia (np. „na 2 placki 30 cm” zamiast „2 porcje”).

## Stany brzegowe
- Pizza: 1 placek 45 cm i 4 placki 25 cm.
- Tortownica 24 → 26 cm i 24 → 18 cm.
- Blacha 25 × 35 → tortownica 26 cm.
- Przepis w trybie „forma” z grupami składników (ciasto, kruszonka, polewa).
- Składnik w sztukach („2 jajka”) przy skali 1,17 — zaokrąglenie jak dziś.

## Poza zakresem
- Przeliczanie czasu i temperatury pieczenia.
- Inna grubość ciasta przy tej samej formie.
- Przeliczanie jednostek (cale ↔ cm). Wymiary podajemy w cm. Cale co najwyżej jako podpis, jeśli uznasz, że pomogą przy pizzy.

## Do oddania
- **Artboardy:** strona przepisu w trybach 2 i 3 (light, dark, 375 px i desktop), formularz w panelu (desktop), stany brzegowe w miarę potrzeby.
- **Handoff:** checklista odbioru, sekcja „Nie zmieniaj” (dzisiejszy licznik porcji dla zwykłych przepisów), model danych (pola i typy), otwarte pytania.
- **Wpis w „Decyzje i backlog.md”.**
