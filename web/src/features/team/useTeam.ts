import { onScopeDispose, ref, shallowRef, watch, type Ref } from 'vue'
import { ApiError, apiGet } from '@/shared/api/http'
import type { SquadPlayer, TeamPage } from '@/shared/api/types'
import type { RequestStatus } from '@/shared/state/screenState'

// The gateway reads the seasons for about a minute after it starts and
// answers 503 meanwhile; the screen asks again this often
const PREPARING_EVERY = 30_000

// useTeam loads a team's page, again whenever the team changes. While the
// seasons are still being read it says so and keeps asking.
export function useTeam(id: Ref<string>) {
  const page = shallowRef<TeamPage | null>(null)
  const request = ref<RequestStatus>('idle')
  const preparing = ref(false)
  const error = shallowRef<unknown>(null)
  const updatedAt = ref<Date | null>(null)
  let loading: AbortController | null = null
  let timer: ReturnType<typeof setTimeout> | undefined

  async function load() {
    clearTimeout(timer)
    loading?.abort()
    const ctrl = new AbortController()
    loading = ctrl
    request.value = 'pending'
    try {
      page.value = await apiGet<TeamPage>(`/teams/${encodeURIComponent(id.value)}`, ctrl.signal)
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
      page.value = null
      void load()
    },
    { immediate: true },
  )
  onScopeDispose(() => {
    clearTimeout(timer)
    loading?.abort()
  })

  return { page, request, preparing, error, updatedAt, retry: load }
}

// useSquad loads a team's squad once its tab is opened: the gateway asks ESPN
// for it, so a visitor who never looks does not cost a request
export function useSquad(id: Ref<string>, wanted: Ref<boolean>) {
  const items = shallowRef<SquadPlayer[]>([])
  const request = ref<RequestStatus>('idle')
  const error = shallowRef<unknown>(null)
  let loading: AbortController | null = null

  async function load() {
    loading?.abort()
    const ctrl = new AbortController()
    loading = ctrl
    request.value = 'pending'
    try {
      items.value = (
        await apiGet<{ items: SquadPlayer[] }>(
          `/teams/${encodeURIComponent(id.value)}/squad`,
          ctrl.signal,
        )
      ).items
      error.value = null
      request.value = 'success'
    } catch (err) {
      if (ctrl.signal.aborted) return
      error.value = err
      request.value = 'error'
    }
  }

  // Another team: forget the squad, and load it again if its tab is open
  watch(id, () => {
    items.value = []
    request.value = 'idle'
  })
  watch(
    [id, wanted],
    () => {
      if (wanted.value && request.value === 'idle') void load()
    },
    { immediate: true },
  )
  onScopeDispose(() => loading?.abort())

  return { items, request, error, retry: load }
}
