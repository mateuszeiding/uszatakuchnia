# Prompt dla Claude Design — kadr zastępczy zdjęcia na karcie przepisu (szkice koncepcji)

## Kontekst
Karta przepisu na liście (`Przepisy.html`, `Poprawki 18.11.html` → komponent Thumb) ma zawsze kadr 1:1 (decyzja 21.11). Przepisy na razie nie mają zdjęć, a własne zdjęcia są odłożone na dłużej, więc kadr zastępczy widać na **każdej** karcie.

Dziś kadr zastępczy to gradient w odcieniu liczonym z id przepisu (`h = id·47 mod 360`) plus „talerz” (dwa koła). Na kadrze w lewym dolnym rogu stoi „wymaga przygotowania”.

## Problem
Właściciel uważa, że obecny kadr wygląda słabo. Dwie rzeczy do zmiany:
1. **Kolor ma wynikać z typu dania**, a nie z id. Typ dania to główny badge na karcie i ma już kolor w DS:
   - śniadanie — `--c-coral`,
   - zupa — `--c-tomato`,
   - danie główne — `--c-terracotta`,
   - sałatka — `--c-lime`,
   - przekąska — `--c-plum`,
   - deser — `--c-pink`,
   - pieczywo i wypieki — `--c-brown`,
   - sos / dodatek — `--c-mustard`,
   - napój — `--c-teal`.

   Przepis może mieć kilka typów dania. Rozstrzygnij, który decyduje o kolorze: pierwszy, główny, a może mieszanka.
2. **Sam wygląd** — „talerz” z dwóch kół nie przekonuje.

## Zadanie: 5 szybkich koncepcji do wyboru
To **jeszcze nie jest pełny design**, tylko szkice do wyboru:
- **5 różnych koncepcji** kadru zastępczego, naprawdę różnych (np. wzór, ikona typu dania, typografia, ilustracja, abstrakcja — twój wybór),
- **jeden widok**: karta na liście, desktop, **jasny motyw**,
- **jedna kolorystyka na koncepcję** — pokaż każdą na 2–3 kartach z różnymi typami dania obok siebie, żeby było widać, jak kolor z kategorii się zmienia,
- **bez** wersji mobilnych, ciemnego motywu, stanów brzegowych, handoffu i checklisty — to dopiero po wyborze,
- przy każdej koncepcji jedno zdanie: idea i ewentualne ryzyko (np. „w ciemnym motywie może świecić”).

Po wyborze właściciela dopracujesz wybraną koncepcję w pełni: telefon, ciemny motyw, „wymaga przygotowania” na kadrze, karta ze zdjęciem obok karty bez, handoff.

## Ograniczenia
- Kadr 1:1, „wymaga przygotowania” w lewym dolnym rogu musi zostać czytelne.
- Hover karty bez zmian: ramka i tytuł w akcencie, nic się nie rusza ani nie skaluje.
- Kadr ma być tłem dla przyszłych zdjęć, więc karta ze zdjęciem obok karty z kadrem zastępczym nie może się gryźć.
- Kolory tylko z tokenów DS (`--c-*`) i ich jaśniejszych lub ciemniejszych wariantów.

## Przy okazji — do zapisania w „Decyzje i backlog.md” (decyzja właściciela 21.11)
**Ułamki w ilościach:** nigdzie nie pokazujemy znaków ½ ⅓ ¼ ⅔ ¾. Połówki i ćwiartki zapisujemy dziesiętnie z przecinkiem („0,25”, „1,5”, „1,75”), tylko tercje jako ułamek („1/3”, „2/3”, „1 1/3”). Obowiązuje wszędzie: lista składników, lista zakupów, znaczniki ilości w krokach. Etykiety części w formularzu („1/2”, „1/4”) zostają z ukośnikiem. W decyzjach i prototypach jest jeszcze wersja z „1/4” albo glifami — popraw ją.
