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

export function scaledAmount(ing: RecipeIngredientDto, units: IUnit[], ratio: number) {
    const unit = units.find((u) => u.id === ing.unitId);
    return scaleQuantity(ing.quantity, ing.unit, unit, ratio);
}

export const DIFFICULTY = ['', 'łatwe', 'średnie', 'trudne'];

export function fmtTime(min: number | null): string | null {
    if (!min) return null;
    if (min < 60) return `${min} min`;
    return min % 60 === 0 ? `${min / 60} h` : `${Math.floor(min / 60)} h ${min % 60} min`;
}

export const servingsLabel = (n: number) =>
    n === 1
        ? 'porcja'
        : n % 10 >= 2 && n % 10 <= 4 && (n % 100 < 10 || n % 100 >= 20)
          ? 'porcje'
          : 'porcji';

/** „12 h” → „12 h wcześniej”, „dzień wcześniej” bez zmian. */
export const leadText = (s: Pick<RecipeStepDto, 'prepLead'>) =>
    s.prepLead ? (s.prepLead.endsWith('wcześniej') ? s.prepLead : `${s.prepLead} wcześniej`) : null;
