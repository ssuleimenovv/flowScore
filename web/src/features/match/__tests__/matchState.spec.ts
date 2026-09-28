// @vitest-environment node
import { describe, expect, it } from 'vitest'
import type { EventList, FlowSeries, Match, MatchEvent, WsMessage } from '@/shared/api/types'
import { applyMessage, fromSnapshot } from '../matchState'

const match = {
  clock: { elapsedSeconds: 120, period: 'first_half', observedAt: '2026-09-28T10:00:00Z' },
  id: 'm1',
  status: 'live',
  score: { home: 1, away: 0 },
  seq: 10,
} as Match

const flow = {
  matchId: 'm1',
  current: { home: 60, away: 20 },
  delta10: { home: 5, away: -3 },
  points: [
    { minute: 2, home: 30, away: 10 },
    { minute: 1, home: 20, away: 10 },
  ],
  updatedAt: '2026-09-28T10:00:00Z',
  seq: 10,
} satisfies FlowSeries

const events = { items: [], seq: 10 } satisfies EventList

const snapshot = () => fromSnapshot(match, flow, events)

function flowUpdate(seq: number, minute: number, home: number, elapsedSeconds = 0): WsMessage {
  return {
    type: 'flow.update',
    matchId: 'm1',
    seq,
    sentAt: '',
    data: {
      current: { home, away: 0 },
      delta10: { home: 0, away: 0 },
      point: { minute, home, away: 0 },
      clock: { elapsedSeconds, period: 'first_half', observedAt: '2026-09-28T10:00:05Z' },
    },
  }
}

function goal(seq: number, side: 'home' | 'away'): WsMessage {
  const data = { id: `g${seq}`, type: 'goal', side, minute: 30 } as MatchEvent
  return { type: 'match.event', matchId: 'm1', seq, sentAt: '', data }
}

describe('matchState', () => {
  it('sorts snapshot points by minute', () => {
    expect(snapshot().points.map((p) => p.minute)).toEqual([1, 2])
  })

  it('skips messages already in the snapshot', () => {
    const state = snapshot()
    expect(applyMessage(state, flowUpdate(10, 3, 99))).toBe(state)
  })

  it('applies a newer flow update and keeps points in order', () => {
    const next = applyMessage(snapshot(), flowUpdate(11, 3, 70))
    expect(next.flow.home).toBe(70)
    expect(next.points.map((p) => p.minute)).toEqual([1, 2, 3])
  })

  it('moves the clock to the one in the flow update', () => {
    const next = applyMessage(snapshot(), flowUpdate(11, 3, 70, 150))
    expect(next.match.clock).toEqual({
      elapsedSeconds: 150,
      period: 'first_half',
      observedAt: '2026-09-28T10:00:05Z',
    })
    expect(next.match.score).toEqual(match.score)
  })

  it('keeps the clock of a newer match part', () => {
    const state = { ...snapshot(), seq: { match: 12, flow: 10, events: 10 } }
    const next = applyMessage(state, flowUpdate(11, 3, 70, 150))
    expect(next.flow.home).toBe(70)
    expect(next.match.clock.elapsedSeconds).toBe(120)
  })

  it('replaces the point of the same minute', () => {
    const next = applyMessage(snapshot(), flowUpdate(11, 2, 45))
    expect(next.points).toHaveLength(2)
    expect(next.points[1]!.home).toBe(45)
  })

  it('adds a newer goal to the timeline and the score', () => {
    const next = applyMessage(snapshot(), goal(11, 'away'))
    expect(next.score).toEqual({ home: 1, away: 1 })
    expect(next.events[0]!.id).toBe('g11')
  })

  it('does not count a goal the snapshot already has', () => {
    const next = applyMessage(snapshot(), goal(9, 'home'))
    expect(next.score).toEqual({ home: 1, away: 0 })
  })
})
