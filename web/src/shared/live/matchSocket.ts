import type { WsMessage } from '@/shared/api/types'
import type { SocketStatus } from '@/shared/state/screenState'
import { backoffDelay } from './backoff'

export interface MatchSocketHandlers {
  onMessage: (message: WsMessage) => void
  // retryAt is when the next attempt starts (ms since epoch), null otherwise
  onStatus: (status: SocketStatus, retryAt: number | null) => void
  // Messages were lost: the caller should refetch the REST resources
  onGap: () => void
}

export interface MatchSocket {
  reconnectNow: () => void
  close: () => void
}

export type SocketFactory = (url: string) => WebSocket

// openMatchSocket keeps one match stream alive: it reconnects with backoff
// after every drop and reports gaps in seq. It knows nothing about Vue.
export function openMatchSocket(
  url: string,
  handlers: MatchSocketHandlers,
  createSocket: SocketFactory = (u) => new WebSocket(u),
): MatchSocket {
  let socket: WebSocket | null = null
  let attempt = 0
  let lastSeq = 0
  let retryTimer: ReturnType<typeof setTimeout> | undefined
  let waiting = false
  let closed = false

  function connect() {
    clearTimeout(retryTimer)
    waiting = false
    handlers.onStatus('connecting', null)

    const ws = createSocket(url)
    socket = ws

    ws.onopen = () => {
      attempt = 0
      handlers.onStatus('open', null)
    }

    ws.onmessage = (e: MessageEvent<string>) => {
      let message: WsMessage
      try {
        message = JSON.parse(e.data)
      } catch {
        return // not our format; ignore rather than break the stream
      }
      if (lastSeq !== 0 && message.seq !== lastSeq + 1) {
        handlers.onGap()
      }
      lastSeq = message.seq
      handlers.onMessage(message)
    }

    ws.onclose = () => {
      if (closed || socket !== ws) return
      scheduleRetry()
    }
  }

  function scheduleRetry() {
    const delay = backoffDelay(attempt)
    attempt++
    waiting = true
    handlers.onStatus('closed', Date.now() + delay)
    retryTimer = setTimeout(connect, delay)
  }

  connect()

  return {
    reconnectNow() {
      if (waiting && !closed) connect()
    },
    close() {
      closed = true
      clearTimeout(retryTimer)
      socket?.close()
    },
  }
}
