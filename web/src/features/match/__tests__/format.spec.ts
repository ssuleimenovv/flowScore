import { describe, expect, it } from 'vitest'
import type { Match } from '@/shared/api/types'
import { flowShare, formatClock, formatDelta, formatMeta, formatMinute } from '../format'

describe('formatMinute', () => {
  it('shows the minute the clock is in', () => {
    expect(formatMinute(72 * 60 + 14, 'second_half')).toBe('72′')
  })

  it('shows added time after 45 and 90', () => {
    expect(formatMinute(45 * 60 + 30, 'first_half')).toBe('45+1′')
    expect(formatMinute(92 * 60 + 5, 'second_half')).toBe('90+3′')
  })
})

describe('formatClock', () => {
  it('pads the seconds', () => {
    expect(formatClock(72 * 60 + 4)).toBe('72:04')
  })
})

describe('formatMeta', () => {
  const base = {
    competition: { id: '2', name: 'Premier League', round: '37' },
    venue: { name: 'Etihad Stadium', city: 'Manchester' },
  } as Match

  it('joins competition, round and venue', () => {
    expect(formatMeta(base)).toBe('Premier League · 37-й тур · Etihad Stadium, Manchester')
  })

  it('skips what is missing', () => {
    expect(formatMeta({ ...base, venue: undefined, competition: { id: '2', name: 'КПЛ' } })).toBe(
      'КПЛ',
    )
  })
})

describe('formatDelta', () => {
  it('marks the direction', () => {
    expect(formatDelta(17.6)).toBe('▲ +18')
    expect(formatDelta(-21.2)).toBe('▼ −21')
    expect(formatDelta(0)).toBe('▲ +0')
  })
})

describe('flowShare', () => {
  it('splits 100 percent between the sides', () => {
    expect(flowShare({ home: 82, away: 37 })).toEqual({ home: 69, away: 31 })
  })

  it('is even when both sides are at zero', () => {
    expect(flowShare({ home: 0, away: 0 })).toEqual({ home: 50, away: 50 })
  })
})
