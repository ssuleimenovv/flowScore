import { apiGet } from '@/shared/api/http'
import type { EventList, FlowSeries, Insight, Match } from '@/shared/api/types'

const path = (id: string) => `/matches/${encodeURIComponent(id)}`

export const matchApi = {
  get: (id: string, signal?: AbortSignal) => apiGet<Match>(path(id), signal),
  flow: (id: string, signal?: AbortSignal) => apiGet<FlowSeries>(`${path(id)}/flow`, signal),
  events: (id: string, signal?: AbortSignal) => apiGet<EventList>(`${path(id)}/events`, signal),
  insight: (id: string, signal?: AbortSignal) => apiGet<Insight>(`${path(id)}/insight`, signal),
}
