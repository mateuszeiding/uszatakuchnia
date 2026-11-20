# Prompt dla Claude Design — tryb gotowania bez przewijania

## Kontekst

Uszata Kuchnia to aplikacja z przepisami właściciela. Strony publiczne muszą działać na telefonie. DS v2, tokeny i komponenty jak w istniejących prototypach.

**Tryb gotowania** to pełny ekran:
- jeden krok naraz,
- „Krok N z M” u góry,
- treść kroku w dużym rozmiarze,
- pod nią wskazówki (żarówka, też duże),
- na kroku 1 blok „Uwagi”,
- duże „Wstecz” / „Dalej” na dole („Dalej” odhacza krok),
- przycisk do arkusza „Sprzęt i składniki”, przeprojektowanego 20.11 na arkusz od dołu,
- ekran nie gaśnie.

20.11 poprawiono błąd: górna część długiego kroku była ucięta. Teraz wysoki krok przewija się od góry.

## Problem (feedback od znajomego, który gotował z aplikacji)

> W trybie gotowania musiałem scrollować w dół, żeby doczytać krok, a gdy ma się łapy brudne od mąki, to się syfi cały telefon. Może krótsze kroki albo trochę mniejszy font, żeby całość mieściła się na ekranie, a „Dalej” to można wcisnąć nawet górną częścią palucha albo nosem.

Tryb gotowania jest od tego, żeby przy kuchence **nie dotykać telefonu poza „Dalej”**. Dziś długi krok z dwiema, trzema wskazówkami nie mieści się na ekranie.

## Cel

1. **Krok mieści się na ekranie bez przewijania** na typowym telefonie: 375 × 667 z paskami przeglądarki, czyli realnie ~375 × 560 na treść. Dotyczy typowych kroków.
2. **„Dalej” da się trafić byle czym:** knykciem, brudnym palcem, bez celowania.

## Do rozstrzygnięcia przez ciebie

- **Dopasowanie treści:**
  - czy tekst kroku zmniejsza się sam, gdy krok jest długi (z dolną granicą czytelności z odległości ~50 cm),
  - czy wskazówki są zwinięte („2 wskazówki”, rozwija się jednym dotknięciem),
  - czy dzieje się jedno i drugie, czy coś innego.
- **Krok, który mimo to się nie zmieści:** co wtedy? Przewijanie jako ostateczność? Wyraźny sygnał „ciąg dalszy niżej”?
- **Strefa „Dalej”:**
  - większy cel, np. cała dolna część ekranu,
  - przesunięcie palcem,
  - dotknięcie w dowolnym miejscu kroku.

  Zważ ryzyko przypadkowego przejścia dalej, bo „Dalej” odhacza krok, oraz to, jak się cofnąć.
- **Uwagi na kroku 1:** blok zabiera miejsce. Jak zmieścić go z krokiem 1?
- **Zasada dla autora:** maksymalna długość kroku (np. liczba znaków) i licznik albo ostrzeżenie w formularzu przepisu, kiedy krok jest za długi do trybu gotowania. Zaproponuj limit.

## Powiązanie z innym featurem w toku

Równolegle projektujesz **ilości składników w treści kroku** (osobny brief). Odwołania wydłużają tekst kroku. Załóż je przy liczeniu miejsca, a jeśli kolejność ci to ułatwia, zaprojektuj oba razem.

## Stany do pokazania
- Krok krótki (1 zdanie).
- Krok typowy (2–3 zdania i 1 wskazówka).
- Krok długi (5 zdań i 3 wskazówki).
- Krok 1 z blokiem „Uwagi”.
- Ostatni krok („Gotowe”).
- Telefon 375 × 667 i 375 × 812, light i dark.
- Telefon w poziomie, czyli oparty o coś na blacie: krótko, jak się zachowuje.

## Poza zakresem
- Sterowanie głosem i gestami bez dotyku.
- Minutniki.
- Desktop: zostaje jak jest, chyba że zmiana wynika naturalnie.

## Do oddania
- **Artboardy:** stany powyżej, plus formularz z ostrzeżeniem o długim kroku (desktop).
- **Handoff:** checklista odbioru, sekcja „Nie zmieniaj”, otwarte pytania.
- **Wpis w „Decyzje i backlog.md”.**
