import { ImageDto } from './ImageDto';
import type { IType } from './IType';

import { IngredientType } from '@/shared/enums/IngredientType';
import type { Macros } from '@/shared/utils/nutrition';

/** Przelicznik jednostki na gramy; inherited = od składnika nadrzędnego. */
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

    /** Kody alergenów UE (patrz IDictionaries.allergens). */
    allergens: string[] = [];
    /** W ilu przepisach użyty (tylko z listy). */
    usedIn: number = 0;

    image: ImageDto = new ImageDto({});

    /** Własne makro na 100 g (brak klucza = dziedziczy z podskładników). */
    nutrition: Macros = {};
    /** Makro na 100 g do wyświetlenia: własne albo odziedziczone. */
    resolvedNutrition: Macros = {};
    /** Przeliczniki: własne + odziedziczone. */
    conversions: IConversion[] = [];

    constructor(obj: Partial<IngredientDto>) {
        Object.assign(this, obj);
        this.image = obj.image ? new ImageDto(obj.image) : new ImageDto({});
    }
}
