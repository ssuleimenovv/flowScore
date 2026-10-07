// @vitest-environment node
import { describe, expect, it } from 'vitest'
import type { MatchSummary } from '@/shared/api/types'
import { dayForecast, flowPeaks, insightLine, liveMinute, matchDay, sparkline } from '../home'

function match(over: Partial<MatchSummary> & { id: string }): MatchSummary {
  return {
    status: 'live',
    kickoffAt: '2026-10-05T18:00:00Z',
    competition: { id: '2', name: 'Premier League' },
    home: { id: '36', code: 'MCI', name: 'Сити' },
    away: { id: '1', code: 'ARS', name: 'Арсенал' },
    score: { home: 0, away: 0 },
    halftimeScore: null,
    clock: { elapsedSeconds: 600, period: 'first_half', observedAt: '2026-10-05T18:10:00Z' },
    stats: [],
    seq: 1,
    flow: { home: 50, away: 50 },
    delta10: { home: 0, away: 0 },
    factors: [],
    points: [],
    prediction: null,
    explanation: null,
    ...over,
  }
}

const points = (home: number[], away: number[]) =>
  home.map((h, i) => ({ minute: i + 1, home: h, away: away[i] ?? 15 }))

describe('matchDay', () => {
  it('splits live from upcoming and leaves finished matches out', () => {
    const day = matchDay(
      [
        match({ id: 'a', status: 'live' }),
        match({ id: 'b', status: 'halftime' }),
        match({ id: 'c', status: 'scheduled' }),
        match({ id: 'd', status: 'finished' }),
      ],
      null,
    )
    expect(day.live.map((m) => m.id)).toEqual(['a', 'b'])
    expect(day.upcoming.map((m) => m.id)).toEqual(['c'])
  })

  it('keeps one competition', () => {
    const day = matchDay(
      [match({ id: 'a' }), match({ id: 'b', competition: { id: '9', name: 'КПЛ' } })],
      'КПЛ',
    )
    expect(day.live.map((m) => m.id)).toEqual(['b'])
  })
})

describe('liveMinute', () => {
  it('ticks on from when the clock was measured', () => {
    const m = match({ id: 'a' })
    expect(liveMinute(m, Date.parse('2026-10-05T18:10:00Z'))).toBe('10′')
    expect(liveMinute(m, Date.parse('2026-10-05T18:12:30Z'))).toBe('12′')
  })
})

describe('flowPeaks', () => {
  it('ranks the biggest rise over five minutes', () => {
    const slow = match({
      id: 'slow',
      points: points([40, 40, 40, 40, 40, 50], [30, 30, 30, 30, 30, 30]),
    })
    const fast = match({
      id: 'fast',
      points: points([20, 20, 20, 20, 20, 20], [20, 30, 40, 50, 60, 80]),
    })
    const peaks = flowPeaks([slow, fast])
    expect(peaks.map((p) => [p.id, p.side, p.rise, p.flow])).toEqual([
      ['fast', 'away', 60, 80],
      ['slow', 'home', 10, 50],
    ])
  })

  it('skips matches without a rise or without five minutes played', () => {
    const flat = match({
      id: 'flat',
      points: points([50, 50, 50, 50, 50, 50], [50, 50, 50, 50, 50, 50]),
    })
    const young = match({ id: 'young', points: points([20, 90], [20, 20]) })
    expect(flowPeaks([flat, young])).toEqual([])
  })
})

describe('insightLine', () => {
  it("shows the agent's headline once it has one", () => {
    const m = match({
      id: 'a',
      explanation: {
        title: 'Арсенал забирает инициативу',
        text: 'Четыре удара за семь минут.',
        minute: 30,
        generatedAt: '2026-10-05T18:31:00Z',
      },
    })
    expect(insightLine(m)).toBe('Арсенал забирает инициативу')
  })

  it('falls back to the headline from Flow', () => {
    expect(insightLine(match({ id: 'a' }))).toBe('Равная игра')
  })
})


describe('sparkline', () => {
  it('maps minutes and Flow into the 300 × 60 box', () => {
    const m = match({ id: 'a', points: [{ minute: 0, home: 100, away: 0 }] })
    expect(sparkline(m, 'home')).toBe('0.0,6.0')
    expect(sparkline(m, 'away')).toBe('0.0,58.0')
  })
})

describe('dayForecast', () => {
  const p = (home: number, draw: number, away: number) => ({
    current: { home, draw, away },
    preMatch: { home, draw, away },
    model: 'Poisson v1',
  })

  it('picks the match the model is surest about', () => {
    const f = dayForecast([
      match({ id: 'close', status: 'scheduled', prediction: p(35, 30, 35) }),
      match({ id: 'clear', status: 'scheduled', prediction: p(20, 26, 54) }),
    ])
    expect(f?.match.id).toBe('clear')
    expect(f?.text).toMatch(/^Арсенал — фаворит/)
  })

  it('calls an even match even', () => {
    const f = dayForecast([match({ id: 'a', status: 'scheduled', prediction: p(35, 30, 35) })])
    expect(f?.text).toMatch(/^Равный матч/)
  })

  it('has nothing to say without chances', () => {
    expect(dayForecast([match({ id: 'a', status: 'scheduled' })])).toBeNull()
  })
})
