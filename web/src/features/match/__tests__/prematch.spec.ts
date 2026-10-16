import { describe, expect, it } from 'vitest'
import type { Match, MatchPreMatch, PastMatch } from '@/shared/api/types'
import { formRows, headToHead, meetingRows, newsRows, shortDate } from '../prematch'

const everton = { id: '368', code: 'EVE', name: 'Everton' }
const chelsea = { id: '363', code: 'CHE', name: 'Chelsea' }
const spurs = { id: '367', code: 'TOT', name: 'Tottenham Hotspur' }
const match = { home: everton, away: chelsea } as Match

const past = (fields: Partial<PastMatch>): PastMatch => ({
  id: 'm',
  date: '2026-09-12T16:30:00Z',
  competition: 'Premier League',
  home: everton,
  away: chelsea,
  score: { home: 0, away: 0 },
  flow: null,
  ...fields,
})

const pre: MatchPreMatch = {
  matchId: 'm1',
  home: {
    matches: [
      { ...past({ id: 'a', home: spurs, away: everton }), result: 'draw' },
      {
        ...past({ id: 'b', score: { home: 1, away: 0 }, competition: 'Carabao Cup' }),
        result: 'win',
      },
    ],
    flow: 44.4,
    flowGames: 1,
  },
  away: { matches: [], flow: null, flowGames: 0 },
  meetings: {
    homeWins: 1,
    draws: 1,
    awayWins: 3,
    matches: [
      past({
        id: 'c',
        date: '2026-03-21T17:30:00Z',
        score: { home: 3, away: 0 },
        flow: { home: 58.2, away: 60.6 },
      }),
      past({ id: 'd', home: chelsea, away: everton, score: { home: 2, away: 0 } }),
    ],
  },
  news: [
    { headline: 'Top 50', url: 'https://www.espn.com/x', publishedAt: '2026-10-16T11:13:00Z' },
  ],
}

describe('pre-match', () => {
  it("shows each team's last matches as В / Н / П with its Flow", () => {
    const [home, away] = formRows(pre, match)
    expect(home).toMatchObject({
      side: 'home',
      code: 'EVE',
      flow: '44',
      flowNote: 'ср. Flow, 1 матч',
    })
    expect(home!.results.map((r) => r.letter)).toEqual(['Н', 'В'])
    expect(home!.results[0]!.title).toBe(
      'Tottenham Hotspur 0:0 Everton, 12 сен 2026 · Premier League',
    )
    // No match the season read: no number
    expect(away).toMatchObject({ flow: null, flowNote: 'ср. Flow, 0 матчей', results: [] })
  })

  it('sums up the meetings for this match’s home and away team', () => {
    expect(headToHead(pre, match)).toEqual({
      home: 'Everton · 1',
      draws: 'Ничьи · 1',
      away: 'Chelsea · 3',
      widths: { home: '20.0%', draw: '20.0%', away: '60.0%' },
      note: 'последние 2 встречи',
    })
  })

  it('lists the meetings, with Flow only for the ones the season read', () => {
    expect(meetingRows(pre)).toEqual([
      {
        id: 'c',
        date: shortDate('2026-03-21T17:30:00Z'),
        title: 'Everton — Chelsea',
        competition: 'Premier League',
        score: '3:0',
        flow: '58–61',
      },
      expect.objectContaining({ id: 'd', title: 'Chelsea — Everton', score: '2:0', flow: null }),
    ])
  })

  it('says how long ago each story came out', () => {
    const at = (iso: string) => newsRows(pre, new Date(iso))[0]!.meta
    expect(at('2026-10-16T11:25:00Z')).toBe('ESPN · 12 мин назад')
    expect(at('2026-10-16T14:20:00Z')).toBe('ESPN · 3 ч назад')
    expect(at('2026-10-18T12:00:00Z')).toBe('ESPN · 2 дня назад')
  })
})
