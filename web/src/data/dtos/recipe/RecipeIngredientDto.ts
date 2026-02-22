export class RecipeIngredientDto {
    sortOrder: number = 0;
    section: string | null = null;

    ingredientId: number = 0;
    ingredientName: string = '';

    amount: number | null = null;
    unit: string | null = null;
    amountText: string | null = null;
    note: string | null = null;

    constructor(obj: Partial<RecipeIngredientDto>) {
        Object.assign(this, obj);
    }
}
