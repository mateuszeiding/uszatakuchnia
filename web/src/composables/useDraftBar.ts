import { onBeforeUnmount, type Ref, shallowRef, watchEffect } from 'vue';

import { isOwner } from '@/auth/createAuth';
import type { RecipeDto } from '@/data/dtos/recipe/RecipeDto';

/** recipe = strona przepisu, list = lista zakupów (inne zdanie w pasku). */
export type DraftBarKind = 'recipe' | 'list';
export type DraftBarState = { recipeId: number; kind: DraftBarKind };

/** Pasek „Szkic” nad nawigacją (App.vue) — ustawia go widok przepisu albo listy zakupów. */
export const draftBar = shallowRef<DraftBarState | null>(null);

/**
 * Pokazuje pasek, gdy wczytany przepis jest szkicem, a ogląda go zalogowany właściciel.
 * Opublikowany przepis (także oglądany przez właściciela) i tryb gotowania — bez paska.
 * `recipe` = tylko pomyślnie wczytany bieżący przepis (w trakcie ładowania i przy błędzie: undefined).
 */
export function useDraftBar(kind: DraftBarKind, recipe: Readonly<Ref<RecipeDto | undefined>>) {
    let mine: DraftBarState | null = null;

    watchEffect(() => {
        const r = recipe.value;
        const next = r && r.status === 'draft' && isOwner.value ? { recipeId: r.id, kind } : null;
        // czyścimy tylko własny pasek — przy przejściu przepis → lista nowy widok ustawia swój
        if (next || draftBar.value === mine) draftBar.value = next;
        mine = next;
    });

    onBeforeUnmount(() => {
        if (mine && draftBar.value === mine) draftBar.value = null;
    });
}
