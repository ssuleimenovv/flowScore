import { describe, expect, it } from 'vitest'
import type { FlowFactor, Match } from '@/shared/api/types'
import { explain } from '../explain'
import type { MatchLive } from '../matchState'

const match = {
  home: { id: '36', code: 'MCI', name: 'Сити' },
  away: { id: '1', code: 'ARS', name: 'Арсенал' },
} as Match

function live(flow: { home: number; away: number }, delta: number, factors: FlowFactor[]) {
  return { match, flow, delta10: { home: delta, away: 0 }, factors } as unknown as MatchLive
}

const factors: FlowFactor[] = [
  { side: 'home', key: 'shots', value: 13.9, count: 4, minutes: 7 },
  { side: 'home', key: 'possession', value: 9, count: 10, minutes: 10, share: 0.68 },
  { side: 'home', key: 'corners', value: 2, count: 1, minutes: 3 },
  { side: 'away', key: 'shots', value: 5, count: 2, minutes: 5 },
]

describe('explain', () => {
  it('names the leader and how its Flow moved', () => {
    const e = explain(live({ home: 82, away: 37 }, 18, factors), 5)
    expect(e.title).toBe('Сити забирает инициативу')
    expect(e.text).toBe('Поток Сити вырос с 64 до 82 за 10 минут.')
  })

  it('calls a close game even', () => {
    expect(explain(live({ home: 50, away: 45 }, 0, []), 5).title).toBe('Равная игра')
  })

  it('lists the factors for the leader, the other side as minuses', () => {
    const lines = explain(live({ home: 82, away: 37 }, 18, factors), 5).lines
    expect(lines.map((l) => [l.text, l.value, l.side])).toEqual([
      ['4 удара за последние 7 минут', '+14', 'home'],
      ['Владение 68% за последние 10 минут', '+9', 'home'],
      ['Арсенал: 2 удара за последние 5 минут', '−5', 'away'],
      ['1 угловой за последние 3 минуты', '+2', 'home'],
    ])
    // 15 points fill a bar, as on the board
    expect(lines[2]!.width).toBe('33.3%')
  })

  it('keeps the strongest when there are more factors than room', () => {
    const lines = explain(live({ home: 82, away: 37 }, 18, factors), 3).lines
    expect(lines.map((l) => l.key)).toEqual(['home-shots', 'home-possession', 'away-shots'])
  })

  it('leaves out what holds the other team back', () => {
    const mirror: FlowFactor = {
      side: 'away',
      key: 'possession',
      value: -5,
      count: 10,
      minutes: 10,
      share: 0.32,
    }
    const keys = explain(live({ home: 82, away: 37 }, 18, [...factors, mirror]), 5).lines.map(
      (l) => l.key,
    )
    expect(keys).not.toContain('away-possession')
  })

  it('says "за последнюю минуту" for one minute', () => {
    const shot: FlowFactor = { side: 'home', key: 'shots', value: 6, count: 1, minutes: 1 }
    expect(explain(live({ home: 60, away: 30 }, 0, [shot]), 5).lines[0]!.text).toBe(
      '1 удар за последнюю минуту',
    )
  })

  it('drops a possession share that contradicts the value', () => {
    const odd: FlowFactor = {
      side: 'home',
      key: 'possession',
      value: 2.8,
      count: 10,
      minutes: 10,
      share: 0.28,
    }
    expect(explain(live({ home: 60, away: 30 }, 0, [odd]), 5).lines[0]!.text).toBe('Владение мячом')
  })
})
