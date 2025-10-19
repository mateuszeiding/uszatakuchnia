import { IngredientType } from '@/shared/enums/IngredientType'
import { EnumBase } from './!EnumBase'

export class IngredientDto extends EnumBase {
  Code: IngredientType = IngredientType.OTHER

  constructor(obj: IngredientDto) {
    super(obj)
  }
}
