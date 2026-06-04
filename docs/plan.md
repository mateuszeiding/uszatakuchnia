# Uszata Kuchnia — plan produktu

## Faza 1 — Przepisy

### Dane przepisu
- Nazwa + tagline
- Zdjęcie
- Kategoria (konfigurowalna przez admina)
- Region/kuchnia (konfigurowalny przez admina)
- Czas w minutach
- Trudność (1–3, stała skala)
- Porcje + kcal na porcję
- Tag "wymaga przygotowania" — boolean, osobny od grup tagów
- Tagi grupy **Diety** (bez glutenu, bez laktozy, keto itd.) — konfigurowalne przez admina
- Tagi grupy **Praktyczne** (1 garnek, meal prep itd.) — konfigurowalne przez admina
- Składniki z ilościami i sekcjami
- Kroki z sekcjami
- Wskazówka/tip
- Status: szkic / opublikowany

### Lista `/recipes`
Grid kart. Sidebar z filtrami: kategoria, czas max, trudność, sekcja "Diety" (chipy), sekcja "Praktyczne" (chipy) — obie grupy wizualnie rozróżnione bo Diety są informacją krytyczną (alergie), Praktyczne to wygoda. Switch "ukryj przepisy wymagające przygotowania". Sortowanie: najnowsze, najszybsze, A–Z. Search po nazwie.

Admin widzi tabs: Wszystkie / Opublikowane / Szkice.

**Uwagi do filtrów (problemy z obecnym designem do rozwiązania):**
- Przycisk "wyczyść filtry" nie może pojawiać się warunkowo nad filtrami — przesuwa UI. Ma być na stałe w sidebarze, nieaktywny gdy żaden filtr nie jest włączony.
- Wizualizacja aktywnych filtrów — użytkownik musi widzieć aktywne filtry bez scrollowania do sidebara. Propozycja: pasek pod nagłówkiem z chipami aktywnych filtrów, każdy chip ma X do usunięcia. **Claude Design przygotowuje 3 warianty.**
- OR / AND per grupa filtrów — wybranie dwóch tagów z tej samej grupy może działać jako OR (którykolwiek) lub AND (oba). Użytkownik ma toggle per grupę chipów. **Claude Design przygotowuje 3 warianty.**

### Karta przepisu
Zdjęcie, nazwa, tagline, badge kategorii, czas, trudność (kropki), badge "wymaga przygotowania" jeśli ustawiony.

### Detal `/recipes/:id`
Hero z zdjęciem, MetaStrip (czas, trudność, porcje, kcal), chipy tagów z zachowaniem podziału na grupy (Diety / Praktyczne), składniki ze skalowaniem porcji (sticky sidebar), kroki z checklistą i progress barem, sekcja powiązanych przepisów na dole (3 karty z tej samej kategorii). Słowa kluczowe w krokach i składnikach linkują do Wiki (faza 3).

### Formularz `/recipes/new` i `/recipes/:id/edit`
Wszystkie pola jak wyżej. Wybór tagów z obu grup. Tag "wymaga przygotowania" jako osobny toggle. Status szkic/opublikowany. Przy dodawaniu składnika — modal z polami: nazwa (wymagana) + typ/kategoria (wymagana), reszta opcjonalna do uzupełnienia później.

### Panel admina `/admin`
Dashboard z kaflem na każdą sekcję: Przepisy, Składniki, Challenges, Wiki. **Claude Design proponuje wygląd.**

### `/admin/recipes`
Lista wszystkich przepisów (szkice + opublikowane), akcje edytuj/usuń.

### `/admin/ingredients`
Lista składników, formularz dodaj/edytuj: nazwa, typ/kategoria, czy alergen, zdjęcie.

### `/admin/settings`
Zarządzanie słownikami: kategorie przepisów (nazwa + kolor badge), regiony/kuchnie, tagi grupy Diety, tagi grupy Praktyczne. **Claude Design proponuje UI.**

---

## Faza 2 — Challenges

### Dane challenge
- Tytuł + opis
- Trudność (1–3)
- Typ (konfigurowalny przez admina: konkretne danie, z podanych składników, tematyka, technika itd.)
- Opcjonalny link do przepisu

### Widok `/challenges`
Master-detail: lewa kolumna — lista challengy z filtrem po trudności i typie, badge "ukończone" na kartach odhaczonych przez użytkownika (localStorage). Prawa kolumna — detal wybranego challenge: tytuł, opis, trudność, typ, opcjonalny link do przepisu, przycisk "Oznacz jako ukończone".

Na mobile: lista pełnoekranowa, kliknięcie zastępuje ją detalem z przyciskiem powrotu.

### `/admin/challenges`
Lista challengy, dodaj/edytuj/usuń. Formularz: tytuł, opis, trudność, typ, opcjonalny link do przepisu.

### `/admin/settings`
Uzupełnić o sekcję typów challengy.

---

## Faza 3 — Wiki

Ostatnia część aplikacji. Szczegóły do doprecyzowania w swoim czasie.

### Ogólny zarys
- **Składniki** — opis, właściwości, przechowywanie, odmiany, zamienniki, przepisy gdzie występuje
- **Techniki gotowania** — opis, zastosowanie, powiązane przepisy
- **Style gotowania**

### Widoki
- `/wiki` — strona główna z podziałem na sekcje (Składniki / Techniki / Style) i searchem
- `/wiki/:slug` — artykuł z treścią i powiązanymi przepisami
- `/admin/wiki` — lista artykułów, formularz: typ, treść, słowa kluczowe

### Powiązanie z przepisami
Słowa kluczowe w treści kroków i składnikach przepisu — podkreślone, hover pokazuje tooltip (nazwa + jedno zdanie opisu), klik przenosi na artykuł Wiki. Admin definiuje słowa kluczowe przy tworzeniu artykułu.

---

## Nawigacja

TopBar: logo, Przepisy, Challenges, Wiki. Jeśli admin zalogowany — link do `/admin`.
