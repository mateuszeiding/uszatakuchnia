import { AromaType } from '@/shared/enums/AromaType'

export class AromaDto {
  name: string = ''
  intensity: number = 0
  code: AromaType = AromaType.OTHER

  constructor(obj: Partial<AromaDto>) {
    Object.assign(this, obj)
  }
}
