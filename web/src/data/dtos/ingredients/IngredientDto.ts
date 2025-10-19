import { IngredientType } from '@/shared/enums/IngredientType'
import type { AromaDto } from './aromaDto'
import type { TasteDto } from './tasteDto'

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
