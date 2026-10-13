// @vitest-environment node
import { describe, expect, it } from 'vitest'
import type { FlowRecord, RoundMatch } from '@/shared/api/types'
import {
  flowText,
  flowWidth,
  goals,
  recordRows,
  roundStatus,
  roundTabs,
  scoreText,
  seasonLine,
  zoneMark,
} from '../league'

const team = (code: string, name: string) => ({ id: code, code, name })

function match(over: Partial<RoundMatch>): RoundMatch {
  return {
    id: '1',
    kickoffAt: '2026-10-17T14:00:00Z',
    state: 'pre',
    clock: '',
    home: team('FUL', 'Fulham'),
    away: team('MNC', 'Manchester City'),
    score: null,
    flow: null,
    ...over,
  }
}

describe('the table', () => {
  it('names the season and the round', () => {
    expect(seasonLine({ season: '2026/27', round: 8 })).toBe('Сезон 2026/27 · 8-й тур')
  })

  it('marks only the places the board marks', () => {
    expect(zoneMark('champions')).toBe('champions')
    expect(zoneMark('europa')).toBe('europa')
    expect(zoneMark('conference')).toBe('none')
    expect(zoneMark('relegation')).toBe('none')
    expect(zoneMark(null)).toBe('none')
  })

  it('writes goals and Flow as the board does', () => {
    expect(goals({ goalsFor: 17, goalsAgainst: 5 })).toBe('17:5')
    expect(flowText(73.6)).toBe('74')
    expect(flowText(null)).toBe('—')
    expect(flowWidth(73.6)).toBe('74%')
    expect(flowWidth(null)).toBe('0%')
  })
})

describe('the rounds', () => {
  it('shows the round between its neighbours, inside the season', () => {
    expect(roundTabs(8, 38)).toEqual([7, 8, 9])
    expect(roundTabs(1, 38)).toEqual([1, 2, 3])
    expect(roundTabs(38, 38)).toEqual([36, 37, 38])
    expect(roundTabs(1, 2)).toEqual([1, 2])
    expect(roundTabs(0, 0)).toEqual([])
  })

  it('writes the score, or a dash before kick-off', () => {
    expect(scoreText(match({ score: { home: 1, away: 3 } }))).toBe('1:3')
    expect(scoreText(match({}))).toBe('—')
  })

  it('says where each match stands', () => {
    const now = new Date('2026-10-17T10:00:00Z')
    const tz = 'Asia/Almaty' // UTC+5: 17 October, 15:00
    expect(roundStatus(match({ state: 'post' }), now, tz)).toBe('ФТ')
    expect(roundStatus(match({ state: 'in', clock: "45'+2'" }), now, tz)).toBe('45′+2′')
    expect(roundStatus(match({ kickoffAt: '2026-10-17T14:00:00Z' }), now, tz)).toBe('19:00')
    expect(roundStatus(match({ kickoffAt: '2026-10-18T13:00:00Z' }), now, tz)).toBe('Завтра')
    expect(roundStatus(match({ kickoffAt: '2026-10-28T19:00:00Z' }), now, tz)).toBe('29 окт')
  })
})

describe('the records', () => {
  it('ranks them as the board does', () => {
    const r: FlowRecord = {
      matchId: '9',
      round: 2,
      home: team('MAN', 'Manchester United'),
      away: team('IPS', 'Ipswich Town'),
      score: { home: 5, away: 2 },
      side: 'home',
      flow: 78.37,
      peak: 99,
    }
    expect(recordRows([r])).toEqual([
      {
        id: '9',
        rank: '01',
        match: 'Manchester United — Ipswich Town',
        round: 'тур 2',
        flow: '78',
      },
    ])
  })
})
