import type { IChallengeFormat } from '../challenge/ChallengeDto';
import type { IType } from '../ingredients/IType';

import type { IngredientType } from '@/shared/enums/IngredientType';

/** Classification axis codes. Axes live in the database, so more may come. */
export type AxisCode = 'dish' | 'base' | 'cuisine' | 'diet' | 'practical' | (string & {});

export interface ITaxonomyValue {
    id: number;
    label: string;
    /** Color token key, e.g. "coral" → var(--c-coral). */
    color?: string;
    /** Dish-type icon key, e.g. "zupa" (DISH_ICONS); missing = no icon. */
    icon?: string;
    /** Number of recipes with this value (dictionaries only). */
    recipeCount?: number;
}

export interface ITaxonomyAxis {
    id: number;
    code: AxisCode;
    label: string;
    multi: boolean;
    required: boolean;
    hasColor: boolean;
    filterMode: 'and' | 'or';
    values: ITaxonomyValue[];
}

export interface IUnit {
    id: number;
    name: string;
    /** Inflected forms [1, 2–4, 5+]; none = doesn't inflect. */
    forms?: [string, string, string];
    noQty: boolean;
    usedIn?: number;
}

export interface IAllergen {
    code: string;
    label: string;
}

export interface IEquipment {
    id: number;
    name: string;
    usedIn: number;
}

export interface IDictionaries {
    axes: ITaxonomyAxis[];
    units: IUnit[];
    allergens: IAllergen[];
    ingredientTypes: IType<IngredientType>[];
    equipment: IEquipment[];
    challengeFormats: IChallengeFormat[];
}

export type RecipeTaxonomy = Partial<Record<AxisCode, ITaxonomyValue[]>>;
