import { RecipeIngredientDto } from './RecipeIngredientDto';
import { RecipeStepDto } from './RecipeStepDto';

type BackendRecipePhoto = { url: string } | null;

export class RecipeBaseDto {
    id: number = 0;
    name: string = '';
    tagline: string | null = null;
    category: string | null = null;
    region: string | null = null;
    timeMinutes: number | null = null;
    difficulty: number | null = null;
    status: string = 'published';
    needsPrep: boolean = false;
    dietTags: string[] = [];
    practicalTags: string[] = [];
    photoUrl: string | null = null;

    constructor(obj: Partial<RecipeBaseDto> & { photo?: BackendRecipePhoto }) {
        Object.assign(this, obj);
        if ('photo' in obj) {
            this.photoUrl = obj.photo?.url ?? null;
        }
    }
}

export class RecipeDto extends RecipeBaseDto {
    servings: number = 1;
    description: string | null = null;
    kcalPerServing: number | null = null;

    steps: RecipeStepDto[] = [];
    ingredients: RecipeIngredientDto[] = [];

    constructor(obj: Partial<RecipeDto> & { photo?: BackendRecipePhoto }) {
        super(obj);
        Object.assign(this, obj);
        this.steps = (obj.steps ?? []).map((s) => new RecipeStepDto(s));
        this.ingredients = (obj.ingredients ?? []).map((i) => new RecipeIngredientDto(i));
    }
}

export interface IUpsertRecipeRequest {
    name: string;
    servings: number;
    description?: string | null;
    tagline?: string | null;
    category?: string | null;
    region?: string | null;
    timeMinutes?: number | null;
    difficulty?: number | null;
    kcalPerServing?: number | null;
    photoUrl?: string | null;
    status: string;
    needsPrep: boolean;
    dietTags: string[];
    practicalTags: string[];
    steps: { stepNo: number; section?: string | null; text: string }[];
    ingredients: {
        sortOrder: number;
        section?: string | null;
        ingredientId: number;
        amount?: number | null;
        unit?: string | null;
        amountText?: string | null;
        note?: string | null;
    }[];
}
