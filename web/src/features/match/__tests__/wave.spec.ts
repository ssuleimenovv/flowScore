import { describe, expect, it } from 'vitest'
import type { MatchEvent } from '@/shared/api/types'
import { buildWave, goalMarks, surname } from '../wave'

const points = [
  { minute: 1, home: 10, away: 0 },
  { minute: 2, home: 0, away: 20 },
  { minute: 3, home: 40, away: 10 },
]

describe('buildWave', () => {
  it('is home minus away with a one-minute window', () => {
    expect(buildWave(points, 1).slice(0, 3)).toEqual([10, -20, 30])
  })

  it('averages the trailing window', () => {
    // minute 3 with a 5-minute window: (10 − 20 + 30) / 3
    expect(buildWave(points, 5)[2]).toBeCloseTo(20 / 3)
  })

  it('leaves unplayed minutes empty up to 90', () => {
    const wave = buildWave(points, 1)
    expect(wave).toHaveLength(90)
    expect(wave[3]).toBeNull()
    expect(wave[89]).toBeNull()
  })

  it('grows past 90 in added time', () => {
    expect(buildWave([{ minute: 93, home: 0, away: 0 }], 1)).toHaveLength(93)
  })
})

describe('goalMarks', () => {
  const goal = (minute: number, side: 'home' | 'away', name: string) =>
    ({ id: `${minute}`, type: 'goal', side, minute, player: { id: '', name } }) as MatchEvent

  it('labels each goal with the score it made, oldest first', () => {
    // The stream keeps events newest first
    const events = [goal(58, 'home', 'Phil Foden'), goal(23, 'away', 'Bukayo Saka')]
    expect(goalMarks(events)).toEqual([
      { minute: 23, side: 'away', label: '0:1 Saka' },
      { minute: 58, side: 'home', label: '1:1 Foden' },
    ])
  })
})

describe('surname', () => {
  it('takes the last word', () => {
    expect(surname('Alexis Alejandro Sánchez Sánchez')).toBe('Sánchez')
    expect(surname('Fernandinho')).toBe('Fernandinho')
  })
})
