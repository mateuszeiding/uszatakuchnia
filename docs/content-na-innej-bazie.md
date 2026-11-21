# Kopia w `content/` na innej bazie — pojęcia ze słowniczka

Przepisy w `content/recipes/*.json` dają się odtworzyć na każdej bazie: składniki, jednostki i klasyfikacja są w nich po nazwie. Wyjątek to pojęcia ze słowniczka. W treści kroku i wskazówki odwołanie ma postać `{{term:<id>|słowo}}`, a `id` pochodzi z bazy, na której przepis dodano.

## Kiedy to ma znaczenie

- **`go run ./cmd/restore-recipe` z kopii tej samej bazy — bezpieczne.** Postgres nie używa id drugi raz, więc id wskazuje to samo pojęcie albo żadne. Brakujące pojęcie strona pokazuje jako samo słowo, a publikacja jest zablokowana do poprawki.
- **Wprowadzanie pliku z `content/` na inną bazę** (nowa baza, świeży Neon, lokalna baza testowa) — **ryzykowne bez mapowania.** Pod tym samym id może tam być inne pojęcie i czytelnik zobaczy cudze wyjaśnienie, a walidacja tego nie złapie.

## Format kopii pojęć

Pojęcia założone przy okazji przepisu Claude zapisuje w `content/glossary/<data>-<temat>.json`, z id z bazy źródłowej:

```json
[
  { "id": 12, "name": "deglasowanie", "explanation": "Wlanie płynu (wina, bulionu, wody) na gorącą patelnię…" }
]
```

## Procedura dla innej bazy

1. **Załóż pojęcia z pliku.** Dla każdego wpisu: `POST /api/glossary` z `{ "name", "explanation" }`. Jeśli pojęcie o tej nazwie już jest (409), weź jego id z `GET /api/glossary`, wyszukując po nazwie bez względu na wielkość liter.
2. **Przemapuj id w przepisie po nazwie.** Zbuduj mapę `stare id` → `nowe id` (para: `id` z pliku i id pojęcia o tej samej nazwie na nowej bazie). W każdym kroku i każdej wskazówce zamień `{{term:<stare>|…}}` na `{{term:<nowe>|…}}`; słowa po `|` nie zmieniaj.
3. **Dopiero wtedy wyślij przepis** (`POST /api/recipes`, jak w CLAUDE.md → „Zapis przez API z sesji Claude Code”).

Odwołanie do id, którego nie ma w mapie, zostaw. API zapisze szkic, ale publikację zablokuje, a strona pokaże samo słowo.
