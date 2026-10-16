import { onScopeDispose, ref, shallowRef } from 'vue'
import { ApiError } from '@/shared/api/http'
import type { MatchPreMatch } from '@/shared/api/types'
import type { RequestStatus } from '@/shared/state/screenState'
import { matchApi } from './api'

// usePreMatch loads "Перед матчем" once: the teams' form and their meetings
// do not change while the match is played. pre stays null for a match the
// source has nothing of the kind for, as a demo one (the API answers 404).
export function usePreMatch(matchId: string) {
  const pre = shallowRef<MatchPreMatch | null>(null)
  const request = ref<RequestStatus>('pending')
  const ctrl = new AbortController()

  matchApi
    .prematch(matchId, ctrl.signal)
    .then((p) => {
      pre.value = p
      request.value = 'success'
    })
    .catch((err: unknown) => {
      if (ctrl.signal.aborted) return
      request.value = err instanceof ApiError && err.status === 404 ? 'success' : 'error'
    })

  onScopeDispose(() => ctrl.abort())
  return { pre, request }
}
