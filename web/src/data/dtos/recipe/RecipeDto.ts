import { RecipeIngredientDto } from "./RecipeIngredientDto"
import { RecipeStepDto } from "./RecipeStepDto"

type BackendRecipePhoto = { url: string } | null

export class RecipeDto {
  id: number = 0
  name: string = ''
  servings: number = 1
  description: string | null = null

  photoUrl: string | null = null

  steps: RecipeStepDto[] = []
  ingredients: RecipeIngredientDto[] = []

  constructor(
    obj: Partial<RecipeDto> & { photo?: BackendRecipePhoto }
  ) {
    Object.assign(this, obj)

    if ('photo' in obj) {
      this.photoUrl = obj.photo?.url ?? null
    }

    this.steps = (obj.steps ?? []).map((s) => new RecipeStepDto(s))
    this.ingredients = (obj.ingredients ?? []).map((i) => new RecipeIngredientDto(i))
  }
}
