import { TasteType } from '@/shared/enums/TasteType'
import type { IType } from './IType'

export class TypeDto implements IType<TasteType> {
  name: string = ''
  intensity: number = 0
  code: TasteType = TasteType.UNKNOWN

  constructor(obj: Partial<TypeDto>) {
    Object.assign(this, obj)
  }
}
