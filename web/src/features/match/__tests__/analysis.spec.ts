// @vitest-environment node
import { describe, expect, it } from 'vitest'
import type { Analysis, Insight, Match, WsMessage } from '@/shared/api/types'
import { applyMessage, fromSnapshot } from '../matchState'

const match = {
  clock: { elapsedSeconds: 0, period: 'first_half', observedAt: '2026-10-07T10:00:00Z' },
  score: { home: 0, away: 0 },
  seq: 10,
} as Match

const flow = {
  matchId: 'm1',
  current: { home: 20, away: 20 },
  delta10: { home: 0, away: 0 },
  factors: [],
  points: [],
  updatedAt: '2026-10-07T10:00:00Z',
  seq: 10,
}
const events = { items: [], seq: 10 }
const prediction = {
  current: { home: 40, draw: 30, away: 30 },
  preMatch: { home: 40, draw: 30, away: 30 },
  model: 'Poisson v1',
}

const analysis = (minute: number): Analysis => ({
  title: `Разбор ${minute}′`,
  text: 'Текст',
  minute,
  generatedAt: '2026-10-07T10:01:00Z',
})

const update = (seq: number, data: Analysis): WsMessage => ({
  type: 'insight.update',
  matchId: 'm1',
  seq,
  sentAt: '2026-10-07T10:01:00Z',
  data,
})

describe('the AI analysis in the match state', () => {
  it('comes from the snapshot, or is null before the agent answers', () => {
    const withText: Insight = { matchId: 'm1', prediction, explanation: analysis(10), seq: 10 }
    expect(fromSnapshot(match, flow, events, withText).analysis?.minute).toBe(10)
    expect(
      fromSnapshot(match, flow, events, { matchId: 'm1', prediction, seq: 10 }).analysis,
    ).toBeNull()
    expect(fromSnapshot(match, flow, events, null).analysis).toBeNull()
  })

  it('is replaced by insight.update', () => {
    const state = fromSnapshot(match, flow, events, null)
    const next = applyMessage(state, update(11, analysis(20)))
    expect(next.analysis?.title).toBe('Разбор 20′')
    expect(next.seq.insight).toBe(11)
  })

  it('skips an update the snapshot already has', () => {
    const state = fromSnapshot(match, flow, events, { matchId: 'm1', prediction, seq: 12 })
    expect(applyMessage(state, update(12, analysis(20)))).toBe(state)
  })
})
