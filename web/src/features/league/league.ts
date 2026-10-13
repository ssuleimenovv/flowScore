import type { FlowRecord, League, RoundMatch, Standing } from '@/shared/api/types'

// "Сезон 2026/27 · 7-й тур", the line over the table
export function seasonLine(l: Pick<League, 'season' | 'round'>): string {
  return `Сезон ${l.season} · ${l.round}-й тур`
}

// The mark beside a place. The board marks the places that lead to the
// Champions League and to the Europa League; every other place is plain.
export type ZoneMark = 'champions' | 'europa' | 'none'

export function zoneMark(zone: Standing['zone']): ZoneMark {
  return zone === 'champions' || zone === 'europa' ? zone : 'none'
}

// "17:5"
export function goals(s: Pick<Standing, 'goalsFor' | 'goalsAgainst'>): string {
  return `${s.goalsFor}:${s.goalsAgainst}`
}

// A team's average Flow as the table shows it, a whole number; "—" before it has played
export function flowText(flow: number | null): string {
  return flow === null ? '—' : String(Math.round(flow))
}

// The width of the Flow bar, on the same 0–100 scale
export function flowWidth(flow: number | null): string {
  return `${Math.round(Math.min(100, Math.max(0, flow ?? 0)))}%`
}

// "Тур 6 · Тур 7 · Тур 8": the round being played between the one before and
// the one after, three in all, kept inside the season
export function roundTabs(round: number, rounds: number): number[] {
  if (rounds < 1) return []
  const first = Math.max(1, Math.min(round - 1, rounds - 2))
  const tabs: number[] = []
  for (let r = first; r <= Math.min(rounds, first + 2); r++) tabs.push(r)
  return tabs
}

// "Фулхэм — Манчестер Сити"
export function pairing(m: Pick<RoundMatch, 'home' | 'away'>): string {
  return `${m.home.name} — ${m.away.name}`
}

// The score under a match: "1:3"; "—" before kick-off
export function scoreText(m: Pick<RoundMatch, 'score'>): string {
  return m.score ? `${m.score.home}:${m.score.away}` : '—'
}

function dayKey(date: Date, timeZone?: string): string {
  return new Intl.DateTimeFormat('en-CA', { timeZone }).format(date)
}

// Where a match of the round stands, as the board writes it: "ФТ" once it is
// over, the minute while it is played ("72′"), and before kick-off the time
// today, "Завтра", or the date ("28 окт")
export function roundStatus(m: RoundMatch, now: Date, timeZone?: string): string {
  if (m.state === 'post') return 'ФТ'
  if (m.state === 'in') return m.clock.replace(/'/g, '′')
  const kickoff = new Date(m.kickoffAt)
  const key = dayKey(kickoff, timeZone)
  if (key === dayKey(now, timeZone)) {
    return new Intl.DateTimeFormat('ru', { hour: '2-digit', minute: '2-digit', timeZone }).format(
      kickoff,
    )
  }
  if (key === dayKey(new Date(now.getTime() + 24 * 3600 * 1000), timeZone)) return 'Завтра'
  return new Intl.DateTimeFormat('ru', { day: 'numeric', month: 'short', timeZone })
    .format(kickoff)
    .replace('.', '')
}

// A Flow record as the board lists it: "01", the match, "тур 5", the value
export interface RecordRow {
  id: string
  rank: string
  match: string
  round: string
  flow: string
}

export function recordRows(records: FlowRecord[]): RecordRow[] {
  return records.map((r, i) => ({
    id: r.matchId,
    rank: String(i + 1).padStart(2, '0'),
    match: pairing(r),
    round: `тур ${r.round}`,
    flow: String(Math.round(r.flow)),
  }))
}
