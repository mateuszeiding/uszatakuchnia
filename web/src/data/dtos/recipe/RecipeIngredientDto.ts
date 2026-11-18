import type { IRecipeNutrition } from '@/shared/utils/nutrition';

/** Przepis, do którego prowadzi składnik (np. bułka → „Bułki do hamburgerów”). */
export interface ILinkedRecipe {
    id: number;
    title: string;
    status: 'draft' | 'published';
    /** Adres zdjęcia z API (przed resolveUrl); null = brak. */
    image: string | null;
}

/** Zamiennik pozycji (dane do alergenów, listy zakupów i makro bez kolejnego zapytania). */
export interface IRecipeSubstitute {
    ingredientId: number;
    name: string;
    /** „Pół na pół” ze składnikiem głównym (po połowie ilości). */
    mix: boolean;
    allergens: string[];
    ingredientType: { code: string; label: string };
    skipShopping: boolean;
    gramsPerUnit: { min: number; max: number } | null;
    /** Składnik ma dane makro (inaczej „bez danych o wartościach odżywczych”). */
    hasNutrition: boolean;
    /** Wkład wariantu w makro na porcję; null = pozycja się nie liczy. */
    nutrition: IRecipeNutrition | null;
}

export class RecipeIngredientDto {
    sortOrder: number = 0;
    section: string | null = null;

    ingredientId: number = 0;
    ingredientName: string = '';

    /** Ilość jako tekst: "2", "½", "1,5", "2–3". */
    quantity: string | null = null;
    unitId: number | null = null;
    /** Nazwa jednostki (forma słownikowa); odmianę robi unitLabel(). */
    unit: string | null = null;
    note: string | null = null;
    optional: boolean = false;
    /** Kody alergenów składnika. */
    allergens: string[] = [];
    /** Dział w sklepie (typ składnika) — lista zakupów. */
    ingredientType: { code: string; label: string } = { code: 'OTHER', label: 'inne' };
    /** Składnik pomijany na liście zakupów (np. woda). */
    skipShopping: boolean = false;
    /** Ile gramów to 1 jednostka tej pozycji; null = nie da się przeliczyć. */
    gramsPerUnit: { min: number; max: number } | null = null;
    /** Link do przepisu; szkic też przychodzi — strona publiczna go nie pokazuje (publicLink). */
    linkedRecipe: ILinkedRecipe | null = null;
    /** Wkład pozycji w makro na porcję (suma wkładów = nutrition przepisu); null = pozycja się nie liczy. */
    nutrition: IRecipeNutrition | null = null;
    substitutes: IRecipeSubstitute[] = [];

    constructor(obj: Partial<RecipeIngredientDto>) {
        Object.assign(this, obj);
        this.substitutes = obj.substitutes ?? [];
    }
}
