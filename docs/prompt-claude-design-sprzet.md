# Prompt dla Claude Design — sprzęt potrzebny do przepisu

## Kontekst

Uszata Kuchnia, strona przepisu `/recipes/:id`, tryb gotowania i formularz `/recipes/:id/edit`. Rozszerzasz istniejące widoki (DS v2, projekt UszataKuchnia). Nie projektujesz ich od nowa.

## Problem

Przepis nie mówi, czego potrzebujesz poza składnikami. Wychodzi to dopiero w krokach: „wyrabiaj hakiem w robocie planetarnym”, „garnek 3–4 l z grubym dnem”, „skrobak do ciasta”, „kratka do studzenia”. Ktoś zaczyna gotować i w połowie okazuje się, że nie ma odpowiedniego garnka albo że bez robota zajmie to dużo dłużej.

## Feature

Przepis ma listę sprzętu.

- Pozycja sprzętu to nazwa z bazy (słownik sprzętu, jak słownik składników), opcjonalnie z doprecyzowaniem w konkretnym przepisie („3–4 l, grube dno”).
- Pozycja może być wymagana albo opcjonalna / zastępowalna („robot planetarny — albo ręcznie”).
- Sprzęt oczywisty (nóż, deska, łyżka) nie trafia na listę. To zasada treści, nie UI, ale zaznacz, jeśli design powinien ją wspierać.

## Widoki do zaprojektowania

### 1. Strona przepisu

Gdzie i jak lista sprzętu pojawia się na stronie, desktop (przyklejony panel składników, karta makro pod krokami) i mobile. Rozstrzygnij, czy to część panelu składników, osobna karta, czy coś w metadanych na górze. Przy decyzji uwzględnij kolejność sekcji z „Decyzje i backlog.md”.

Do pokazania: 1 pozycja, 5–6 pozycji, pozycja z doprecyzowaniem, pozycja opcjonalna z alternatywą.

### 2. Tryb gotowania

Czy sprzęt jest tu widoczny (np. na starcie, jako „przygotuj”), czy tylko na stronie przepisu.

### 3. Formularz przepisu

Dodawanie, kolejność i usuwanie pozycji sprzętu, wybór ze słownika z wyszukiwaniem, pole doprecyzowania, przełącznik wymagany / opcjonalny. Gdzie w formularzu to trafia względem składników i kroków.

### 4. Słownik sprzętu (admin)

Lista sprzętu w panelu admina, wzorem listy składników i jednostek: dodawanie, zmiana nazwy, usuwanie pozycji używanej w przepisach.

### 5. Stany brzegowe

- Przepis bez sprzętu: sekcja ukryta czy z komunikatem.
- Bardzo długa nazwa albo doprecyzowanie.

## Poza zakresem (nie projektuj)

- Filtrowanie listy przepisów po sprzęcie („mam tylko mikrofalę”). Możesz zaznaczyć, gdzie by kiedyś pasowało.
- Ikony / zdjęcia sprzętu.
- Linki do sprzętu w treści kroków.

## Do oddania

Artboardy: light, dark, 375 px, dla widoków 1 i 3 (oraz 4 w wersji desktop), handoff z checklistą odbioru i wpisem w „Decyzje i backlog.md”.
