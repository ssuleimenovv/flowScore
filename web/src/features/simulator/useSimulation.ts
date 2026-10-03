import { onScopeDispose, shallowRef, watch, type Ref } from 'vue'
import type { SimulatedRed, Simulation } from '@/shared/api/types'
import { simulatorApi } from './api'

// useSimulation asks the model for the scenario whenever it changes and
// whenever the match moves on: "now" is a moment, and a red at the 60th
// minute means less at 55′ than at 40′. A newer request cancels the older,
// so a slow answer never replaces a fresh one.
export function useSimulation(
  matchId: string,
  reds: Ref<SimulatedRed[]>,
  moment: Ref<unknown>, // anything that changes when the match does
) {
  const simulation = shallowRef<Simulation | null>(null)
  const error = shallowRef<unknown>(null)
  let running: AbortController | null = null

  async function run() {
    running?.abort()
    const ctrl = new AbortController()
    running = ctrl
    try {
      simulation.value = await simulatorApi.simulate(matchId, { reds: reds.value }, ctrl.signal)
      error.value = null
    } catch (err) {
      if (!ctrl.signal.aborted) error.value = err
    } finally {
      if (running === ctrl) running = null
    }
  }

  watch([reds, moment], run, { immediate: true, deep: true })
  onScopeDispose(() => running?.abort())

  return { simulation, error, retry: run }
}
