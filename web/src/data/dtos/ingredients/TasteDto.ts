import { TasteType } from '@/shared/enums/TasteType'

export class TasteDto {
  name: string = ''
  intensity: number = 0
  code: TasteType = TasteType.UNKNOWN

  constructor(obj: Partial<TasteDto>) {
    Object.assign(this, obj)
  }
}
