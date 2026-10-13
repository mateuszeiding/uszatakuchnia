import type { IType } from '../ingredients/IType';

import type { IngredientType } from '@/shared/enums/IngredientType';

/** Kody osi klasyfikacji. Osie są w bazie — mogą dojść kolejne. */
export type AxisCode = 'dish' | 'base' | 'cuisine' | 'diet' | 'practical' | (string & {});

export interface ITaxonomyValue {
    id: number;
    label: string;
    /** Klucz tokenu koloru, np. "coral" → var(--c-coral). */
    color?: string;
    /** Liczba przepisów z tą wartością (tylko w słownikach). */
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
    /** Formy odmiany [1, 2–4, 5+]; brak = nieodmienna. */
    forms?: [string, string, string];
    noQty: boolean;
    /** W ilu składnikach przepisów użyta. */
    usedIn?: number;
}

export interface IAllergen {
    code: string;
    label: string;
}

export interface IDictionaries {
    axes: ITaxonomyAxis[];
    units: IUnit[];
    allergens: IAllergen[];
    ingredientTypes: IType<IngredientType>[];
}

/** Wartości klasyfikacji przepisu pogrupowane po kodzie osi. */
export type RecipeTaxonomy = Partial<Record<AxisCode, ITaxonomyValue[]>>;
