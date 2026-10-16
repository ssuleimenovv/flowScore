import type { LineupMinute, LineupPlayer, MatchEvent, Side, TeamLineup } from '@/shared/api/types'
import { surname } from './wave'

// "Составы" from the MatchExtended boards: the eleven on the pitch, the same
// as a list with what each player did, and the bench.

export interface PitchPlayer {
  id: string
  number: string // "47", "" if the source does not say
  name: string // "Foden"
  // Across the pitch, on a desktop: the home team attacks to the right
  left: string // "46%"
  top: string
  // Down the pitch, on a phone: the home team at the top attacks downwards
  downLeft: string
  downTop: string
  goal: boolean
  card: 'yellow' | 'red' | null
  off: string | null // "64′", when the player went off
}

export interface ListRow {
  id: string
  number: string
  name: string // "Phil Foden"
  events: string // "Гол · ЖК ↓ 64′"
  flow: string // "+9", "−3", "0": the player's share of the team's Flow
}

export interface BenchRow {
  id: string
  number: string
  name: string
  on: string | null // "↑ 64′"
}

// "64′", "90+3′"
export function minuteText(m: LineupMinute): string | null {
  if (!m) return null
  return m.addedTime ? `${m.minute}+${m.addedTime}′` : `${m.minute}′`
}

// What each player of a team did, from the match events: the source gives
// the event's player by the same id as the lineup
interface Deeds {
  goals: number
  yellow: boolean
  red: boolean
  flow: number
}

function deeds(events: MatchEvent[], side: Side): Map<string, Deeds> {
  const out = new Map<string, Deeds>()
  for (const e of events) {
    // An own goal is the other team's, and its player is not one of theirs
    if (e.side !== side || !e.player?.id) continue
    const d = out.get(e.player.id) ?? { goals: 0, yellow: false, red: false, flow: 0 }
    if (e.type === 'goal') d.goals++
    if (e.type === 'yellow_card') d.yellow = true
    if (e.type === 'red_card') d.red = true
    d.flow += contribution(e)
    out.set(e.player.id, d)
  }
  return out
}

// A player's part in the team's Flow (docs/FLOW.md, section 4): what each of
// the player's moments added, a yellow card taken off since it helps the
// opponent. A substitution is the coach's, not the player's.
function contribution(e: MatchEvent): number {
  const w = e.flowImpact ?? 0
  if (e.type === 'substitution') return 0
  return e.type === 'yellow_card' ? -w : w
}

function flowText(n: number): string {
  const r = Math.round(n)
  if (r > 0) return `+${r}`
  if (r < 0) return `−${Math.abs(r)}`
  return '0'
}

// The eleven where they stand, each team in its own half. The server gives
// every place attacking to the right; the away team is turned half a circle
// round to its end. Down the pitch everything turns a quarter: a team facing
// down has its left wing on the right of the screen.
export function pitchPlayers(team: TeamLineup, side: Side, events: MatchEvent[]): PitchPlayer[] {
  const did = deeds(events, side)
  return team.starters
    .filter((p): p is LineupPlayer & { position: { x: number; y: number } } => !!p.position)
    .map((p) => {
      const d = did.get(p.id)
      const home = side === 'home'
      return {
        id: p.id,
        number: p.number ?? '',
        name: surname(p.name),
        left: `${home ? p.position.x : 100 - p.position.x}%`,
        top: `${home ? p.position.y : 100 - p.position.y}%`,
        downLeft: `${home ? 100 - p.position.y : p.position.y}%`,
        downTop: `${home ? p.position.x : 100 - p.position.x}%`,
        goal: (d?.goals ?? 0) > 0,
        card: d?.red ? 'red' : d?.yellow ? 'yellow' : null,
        off: minuteText(p.off),
      }
    })
}

// The eleven as a list: what each one did and their share of the Flow
export function listRows(team: TeamLineup, side: Side, events: MatchEvent[]): ListRow[] {
  const did = deeds(events, side)
  return team.starters.map((p) => {
    const d = did.get(p.id)
    const done: string[] = []
    if (d?.goals) done.push(d.goals === 1 ? 'Гол' : `${d.goals} гола`)
    if (d?.red) done.push('КК')
    else if (d?.yellow) done.push('ЖК')
    const off = minuteText(p.off)
    return {
      id: p.id,
      number: p.number ?? '',
      name: p.name,
      events: [done.join(' · '), off && `↓ ${off}`].filter(Boolean).join(' '),
      flow: flowText(d?.flow ?? 0),
    }
  })
}

export function benchRows(team: TeamLineup): BenchRow[] {
  return team.bench.map((p) => {
    const on = minuteText(p.on)
    return { id: p.id, number: p.number ?? '', name: p.name, on: on && `↑ ${on}` }
  })
}

// "Manchester City 4-4-2 · Arsenal 4-2-3-1": a team the source gives no
// formation for is named alone
export function formationsLine(home: string, away: string, h: TeamLineup, a: TeamLineup): string {
  const team = (name: string, t: TeamLineup) => (t.formation ? `${name} ${t.formation}` : name)
  return `${team(home, h)} · ${team(away, a)}`
}
