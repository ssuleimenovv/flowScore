import { computed, onScopeDispose, ref, shallowRef, triggerRef } from 'vue'
import type { FlowPoint, FlowValues, MatchEvent, WsMessage } from '@/shared/api/types'
import type { SocketStatus } from '@/shared/state/screenState'
import { openMatchSocket } from './matchSocket'

// How many timeline events to keep; older ones are dropped.
const MAX_EVENTS = 100

export interface StreamEvent extends MatchEvent {
  seq: number // unique per message, safe as a v-for key
}

// useMatchStream exposes one match's live stream as reactive state.
// The socket closes when the calling component unmounts.
export function useMatchStream(matchId: string) {
  const status = ref<SocketStatus>('connecting')
  const retryAt = ref<number | null>(null)
  const now = ref(Date.now())

  const flow = shallowRef<FlowValues | null>(null)
  const delta10 = shallowRef<FlowValues | null>(null)
  const points = shallowRef(new Map<number, FlowPoint>())
  const events = shallowRef<StreamEvent[]>([])
  const gaps = ref(0)

  const secondsToRetry = computed(() =>
    retryAt.value === null ? null : Math.max(0, Math.ceil((retryAt.value - now.value) / 1000)),
  )
  const wave = computed(() => [...points.value.values()].sort((a, b) => a.minute - b.minute))

  let countdown: ReturnType<typeof setInterval> | undefined

  const socket = openMatchSocket(streamUrl(matchId), {
    onMessage: handle,
    onStatus(next, at) {
      status.value = next
      retryAt.value = at
      clearInterval(countdown)
      if (at !== null) {
        now.value = Date.now()
        countdown = setInterval(() => (now.value = Date.now()), 250)
      }
    },
    onGap() {
      gaps.value++ // later: refetch REST resources here
    },
  })

  function handle(message: WsMessage) {
    switch (message.type) {
      case 'flow.update':
        flow.value = message.data.current
        delta10.value = message.data.delta10
        points.value.set(message.data.point.minute, message.data.point)
        triggerRef(points)
        break
      case 'match.event':
        events.value = [{ ...message.data, seq: message.seq }, ...events.value].slice(0, MAX_EVENTS)
        break
    }
  }

  const onOnline = () => socket.reconnectNow()
  window.addEventListener('online', onOnline)

  onScopeDispose(() => {
    socket.close()
    clearInterval(countdown)
    window.removeEventListener('online', onOnline)
  })

  return {
    status,
    secondsToRetry,
    flow,
    delta10,
    wave,
    events,
    gaps,
    reconnectNow: socket.reconnectNow,
  }
}

// The stream goes through the same host as the page; in development
// Vite proxies /ws to the Go gateway.
function streamUrl(matchId: string): string {
  const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:'
  return `${protocol}//${location.host}/ws/matches/${encodeURIComponent(matchId)}/stream`
}
