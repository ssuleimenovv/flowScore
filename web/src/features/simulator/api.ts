import { apiPost } from '@/shared/api/http'
import type { Simulation, SimulationRequest } from '@/shared/api/types'

export const simulatorApi = {
  simulate: (id: string, body: SimulationRequest, signal?: AbortSignal) =>
    apiPost<Simulation>(`/matches/${encodeURIComponent(id)}/simulate`, body, signal),
}
