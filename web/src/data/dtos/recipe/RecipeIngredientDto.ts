import type { IRecipeNutrition } from '@/shared/utils/nutrition';

/** Recipe this ingredient links to (e.g. bun → "Bułki do hamburgerów"). */
export interface ILinkedRecipe {
    id: number;
    title: string;
    status: 'draft' | 'published';
    /** Photo URL from the API (before resolveUrl); null = none. */
    image: string | null;
}

/** Row substitute (data for allergens, the shopping list and macros without another request). */
export interface IRecipeSubstitute {
    ingredientId: number;
    name: string;
    /** "Half and half" with the main ingredient (half the quantity each). */
    mix: boolean;
    allergens: string[];
    ingredientType: { code: string; label: string };
    skipShopping: boolean;
    gramsPerUnit: { min: number; max: number } | null;
    /** The ingredient has macro data (otherwise "no nutrition data"). */
    hasNutrition: boolean;
    /** The variant's share of the per-serving macros; null = the row doesn't count. */
    nutrition: IRecipeNutrition | null;
}

export class RecipeIngredientDto {
    sortOrder: number = 0;
    section: string | null = null;

    ingredientId: number = 0;
    ingredientName: string = '';

    /** Quantity as text: "2", "½", "1,5", "2–3". */
    quantity: string | null = null;
    unitId: number | null = null;
    /** Unit name (dictionary form); unitLabel() inflects it. */
    unit: string | null = null;
    note: string | null = null;
    optional: boolean = false;
    allergens: string[] = [];
    /** Store aisle (ingredient type), for the shopping list. */
    ingredientType: { code: string; label: string } = { code: 'OTHER', label: 'inne' };
    /** Skipped on the shopping list (e.g. water). */
    skipShopping: boolean = false;
    /** Grams per 1 unit of this row; null = can't convert. */
    gramsPerUnit: { min: number; max: number } | null = null;
    /** Recipe link; drafts come too, and the public page hides them (publicLink). */
    linkedRecipe: ILinkedRecipe | null = null;
    /** This row's share of the per-serving macros (shares add up to the recipe's nutrition); null = the row doesn't count. */
    nutrition: IRecipeNutrition | null = null;
    substitutes: IRecipeSubstitute[] = [];

    constructor(obj: Partial<RecipeIngredientDto>) {
        Object.assign(this, obj);
        this.substitutes = obj.substitutes ?? [];
    }
}
