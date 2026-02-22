export class RecipeStepDto {
  stepNo: number = 0
  section: string | null = null
  text: string = ''

  constructor(obj: Partial<RecipeStepDto>) {
    Object.assign(this, obj)
  }
}
