import { RecipeIngredientDto } from './RecipeIngredientDto';
import { RecipeStepDto } from './RecipeStepDto';

import { API } from '@/data/api/API';
import type { ITaxonomyValue, RecipeTaxonomy } from '@/data/dtos/taxonomy/DictionariesDto';

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
    timeMinutes: number | null = null;
    difficulty: number | null = null;
    status: RecipeStatus = 'draft';
    /** Wyliczane: przepis ma krok oznaczony jako przygotowanie wcześniej. */
    needsPrep: boolean = false;
    taxonomy: RecipeTaxonomy = {};
    /** Adres do wyświetlenia (lokalnie z hostem API). */
    photoUrl: string | null = null;
    /** Adres zapisany w bazie (np. /api/photos/1) — do odesłania przy zapisie. */
    photoPath: string | null = null;
    /** Nazwy składników (lista — wyszukiwanie po składnikach). */
    ingredientNames: string[] = [];

    constructor(obj: Partial<RecipeBaseDto> & { photo?: BackendRecipePhoto }) {
        Object.assign(this, obj);
        this.taxonomy = obj.taxonomy ?? {};
        if ('photo' in obj) {
            this.photoPath = obj.photo?.url ?? null;
            this.photoUrl = API.resolveUrl(this.photoPath);
        }
    }
}

/** Pierwszy typ dania — główny badge na karcie. API zwraca zwykłe obiekty, więc helper zamiast gettera. */
export const mainDish = (r: Pick<RecipeBaseDto, 'taxonomy'>): ITaxonomyValue | null =>
    r.taxonomy.dish?.[0] ?? null;

export const cuisineOf = (r: Pick<RecipeBaseDto, 'taxonomy'>): ITaxonomyValue | null =>
    r.taxonomy.cuisine?.[0] ?? null;

export const taxonomyLabels = (r: Pick<RecipeBaseDto, 'taxonomy'>, axis: string): string[] =>
    (r.taxonomy[axis] ?? []).map((v) => v.label);

export class RecipeDto extends RecipeBaseDto {
    servings: number = 1;
    description: string | null = null;
    tip: string | null = null;
    kcalPerServing: number | null = null;
    /** Kody alergenów wyliczone ze składników. */
    allergens: string[] = [];

    steps: RecipeStepDto[] = [];
    ingredients: RecipeIngredientDto[] = [];

    constructor(obj: Partial<RecipeDto> & { photo?: BackendRecipePhoto }) {
        super(obj);
        Object.assign(this, obj);
        this.taxonomy = obj.taxonomy ?? {};
        this.allergens = obj.allergens ?? [];
        this.steps = (obj.steps ?? []).map((s) => new RecipeStepDto(s));
        this.ingredients = (obj.ingredients ?? []).map((i) => new RecipeIngredientDto(i));
    }
}

export interface IUpsertRecipeRequest {
    name: string;
    servings: number;
    description?: string | null;
    tagline?: string | null;
    tip?: string | null;
    timeMinutes?: number | null;
    difficulty?: number | null;
    kcalPerServing?: number | null;
    photoUrl?: string | null;
    status: RecipeStatus;
    taxonomyValueIds: number[];
    steps: {
        stepNo: number;
        section?: string | null;
        text: string;
        prep: boolean;
        prepLead?: string | null;
    }[];
    ingredients: {
        sortOrder: number;
        section?: string | null;
        ingredientId: number;
        quantity?: string | null;
        unitId?: number | null;
        note?: string | null;
        optional: boolean;
    }[];
}

/** Błąd walidacji z API (422). */
export interface IFieldError {
    field: string;
    message: string;
}
