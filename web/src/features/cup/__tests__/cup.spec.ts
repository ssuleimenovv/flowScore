// @vitest-environment node
import { describe, expect, it } from 'vitest'
import type { Stage, Tie, TieLeg, TieSide } from '@/shared/api/types'
import {
  isTwoLegged,
  legsNote,
  liveClock,
  matchesNote,
  openingStage,
  tieCells,
  tieColumns,
  tieFootnote,
  tieMark,
  tieNote,
} from '../cup'

const team = (code: string, name: string) => ({ id: code, code, name })

function side(over: Partial<TieSide>): TieSide {
  return {
    team: team('BAR', 'Барселона'),
    goals: [null],
    total: null,
    shootout: null,
    through: null,
    ...over,
  }
}

function leg(over: Partial<TieLeg>): TieLeg {
  return { matchId: '1', kickoffAt: '2026-10-17T14:00:00Z', state: 'post', clock: '', ...over }
}

function tie(over: Partial<Tie>): Tie {
  return {
    home: side({}),
    away: side({ team: team('BEN', 'Бенфика') }),
    legs: [leg({})],
    finish: null,
    ...over,
  }
}

const stages: Stage[] = [
  { id: 'knockout-round-playoffs', name: 'Стыки', ties: 8 },
  { id: 'round-of-16', name: '1/8', ties: 8 },
  { id: 'quarterfinals', name: '1/4', ties: 4 },
]

describe('the stages', () => {
  it('opens on the stage being played, else on the last', () => {
    expect(openingStage(stages, 'round-of-16')).toBe('round-of-16')
    expect(openingStage(stages, '')).toBe('quarterfinals')
    expect(openingStage(stages, 'final')).toBe('quarterfinals')
    expect(openingStage([], '')).toBe('')
  })

  it('tells two legs from one match', () => {
    expect(isTwoLegged([tie({ legs: [leg({}), leg({ matchId: '2' })] })])).toBe(true)
    expect(isTwoLegged([tie({})])).toBe(false)
    expect(isTwoLegged([])).toBe(false)
  })

  it('heads the columns as the board does', () => {
    expect(tieColumns(true)).toEqual(['1-й', '2-й', 'Σ'])
    expect(tieColumns(false)).toEqual(['Счёт'])
  })

  it('writes the side notes and the footnote', () => {
    expect(legsNote(true)).toBe('два матча')
    expect(legsNote(false)).toBe('один матч')
    expect(tieFootnote(true)).toBe('Зелёная метка — прошёл дальше. Сумма двух матчей справа.')
    expect(tieFootnote(false)).toBe('Один матч. При ничьей — дополнительное время и пенальти.')
  })

  it('counts the matches of a stage in Russian', () => {
    expect(matchesNote({ ties: 32 }, [])).toBe('32 матча')
    expect(matchesNote({ ties: 21 }, [])).toBe('21 матч')
    expect(matchesNote({ ties: 16 }, [])).toBe('16 матчей')
    expect(matchesNote({ ties: 12 }, [])).toBe('12 матчей')
    expect(matchesNote({ ties: 1 }, [])).toBe('1 матч')
    // two legs a tie: once the ties come, every leg counts
    expect(
      matchesNote({ ties: 2 }, [
        tie({ legs: [leg({}), leg({})] }),
        tie({ legs: [leg({}), leg({})] }),
      ]),
    ).toBe('4 матча')
  })
})

describe('a tie', () => {
  it('marks who went through, who went out and an open tie', () => {
    expect(tieMark({ through: true })).toBe('through')
    expect(tieMark({ through: false })).toBe('out')
    expect(tieMark({ through: null })).toBe('open')
  })

  it('writes each leg, then the sum in bold', () => {
    expect(tieCells(side({ goals: [2, 1], total: 3 }), true)).toEqual(['2', '1', '3'])
    // "Ливерпуль 3 – 3": the second leg is to come
    expect(tieCells(side({ goals: [3, null], total: 3 }), true)).toEqual(['3', '–', '3'])
    expect(tieCells(side({ goals: [null, null], total: null }), true)).toEqual(['–', '–', '–'])
  })

  it('writes a single match as its score', () => {
    expect(tieCells(side({ goals: [3], total: 3 }), false)).toEqual(['3'])
    expect(tieCells(side({ goals: [null], total: null }), false)).toEqual(['–'])
  })

  it('shows the minute, and extra time past the 90', () => {
    expect(liveClock("72'")).toBe('72′')
    expect(liveClock("90'+3'")).toBe('90′+3′')
    expect(liveClock("104'")).toBe('ДВ 104′')
    expect(liveClock('HT')).toBe('HT')
  })
})

describe('the note under a tie', () => {
  const now = new Date('2026-10-17T10:00:00Z')
  const tz = 'Asia/Almaty' // UTC+5: 17 October, 15:00
  const two = (a: Partial<TieLeg>, b: Partial<TieLeg>) => [leg(a), leg({ matchId: '2', ...b })]

  it('is the pill while a leg is played', () => {
    expect(tieNote(tie({ legs: [leg({ state: 'in', clock: "104'" })] }), now, tz)).toEqual({
      live: true,
      text: 'ДВ 104′',
    })
    expect(tieNote(tie({ legs: two({}, { state: 'in', clock: "72'" }) }), now, tz)).toEqual({
      live: true,
      text: '72′',
    })
  })

  it('says how the tie was settled, the first home team first', () => {
    const pens = (legs: TieLeg[]) =>
      tie({
        legs,
        finish: 'pens',
        home: side({ shootout: 4, through: true }),
        away: side({ shootout: 3, through: false }),
      })
    expect(tieNote(pens(two({}, {})), now, tz)?.text).toBe('после ДВ · пен. 4:3')
    expect(tieNote(pens([leg({})]), now, tz)?.text).toBe('пен. 4:3')
    expect(tieNote(tie({ finish: 'aet' }), now, tz)).toEqual({ live: false, text: 'после ДВ' })
    expect(tieNote(tie({ finish: 'pens' }), now, tz)?.text).toBe('пен.')
  })

  it('says nothing under a tie settled in the 90 minutes', () => {
    expect(tieNote(tie({}), now, tz)).toBeNull()
    expect(tieNote(tie({ legs: two({}, {}) }), now, tz)).toBeNull()
  })

  it('says when the next leg is', () => {
    expect(
      tieNote(tie({ legs: two({}, { state: 'pre', kickoffAt: '2026-10-28T19:00:00Z' }) }), now, tz)
        ?.text,
    ).toBe('ответный 29 окт')
    expect(
      tieNote(tie({ legs: two({}, { state: 'pre', kickoffAt: '2026-10-18T13:00:00Z' }) }), now, tz)
        ?.text,
    ).toBe('ответный завтра')
    expect(
      tieNote(
        tie({
          legs: two(
            { state: 'pre', kickoffAt: '2026-10-21T19:00:00Z' },
            { state: 'pre', kickoffAt: '2026-10-28T19:00:00Z' },
          ),
        }),
        now,
        tz,
      )?.text,
    ).toBe('первый 22 окт')
  })

  it('says when a single match kicks off, as a round does', () => {
    const pre = (kickoffAt: string) => tie({ legs: [leg({ state: 'pre', kickoffAt })] })
    expect(tieNote(pre('2026-10-17T14:00:00Z'), now, tz)?.text).toBe('19:00')
    expect(tieNote(pre('2026-10-18T13:00:00Z'), now, tz)?.text).toBe('Завтра')
    expect(tieNote(pre('2026-10-28T19:00:00Z'), now, tz)?.text).toBe('29 окт')
  })
})
