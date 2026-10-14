import type { MatchEvent, Side } from '@/shared/api/types'
import { surname } from './wave'

// "Карта ударов" from the MatchPhases boards: every shot where it was taken,
// the home team attacking to the right and the away team to the left.

export type ShotOutcome = 'goal' | 'on' | 'block' | 'off' // a blocked shot is drawn as one on target

export interface Shot {
  id: string
  side: Side
  outcome: ShotOutcome
  left: string // "88%" from the left edge of the pitch
  top: string // "46%" from the top
  size: string // "30px": the bigger the circle, the bigger the chance
  label: string // "Foden 31′ · xG 0.11"
}

export interface ShotTotals {
  home: { xg: string; shots: number } // "1.84", 8
  away: { xg: string; shots: number }
}

const OUTCOMES: Partial<Record<MatchEvent['type'], ShotOutcome>> = {
  goal: 'goal',
  shot_on_target: 'on',
  shot_blocked: 'block',
  shot_off_target: 'off',
}

// Our model's xG when there is one (docs/XG.md), else the source's. An own
// goal has neither, so it is not a shot of the team it counts for.
function xgOf(e: MatchEvent): number | null {
  if (typeof e.modelXG === 'number') return e.modelXG
  if (typeof e.xG === 'number') return e.xG
  return null
}

function isShot(e: MatchEvent): e is MatchEvent & { side: Side } {
  return e.type in OUTCOMES && e.side !== null && xgOf(e) !== null
}

// The circles, oldest first so a later shot is drawn over an earlier one.
// The source gives every position attacking to the right, so the away team's
// shots are turned half a circle round to its end of the pitch.
export function shots(events: MatchEvent[]): Shot[] {
  const out: Shot[] = []
  for (let i = events.length - 1; i >= 0; i--) {
    const e = events[i]!
    const pos = e.position
    if (!isShot(e) || pos?.x === undefined || pos.y === undefined) continue
    const home = e.side === 'home'
    const xg = xgOf(e)!
    const minute = e.addedTime ? `${e.minute}+${e.addedTime}′` : `${e.minute}′`
    const who = e.player ? `${surname(e.player.name)} ` : ''
    out.push({
      id: e.id,
      side: e.side,
      outcome: OUTCOMES[e.type]!,
      left: `${home ? pos.x : 100 - pos.x}%`,
      top: `${home ? pos.y : 100 - pos.y}%`,
      size: `${Math.round(12 + xg * 44)}px`,
      label: `${who}${minute} · xG ${xg.toFixed(2)}`,
    })
  }
  return out
}

// "xG 1.84 · 8 ударов": every shot counts here, a shot the source gave no
// position too
export function shotTotals(events: MatchEvent[]): ShotTotals {
  const totals = { home: { xg: 0, shots: 0 }, away: { xg: 0, shots: 0 } }
  for (const e of events) {
    if (!isShot(e)) continue
    totals[e.side].xg += xgOf(e)!
    totals[e.side].shots++
  }
  return {
    home: { xg: totals.home.xg.toFixed(2), shots: totals.home.shots },
    away: { xg: totals.away.xg.toFixed(2), shots: totals.away.shots },
  }
}

// "8 ударов", "1 удар", "3 удара"
export function shotsWord(n: number): string {
  const ten = n % 10
  const hundred = n % 100
  if (ten === 1 && hundred !== 11) return `${n} удар`
  if (ten >= 2 && ten <= 4 && (hundred < 12 || hundred > 14)) return `${n} удара`
  return `${n} ударов`
}
