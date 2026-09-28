import { describe, expect, it } from 'vitest'
import type { Match, MatchEvent } from '@/shared/api/types'
import { chronicle } from '../chronicle'

const match = {
  home: { id: '36', code: 'MCI', name: 'Manchester City' },
  away: { id: '1', code: 'ARS', name: 'Arsenal' },
} as Match

const event = (fields: Partial<MatchEvent>) =>
  ({ side: 'home', minute: 10, ...fields }) as MatchEvent

describe('chronicle', () => {
  // Newest first, as the stream keeps them
  const events = [
    event({ id: 'e4', type: 'foul', minute: 60 }),
    event({ id: 'e3', type: 'yellow_card', side: 'away', minute: 52, flowImpact: 5 }),
    event({
      id: 'e2',
      type: 'goal',
      minute: 45,
      addedTime: 2,
      flowImpact: 23,
      player: { id: '', name: 'Phil Foden' },
      assist: { id: '', name: 'Kevin De Bruyne' },
    }),
    event({ id: 'e1', type: 'goal', side: 'away', minute: 23, flowImpact: 23 }),
  ]
  const items = chronicle(events, match)

  it('skips the events that do not move the flow', () => {
    expect(items.map((i) => i.id)).toEqual(['e3', 'e2', 'e1'])
  })

  it('describes a goal with the score it made', () => {
    expect(items[1]).toMatchObject({
      minute: '45+2′',
      code: 'ГОЛ',
      title: 'Гол! Foden',
      detail: 'Ассист Bruyne · 1:1',
      goal: 'home',
      impact: { text: '+23', side: 'home' },
    })
  })

  it('credits a card to the opponent', () => {
    expect(items[0]!.impact).toEqual({ text: '+5', side: 'home' })
  })
})
