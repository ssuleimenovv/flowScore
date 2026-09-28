import { describe, expect, it } from 'vitest'
import { statLines } from '../stats'

describe('statLines', () => {
  it('formats each row the way the board shows it', () => {
    const lines = statLines([
      { key: 'possession', home: 64, away: 36 },
      { key: 'xg', home: 1.84, away: 0.6 },
      { key: 'shots', home: 14, away: 6 },
    ])
    expect(lines.map((l) => [l.label, l.home, l.away])).toEqual([
      ['Владение', '64%', '36%'],
      ['xG', '1.84', '0.60'],
      ['Удары', '14', '6'],
    ])
    expect(lines[2]!.homeWidth).toBe('70.0%')
  })

  it('draws empty bars while both sides are at zero', () => {
    const [line] = statLines([{ key: 'shots', home: 0, away: 0 }])
    expect([line!.homeWidth, line!.awayWidth]).toEqual(['0%', '0%'])
  })

  it('leaves out a key it has no label for', () => {
    expect(statLines([{ key: 'dribbles', home: 3, away: 1 }])).toEqual([])
  })
})
