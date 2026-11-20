import { ImageDto } from './ImageDto';
import type { IType } from './IType';

import { IngredientType } from '@/shared/enums/IngredientType';
import type { Macros } from '@/shared/utils/nutrition';

/** Unit-to-grams conversion; inherited = from the parent ingredient. */
export interface IConversion {
    unitId: number;
    unitName: string;
    gramsMin: number;
    gramsMax: number | null;
    inherited: boolean;
}

export class IngredientDto {
    id: number = 0;
    name: string = '';
    parentId: number | null = null;

    type: IType<IngredientType> = {
        name: '',
        code: IngredientType.OTHER,
    };

    skipShopping: boolean = false;
    /** EU allergen codes (see IDictionaries.allergens). */
    allergens: string[] = [];
    /** Number of recipes using it (list only). */
    usedIn: number = 0;

    image: ImageDto = new ImageDto({});

    /** Own macros per 100 g (missing key = inherited from sub-ingredients). */
    nutrition: Macros = {};
    /** Macros per 100 g for display: own or inherited. */
    resolvedNutrition: Macros = {};
    conversions: IConversion[] = [];

    constructor(obj: Partial<IngredientDto>) {
        Object.assign(this, obj);
        this.image = obj.image ? new ImageDto(obj.image) : new ImageDto({});
    }
}
