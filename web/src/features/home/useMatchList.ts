import { onScopeDispose, ref, shallowRef } from 'vue'
import { apiGet } from '@/shared/api/http'
import type { MatchList, MatchSummary } from '@/shared/api/types'
import type { RequestStatus } from '@/shared/state/screenState'

// How often the home screen asks for the list. A card shows the minute and
// the score, and the clock ticks on between polls, so 5 seconds is plenty.
const EVERY = 5000

// useMatchList polls GET /matches while the page is open and visible. A failed
// poll keeps the last list on screen: the next one usually works.
export function useMatchList() {
  const items = shallowRef<MatchSummary[]>([])
  const request = ref<RequestStatus>('idle')
  const error = shallowRef<unknown>(null)
  let timer: ReturnType<typeof setTimeout> | undefined
  let loading: AbortController | null = null

  async function load() {
    clearTimeout(timer)
    loading?.abort()
    const ctrl = new AbortController()
    loading = ctrl
    if (request.value !== 'success') request.value = 'pending'
    try {
      items.value = (await apiGet<MatchList>('/matches', ctrl.signal)).items
      error.value = null
      request.value = 'success'
    } catch (err) {
      if (ctrl.signal.aborted) return
      error.value = err
      if (items.value.length === 0) request.value = 'error'
    }
    // A hidden tab does not poll; it loads again when it is shown
    if (document.visibilityState === 'visible') timer = setTimeout(load, EVERY)
  }

  const onVisible = () => {
    if (document.visibilityState === 'visible') void load()
  }
  document.addEventListener('visibilitychange', onVisible)
  void load()

  onScopeDispose(() => {
    clearTimeout(timer)
    loading?.abort()
    document.removeEventListener('visibilitychange', onVisible)
  })

  return { items, request, error, retry: load }
}
