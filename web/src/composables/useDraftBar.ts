import { onBeforeUnmount, type Ref, shallowRef, watchEffect } from 'vue';

import { isOwner } from '@/auth/createAuth';
import type { RecipeDto } from '@/data/dtos/recipe/RecipeDto';

/** recipe = recipe page, list = shopping list (a different sentence in the bar). */
export type DraftBarKind = 'recipe' | 'list';
export type DraftBarState = { recipeId: number; kind: DraftBarKind };

/** The "Szkic" bar above the nav (App.vue); set by the recipe view or the shopping list view. */
export const draftBar = shallowRef<DraftBarState | null>(null);

/**
 * Shows the bar when the loaded recipe is a draft and the logged-in owner is viewing it.
 * A published recipe (also when the owner views it) and cooking mode: no bar.
 * `recipe` = only a successfully loaded current recipe (undefined while loading and on error).
 */
export function useDraftBar(kind: DraftBarKind, recipe: Readonly<Ref<RecipeDto | undefined>>) {
    let mine: DraftBarState | null = null;

    watchEffect(() => {
        const r = recipe.value;
        const next = r && r.status === 'draft' && isOwner.value ? { recipeId: r.id, kind } : null;
        // clear only our own bar: on recipe → list navigation the new view sets its own
        if (next || draftBar.value === mine) draftBar.value = next;
        mine = next;
    });

    onBeforeUnmount(() => {
        if (mine && draftBar.value === mine) draftBar.value = null;
    });
}
