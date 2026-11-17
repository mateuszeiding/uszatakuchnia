# Prompt dla Claude Design — linki do innych przepisów

## Kontekst

Uszata Kuchnia, strona przepisu `/recipes/:id` i formularz `/recipes/:id/edit`. Rozszerzasz istniejące widoki (DS v2, projekt UszataKuchnia). Nie projektujesz ich od nowa.

## Problem

Część przepisów korzysta z innych przepisów. Przykład: smashburger ma w składnikach bułki, a w aplikacji jest osobny przepis „Bułki do hamburgerów”. Teraz nie da się przejść od jednego do drugiego, więc trzeba szukać na liście.

## Feature

Składnik w przepisie może wskazywać inny przepis z aplikacji. Wtedy przy tym składniku widać link do tego przepisu.

- Link jest opcjonalny, na poziomie wiersza składnika w konkretnym przepisie. Nie dotyczy składnika globalnie: „bułka” w jednym przepisie może linkować do bułek, a w innym nie.
- Jeden wiersz to najwyżej jeden link.
- Składnik zostaje zwykłym składnikiem: ilość, jednostka, notatka, makro i lista zakupów działają jak dotąd. Link tylko dochodzi do wiersza.

## Widoki do zaprojektowania

### 1. Strona przepisu — panel składników

Wiersz składnika z linkiem do przepisu, desktop (przyklejony panel) i mobile.

**Ograniczenie:** cały wiersz składnika to dziś przycisk do odhaczania (checkbox). Link nie może być w środku przycisku, a kliknięcie w link nie może odhaczać składnika. Zaprojektuj tak, żeby odhaczanie i przejście do przepisu były osobnymi, wyraźnymi celami kliknięcia, także na telefonie (min. 44 px).

Zaprojektuj **2–3 warianty** oznaczenia linku, np. ikona/strzałka na końcu wiersza, podkreślona nazwa składnika, mały chip „przepis”.

Do pokazania razem z linkiem: odhaczony wiersz (przekreślenie), wiersz z notatką, wiersz z „(opcjonalnie)” i „też: …”, długa nazwa przepisu.

### 2. Tryb gotowania

Składniki w trybie gotowania mają większe wiersze (ten sam komponent). Czy link ma tu być? Jeśli tak, to jak, żeby przypadkiem nie wyjść z gotowania (np. otwiera w nowej karcie albo prosi o potwierdzenie). Rozstrzygnij.

### 3. Formularz przepisu — ustawianie linku

W edytorze składników (`IngredientsEditor`, wiersz z autosuggestem składnika, ilością, jednostką, notatką, „opcjonalnie”) trzeba móc:
- wybrać przepis do podlinkowania, z wyszukiwaniem po nazwie (przepisów będą dziesiątki),
- zobaczyć, który przepis jest wybrany,
- usunąć link.

Wiersz edytora już jest gęsty, zwłaszcza na 375 px. Zaprojektuj, gdzie to pole trafia, np. rozwijane „więcej” przy wierszu albo osobna linia pod wierszem.

### 4. Stany brzegowe

- **Podlinkowany przepis jest szkicem.** Szkice nie są widoczne na publicznej liście, ale da się je otworzyć po ID. Pokazywać link, ukryć go, czy oznaczyć „szkic” (np. tylko w formularzu)?
- **Podlinkowany przepis usunięto.** Link znika z przepisu. Czy formularz ma coś wtedy pokazać?
- **Przepis linkuje sam do siebie.** Nie pozwalamy, więc na liście wyboru go nie ma.

## Poza zakresem (nie projektuj)

- „Przepis w przepisie”, czyli wstawianie całego przepisu jako składnika z przeliczaniem ilości i makro. Może kiedyś, na razie nie ma przykładu.
- Linki w treści kroków.
- Sekcja „Powiązane przepisy” / „Zobacz też”.
- Odwrotne linki („Użyte w: Smashburger”). Możesz zaznaczyć, gdzie by kiedyś pasowały, ale bez projektu.

## Do oddania

Artboardy: light, dark, 375 px, dla widoków 1–3, oraz handoff z checklistą odbioru jak przy poprzednich feature'ach.
