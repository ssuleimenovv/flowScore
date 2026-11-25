// @vitest-environment node
import { describe, expect, it } from 'vitest'
import type { PlayerMatch } from '@/shared/api/types'
import { clubLeague, finishing, logRow, per90, seasonTiles, totals } from '../player'

const team = (code: string) => ({ id: code, code, name: code })

function match(over: Partial<PlayerMatch>): PlayerMatch {
  return {
    id: '1',
    leagueId: 'eng.1',
    league: 'Премьер-лига',
    kickoffAt: '2026-11-21T12:30:00Z',
    home: team('MNC'),
    away: team('FUL'),
    score: { home: 3, away: 1 },
    side: 'home',
    starter: true,
    minutes: 90,
    goals: 0,
    assists: 0,
    shots: 0,
    onTarget: 0,
    xg: 0,
    ...over,
  }
}

const season = [
  match({ goals: 2, assists: 1, shots: 5, onTarget: 3, xg: 1.2 }),
  match({
    id: '2',
    side: 'away',
    home: team('LIV'),
    away: team('MNC'),
    starter: false,
    minutes: 30,
    shots: 1,
    xg: 0.3,
  }),
]

describe('player', () => {
  it('sums the season and counts per 90 minutes', () => {
    const t = totals(season)
    expect(t).toMatchObject({
      matches: 2,
      starts: 1,
      minutes: 120,
      goals: 2,
      shots: 6,
      onTarget: 3,
    })
    expect(per90(t.goals, t.minutes)).toBe('1.50')
    expect(per90(1, 0)).toBe('—')
    expect(seasonTiles(t).map((x) => x.value)).toEqual(['2', '120', '2', '1', '1.5', '6'])
  })

  it('weighs the goals against the xG', () => {
    const f = finishing(totals(season))
    expect(f.diff).toBe('+0.5')
    expect(f.goalsWidth).toBe('100%')
    expect(f.xgWidth).toBe('75%')
    expect(finishing(totals([match({ shots: 4, xg: 2 })])).diff).toBe('−2.0')
  })

  it('leaves the matches without an xG out of the weighing', () => {
    // A cup match: two goals, and no places of the shots at ESPN
    const cup = match({ id: '3', goals: 2, shots: 2, onTarget: 2, xg: null })
    const t = totals([...season, cup])
    expect(t).toMatchObject({ goals: 4, xg: 1.5, xgGoals: 2, xgMinutes: 120, noXG: 1 })
    const f = finishing(t)
    expect(f.goals).toBe(2)
    expect(f.diff).toBe('+0.5')
    expect(f.note).toBe('Не считая 1 матча: в них источник не даёт мест ударов.')
    expect(logRow(cup).xg).toBe('—')
    expect(finishing(totals(season)).note).toBe('')

    const cupsOnly = totals([cup])
    expect(seasonTiles(cupsOnly)[4]!.value).toBe('—')
    expect(finishing(cupsOnly)).toMatchObject({
      note: '',
      text: 'Для этих матчей источник не даёт мест ударов.',
    })
  })

  it('finds the player’s league in his log, past the cups', () => {
    const cup = match({ id: '3', leagueId: 'ger.dfb_pokal' })
    const ucl = match({ id: '4', leagueId: 'uefa.champions' })
    expect(clubLeague([cup, ucl, match({ leagueId: 'ger.1' })])).toBe('ger.1')
    expect(clubLeague([cup])).toBeNull()
  })

  it('writes a match from the side of the player', () => {
    const row = logRow(season[1]!, 'UTC')
    expect(row).toMatchObject({ opponent: 'LIV', venue: 'Г', score: '1:3', date: '21 нояб' })
  })
})
