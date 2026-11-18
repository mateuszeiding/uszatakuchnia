# Sprzęt i zamienniki — plan wdrożenia (18.11.2026)

Źródło: artefakt Claude Design — `Handoff - Claude Code.md` (sekcje „Sprzęt potrzebny do przepisu”, „Zamienniki składnika w przepisie”), `Decyzje i backlog.md`, prototypy `sp-main.jsx`, `zm-main.jsx`.
Baza ma prawdziwe dane: nowe tabele zakłada AutoMigrate, reguły FK / indeksy — wersjonowane migracje (`0007`, `0008`). Nigdy `resetdb`.
Testy: lokalny embedded Postgres z seedem (scratchpad), zrzuty light / dark / 375 px.

## Etap 1 — Sprzęt (commit na etap)

| # | Punkt checklisty | Zadanie |
|---|---|---|
| S1 | Migracje `equipment`, `recipe_equipment`; DTO; CRUD słownika z zamianą przy usuwaniu | Encje `Equipment`, `RecipeEquipment`; migracja 0007 (unikalne `lower(name)`, FK CASCADE / RESTRICT); `RecipeDto.equipment`; zapis w `ReplaceRecipeChildren` + walidacja; słownik w `GET /api/taxonomy/list` z liczbą użyć; CRUD `/api/taxonomy/equipment` (`replaceWith` / `removeFromRecipes`) |
| S2 | Karta „Sprzęt” nad krokami / między składnikami a krokami; 1 i 6 pozycji, opcjonalna, długie nazwy — light i dark | `EquipmentCard.vue`, siatka `RecipeBody` (Sprzęt → Kroki → Wartości odżywcze) |
| S3 | Przepis bez sprzętu: brak karty | `v-if` na niepustej liście |
| S4 | Tryb gotowania: „Sprzęt i składniki”, sekcja sprzętu na początku | `CookMode.vue` (ta sama lista, jedna kolumna, tryb compact) |
| S5 | Formularz: sekcja 03, podpowiedzi, dodanie do słownika, kolejność, opcjonalny, usuwanie, zapis i odczyt; Kroki jako 04 | `EquipmentEditor.vue` + `EquipmentSuggest.vue`, model w `recipeForm.ts` |
| S6 | Panel: zakładka „Sprzęt”, dodawanie, zmiana nazwy, usuwanie z zamianą / usunięciem z przepisów | `EquipmentTable.vue`, `DeleteEquipmentModal.vue`, zakładka w `AdminTaxonomy.vue` |
| S7 | Nie zmieniaj: przyklejony panel składników, kolejność Kroki → Wartości odżywcze | Sprawdzenie na zrzutach (desktop, przewinięty) |

## Etap 2 — Zamienniki (commit na etap)

| # | Punkt checklisty | Zadanie |
|---|---|---|
| Z1 | Migracja tabeli, walidacja (≠ główny, duplikaty), blokada usunięcia składnika używanego jako zamiennik | Encja `RecipeIngredientSubstitute`; migracja 0008 (FK RESTRICT do `ingredients`, CASCADE do wiersza przepisu, unikalność); 422 w `ValidateRecipe`; licznik użyć przy usuwaniu składnika |
| Z2 | DTO z danymi zamienników + wkłady makro pozycji; bez wyboru makro = dzisiejsze | `RecipeIngredientDto.substitutes[]` + `nutrition` (wkład na porcję) dla głównego i każdego zamiennika; `Compute` złożone z `ComputeLine` (wynik bez zmian — test) |
| Z3 | Wiersz: ikona ⇄ obok ikony linku, osobny cel, 44 px na telefonie; stan wybrany | `IngRow.vue` |
| Z4 | Popover (desktop) i arkusz (375 px): radio, ilości, „zawiera: …”, „bez danych…”, Esc, fokus — light i dark | `SubstituteChooser.vue` (Teleport — lista składników przewija się w panelu) |
| Z5 | Karta „Ikony przy składnikach”: tylko obecne ikony, bez ikon — bez karty, „Ukryj” zapamiętane; desktop (panel przyklejony) i 375 px | `IconLegend.vue`, kolumna lewa w `RecipeBody` |
| Z6 | „Z zamiennikami: n · Przywróć przepis”; zapamiętanie w `uk:progress` | `useCookingProgress` (`subs`, nieaktualne wpisy ignorowane), `IngredientsPanel` |
| Z7 | Alergeny, makro (zdanie w karcie) i lista zakupów (mieszanka = 2 pozycje) wg wyboru | `substitutes.ts` (efektywne pozycje + suma wkładów makro), `NutritionCard`, `shoppingList.ts`, `ShoppingList` |
| Z8 | Tryb gotowania: wybór w składnikach | `CookMode.vue` — ta sama ikona i wybór, stan wspólny |
| Z9 | Formularz: „zamiast” / „pół na pół”, usuwanie, zapis i odczyt; chipy zamienników i linku w jednej siatce | `SubstitutePicker.vue`, `IngredientsEditor.vue`, `recipeForm.ts` |
| Z10 | Nie zmieniaj: odhaczanie, skalowanie, link do przepisu, „Przydatne przepisy” | Testy jednostkowe + zrzuty |
