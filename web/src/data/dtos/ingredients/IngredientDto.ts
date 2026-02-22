import { IngredientType } from '@/shared/enums/IngredientType'
import { ImageDto } from './ImageDto'
import type { IType } from './IType'

export class IngredientDto {
  id: number = 0
  name: string = ''
  parentId: number | null = null

  type: IType<IngredientType> = {
    name: '',
    code: IngredientType.OTHER,
  }

  isAllergen: boolean = false

  image: ImageDto = new ImageDto({})

  constructor(obj: Partial<IngredientDto>) {
    Object.assign(this, obj)
    this.image = obj.image ? new ImageDto(obj.image) : new ImageDto({})
  }
}
