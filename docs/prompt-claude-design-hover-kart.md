# Prompt dla Claude Design — efekt najechania na karty (rozmyty tekst)

## Problem

Karty przepisów i karty składników po najechaniu kursorem podjeżdżają o 2 px do góry i dostają cień (`transform: translateY(-2px)` + `--shadow-card`). Właściciel zgłosił, że **tekst na karcie się wtedy rozmywa**.

Przyczyna jest techniczna:
- na czas animacji przeglądarka (Chrome) traktuje kartę jak obrazek i przesuwa go jako całość,
- przy skalowaniu ekranu w Windows 125% lub 150% przesunięcie o 2 px wypada na pół fizycznego piksela,
- przez to tekst traci ostrość, dopóki kursor jest na karcie.

Tego nie da się pewnie naprawić samym CSS przy zachowaniu przesunięcia, więc problem leży w projekcie efektu, a nie w implementacji.

Gdzie dziś występuje:
- **karta przepisu** na liście przepisów (`Przepisy.html`): podjazd 2 px + cień,
- **karta składnika** na liście składników: podjazd 2 px + cień,
- **przycisk akcentowy** (`btn--accent`): podjazd 1 px + cień akcentowy; tekst w przycisku też może się lekko rozmywać,
- **karta wyzwania** (`Wyzwania.html`): **bez podjazdu**, tylko ramka ciemnieje (`--rule` → `--rule-strong`). Tu problemu nie ma.

## Zadanie

Zaproponuj **2–3 warianty** efektu najechania dla kart (przepis, składnik, a dla spójności także wyzwanie), które:
- **nie przesuwają ani nie skalują** treści z tekstem: żadnego `translate`, `scale` ani zmian rozmiaru,
- nadal jasno mówią „to da się kliknąć”, także w trybie ciemnym (cień w dark mode słabo widać),
- są spójne między trzema typami kart, a jeśli się różnią, to celowo i z uzasadnieniem,
- są subtelne, w tonie DS v2.

Przykładowe kierunki, nie narzucam: zmiana ramki, cień bez ruchu, zmiana tła karty, akcent na tytule (kolor / podkreślenie), delikatny ruch samego zdjęcia wewnątrz kadru zamiast całej karty.

Osobno, krótko: **przycisk akcentowy**. Zostawić podjazd 1 px czy zamienić go na coś bez ruchu?

## Stany do pokazania

- Spoczynek, najechanie i fokus z klawiatury (fokus musi być widoczny niezależnie od hovera).
- Light i dark.
- Telefon: hover tam nie istnieje. Upewnij się, że nic nie zostaje „przyklejone” po tapnięciu.
- `prefers-reduced-motion`: jeśli wariant ma przejście, co się dzieje przy tym ustawieniu.

## Do oddania

- Warianty obok siebie (artboard) z rekomendacją.
- Po wyborze właściciela:
  - handoff z checklistą odbioru (które klasy i komponenty się zmieniają, tokeny),
  - wpis w „Decyzje i backlog.md”: zasada „hover bez przesuwania treści z tekstem”, żeby nie wróciła w nowych ekranach.
