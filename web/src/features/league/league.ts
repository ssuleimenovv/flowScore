import type {
  FlowRecord,
  League,
  LeagueItem,
  RoundMatch,
  Stage,
  Standing,
} from '@/shared/api/types'

// "1/8" is written "1/8 финала"; "Финал", "Стыки" and "3-й раунд" stay as they are
export function stageTitle(name: string): string {
  return /^1\/\d+$/.test(name) ? `${name} финала` : name
}

// "Сезон 2026/27 · 7-й тур", the line over the table. A cup in its knockouts
// names the stage ("1/8 финала"), the Super Cup its one match, and a cup
// before its first stage only the season.
export function seasonLine(
  l: Pick<League, 'season' | 'round' | 'rounds' | 'stage' | 'stages'>,
): string {
  const stage = l.stages.find((s) => s.id === l.stage)
  if (stage) return `Сезон ${l.season} · ${stageTitle(stage.name)}`
  if (l.rounds === 1) return `Сезон ${l.season} · Матч за трофей`
  if (l.round > 0) return `Сезон ${l.season} · ${l.round}-й тур`
  return `Сезон ${l.season}`
}

// The mark beside a place. The League board marks the places that lead to
// the Champions League and to the Europa League; the cup board marks the
// round of 16 the same as the first and the playoffs the same as the second.
// Every other place is plain.
export type ZoneMark = 'champions' | 'europa' | 'none'

export function zoneMark(zone: Standing['zone']): ZoneMark {
  if (zone === 'champions' || zone === 'last16') return 'champions'
  if (zone === 'europa' || zone === 'playoff') return 'europa'
  return 'none'
}

// The Champions, Europa and Conference League: a league phase of 36, not a league
export function isEuroCup(leagueId: string | null): boolean {
  return leagueId?.startsWith('uefa.') ?? false
}

// What a competition's page shows. A league ("eng.1", a country and a tier)
// and a European cup have a table; a cup has knockout stages and no table.
// The Super Cup is one match: a cup, though it sits among the European cups.
export type CompetitionKind = 'league' | 'euro' | 'cup'

export function competitionKind(leagueId: string | null): CompetitionKind {
  if (leagueId === null || /^[a-z]+\.\d+$/.test(leagueId)) return 'league'
  return isEuroCup(leagueId) && leagueId !== 'uefa.super_cup' ? 'euro' : 'cup'
}

// The chips in groups, as on the CompetitionsExtended board: the leagues, the
// European cups and the national cups, each in the gateway's order; a group
// with nothing in it is left out
const GROUPS = ['Лиги', 'Еврокубки', 'Кубки'] as const

export function competitionGroup(leagueId: string): (typeof GROUPS)[number] {
  if (isEuroCup(leagueId)) return 'Еврокубки'
  return competitionKind(leagueId) === 'league' ? 'Лиги' : 'Кубки'
}

export function competitionGroups(items: LeagueItem[]): { name: string; items: LeagueItem[] }[] {
  return GROUPS.map((name) => ({
    name,
    items: items.filter((l) => competitionGroup(l.id) === name),
  })).filter((g) => g.items.length > 0)
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
// today, "Завтра", or the date ("28 окт"). A leg of a cup tie is written the same way.
export function roundStatus(
  m: Pick<RoundMatch, 'state' | 'clock' | 'kickoffAt'>,
  now: Date,
  timeZone?: string,
): string {
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

// A Flow record as the board lists it: "01", the match, "тур 5", the value.
// A match of a cup's knockouts names its stage instead: "3-й раунд".
export interface RecordRow {
  id: string
  rank: string
  match: string
  round: string
  flow: string
}

export function recordRows(records: FlowRecord[], stages: Stage[] = []): RecordRow[] {
  return records.map((r, i) => {
    const stage = stages.find((s) => s.id === r.stage)
    return {
      id: r.matchId,
      rank: String(i + 1).padStart(2, '0'),
      match: pairing(r),
      round: stage ? stageTitle(stage.name) : r.round > 0 ? `тур ${r.round}` : '',
      flow: String(Math.round(r.flow)),
    }
  })
}
