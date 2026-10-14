import { computed, onScopeDispose, ref, shallowRef, watch, type Ref } from 'vue'
import { ApiError } from '@/shared/api/http'
import type { MatchLineups } from '@/shared/api/types'
import type { RequestStatus } from '@/shared/state/screenState'
import { matchApi } from './api'
import type { MatchLive } from './matchState'

// Until the match is over the lineups change: the source puts them up about
// an hour before kick-off and adds every substitution
const REFRESH_MS = 60_000

// useLineups loads both teams' lineups and loads them again on every
// substitution the stream brings, and once a minute until the final whistle.
// lineups stays null while the source has none (the API answers 404).
export function useLineups(matchId: string, live: Ref<MatchLive | null>) {
  const lineups = shallowRef<MatchLineups | null>(null)
  const request = ref<RequestStatus>('idle')
  let loading: AbortController | null = null

  async function load() {
    loading?.abort()
    const ctrl = new AbortController()
    loading = ctrl
    if (request.value !== 'success') request.value = 'pending'
    try {
      lineups.value = await matchApi.lineups(matchId, ctrl.signal)
      request.value = 'success'
    } catch (err) {
      if (ctrl.signal.aborted) return
      if (err instanceof ApiError && err.status === 404) {
        lineups.value = null
        request.value = 'success'
      } else if (request.value !== 'success') {
        request.value = 'error' // a failed refresh keeps what is on screen
      }
    } finally {
      if (loading === ctrl) loading = null
    }
  }

  const status = computed(() => live.value?.match.status ?? null)
  const substitutions = computed(
    () => live.value?.events.filter((e) => e.type === 'substitution').length ?? 0,
  )
  // A new kick-off of a replay starts it again without them, and the final
  // whistle shows all of them
  watch([substitutions, status], () => void load())

  void load()
  const timer = setInterval(() => {
    if (status.value !== 'finished') void load()
  }, REFRESH_MS)

  onScopeDispose(() => {
    clearInterval(timer)
    loading?.abort()
  })

  return { lineups, request }
}
