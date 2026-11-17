import { confirmAction } from './useConfirm';

import { API } from '@/data/api/API';
import type { RecipeDto, RecipeStatus } from '@/data/dtos/recipe/RecipeDto';

/**
 * Potwierdzenie usunięcia przepisu (formularz i Panel → Przepisy).
 * Jeśli prowadzą do niego składniki innych przepisów — akapit i lista „Przepis · składnik”
 * (linki znikną przez ON DELETE SET NULL, składniki zostaną).
 */
export async function confirmRecipeDelete(r: {
    id: number;
    name: string;
    status: RecipeStatus;
}): Promise<boolean> {
    // Świeże dane bez cache — lista linków ma być aktualna w chwili usuwania.
    const raw = await API.Client.get<Pick<RecipeDto, 'linkedFrom'>>(`recipes/${r.id}`);
    const from = raw.linkedFrom ?? [];
    const recipes = new Set(from.map((x) => x.recipeId)).size;

    const lead =
        r.status === 'published'
            ? 'Przepis jest opublikowany i zniknie ze strony.'
            : 'Przepis zniknie z panelu.';
    const links = from.length
        ? ` Prowadzą do niego składniki w ${recipes} ${recipes === 1 ? 'przepisie' : 'przepisach'} — te linki znikną, składniki zostaną:`
        : '';

    return confirmAction({
        title: `Usunąć „${r.name || 'ten przepis'}”?`,
        body: lead + links,
        items: from.map((x) => `${x.recipeName} · ${x.ingredientName.toLowerCase()}`),
        note: 'Tej operacji nie da się cofnąć.',
        confirmLabel: 'Usuń przepis',
        danger: true,
    });
}
