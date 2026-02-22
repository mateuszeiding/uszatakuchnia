import { useQueryClient, type FetchQueryOptions } from '@tanstack/vue-query'
import { API } from '../API'
import type { EnumBase } from '@/data/dtos/enums/!EnumBase'

type RequestMap = {
  ingredient: EnumBase
}

export const fetchEnums = async <E extends keyof RequestMap>(endpoint: E) => {
  const qc = useQueryClient()
  const options: FetchQueryOptions<RequestMap[E][]> = {
    queryKey: [endpoint],
    queryFn: () => API.Client.get<RequestMap[E][]>(`enums/${endpoint}`),
  }
  return await qc.fetchQuery(options)
}
