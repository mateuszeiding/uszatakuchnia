import { API } from '@/data/api/API';
import type { ILinkedRecipe, RecipeIngredientDto } from '@/data/dtos/recipe/RecipeIngredientDto';
import type { RecipeStepDto } from '@/data/dtos/recipe/RecipeStepDto';
import type { IUnit } from '@/data/dtos/taxonomy/DictionariesDto';
import { scaleQuantity } from '@/shared/utils/quantity';

export interface ISection<T> {
    label: string | null;
    items: T[];
}

/** Consecutive items with the same section → one group (API order). */
export function bySection<T extends { section: string | null }>(items: T[]): ISection<T>[] {
    const out: ISection<T>[] = [];
    for (const it of items) {
        const last = out[out.length - 1];
        if (last && last.label === it.section) last.items.push(it);
        else out.push({ label: it.section, items: [it] });
    }
    return out;
}

/** Groups to display: >1 group → headers; one named group → its name next to the title; no groups → a plain list. */
export function ingredientGroups(items: RecipeIngredientDto[]) {
    const sections = bySection(items);
    const single = sections.length === 1 ? sections[0]!.label : null;
    return { sections, grouped: sections.length > 1, single };
}

/** Recipe link visible on the public page; drafts aren't shown. */
export const publicLink = (ing: RecipeIngredientDto): ILinkedRecipe | null =>
    ing.linkedRecipe?.status === 'published' ? ing.linkedRecipe : null;

export interface IUsefulRecipe {
    id: number;
    title: string;
    image: string | null;
    /** Names of the ingredients that link to it ("do: bułka do hamburgera"). */
    for: string[];
}

/** "Przydatne przepisy": published links from the ingredients; the same recipe from several rows → one card. */
export function usefulRecipes(items: RecipeIngredientDto[]): IUsefulRecipe[] {
    const out: IUsefulRecipe[] = [];
    for (const ing of items) {
        const r = publicLink(ing);
        if (!r) continue;
        const name = ing.ingredientName.toLowerCase();
        const seen = out.find((x) => x.id === r.id);
        if (seen) {
            if (!seen.for.includes(name)) seen.for.push(name);
        } else out.push({ id: r.id, title: r.title, image: API.resolveUrl(r.image), for: [name] });
    }
    return out;
}

export function scaledAmount(ing: RecipeIngredientDto, units: IUnit[], ratio: number) {
    const unit = units.find((u) => u.id === ing.unitId);
    return scaleQuantity(ing.quantity, ing.unit, unit, ratio);
}

export const DIFFICULTY = ['', 'łatwe', 'średnie', 'trudne'];

export const servingsLabel = (n: number) =>
    n === 1
        ? 'porcja'
        : n % 10 >= 2 && n % 10 <= 4 && (n % 100 < 10 || n % 100 >= 20)
          ? 'porcje'
          : 'porcji';

/** "12 h" → "12 h wcześniej"; "dzień wcześniej" stays as is. */
export const leadText = (s: Pick<RecipeStepDto, 'prepLead'>) =>
    s.prepLead ? (s.prepLead.endsWith('wcześniej') ? s.prepLead : `${s.prepLead} wcześniej`) : null;
