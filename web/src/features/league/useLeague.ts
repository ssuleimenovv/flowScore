import { onScopeDispose, ref, shallowRef, watch, type Ref } from 'vue'
import { ApiError, apiGet } from '@/shared/api/http'
import type {
  League,
  LeagueItem,
  LeagueList,
  RoundList,
  RoundMatch,
  Standing,
  StandingList,
} from '@/shared/api/types'
import type { RequestStatus } from '@/shared/state/screenState'

// The gateway reads a league's season for about a minute after it starts and
// answers 503 meanwhile; the screen asks again this often
const PREPARING_EVERY = 30_000

// A round with a match under way is asked again this often, for the minute and the score
const LIVE_EVERY = 30_000

// useLeagues loads the leagues the gateway follows, once.
export function useLeagues() {
  const items = shallowRef<LeagueItem[]>([])
  const request = ref<RequestStatus>('idle')
  const error = shallowRef<unknown>(null)

  async function load() {
    request.value = 'pending'
    try {
      items.value = (await apiGet<LeagueList>('/leagues')).items
      error.value = null
      request.value = 'success'
    } catch (err) {
      error.value = err
      request.value = 'error'
    }
  }
  void load()

  return { items, request, error, retry: load }
}

// useLeague loads a league's header and table, again whenever the league
// changes. While the season is still being read it says so and keeps asking.
export function useLeague(id: Ref<string | null>) {
  const league = shallowRef<League | null>(null)
  const standings = shallowRef<Standing[]>([])
  const request = ref<RequestStatus>('idle')
  const preparing = ref(false)
  const error = shallowRef<unknown>(null)
  const updatedAt = ref<Date | null>(null)
  let loading: AbortController | null = null
  let timer: ReturnType<typeof setTimeout> | undefined

  async function load() {
    clearTimeout(timer)
    loading?.abort()
    const leagueId = id.value
    if (!leagueId) return
    const ctrl = new AbortController()
    loading = ctrl
    request.value = 'pending'
    try {
      const path = `/leagues/${encodeURIComponent(leagueId)}`
      const [l, s] = await Promise.all([
        apiGet<League>(path, ctrl.signal),
        apiGet<StandingList>(`${path}/standings`, ctrl.signal),
      ])
      league.value = l
      standings.value = s.items
      preparing.value = false
      error.value = null
      updatedAt.value = new Date()
      request.value = 'success'
    } catch (err) {
      if (ctrl.signal.aborted) return
      preparing.value = err instanceof ApiError && err.status === 503
      error.value = err
      request.value = 'error'
      if (preparing.value) timer = setTimeout(load, PREPARING_EVERY)
    }
  }

  watch(
    id,
    () => {
      league.value = null
      standings.value = []
      void load()
    },
    { immediate: true },
  )
  onScopeDispose(() => {
    clearTimeout(timer)
    loading?.abort()
  })

  return { league, standings, request, preparing, error, updatedAt, retry: load }
}

// useRound loads the matches of one round. While one of them is under way it
// asks again, as long as the tab is visible.
export function useRound(id: Ref<string | null>, round: Ref<number | null>) {
  const items = shallowRef<RoundMatch[]>([])
  const request = ref<RequestStatus>('idle')
  const error = shallowRef<unknown>(null)
  let loading: AbortController | null = null
  let timer: ReturnType<typeof setTimeout> | undefined

  async function load() {
    clearTimeout(timer)
    loading?.abort()
    if (!id.value || !round.value) return
    const ctrl = new AbortController()
    loading = ctrl
    if (request.value !== 'success') request.value = 'pending'
    try {
      const path = `/leagues/${encodeURIComponent(id.value)}/rounds/${round.value}`
      items.value = (await apiGet<RoundList>(path, ctrl.signal)).items
      error.value = null
      request.value = 'success'
    } catch (err) {
      if (ctrl.signal.aborted) return
      error.value = err
      if (items.value.length === 0) request.value = 'error'
    }
    if (items.value.some((m) => m.state === 'in') && document.visibilityState === 'visible') {
      timer = setTimeout(load, LIVE_EVERY)
    }
  }

  const onVisible = () => {
    if (document.visibilityState === 'visible') void load()
  }
  document.addEventListener('visibilitychange', onVisible)

  watch(
    [id, round],
    () => {
      items.value = []
      request.value = 'idle'
      void load()
    },
    { immediate: true },
  )
  onScopeDispose(() => {
    clearTimeout(timer)
    loading?.abort()
    document.removeEventListener('visibilitychange', onVisible)
  })

  return { items, request, error, retry: load }
}
