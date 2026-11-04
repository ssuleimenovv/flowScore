// @vitest-environment node
import { describe, expect, it } from 'vitest'
import type { FlowRecord, LeagueItem, RoundMatch } from '@/shared/api/types'
import {
  competitionGroups,
  competitionKind,
  flowText,
  flowWidth,
  goals,
  isEuroCup,
  recordRows,
  roundStatus,
  roundTabs,
  scoreText,
  seasonLine,
  stageTitle,
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
    const l = { season: '2026/27', round: 8, rounds: 38, stage: '', stages: [] }
    expect(seasonLine(l)).toBe('Сезон 2026/27 · 8-й тур')
  })

  it('marks only the places the board marks', () => {
    expect(zoneMark('champions')).toBe('champions')
    expect(zoneMark('europa')).toBe('europa')
    expect(zoneMark('conference')).toBe('none')
    expect(zoneMark('relegation')).toBe('none')
    expect(zoneMark(null)).toBe('none')
    expect(zoneMark('last16')).toBe('champions')
    expect(zoneMark('playoff')).toBe('europa')
  })

  it('tells a European cup from a league', () => {
    expect(isEuroCup('uefa.champions')).toBe(true)
    expect(isEuroCup('uefa.europa.conf')).toBe(true)
    expect(isEuroCup('eng.1')).toBe(false)
    expect(isEuroCup(null)).toBe(false)
  })

  it('writes goals and Flow as the board does', () => {
    expect(goals({ goalsFor: 17, goalsAgainst: 5 })).toBe('17:5')
    expect(flowText(73.6)).toBe('74')
    expect(flowText(null)).toBe('—')
    expect(flowWidth(73.6)).toBe('74%')
    expect(flowWidth(null)).toBe('0%')
  })
})

describe('the cups', () => {
  const stages = [
    { id: 'first-round', name: '1-й раунд', ties: 16 },
    { id: 'round-of-16', name: '1/8', ties: 8 },
    { id: 'final', name: 'Финал', ties: 1 },
  ]

  it('tells a league, a European cup and a cup apart', () => {
    expect(competitionKind('eng.1')).toBe('league')
    expect(competitionKind('uefa.champions')).toBe('euro')
    expect(competitionKind('uefa.europa.conf')).toBe('euro')
    expect(competitionKind('uefa.super_cup')).toBe('cup')
    expect(competitionKind('eng.fa')).toBe('cup')
    expect(competitionKind('ita.coppa_italia')).toBe('cup')
    expect(competitionKind(null)).toBe('league')
  })

  it('groups the chips as the board does, in the gateway order', () => {
    const item = (id: string): LeagueItem => ({ id, name: id, ready: true })
    const ids = ['eng.1', 'esp.1', 'uefa.champions', 'uefa.super_cup', 'eng.fa', 'ita.coppa_italia']
    expect(competitionGroups(ids.map(item)).map((g) => [g.name, g.items.map((l) => l.id)])).toEqual(
      [
        ['Лиги', ['eng.1', 'esp.1']],
        ['Еврокубки', ['uefa.champions', 'uefa.super_cup']],
        ['Кубки', ['eng.fa', 'ita.coppa_italia']],
      ],
    )
    expect(competitionGroups([item('eng.1')]).map((g) => g.name)).toEqual(['Лиги'])
    expect(competitionGroups([])).toEqual([])
  })

  it('writes a stage as the season line does', () => {
    expect(stageTitle('1/8')).toBe('1/8 финала')
    expect(stageTitle('1/16')).toBe('1/16 финала')
    expect(stageTitle('Финал')).toBe('Финал')
    expect(stageTitle('3-й раунд')).toBe('3-й раунд')
    expect(stageTitle('Стыки')).toBe('Стыки')
  })

  it('names the stage being played over the bracket', () => {
    const cup = { season: '2026/27', round: 0, rounds: 0, stages }
    expect(seasonLine({ ...cup, stage: 'first-round' })).toBe('Сезон 2026/27 · 1-й раунд')
    expect(seasonLine({ ...cup, stage: 'round-of-16' })).toBe('Сезон 2026/27 · 1/8 финала')
    expect(seasonLine({ ...cup, stage: '', stages: [] })).toBe('Сезон 2026/27')
    const euro = { season: '2026/27', round: 4, rounds: 8, stage: '', stages: [] }
    expect(seasonLine(euro)).toBe('Сезон 2026/27 · 4-й тур')
    expect(seasonLine({ ...euro, stage: 'round-of-16', stages })).toBe('Сезон 2026/27 · 1/8 финала')
    expect(seasonLine({ season: '2026/27', round: 1, rounds: 1, stage: '', stages: [] })).toBe(
      'Сезон 2026/27 · Матч за трофей',
    )
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
      stage: '',
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

  it('names a knockout match by its stage', () => {
    const r: FlowRecord = {
      matchId: '7',
      round: 0,
      stage: 'third-round',
      home: team('WRE', 'Wrexham'),
      away: team('NFO', 'Nottingham Forest'),
      score: { home: 3, away: 3 },
      side: 'home',
      flow: 84.2,
      peak: 97,
    }
    const stages = [
      { id: 'third-round', name: '3-й раунд', ties: 32 },
      { id: 'round-of-16', name: '1/8', ties: 8 },
    ]
    expect(recordRows([r], stages)[0]?.round).toBe('3-й раунд')
    expect(recordRows([{ ...r, stage: 'round-of-16' }], stages)[0]?.round).toBe('1/8 финала')
    expect(recordRows([r])[0]?.round).toBe('')
  })
})
