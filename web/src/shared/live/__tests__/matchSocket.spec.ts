// @vitest-environment node
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { openMatchSocket, type MatchSocketHandlers } from '../matchSocket'

// FakeSocket stands in for the browser WebSocket; tests drive it by hand.
class FakeSocket {
  onopen: (() => void) | null = null
  onmessage: ((e: { data: string }) => void) | null = null
  onclose: (() => void) | null = null
  closed = false

  close() {
    this.closed = true
  }

  open() {
    this.onopen?.()
  }

  receive(seq: number) {
    this.onmessage?.({ data: JSON.stringify({ type: 'flow.update', seq }) })
  }

  drop() {
    this.onclose?.()
  }
}

function setup() {
  const sockets: FakeSocket[] = []
  const handlers = {
    onMessage: vi.fn<MatchSocketHandlers['onMessage']>(),
    onStatus: vi.fn<MatchSocketHandlers['onStatus']>(),
    onGap: vi.fn<MatchSocketHandlers['onGap']>(),
  } satisfies MatchSocketHandlers

  const conn = openMatchSocket('ws://test', handlers, () => {
    const s = new FakeSocket()
    sockets.push(s)
    return s as unknown as WebSocket
  })
  return { sockets, handlers, conn }
}

describe('openMatchSocket', () => {
  beforeEach(() => vi.useFakeTimers())
  afterEach(() => vi.useRealTimers())

  it('reports a gap in seq', () => {
    const { sockets, handlers } = setup()
    sockets[0]!.open()

    sockets[0]!.receive(1)
    sockets[0]!.receive(2)
    sockets[0]!.receive(4)

    expect(handlers.onMessage).toHaveBeenCalledTimes(3)
    expect(handlers.onGap).toHaveBeenCalledTimes(1)
  })

  it('reconnects after a drop with a countdown', () => {
    const { sockets, handlers } = setup()
    sockets[0]!.open()

    sockets[0]!.drop()

    const [status, retryAt] = handlers.onStatus.mock.lastCall!
    expect(status).toBe('closed')
    if (retryAt === null) throw new Error('no retryAt')
    expect(retryAt).toBeGreaterThan(Date.now())

    vi.advanceTimersByTime(retryAt - Date.now())
    expect(sockets).toHaveLength(2)
  })

  it('reconnects at once on reconnectNow', () => {
    const { sockets, conn } = setup()
    sockets[0]!.drop()

    conn.reconnectNow()

    expect(sockets).toHaveLength(2)
  })

  it('stops for good after close', () => {
    const { sockets, conn } = setup()
    conn.close()
    sockets[0]!.drop()

    vi.advanceTimersByTime(60_000)

    expect(sockets[0]!.closed).toBe(true)
    expect(sockets).toHaveLength(1)
  })
})
