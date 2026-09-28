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
  minute: number
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
    return { minute: goal.minute, side: goal.side, label: `${home}:${away}${scorer}` }
  })
}

// "Alexis Alejandro Sánchez Sánchez" → "Sánchez". A guess: the data has full
// names only, until the backend sends short ones.
export function surname(name: string): string {
  const words = name.trim().split(/\s+/)
  return words[words.length - 1] ?? name
}
