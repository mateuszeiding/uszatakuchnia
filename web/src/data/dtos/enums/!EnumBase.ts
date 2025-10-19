import type { AromaType } from '@/shared/enums/AromaType'
import type { IngredientType } from '@/shared/enums/IngredientType'
import type { TasteType } from '@/shared/enums/TasteType'

export class EnumBase {
  ID: number = 0
  Code!: AromaType | IngredientType | TasteType
  Name: string = ''

  protected constructor(obj: EnumBase) {
    Object.assign(this, obj)
  }
}
