import { RecipeIngredientDto } from './RecipeIngredientDto';
import { RecipeStepDto } from './RecipeStepDto';

import { API } from '@/data/api/API';
import type { ITaxonomyValue, RecipeTaxonomy } from '@/data/dtos/taxonomy/DictionariesDto';
import type { IRecipeNutrition, SumRange } from '@/shared/utils/nutrition';

type BackendRecipePhoto = { url: string } | null;

export type RecipeStatus = 'draft' | 'published';

export const STATUS_LABEL: Record<RecipeStatus, string> = {
    draft: 'Szkic',
    published: 'Opublikowany',
};

export class RecipeBaseDto {
    id: number = 0;
    name: string = '';
    tagline: string | null = null;
    /** Time in minutes; 0 = not given (draft). */
    timeMin: number = 0;
    timeMax: number | null = null;
    difficulty: number | null = null;
    status: RecipeStatus = 'draft';
    /** Computed: the recipe has a step marked as prep ahead. */
    needsPrep: boolean = false;
    taxonomy: RecipeTaxonomy = {};
    /** URL for display (locally with the API host). */
    photoUrl: string | null = null;
    /** URL stored in the database (e.g. /api/photos/1), sent back on save. */
    photoPath: string | null = null;
    kcal: SumRange | null = null;
    /** Ingredient names (list: search by ingredient). */
    ingredientNames: string[] = [];
    /** Substitute names (list search: "pancetta" finds the recipe with bacon). */
    substituteNames: string[] = [];

    constructor(obj: Partial<RecipeBaseDto> & { photo?: BackendRecipePhoto }) {
        Object.assign(this, obj);
        this.taxonomy = obj.taxonomy ?? {};
        if ('photo' in obj) {
            this.photoPath = obj.photo?.url ?? null;
            this.photoUrl = API.resolveUrl(this.photoPath);
        }
    }
}

/** First dish type: the main badge on the card. The API returns plain objects, hence a helper instead of a getter. */
export const mainDish = (r: Pick<RecipeBaseDto, 'taxonomy'>): ITaxonomyValue | null =>
    r.taxonomy.dish?.[0] ?? null;

export const cuisineOf = (r: Pick<RecipeBaseDto, 'taxonomy'>): ITaxonomyValue | null =>
    r.taxonomy.cuisine?.[0] ?? null;

export const taxonomyLabels = (r: Pick<RecipeBaseDto, 'taxonomy'>, axis: string): string[] =>
    (r.taxonomy[axis] ?? []).map((v) => v.label);

export interface IRecipeEquipment {
    equipmentId: number;
    name: string;
    spec: string | null;
    optional: boolean;
}

export class RecipeDto extends RecipeBaseDto {
    servings: number = 1;
    description: string | null = null;
    tip: string | null = null;
    cautions: string[] = [];
    nutrition: IRecipeNutrition | null = null;
    allergens: string[] = [];

    steps: RecipeStepDto[] = [];
    ingredients: RecipeIngredientDto[] = [];
    equipment: IRecipeEquipment[] = [];
    /** Ingredient rows of other recipes that link to this one (delete modal). */
    linkedFrom: { recipeId: number; recipeName: string; ingredientName: string }[] = [];

    constructor(obj: Partial<RecipeDto> & { photo?: BackendRecipePhoto }) {
        super(obj);
        Object.assign(this, obj);
        this.taxonomy = obj.taxonomy ?? {};
        this.allergens = obj.allergens ?? [];
        this.cautions = obj.cautions ?? [];
        this.steps = (obj.steps ?? []).map((s) => new RecipeStepDto(s));
        this.ingredients = (obj.ingredients ?? []).map((i) => new RecipeIngredientDto(i));
        this.linkedFrom = obj.linkedFrom ?? [];
        this.equipment = obj.equipment ?? [];
    }
}

export interface IUpsertRecipeRequest {
    name: string;
    servings: number;
    description?: string | null;
    tagline?: string | null;
    tip?: string | null;
    /** Notes without empty items; missing field = unchanged. */
    cautions?: string[];
    timeMin: number;
    timeMax?: number | null;
    difficulty?: number | null;
    photoUrl?: string | null;
    status: RecipeStatus;
    taxonomyValueIds: number[];
    steps: {
        stepNo: number;
        section?: string | null;
        text: string;
        prep: boolean;
        prepLead?: string | null;
        hints: string[];
    }[];
    ingredients: {
        sortOrder: number;
        /** Stable row key (kept on edit, new rows get a random one). */
        key: string;
        section?: string | null;
        ingredientId: number;
        quantity?: string | null;
        unitId?: number | null;
        note?: string | null;
        optional: boolean;
        linkedRecipeId?: number | null;
        substitutes: { ingredientId: number; mix: boolean }[];
    }[];
    equipment: { equipmentId: number; spec: string | null; optional: boolean }[];
}

/** Validation error from the API (422). */
export interface IFieldError {
    field: string;
    message: string;
}
