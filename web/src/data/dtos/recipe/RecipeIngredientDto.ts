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

    constructor(obj: Partial<RecipeIngredientDto>) {
        Object.assign(this, obj);
    }
}
