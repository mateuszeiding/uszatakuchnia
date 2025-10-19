import { AromaType } from '@/shared/enums/AromaType'
import { EnumBase } from './!EnumBase'

export class AromaDto extends EnumBase {
  Code: AromaType = AromaType.OTHER

  constructor(obj: AromaDto) {
    super(obj)
  }
}
