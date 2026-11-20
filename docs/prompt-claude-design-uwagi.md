# Prompt dla Claude Design — Uwagi (przepis i wyzwanie)

## Kontekst

Uszata Kuchnia ma dwa filary: **przepisy** i **wyzwania**. Oba masz w projekcie:
- strona przepisu: `Przepis.html` (wraz z `Wskazowki.html`, `Sprzet.html`, `Zamienniki.html`, `Lista zakupow.html`),
- wyzwania: `Wyzwania.html`.

Aplikacja jest publiczna do czytania, bez kont użytkowników. Admin jest jeden i pracuje na desktopie, a strony publiczne muszą działać na telefonie. DS v2, tokeny i komponenty jak w istniejących prototypach.

## Problem

Przy części przepisów i wyzwań jest rzecz, której **nie wolno zawalić**, a którą łatwo przeoczyć w opisie albo w krokach. Przykłady:
- tatar wołowy: mięso tylko świeże od rzeźnika, nie paczkowane z marketu,
- surowe żółtko: jajka tylko pewne, świeże,
- ryba na surowo: tylko „sushi grade”, mrożona wcześniej,
- „zarezerwuj 2 dni — ciasto leżakuje w lodówce”.

Dziś takie rzeczy toną w tekście. Właściciel chce, żeby były widoczne **zanim** ktoś zacznie zakupy albo gotowanie, i to wprost, bez owijania w bawełnę.

## Feature: „Uwagi”

Przepis i wyzwanie mogą mieć **0–n uwag**. Każda uwaga to krótki tekst: jedno zdanie, góra dwa. Uwagi są **listą punktów**, nie jednym blokiem tekstu: właściciel chce, żeby było widać, że rzeczy jest kilka. Kolejność ustala autor.

### Czym uwaga różni się od wskazówki — ważne
W przepisach są już **wskazówki przy krokach** (decyzja „Wskazówki”, wariant A): „jeśli…”, „możesz też…”, „najłatwiej…”. Wtedy **odrzuciliśmy** stałą sekcję „Zanim zaczniesz” na poziomie przepisu, bo dwa miejsca na rady to szum.

Uwagi to **inna kategoria** i nie mogą tamtej decyzji obejść tylnymi drzwiami:
- **wskazówka** = jak zrobić lepiej, dotyczy kroku, opcjonalna do przeczytania,
- **uwaga** = czego nie zawalić, dotyczy całości (zakupy, bezpieczeństwo, czas, planowanie), trzeba ją przeczytać przed startem.

Wygląd musi to rozróżnienie oddawać: uwaga jest mocniejsza wizualnie niż wskazówka, ale to nie jest alarm na pół ekranu. W formularzu przyda się krótka podpowiedź dla autora, kiedy coś jest uwagą, a kiedy wskazówką do kroku.

## Gdzie i jak — do rozstrzygnięcia przez ciebie

1. **Strona wyzwania.** Gdzie stoją uwagi względem metadanych, opisu i karty „Składniki”? Telefon i desktop.
2. **Strona przepisu.** Gdzie stoją uwagi przy dzisiejszym układzie? Ograniczenia:
   - desktop ma przyklejony panel składników o ograniczonej wysokości, więc tam nie wkładamy,
   - prawa kolumna ma kolejność Sprzęt → Kroki → Wartości odżywcze,
   - telefon ma sekcje jedna pod drugą.
   Uwaga ma być widoczna, zanim ktoś zacznie składać zakupy.
3. **Lista zakupów przepisu.** Typowa uwaga dotyczy zakupów („kup porządne mięso”). Rozstrzygnij, czy uwagi pokazujemy też na liście zakupów (ekran i wydruk). Jeśli tak, to wszystkie, czy tylko część? Jeśli uważasz, że potrzebny jest podział na „uwagi zakupowe” i resztę, zaproponuj go, ale prostsze rozwiązanie wygrywa.
4. **Tryb gotowania przepisu.** Czy i gdzie uwagi się pojawiają (np. przed pierwszym krokiem)?
5. **Lista przepisów i lista wyzwań (karty).** Czy karta sygnalizuje, że są uwagi? Domyślnie nie. Zaproponuj sygnał tylko wtedy, gdy uważasz, że jest potrzebny.
6. **Panel admina (desktop):**
   - formularz wyzwania: dziś ma 01 Podstawowe, 02 Zadanie, 03 Składniki,
   - formularz przepisu: dziś ma m.in. Składniki, 03 Sprzęt, 04 Kroki,
   - w obu: gdzie wchodzą uwagi, pola uwag (✕, przeciąganie), „+ uwaga” i podpowiedź z sekcji „Czym uwaga różni się od wskazówki”.

## Stany brzegowe

- Brak uwag: nic się nie pokazuje, układ jak dziś.
- Jedna uwaga.
- 4–5 uwag.
- Bardzo długa uwaga.
- Wyzwanie z uwagami bez składników.
- Przepis z uwagami, sprzętem i objaśnieniem ikon jednocześnie: sprawdź, czy góra strony nie robi się ścianą kart.

## Poza zakresem

- Poziomy ważności, kolory per uwaga, ikony do wyboru.
- Powiązanie uwagi z konkretnym składnikiem lub krokiem: do kroku służą wskazówki.
- „Przeczytałem / ukryj” zapamiętane na urządzeniu, chyba że uznasz to za konieczne. Wtedy uzasadnij.
- Zmiany wskazówek przy krokach.

## Do oddania

- Artboardy: light, dark i 375 px dla strony wyzwania, strony przepisu, listy zakupów (jeśli dotyczy) i trybu gotowania (jeśli dotyczy).
- Desktop dla formularzy w panelu.
- Handoff:
  - model danych: lista tekstów z kolejnością, osobno dla przepisu i wyzwania,
  - checklista odbioru,
  - sekcja „Nie zmieniaj”,
  - otwarte pytania.
- Wpis w „Decyzje i backlog.md” z odniesieniem do decyzji „Wskazówki” (dlaczego to nie jest odrzucone „Zanim zaczniesz”).
