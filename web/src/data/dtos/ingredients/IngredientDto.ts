import { ImageDto } from './ImageDto';
import type { IType } from './IType';

import { IngredientType } from '@/shared/enums/IngredientType';

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

    constructor(obj: Partial<IngredientDto>) {
        Object.assign(this, obj);
        this.image = obj.image ? new ImageDto(obj.image) : new ImageDto({});
    }
}
