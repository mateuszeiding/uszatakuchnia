# Prompt dla Claude Design — zamienniki składnika w przepisie

## Kontekst

Uszata Kuchnia, strona przepisu `/recipes/:id`, tryb gotowania i formularz `/recipes/:id/edit`. Rozszerzasz istniejące widoki (DS v2, projekt UszataKuchnia). Nie projektujesz ich od nowa.

## Problem

W wielu przepisach jeden składnik da się wymienić na inny bez zmiany przepisu. Przykład: ragù alla bolognese robi się na mielonej wołowinie, ale równie dobrze na wieprzowinie albo pół na pół. Dziś trzeba wybrać jedno i resztę dopisać w notatce („lub wieprzowina”), więc zamiennik nie jest prawdziwym składnikiem: nie ma makro, nie trafia do listy zakupów, nie da się go kliknąć.

## Feature

Wiersz składnika w konkretnym przepisie może mieć zamienniki: inne składniki z bazy, którymi można go zastąpić.

- Zamienniki są na poziomie wiersza w przepisie, nie globalnie. Wołowina w ragù ma zamiennik „wieprzowina”, w tatarze nie ma żadnego.
- Jeden składnik jest główny (ten, którego używa autor). Zamienników może być kilka.
- Wariant mieszany („pół na pół”) jest częsty przy mięsie. Rozstrzygnij, czy to osobny zamiennik, czy coś innego.
- Ilość i jednostka: na start zamiennik ma tę samą ilość co składnik główny. Jeśli widzisz przypadek, w którym to nie wystarczy, opisz go, ale go nie projektuj.

## Widoki do zaprojektowania

### 1. Strona przepisu — panel składników

Jak pokazać, że składnik ma zamienniki, desktop (przyklejony panel) i mobile.

**Ograniczenie:** cały wiersz składnika to dziś przycisk do odhaczania. Zamienniki nie mogą z tym kolidować (osobne cele kliknięcia, min. 44 px na telefonie).

Rozstrzygnij:
- Czy użytkownik może wybrać zamiennik (i wtedy zmienia się makro i lista zakupów), czy zamienniki to tylko informacja?
- Jeśli wybór: czy jest zapamiętywany i gdzie widać, że przepis jest „w wariancie”.

Do pokazania: wiersz z jednym i z kilkoma zamiennikami, wiersz z notatką, odhaczony wiersz, długie nazwy.

### 2. Makro i lista zakupów

Karta wartości odżywczych liczy makro ze składników. Rozstrzygnij, co pokazuje, gdy składnik ma zamienniki: tylko składnik główny, wybrany wariant, czy zakres. Lista zakupów (osobny brief) bierze wybrany wariant albo główny.

### 3. Tryb gotowania

Czy zamienniki są tu widoczne i czy da się je zmienić w trakcie gotowania.

### 4. Formularz przepisu — dodawanie zamienników

W edytorze składników (`IngredientsEditor`: autosuggest składnika, ilość, jednostka, notatka, „opcjonalnie”) trzeba móc dodać, zobaczyć i usunąć zamienniki wiersza. Wiersz już jest gęsty, zwłaszcza na 375 px. Jeśli w tym samym miejscu ląduje link do innego przepisu (osobny feature „linki do przepisów”), zaprojektuj oba razem.

### 5. Stany brzegowe

- Zamiennik bez danych makro.
- Zamiennik z innym alergenem niż główny składnik (np. główny bez glutenu, zamiennik z glutenem). Czy alergeny przepisu to suma, czy zależą od wybranego wariantu?
- Składnik główny usunięty z bazy.

## Poza zakresem (nie projektuj)

- Globalne zamienniki na stronie składnika (planowane osobno w opisie składnika).
- Inna ilość dla zamiennika i przeliczniki między nimi.
- Wymiana całych grup składników („wersja wegetariańska”).

## Do oddania

Artboardy: light, dark, 375 px, dla widoków 1, 2 i 4, oraz handoff z checklistą odbioru i wpisem w „Decyzje i backlog.md”.
