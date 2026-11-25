// @vitest-environment node
import { describe, expect, it } from 'vitest'
import type { MatchSummary } from '@/shared/api/types'
import {
  competitions,
  dayForecast,
  flowPeaks,
  groups,
  inGroup,
  insightLine,
  liveMinute,
  matchDay,
  sparkline,
} from '../home'

function match(over: Partial<MatchSummary> & { id: string }): MatchSummary {
  return {
    source: 'espn',
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
  // Thursday evening in Almaty, UTC+5
  const now = new Date('2026-10-08T15:00:00Z')
  const tz = 'Asia/Almaty'
  const at = (id: string, kickoffAt: string) => match({ id, status: 'scheduled', kickoffAt })

  it('splits live from upcoming and leaves finished matches out', () => {
    const day = matchDay(
      [
        match({ id: 'a', status: 'live' }),
        match({ id: 'b', status: 'halftime' }),
        at('c', '2026-10-08T17:00:00Z'),
        match({ id: 'd', status: 'finished' }),
      ],
      null,
      now,
      tz,
    )
    expect(day.live.map((m) => m.id)).toEqual(['a', 'b'])
    expect(day.days.map((d) => d.matches.map((m) => m.id))).toEqual([['c']])
  })

  it("puts the real matches under their day, in the viewer's time zone", () => {
    const day = matchDay(
      [
        at('today', '2026-10-08T17:00:00Z'),
        at('friday', '2026-10-09T18:30:00Z'), // 23:30 on Friday in Almaty
        at('midnight', '2026-10-09T19:00:00Z'), // 00:00 on Saturday there
        at('saturday', '2026-10-10T11:30:00Z'),
      ],
      null,
      now,
      tz,
    )
    expect(day.days.map((d) => [d.title, d.matches.map((m) => m.id)])).toEqual([
      ['Сегодня', ['today']],
      ['Завтра', ['friday']],
      ['Суббота, 10 октября', ['midnight', 'saturday']],
    ])
  })

  it('keeps the demo replays apart from the real matches', () => {
    const day = matchDay(
      [
        match({ id: 'real' }),
        match({ id: 'replay', source: 'replay' }),
        match({ id: 'next', source: 'replay', status: 'scheduled' }),
      ],
      null,
      now,
      tz,
    )
    expect(day.live.map((m) => m.id)).toEqual(['real'])
    expect(day.demo.live.map((m) => m.id)).toEqual(['replay'])
    expect(day.demo.upcoming.map((m) => m.id)).toEqual(['next'])
    expect(day.days).toEqual([])
  })

  it('keeps one competition', () => {
    const day = matchDay(
      [match({ id: 'a' }), match({ id: 'b', competition: { id: '9', name: 'КПЛ' } })],
      '9',
    )
    expect(day.live.map((m) => m.id)).toEqual(['b'])
  })
})

describe('competitions', () => {
  const cup = match({ id: 'a', competition: { id: 'ger.dfb_pokal', name: 'Кубок Германии' } })
  const ucl = match({ id: 'c', competition: { id: 'uefa.champions', name: 'Лига чемпионов' } })
  const liga = match({ id: 'd', competition: { id: 'esp.1', name: 'Ла Лига' } })
  const demo = match({
    id: 'r',
    source: 'replay',
    competition: { id: '2', name: 'Premier League' },
  })

  it('lists the leagues, then the European cups, then the cups, and no chip for the demo', () => {
    const list = competitions([
      cup,
      demo,
      ucl,
      liga,
      match({ id: 'e', competition: liga.competition }),
    ])
    expect(list.map((c) => c.id)).toEqual(['esp.1', 'uefa.champions', 'ger.dfb_pokal'])
    expect(groups(list)).toEqual(['Все', 'Лиги', 'Еврокубки', 'Кубки'])
    expect(groups([list[1]!])).toEqual(['Все', 'Еврокубки'])
  })

  it('keeps a group’s matches, and everything, the demo too, for «Все»', () => {
    const items = [cup, demo, ucl, liga]
    expect(inGroup(items, 'Еврокубки').map((m) => m.id)).toEqual(['c'])
    expect(inGroup(items, 'Кубки').map((m) => m.id)).toEqual(['a'])
    expect(inGroup(items, 'Все')).toHaveLength(4)
  })

  it('puts a day’s matches in a card per competition, in the same order at every poll', () => {
    const at = (id: string, kickoffAt: string, comp: MatchSummary['competition'], home: string) =>
      match({
        id,
        status: 'scheduled',
        kickoffAt,
        competition: comp,
        home: { id, code: id, name: home },
      })
    const now = new Date('2026-11-25T10:00:00Z')
    const items = [
      at('1', '2026-11-25T20:00:00Z', ucl.competition, 'Sporting CP'),
      at('2', '2026-11-25T17:45:00Z', cup.competition, 'Hertha'),
      at('3', '2026-11-25T20:00:00Z', ucl.competition, 'Club Brugge'),
      at('4', '2026-11-25T15:00:00Z', liga.competition, 'Sevilla'),
    ]
    for (const list of [items, [...items].reverse()]) {
      const [d] = matchDay(list, null, now, 'UTC').days
      expect(d!.competitions.map((c) => [c.id, c.matches.map((m) => m.id)])).toEqual([
        ['esp.1', ['4']],
        ['uefa.champions', ['3', '1']],
        ['ger.dfb_pokal', ['2']],
      ])
      expect(d!.matches.map((m) => m.id)).toEqual(['4', '2', '3', '1'])
    }
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
