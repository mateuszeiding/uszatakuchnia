import { IngredientType } from '@/shared/enums/IngredientType'
import type { AromaDto } from './AromaDto'
import type { TasteDto } from './TasteDto'

export class IngredientDto {
  id: number = 0
  name: string = ''
  parentId: number | null = null
  aromas: AromaDto[] = []
  tastes: TasteDto[] = []
  type: IngredientType = IngredientType.OTHER
  isAllergen: boolean = false

  constructor(obj: Partial<IngredientDto>) {
    Object.assign(this, obj)
  }
}
