import { TasteType } from '@/shared/enums/TasteType'
import { EnumBase } from './!EnumBase'

export class TasteDto extends EnumBase {
  Code: TasteType = TasteType.UNKNOWN

  constructor(obj: TasteDto) {
    super(obj)
  }
}
