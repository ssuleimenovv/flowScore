import { onScopeDispose, ref, shallowRef, watch, type Ref } from 'vue'
import { apiGet } from '@/shared/api/http'
import type { PlayerPage } from '@/shared/api/types'
import type { RequestStatus } from '@/shared/state/screenState'

// The seasons are read for about a minute after the gateway starts, and the
// page has the matches read so far: it is asked again now and then
const AGAIN_EVERY = 5 * 60_000

// usePlayer loads a player's page, again whenever the player changes, and
// keeps it fresh while the tab is visible
export function usePlayer(id: Ref<string>) {
  const page = shallowRef<PlayerPage | null>(null)
  const request = ref<RequestStatus>('idle')
  const error = shallowRef<unknown>(null)
  const updatedAt = ref<Date | null>(null)
  let loading: AbortController | null = null
  let timer: ReturnType<typeof setTimeout> | undefined

  async function load() {
    clearTimeout(timer)
    loading?.abort()
    const ctrl = new AbortController()
    loading = ctrl
    // Asked again with the page on screen: it stays, no skeleton
    if (!page.value) request.value = 'pending'
    try {
      page.value = await apiGet<PlayerPage>(`/players/${encodeURIComponent(id.value)}`, ctrl.signal)
      error.value = null
      updatedAt.value = new Date()
      request.value = 'success'
    } catch (err) {
      if (ctrl.signal.aborted) return
      error.value = err
      if (!page.value) request.value = 'error'
    }
    if (page.value && document.visibilityState === 'visible') {
      timer = setTimeout(load, AGAIN_EVERY)
    }
  }

  const onVisible = () => {
    if (document.visibilityState === 'visible') void load()
  }
  document.addEventListener('visibilitychange', onVisible)

  watch(
    id,
    () => {
      page.value = null
      void load()
    },
    { immediate: true },
  )
  onScopeDispose(() => {
    clearTimeout(timer)
    loading?.abort()
    document.removeEventListener('visibilitychange', onVisible)
  })

  return { page, request, error, updatedAt, retry: load }
}
