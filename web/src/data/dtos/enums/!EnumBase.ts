export class EnumBase<T = unknown> {
  ID: number = 0
  Code!: T
  Name: string = ''

  protected constructor(obj: EnumBase) {
    Object.assign(this, obj)
  }
}
