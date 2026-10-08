// @vitest-environment node
import { describe, expect, it } from 'vitest'
import type { Insight, Match, Prediction, WsMessage } from '@/shared/api/types'
import { applyMessage, fromSnapshot } from '../matchState'
import { finalColumns, finalVerdict, formatPreMatch, outcomeColumns } from '../outcome'

const match = {
  home: { id: '36', code: 'MCI', name: 'Сити' },
  away: { id: '1', code: 'ARS', name: 'Арсенал' },
  clock: { elapsedSeconds: 0, period: 'first_half', observedAt: '2026-10-02T10:00:00Z' },
  score: { home: 0, away: 0 },
  seq: 10,
} as Match

const prediction: Prediction = {
  current: { home: 52, draw: 25, away: 23 },
  preMatch: { home: 38, draw: 30, away: 32 },
  model: 'Poisson v1',
}

describe('outcomeColumns', () => {
  it('shows each chance and how it moved since kick-off', () => {
    expect(outcomeColumns(match, prediction)).toEqual([
      { key: 'home', label: 'Победа Сити', value: 52, change: '▲ 14 с начала' },
      { key: 'draw', label: 'Ничья', value: 25, change: '▼ 5' },
      { key: 'away', label: 'Победа Арсенал', value: 23, change: '▼ 9' },
    ])
  })

  it('marks no change before anything happened', () => {
    const still = { ...prediction, current: prediction.preMatch }
    expect(outcomeColumns(match, still).map((c) => c.change)).toEqual([
      '= 0 с начала',
      '= 0',
      '= 0',
    ])
  })

  it('writes the pre-match line', () => {
    expect(formatPreMatch(prediction)).toBe('До матча: 38 · 30 · 32')
  })
})

describe('after the final whistle', () => {
  const over = { ...prediction, current: { home: 0, draw: 0, away: 100 } }

  it('shows the pre-match chances and marks the result', () => {
    const columns = finalColumns(match, over, { home: 4, away: 5 })
    expect(columns.map((c) => [c.value, c.change])).toEqual([
      [38, ''],
      [30, ''],
      [32, '✓ итог'],
    ])
  })

  it("says the model's favourite won", () => {
    expect(finalVerdict(match, over, { home: 2, away: 0 })).toBe(
      'Победа Сити 2:0. Модель это ждала: 38% до матча.',
    )
  })

  it('says the result surprised the model', () => {
    expect(finalVerdict(match, over, { home: 1, away: 1 })).toBe(
      'Ничья 1:1. Неожиданно: до матча модель давала этому 30%.',
    )
  })
})



describe('prediction in the match state', () => {
  const flow = {
    matchId: 'm1',
    current: { home: 20, away: 20 },
    delta10: { home: 0, away: 0 },
    factors: [],
    points: [],
    updatedAt: '2026-10-02T10:00:00Z',
    seq: 10,
  }
  const events = { items: [], seq: 10 }
  const insight: Insight = { matchId: 'm1', prediction, seq: 10 }

  const update = (seq: number, home: number): WsMessage => ({
    type: 'prediction.update',
    matchId: 'm1',
    seq,
    sentAt: '2026-10-02T10:01:00Z',
    data: { home, draw: 100 - home - 20, away: 20 },
  })

  it('moves the current chances and keeps the pre-match ones', () => {
    const next = applyMessage(fromSnapshot(match, flow, events, insight), update(11, 60))
    expect(next.prediction?.current).toEqual({ home: 60, draw: 20, away: 20 })
    expect(next.prediction?.preMatch).toEqual(prediction.preMatch)
  })

  it('skips an update the snapshot already has', () => {
    const state = fromSnapshot(match, flow, events, insight)
    expect(applyMessage(state, update(10, 60))).toBe(state)
  })

  it('ignores updates without a model', () => {
    const state = fromSnapshot(match, flow, events, null)
    expect(applyMessage(state, update(11, 60)).prediction).toBeNull()
  })
})
