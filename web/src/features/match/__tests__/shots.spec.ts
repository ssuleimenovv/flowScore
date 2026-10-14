import { describe, expect, it } from 'vitest'
import type { MatchEvent } from '@/shared/api/types'
import { shots, shotsWord, shotTotals } from '../shots'

const event = (fields: Partial<MatchEvent>) =>
  ({ side: 'home', minute: 10, xG: null, modelXG: null, ...fields }) as MatchEvent

describe('shots', () => {
  // Newest first, as the stream keeps them
  const events = [
    event({ id: 'og', type: 'goal', side: 'away', minute: 70, position: { x: 95, y: 50 } }),
    event({ id: 'e5', type: 'corner', minute: 66 }),
    event({ id: 'e4', type: 'shot_off_target', side: 'away', minute: 60, xG: 0.04 }),
    event({
      id: 'e3',
      type: 'shot_blocked',
      side: 'away',
      minute: 52,
      modelXG: 0.12,
      position: { x: 82, y: 30 },
    }),
    event({
      id: 'e2',
      type: 'goal',
      minute: 45,
      addedTime: 2,
      xG: 0.4,
      modelXG: 0.5,
      player: { id: '', name: 'Phil Foden' },
      position: { x: 88, y: 46 },
    }),
  ]

  it('draws the shots with a position, oldest first', () => {
    expect(shots(events).map((s) => s.id)).toEqual(['e2', 'e3'])
  })

  it('sizes a shot by our xG and names who took it', () => {
    expect(shots(events)[0]).toEqual({
      id: 'e2',
      side: 'home',
      outcome: 'goal',
      left: '88%',
      top: '46%',
      size: '34px',
      label: 'Foden 45+2′ · xG 0.50',
    })
  })

  it('turns the away team round to attack the left goal', () => {
    expect(shots(events)[1]).toMatchObject({
      outcome: 'block',
      left: '18%',
      top: '70%',
      size: '17px',
    })
  })

  it('counts every shot in the totals, one without a position too, but not an own goal', () => {
    expect(shotTotals(events)).toEqual({
      home: { xg: '0.50', shots: 1 },
      away: { xg: '0.16', shots: 2 },
    })
  })

  it('says how many shots in Russian', () => {
    expect([0, 1, 2, 5, 11, 12, 21, 24, 112].map(shotsWord)).toEqual([
      '0 ударов',
      '1 удар',
      '2 удара',
      '5 ударов',
      '11 ударов',
      '12 ударов',
      '21 удар',
      '24 удара',
      '112 ударов',
    ])
  })
})
