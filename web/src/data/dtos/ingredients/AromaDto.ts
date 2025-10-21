import { AromaType } from '@/shared/enums/AromaType'
import type { IType } from './IType'

export class AromaDto implements IType<AromaType> {
  name: string = ''
  intensity: number = 0
  code: AromaType = AromaType.OTHER

  constructor(obj: Partial<AromaDto>) {
    Object.assign(this, obj)
  }
}
