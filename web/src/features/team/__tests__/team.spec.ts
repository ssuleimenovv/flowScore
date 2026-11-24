// @vitest-environment node
import { describe, expect, it } from 'vitest'
import type { TeamMatch } from '@/shared/api/types'
import {
  averageFlow,
  form,
  kickoffLine,
  lastFlows,
  leagueLine,
  opponentLine,
  opponentOf,
  results,
  scoreFor,
  upcoming,
  venueOf,
} from '../team'

const team = (code: string) => ({ id: code, code, name: code })

function match(over: Partial<TeamMatch>): TeamMatch {
  return {
    id: '1',
    leagueId: 'eng.1',
    league: 'Премьер-лига',
    kickoffAt: '2026-10-17T14:00:00Z',
    state: 'post',
    clock: '',
    home: team('MNC'),
    away: team('ARS'),
    score: { home: 1, away: 0 },
    flow: { home: 70, away: 30 },
    result: 'win',
    ...over,
  }
}

// City's season so far: two read, one not read yet, one under way, one to come
const season = [
  match({ id: 'a', result: 'win' }),
  match({
    id: 'b',
    home: team('LIV'),
    away: team('MNC'),
    flow: { home: 60, away: 50 },
    result: 'loss',
  }),
  match({ id: 'c', flow: null, result: 'draw' }),
  match({ id: 'd', state: 'in', result: '', flow: null }),
  match({ id: 'e', state: 'pre', result: '', score: null, flow: null }),
]

describe('team', () => {
  it('reads each match from the side of the team', () => {
    expect(opponentOf(season[1]!, 'MNC').code).toBe('LIV')
    expect(lastFlows(season, 'MNC').map((b) => b.flow)).toEqual([70, 50])
    expect(averageFlow(season, 'MNC')).toBe(60)
    expect(averageFlow([], 'MNC')).toBeNull()
  })

  it('takes the form from the matches over, oldest first', () => {
    expect(form(season)).toEqual(['win', 'loss', 'draw'])
    expect(form(season, 2)).toEqual(['loss', 'draw'])
  })

  it('splits the season into the matches to come and the results', () => {
    expect(upcoming(season).map((m) => m.id)).toEqual(['d', 'e'])
    expect(results(season, 2).map((m) => m.id)).toEqual(['c', 'b'])
  })

  it('names the league, its country and the place', () => {
    const league = { id: 'eng.1', name: 'Премьер-лига', rank: 3 }
    expect(leagueLine(league)).toBe('Премьер-лига · Англия · 3-е место')
    expect(leagueLine(league, false)).toBe('Премьер-лига · 3-е место')
    expect(leagueLine({ id: 'uefa.champions', name: 'Лига чемпионов', rank: 12 })).toBe(
      'Лига чемпионов · 12-е место',
    )
    expect(leagueLine(null)).toBe('')
  })

  it('writes a match from the side of the team', () => {
    const away = season[1]! // LIV 60:50 MNC... the score is 1:0 for the home side
    expect(scoreFor(away, 'MNC')).toBe('0:1')
    expect(venueOf(away, 'MNC')).toBe('В гостях')
    expect(opponentLine(away, 'MNC', 'eng.1')).toBe('LIV')
    expect(
      opponentLine(
        { ...away, leagueId: 'uefa.champions', league: 'Лига чемпионов' },
        'MNC',
        'eng.1',
      ),
    ).toBe('LIV · Лига чемпионов')
  })

  it('names the day of a kick-off', () => {
    const m = match({ kickoffAt: '2026-11-24T20:00:00Z' })
    expect(kickoffLine(m, new Date('2026-11-24T10:00:00Z'), 'UTC')).toBe('Сегодня, 20:00')
    expect(kickoffLine(m, new Date('2026-11-20T10:00:00Z'), 'UTC')).toBe('24 нояб, 20:00')
  })
})
