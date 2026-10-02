import { computed, onScopeDispose, ref, shallowRef } from 'vue'
import { ApiError } from '@/shared/api/http'
import type { WsMessage } from '@/shared/api/types'
import { openMatchSocket } from '@/shared/live/matchSocket'
import type { RequestStatus, SocketStatus } from '@/shared/state/screenState'
import { matchApi } from './api'
import { applyMessage, fromSnapshot, type MatchLive } from './matchState'

// Messages kept while a snapshot loads; if more arrive, the next reload covers them.
const MAX_BUFFER = 1000

// useMatchLive loads the match snapshot over REST and keeps it current
// with the WebSocket stream. It reloads the snapshot after every
// (re)connect and every gap in seq.
export function useMatchLive(matchId: string) {
  const state = shallowRef<MatchLive | null>(null)
  const request = ref<RequestStatus>('idle')
  const error = shallowRef<unknown>(null)
  const socketStatus = ref<SocketStatus>('connecting')
  const retryAt = ref<number | null>(null)
  const now = ref(Date.now())
  // When the data last changed; the offline banner shows it ("данные на 21:14")
  const updatedAt = shallowRef<Date | null>(null)

  const secondsToRetry = computed(() =>
    retryAt.value === null ? null : Math.max(0, Math.ceil((retryAt.value - now.value) / 1000)),
  )

  let buffer: WsMessage[] = []
  let loading: AbortController | null = null
  let countdown: ReturnType<typeof setInterval> | undefined

  async function loadSnapshot() {
    loading?.abort()
    const ctrl = new AbortController()
    loading = ctrl
    request.value = 'pending'

    try {
      const [match, flow, events, insight] = await Promise.all([
        matchApi.get(matchId, ctrl.signal),
        matchApi.flow(matchId, ctrl.signal),
        matchApi.events(matchId, ctrl.signal),
        matchApi.insight(matchId, ctrl.signal).catch(withoutInsight),
      ])
      let next = fromSnapshot(match, flow, events, insight)
      for (const message of buffer) next = applyMessage(next, message)
      buffer = []

      state.value = next
      updatedAt.value = new Date()
      error.value = null
      request.value = 'success'
    } catch (err) {
      if (ctrl.signal.aborted) return // a newer load replaced this one
      error.value = err
      request.value = 'error'
    } finally {
      if (loading === ctrl) loading = null
    }
  }

  const socket = openMatchSocket(streamUrl(matchId), {
    onMessage(message) {
      if (loading || !state.value) {
        buffer.push(message)
        if (buffer.length > MAX_BUFFER) buffer.shift()
        return
      }
      state.value = applyMessage(state.value, message)
      updatedAt.value = new Date()
    },
    onStatus(next, at) {
      socketStatus.value = next
      retryAt.value = at
      clearInterval(countdown)
      if (at !== null) {
        now.value = Date.now()
        countdown = setInterval(() => (now.value = Date.now()), 250)
      }
      if (next === 'open') void loadSnapshot()
    },
    onGap() {
      void loadSnapshot()
    },
  })

  const onOnline = () => socket.reconnectNow()
  window.addEventListener('online', onOnline)

  onScopeDispose(() => {
    socket.close()
    loading?.abort()
    clearInterval(countdown)
    window.removeEventListener('online', onOnline)
  })

  return {
    state,
    request,
    error,
    updatedAt,
    socketStatus,
    retryAt,
    secondsToRetry,
    reconnectNow: socket.reconnectNow,
    retry: loadSnapshot,
  }
}

// A match without a prediction answers 404 (the gateway runs without an
// outcome model): the screen works without the card. Any other error is real.
function withoutInsight(err: unknown): null {
  if (err instanceof ApiError && err.status === 404) return null
  throw err
}

// The stream goes through the same host as the page; in development
// Vite proxies /ws to the Go gateway.
function streamUrl(matchId: string): string {
  const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:'
  return `${protocol}//${location.host}/ws/matches/${encodeURIComponent(matchId)}/stream`
}
