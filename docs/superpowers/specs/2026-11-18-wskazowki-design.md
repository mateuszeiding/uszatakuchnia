# Wskazówki przy krokach — projekt (krótki)

Zaakceptowane 18.11.2026. Źródło: Claude Design, `Wskazowki.html` → `wk-main.jsx` (artboardy „A”), handoff „Wskazówki przy krokach”, decyzja „Wskazówki” w `Decyzje i backlog.md`.

## Po co
Wskazówka na końcu przepisu przychodzi po ugotowaniu. Rady („jeśli…”, „możesz też…”) mają stać przy kroku, którego dotyczą (0–n na krok).

## Cele / nie-cele
- Cel: `recipe_steps.hints`, wyświetlanie pod krokiem (strona, tryb gotowania), edycja w formularzu, przeniesienie starej `recipes.tip` do kroku.
- Nie-cel: zmiany odhaczania kroków, „przygotuj wcześniej”, grup kroków. Usunięcie kolumny `tip` (później, osobną migracją, gdy wszystkie przeniesione).
- Anty-cel: utrata istniejących wskazówek (`tip`) — zostaje w bazie, dopóki admin jej nie przeniesie.

## Podejście
Jedno sensowne: lista tekstów w jsonb na kroku (kroki są zapisywane w całości przy upsercie, osobna tabela nic nie daje). Alternatywy (osobna tabela `recipe_step_hints`, tekst z separatorem) odrzucone — więcej kodu albo kruche.

## Projekt
1. **Dane:** typ `StringList []string` (Scan/Value jsonb), kolumna `hints jsonb NOT NULL DEFAULT '[]'` przez AutoMigrate; migracja `0006_recipe_step_hints` dopisuje `[]` tam, gdzie NULL (HasColumn). Zapis: trim, puste pomijane po cichu, > 300 znaków → 422 z polem `steps[i].hints[j]`.
2. **API:** `RecipeStepDto.hints: string[]` (zawsze tablica), `UpsertStepRequest.hints`.
3. **Strona:** `StepHints.vue` (rozmiar normalny / duży, stan `done`) pod treścią kroku w `StepsList.vue` i `CookMode.vue`. `tip-box` usunięty; zamiast niego, gdy `tip` niepusty, blok „Zanim zaczniesz” nad krokami.
4. **Formularz:** sekcja „04 Wskazówka” usunięta (decyzja właściciela 18.11). W `StepsEditor` pod treścią kroku wiersze wskazówek + „+ wskazówka”. Gdy `tip` niepusty — nad krokami treść + select „Przenieś do kroku…”: dopisuje do kroku i czyści `tip` w stanie formularza; w bazie po „Zapisz”.

## Decyzje Claude (do potwierdzenia)
- Puste wskazówki pomijane po cichu (nie błąd); za długie → błąd walidacji.
- Formularz na telefonie bez osobnego projektu — ten sam układ, responsywny.
- Testy na lokalnym Postgresie, nie na produkcji.
