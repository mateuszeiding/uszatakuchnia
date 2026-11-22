import { computed, type Ref } from 'vue';

import type { ICookingProgress } from '@/composables/useCookingProgress';
import {
    defaultSize,
    effectiveMode,
    type IPortionRecipe,
    portionScale,
    resolveSize,
    sameSize,
    sizeLabel,
} from '@/shared/utils/portion';

/**
 * The scale every view reads (ingredient list, step quantities, cooking mode, shopping list, macros), from the
 * shared uk:progress:<id> state: servings in servings mode, size otherwise.
 */
export function usePortion(recipe: () => IPortionRecipe, progress: Ref<ICookingProgress>) {
    const mode = computed(() => effectiveMode(recipe()));
    const size = computed(() => resolveSize(recipe(), progress.value.size));
    const scale = computed(() => portionScale(recipe(), progress.value.servings, size.value));
    /** The reader's size differs from the recipe's ("Przeliczone: …"). */
    const changed = computed(
        () => mode.value !== 'servings' && !sameSize(size.value, defaultSize(recipe()))
    );
    const label = computed(() => sizeLabel(recipe(), size.value));
    return { mode, size, scale, changed, label };
}
