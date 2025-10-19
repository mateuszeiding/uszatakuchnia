import { useQueryClient, type FetchQueryOptions } from '@tanstack/vue-query'
import { API } from '../API'
import type { TasteDto } from '@/data/dtos/enums/TasteDto'
import type { AromaDto } from '@/data/dtos/enums/AromaDto'
import type { IngredientDto } from '@/data/dtos/enums/IngredientDto'

type RequestMap = {
  aroma: AromaDto
  ingredient: IngredientDto
  taste: TasteDto
}

export const fetchEnums = async <E extends keyof RequestMap>(endpoint: E) => {
  const qc = useQueryClient()
  const options: FetchQueryOptions<RequestMap[E][]> = {
    queryKey: [endpoint],
    queryFn: () => API.Client.get<RequestMap[E][]>(`enums/${endpoint}`),
  }
  return await qc.fetchQuery(options)
}
