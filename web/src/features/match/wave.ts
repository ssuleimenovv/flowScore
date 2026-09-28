import type { FlowPoint, MatchEvent, Side } from '@/shared/api/types'

// The wave is drawn for at least a full match, so the unplayed part shows too.
export const MATCH_MINUTES = 90

// One value per minute: home Flow minus away Flow (−100…100), averaged over
// the last `window` minutes. null marks a minute that is not played yet.
// `points` must be sorted by minute, as MatchLive keeps them.
export function buildWave(points: FlowPoint[], window: number): Array<number | null> {
  const raw = new Map(points.map((p) => [p.minute, p.home - p.away]))
  const last = points[points.length - 1]?.minute ?? 0
  const total = Math.max(MATCH_MINUTES, last)

  const wave: Array<number | null> = []
  for (let minute = 1; minute <= total; minute++) {
    if (minute > last) {
      wave.push(null)
      continue
    }
    let sum = 0
    let count = 0
    for (let m = Math.max(1, minute - window + 1); m <= minute; m++) {
      const value = raw.get(m)
      if (value !== undefined) {
        sum += value
        count++
      }
    }
    wave.push(count > 0 ? sum / count : 0)
  }
  return wave
}

export interface GoalMark {
  column: number // the wave's minute column the goal falls in, from 1
  side: Side
  label: string // "0:1 Сака"
}

// Goals in match order, each labelled with the score it made.
export function goalMarks(events: MatchEvent[]): GoalMark[] {
  const goals = events
    .filter((e): e is MatchEvent & { side: Side } => e.type === 'goal' && e.side !== null)
    .sort((a, b) => a.minute - b.minute || (a.addedTime ?? 0) - (b.addedTime ?? 0))

  let home = 0
  let away = 0
  return goals.map((goal) => {
    if (goal.side === 'home') home++
    else away++
    const scorer = goal.player ? ` ${surname(goal.player.name)}` : ''
    return { column: column(goal), side: goal.side, label: `${home}:${away}${scorer}` }
  })
}

// A goal at 7:30 is 7′, but the wave keeps it in column 8, the one that ends
// at 8:00. In added time 45+2 is the 47th minute played: column 47.
function column(goal: MatchEvent): number {
  return goal.minute + (goal.addedTime ?? 1)
}

// Words that belong to the surname in front of it: "De Bruyne", "van Dijk"
const PARTICLES = new Set([
  'de',
  'da',
  'di',
  'del',
  'della',
  'dos',
  'du',
  'van',
  'von',
  'der',
  'den',
  'ter',
  'le',
  'la',
])

// "Sergio Agüero" → "Agüero", "Kevin De Bruyne" → "De Bruyne". The backend
// sends the name fans know, so the last word is the surname.
export function surname(name: string): string {
  const words = name.trim().split(/\s+/)
  let start = words.length - 1
  while (start > 1 && PARTICLES.has(words[start - 1]!.toLowerCase())) start--
  return words.slice(start).join(' ') || name
}
