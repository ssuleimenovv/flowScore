import { describe, expect, it } from 'vitest'
import type { LineupPlayer, MatchEvent, TeamLineup } from '@/shared/api/types'
import { benchRows, formationsLine, listRows, minuteText, pitchPlayers } from '../lineups'

const player = (fields: Partial<LineupPlayer>) =>
  ({
    id: '',
    name: '',
    number: null,
    position: null,
    on: null,
    off: null,
    ...fields,
  }) as LineupPlayer

const city: TeamLineup = {
  formation: '4-4-2',
  starters: [
    player({ id: '10', name: 'Sergio Agüero', number: '10', position: { x: 58, y: 33 } }),
    player({
      id: '72',
      name: 'Kelechi Iheanacho',
      number: '72',
      position: { x: 58, y: 67 },
      off: { minute: 70, addedTime: null },
    }),
  ],
  bench: [
    player({ id: '7', name: 'Raheem Sterling', number: '7', on: { minute: 70, addedTime: null } }),
    player({ id: '13', name: 'Willy Caballero', number: '13' }),
  ],
}

const event = (fields: Partial<MatchEvent>) =>
  ({ side: 'home', minute: 10, flowImpact: 0, ...fields }) as MatchEvent

// Newest first, as the stream keeps them
const events = [
  event({ id: 'e5', type: 'substitution', player: { id: '72', name: '' }, flowImpact: 5 }),
  event({ id: 'e4', type: 'yellow_card', player: { id: '72', name: '' }, flowImpact: 4 }),
  event({ id: 'e3', type: 'goal', player: { id: '10', name: '' }, flowImpact: 23 }),
  event({ id: 'e2', type: 'shot_off_target', player: { id: '10', name: '' }, flowImpact: 6.4 }),
  // An own goal of the away team's player, which counts for the home team
  event({ id: 'e1', type: 'goal', side: 'away', player: { id: '72', name: '' }, flowImpact: 23 }),
]

describe('lineups', () => {
  it('puts the home team on its half and marks who scored and who was booked', () => {
    expect(pitchPlayers(city, 'home', events)).toEqual([
      {
        id: '10',
        number: '10',
        name: 'Agüero',
        left: '58%',
        top: '33%',
        downLeft: '67%',
        downTop: '58%',
        goal: true,
        card: null,
        off: null,
      },
      {
        id: '72',
        number: '72',
        name: 'Iheanacho',
        left: '58%',
        top: '67%',
        downLeft: '33%',
        downTop: '58%',
        goal: false,
        card: 'yellow',
        off: '70′',
      },
    ])
  })

  it('turns the away team round to its end of the pitch', () => {
    expect(pitchPlayers(city, 'away', [])[0]).toMatchObject({ left: '42%', top: '67%' })
  })

  it('stands the pitch upright on a phone: the home team at the top, its left wing on the right', () => {
    // Agüero, on the left of the home attack: down the pitch at the top, right of centre
    expect(pitchPlayers(city, 'home', [])[0]).toMatchObject({ downLeft: '67%', downTop: '58%' })
    // The same place for the away team: at the bottom, left of centre
    expect(pitchPlayers(city, 'away', [])[0]).toMatchObject({ downLeft: '33%', downTop: '42%' })
  })

  it("lists what each player did and their part in the team's Flow", () => {
    expect(listRows(city, 'home', events)).toEqual([
      { id: '10', number: '10', name: 'Sergio Agüero', events: 'Гол', flow: '+29' },
      // The yellow card helped the opponent; the substitution is the coach's
      { id: '72', number: '72', name: 'Kelechi Iheanacho', events: 'ЖК ↓ 70′', flow: '−4' },
    ])
  })

  it('shows on the bench who came on and when', () => {
    expect(benchRows(city)).toEqual([
      { id: '7', number: '7', name: 'Raheem Sterling', on: '↑ 70′' },
      { id: '13', number: '13', name: 'Willy Caballero', on: null },
    ])
  })

  it('writes minutes and formations the way the board does', () => {
    expect(minuteText({ minute: 90, addedTime: 3 })).toBe('90+3′')
    expect(minuteText(null)).toBeNull()
    expect(formationsLine('Сити', 'Арсенал', city, { ...city, formation: null })).toBe(
      'Сити 4-4-2 · Арсенал',
    )
  })
})
