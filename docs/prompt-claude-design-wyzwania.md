# Prompt dla Claude Design — Wyzwania

## Kontekst

Uszata Kuchnia, nowa sekcja aplikacji: **Wyzwania**. To drugi filar produktu obok przepisów. Aplikacja jest publiczna do czytania, bez kont użytkowników. Admin jest jeden (właściciel) i pracuje tylko na desktopie. Strony publiczne muszą działać na telefonie. DS v2, tokeny i komponenty jak w istniejących prototypach.

Dziś w aplikacji nie ma żadnej trasy ani widoku wyzwań, więc projektujesz je od zera, w stylu istniejących ekranów (`Przepisy.html`, `Przepis.html`, `Panel.html`).

## Problem

Właściciel chce trenować gotowanie na zadaniach w stylu kulinarnych programów konkursowych, a przy okazji mieć zabawę: wylosować coś i ugotować. Dziś trzyma takie zadania w głowie albo nigdzie.

## Feature

Wyzwanie to samodzielne zadanie kulinarne, **niezwiązane z przepisami** (bez linku do przepisu). Zestaw formatów ma być ogólny. Nie projektuj pod konkretny program: zadania mogą pochodzić skądkolwiek albo być wymyślone.

Dane wyzwania:
- **tytuł**,
- **format**: wartość ze słownika edytowalnego w panelu. Na start np. „z podanych składników”, „na temat”, „odtwórz danie”. Lista może się zmieniać, więc UI nie może zakładać konkretnych wartości,
- **opis zadania**: tekst, kilka zdań,
- **limit czasu**: liczba minut, tylko jako informacja. W aplikacji **nie ma** zegara ani odliczania,
- **składniki**: lista do dyspozycji lub obowiązkowych, opcjonalna (przy części formatów pusta). Na start sam tekst pozycji, bez powiązania z bazą składników — rozstrzygnij w handoffie, czy powiązanie ma sens,
- **trudność w skali 1–5**: ma być atrakcyjna wizualnie („fajne ikony”). Przepisy mają trudność 1–3 w kropkach, więc wyzwania mogą wyglądać inaczej. Zaproponuj 2–3 warianty,
- **źródło**: opcjonalny krótki tekst (skąd pochodzi zadanie).

Zachowania:
- **Lista wyzwań**: główny widok. Filtrów na razie **nie ma**, dojdą, kiedy wyzwań będzie dużo. Zaznacz, gdzie kiedyś by stanęły, żeby potem nie przebudowywać układu.
- **„Losuj wyzwanie”**: pokazuje jedno wylosowane wyzwanie. Można zawęzić losowanie trudnością. Losowanie pomija wyzwania oznaczone jako zrobione. Rozstrzygnij, co jeśli zostały same zrobione.
- **„Zrobione”**: oznaczenie zapamiętane na urządzeniu (localStorage, jak postęp gotowania), bez kont. Zrobione zostaje na liście na swoim miejscu, tylko oznaczone. Odznaczenie ma być łatwe.
- **Wpisywanie wyzwań**: głównie przez Claude, przez API (tak jak dziś przepisy). W panelu admina lista wyzwań z edycją i usuwaniem (usunięcie z potwierdzeniem, jak wszędzie), formularz wyzwania i słownik formatów w Panel → Klasyfikacja (wzorem „Jednostki” / „Sprzęt”).

## Widoki do zaprojektowania

1. **Wejście**: pozycja „Wyzwania” w górnym pasku obok „Przepisy”. Trasa: zaproponuj (np. `/wyzwania`).
2. **Lista wyzwań**: telefon i desktop. Co widać na pozycji listy (tytuł, format, trudność, czas, „zrobione”), gdzie stoi „Losuj”.
3. **Wyzwanie**: pełny widok jednego wyzwania (osobna strona czy rozwinięcie na liście — rozstrzygnij). Opis, składniki, czas, trudność, źródło, przycisk „Zrobione”.
4. **Losowanie**: wybór trudności, wynik, „losuj jeszcze raz”.
5. **Panel admina (desktop)**: lista wyzwań, formularz dodaj / edytuj, zakładka formatów w Klasyfikacji (dodawanie, zmiana nazwy, usuwanie używanego formatu — przeniesienie albo blokada, jak przy innych słownikach).

## Stany brzegowe

- Brak wyzwań (pusta lista), jedno wyzwanie.
- Wszystkie zrobione, a użytkownik losuje.
- Wyzwanie bez składników, bez źródła, z bardzo długim opisem, z bardzo długą listą składników.
- Losowanie z trudnością, w której nie ma żadnego wyzwania.

## Poza zakresem (nie projektuj)

- Zegar / odliczanie czasu, tryb „w trakcie wyzwania”.
- Link do przepisu.
- Zdjęcia (zdjęcia w aplikacji są odłożone).
- Punkty, ranking, konta, udostępnianie wyników.
- Wyzwania zespołowe.
- Filtry listy (tylko zaznacz miejsce).

## Do oddania

Artboardy: light, dark, 375 px dla widoków 2–4, desktop dla 2–5. Do tego handoff z modelem danych (pola i typy; jeśli masz wgląd w repo, nazwy wzorem `db/entities`), checklistą odbioru, sekcją „Nie zmieniaj” i otwartymi pytaniami. Wpis w „Decyzje i backlog.md”.
