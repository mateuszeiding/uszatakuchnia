# Prompt dla Claude Design — Faza 1: Przepisy

## Kontekst

Pracujemy nad aplikacją Uszata Kuchnia — personal recipe management app. Masz zaprojektować widoki dla Fazy 1 (Przepisy + Panel admina). Faza 2 (Challenges) i Faza 3 (Wiki) są poza scope — nie projektuj ich.

## Stack i Design System

- Frontend: Vue 3 + TypeScript
- Design System v2 (DS v2): tokeny w `web/src/scss/config/_tokens.scss`, komponenty w `web/src/scss/config/_components.scss`
- Kolory: ghost white #F8F8FF (tło), accent purple #6B4EE6, Geist + Geist Mono
- Klasy DS: `btn`, `badge`, `card`, `chip`, `input`, `dots` (trudność), `cat-*` (badge kategorii)
- Istniejące designy znajdziesz w `C:\Users\eidin\.claude\projects\C--Users-eidin-Repos-uszatakuchnia\design-extracted\uszatakuchnia\project\`

## Co już istnieje (nie przeprojektowuj od zera)

- `/recipes` — lista z sidebar filtrami, grid kart, results bar — **wymaga poprawek opisanych niżej**
- `/recipes/:id` — detal przepisu: hero, MetaStrip, składniki ze skalowaniem, kroki z checklistą i progress barem
- `/recipes/new` i `/recipes/:id/edit` — formularz tworzenia/edycji

## Widoki do zaprojektowania

### 1. Poprawki do `/recipes` — filtry

Obecny sidebar filtrów ma trzy problemy do rozwiązania:

**Problem 1 — przycisk "wyczyść filtry"** pojawia się warunkowo nad filtrami i przesuwa UI. Rozwiązanie: przycisk ma być na stałe w sidebarze (np. na dole lub obok nagłówka), nieaktywny (wyszarzony) gdy żaden filtr nie jest włączony, aktywny gdy jest.

**Problem 2 — brak wizualizacji aktywnych filtrów.** Użytkownik po wybraniu filtrów nie widzi co ma włączone bez scrollowania do sidebara. Potrzebny jest pasek/sekcja pod nagłówkiem strony z chipami aktywnych filtrów — każdy chip pokazuje nazwę filtru i ma X do usunięcia. **Zaprojektuj 3 warianty tego rozwiązania.**

**Problem 3 — OR / AND per grupa filtrów.** Wybranie dwóch tagów z tej samej grupy (np. "bez glutenu" + "keto") może działać jako OR (pokaż przepisy które mają którykolwiek z nich) albo AND (pokaż tylko te które mają oba). Użytkownik powinien móc to przełączyć per grupa. **Zaprojektuj 3 warianty togglea OR/AND.**

### 2. Rozszerzenie filtrów w sidebarze

Do istniejących filtrów (kategoria, czas, trudność) dodaj:

- Sekcja **Diety** — chipy tagów (bez glutenu, bez laktozy, keto itd.). Diety to informacja krytyczna (alergie, styl życia) — muszą być wizualnie wyróżnione względem grupy Praktyczne.
- Sekcja **Praktyczne** — chipy tagów (1 garnek, meal prep itd.). To informacja o wygodzie, nie krytyczna.
- Switch **"ukryj przepisy wymagające przygotowania"** — dla przepisów gdzie np. marynata robi się dzień wcześniej.

Obie grupy tagów muszą być wizualnie rozróżnione — inny kolor chipów, może inna forma.

### 3. Rozszerzenie karty przepisu

Karta ma pokazywać badge "wymaga przygotowania" jeśli przepis ma ten tag ustawiony.

### 4. Rozszerzenie detalu `/recipes/:id`

- Chipy tagów Diety i Praktyczne pod MetaStripem, z zachowaniem wizualnego rozróżnienia grup
- Badge "wymaga przygotowania" jeśli ustawiony
- Sekcja powiązanych przepisów na dole strony — 3 karty z tej samej kategorii

### 5. Rozszerzenie formularza `/recipes/new` i `/recipes/:id/edit`

- Wybór tagów z grupy Diety i Praktyczne (chipy, multi-select)
- Toggle "wymaga przygotowania"
- Przełącznik statusu: Szkic / Opublikowany
- Modal "Dodaj nowy składnik" — otwiera się gdy użytkownik przy dodawaniu składników do przepisu chce dodać składnik którego nie ma w bazie. Modal zawiera: nazwa (wymagana), typ/kategoria (wymagana), reszta opcjonalna. Prosty, minimalny.

### 6. Panel admina `/admin` — dashboard

Strona główna panelu admina. Zawiera kafle/sekcje na: Przepisy, Składniki, Challenges, Wiki. **Zaproponuj wygląd — 3 warianty.**

### 7. `/admin/recipes`

Lista wszystkich przepisów (szkice + opublikowane) z akcjami edytuj/usuń. Widoczny status każdego przepisu.

### 8. `/admin/ingredients`

Lista składników z akcjami edytuj/usuń. Formularz dodaj/edytuj: nazwa, typ/kategoria, czy alergen (toggle), zdjęcie.

### 9. `/admin/settings`

Zarządzanie słownikami aplikacji. Sekcje:
- Kategorie przepisów — każda kategoria ma nazwę i kolor badge
- Regiony/kuchnie — lista wartości
- Tagi grupy Diety — lista tagów
- Tagi grupy Praktyczne — lista tagów

**Zaproponuj UI do zarządzania tymi listami — 3 warianty.**

## Czego NIE projektować

- `/challenges` i wszystko z tym związane
- `/wiki` i wszystko z tym związane
- Linki z przepisów do Wiki (tooltopy, podkreślone słowa kluczowe) — to faza 3
- Przyciski "Oceń przepis" i "Dodaj komentarz" w DoneBlock — zostają martwe
