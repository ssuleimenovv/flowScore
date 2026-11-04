import { onScopeDispose, ref, shallowRef, watch, type Ref } from 'vue'
import { apiGet } from '@/shared/api/http'
import type { Tie, TieList } from '@/shared/api/types'
import type { RequestStatus } from '@/shared/state/screenState'

// A stage with a leg under way is asked again this often, for the minute and the score
const LIVE_EVERY = 30_000

// useTies loads the ties of one stage of a cup, again whenever the cup or the
// stage changes. While a leg is under way it asks again, as long as the tab is visible.
export function useTies(id: Ref<string | null>, stage: Ref<string | null>) {
  const items = shallowRef<Tie[]>([])
  const request = ref<RequestStatus>('idle')
  const error = shallowRef<unknown>(null)
  let loading: AbortController | null = null
  let timer: ReturnType<typeof setTimeout> | undefined

  async function load() {
    clearTimeout(timer)
    loading?.abort()
    if (!id.value || !stage.value) return
    const ctrl = new AbortController()
    loading = ctrl
    if (request.value !== 'success') request.value = 'pending'
    try {
      const path = `/leagues/${encodeURIComponent(id.value)}/stages/${encodeURIComponent(stage.value)}/ties`
      items.value = (await apiGet<TieList>(path, ctrl.signal)).items
      error.value = null
      request.value = 'success'
    } catch (err) {
      if (ctrl.signal.aborted) return
      error.value = err
      if (items.value.length === 0) request.value = 'error'
    }
    const live = items.value.some((t) => t.legs.some((l) => l.state === 'in'))
    if (live && document.visibilityState === 'visible') timer = setTimeout(load, LIVE_EVERY)
  }

  const onVisible = () => {
    if (document.visibilityState === 'visible') void load()
  }
  document.addEventListener('visibilitychange', onVisible)

  watch(
    [id, stage],
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
