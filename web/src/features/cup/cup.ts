import { roundStatus } from '@/features/league/league'
import type { Stage, Tie, TieSide } from '@/shared/api/types'

// The stage the card opens on: the one being played, else the last one
export function openingStage(stages: Stage[], current: string): string {
  return stages.some((s) => s.id === current) ? current : (stages[stages.length - 1]?.id ?? '')
}

// A stage played over two legs, as the European cups' are
export function isTwoLegged(ties: Tie[]): boolean {
  return ties.some((t) => t.legs.length > 1)
}

// "1-й · 2-й · Σ" over a two-legged tie, "Счёт" over a single match
export function tieColumns(twoLegs: boolean): string[] {
  return twoLegs ? ['1-й', '2-й', 'Σ'] : ['Счёт']
}

// The European card's side note: "два матча"
export function legsNote(twoLegs: boolean): string {
  return twoLegs ? 'два матча' : 'один матч'
}

// The cup page's side note: "32 матча". Before the ties come, a match a tie
export function matchesNote(stage: Pick<Stage, 'ties'>, ties: Tie[]): string {
  const n = ties.length ? ties.reduce((sum, t) => sum + t.legs.length, 0) : stage.ties
  return `${n} ${plural(n, 'матч', 'матча', 'матчей')}`
}

// The line under the ties
export function tieFootnote(twoLegs: boolean): string {
  return twoLegs
    ? 'Зелёная метка — прошёл дальше. Сумма двух матчей справа.'
    : 'Один матч. При ничьей — дополнительное время и пенальти.'
}

// The mark beside a team: green once it is through, faint once it is out,
// the bare track while the tie is open
export type TieMark = 'through' | 'out' | 'open'

export function tieMark(side: Pick<TieSide, 'through'>): TieMark {
  if (side.through === null) return 'open'
  return side.through ? 'through' : 'out'
}

// A team's cells on the right: "2 · 1 · 3", its goals in each leg and the sum,
// or its score in a single match; "–" before a leg kicks off. The last is bold.
export function tieCells(side: Pick<TieSide, 'goals' | 'total'>, twoLegs: boolean): string[] {
  const goal = (n: number | null | undefined) => (n === null || n === undefined ? '–' : String(n))
  return twoLegs
    ? [goal(side.goals[0]), goal(side.goals[1]), goal(side.total)]
    : [goal(side.goals[0])]
}

// The minute in the pill: "72′", and past the 90 "ДВ 104′"
export function liveClock(clock: string): string {
  const minute = Number.parseInt(clock, 10)
  const text = clock.replace(/'/g, '′')
  return minute > 90 ? `ДВ ${text}` : text
}

// What the board writes under a tie: the pill while a leg is played, how the
// tie was settled ("после ДВ · пен. 4:3", the home team of the first leg
// first), or when the next leg is ("ответный 11 мар").
//
// The API says only that the deciding match went to penalties, not whether
// extra time came first. A two-legged tie always has it (UEFA's rule), so it
// reads "после ДВ · пен. 4:3"; a single match may go straight to penalties
// (the League Cup, Coppa Italia's early rounds), so it reads "пен. 3:0",
// which is never wrong.
export type TieNote = { live: boolean; text: string }

export function tieNote(tie: Tie, now: Date, timeZone?: string): TieNote | null {
  const playing = tie.legs.find((l) => l.state === 'in')
  if (playing) return { live: true, text: liveClock(playing.clock) }

  const twoLegs = tie.legs.length > 1
  if (tie.finish === 'aet') return { live: false, text: 'после ДВ' }
  if (tie.finish === 'pens') {
    const { shootout: home } = tie.home
    const { shootout: away } = tie.away
    const pens = home === null || away === null ? 'пен.' : `пен. ${home}:${away}`
    return { live: false, text: twoLegs ? `после ДВ · ${pens}` : pens }
  }

  const next = tie.legs.findIndex((l) => l.state === 'pre')
  const leg = tie.legs[next]
  if (!leg) return null
  const when = roundStatus(leg, now, timeZone)
  if (!twoLegs) return { live: false, text: when }
  return { live: false, text: `${next === 0 ? 'первый' : 'ответный'} ${lowerFirst(when)}` }
}

// Russian plural: 1 матч, 2 матча, 5 матчей, 21 матч
function plural(n: number, one: string, few: string, many: string): string {
  const tens = n % 100
  if (tens >= 11 && tens <= 14) return many
  if (n % 10 === 1) return one
  if (n % 10 >= 2 && n % 10 <= 4) return few
  return many
}

function lowerFirst(s: string): string {
  return s.charAt(0).toLowerCase() + s.slice(1)
}
