import type { IngredientType } from '@/shared/enums/IngredientType';

export class EnumBase {
    ID: number = 0;
    Code!: IngredientType;
    Name: string = '';

    protected constructor(obj: EnumBase) {
        Object.assign(this, obj);
    }
}
