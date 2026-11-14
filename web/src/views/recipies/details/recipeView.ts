import type { RecipeIngredientDto } from '@/data/dtos/recipe/RecipeIngredientDto';
import type { RecipeStepDto } from '@/data/dtos/recipe/RecipeStepDto';
import type { IUnit } from '@/data/dtos/taxonomy/DictionariesDto';
import { scaleQuantity } from '@/shared/utils/quantity';

export interface ISection<T> {
    label: string | null;
    items: T[];
}

/** Kolejne elementy z tą samą sekcją → jedna grupa (kolejność z API). */
export function bySection<T extends { section: string | null }>(items: T[]): ISection<T>[] {
    const out: ISection<T>[] = [];
    for (const it of items) {
        const last = out[out.length - 1];
        if (last && last.label === it.section) last.items.push(it);
        else out.push({ label: it.section, items: [it] });
    }
    return out;
}

/** Inne grupy, w których występuje ten sam składnik (małymi literami) — dopisek „też: …”. */
export function alsoIn(ing: RecipeIngredientDto, all: RecipeIngredientDto[]): string[] {
    if (!ing.section) return [];
    const out = all
        .filter(
            (x) => x.ingredientId === ing.ingredientId && x.section && x.section !== ing.section
        )
        .map((x) => x.section!.toLowerCase());
    return [...new Set(out)];
}

/** Grupy do wyświetlenia: >1 grupa → nagłówki; jedna nazwana → nazwa przy tytule; bez grup → sama lista. */
export function ingredientGroups(items: RecipeIngredientDto[]) {
    const sections = bySection(items);
    const single = sections.length === 1 ? sections[0]!.label : null;
    return { sections, grouped: sections.length > 1, single };
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

/** „12 h” → „12 h wcześniej”, „dzień wcześniej” bez zmian. */
export const leadText = (s: Pick<RecipeStepDto, 'prepLead'>) =>
    s.prepLead ? (s.prepLead.endsWith('wcześniej') ? s.prepLead : `${s.prepLead} wcześniej`) : null;
