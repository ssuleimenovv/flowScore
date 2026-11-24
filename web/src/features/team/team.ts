import type { Team, TeamLeague, TeamMatch } from '@/shared/api/types'

// The page is the team's, so each match is read from its side: its Flow,
// its opponent, its result. The matches come by kick-off, oldest first.

export function isHome(m: TeamMatch, teamId: string): boolean {
  return m.home.id === teamId
}

export function opponentOf(m: TeamMatch, teamId: string): Team {
  return isHome(m, teamId) ? m.away : m.home
}

// The team's own Flow in a match; null until the match is over and read
export function flowOf(m: TeamMatch, teamId: string): number | null {
  if (!m.flow) return null
  return isHome(m, teamId) ? m.flow.home : m.flow.away
}

export const RESULT_LETTER = { win: 'В', draw: 'Н', loss: 'П' } as const
export type Result = keyof typeof RESULT_LETTER

// For a screen reader: «В» alone says nothing
export const RESULT_WORD = { win: 'победа', draw: 'ничья', loss: 'поражение' } as const

// The country of a domestic league, as on the board: «Премьер-лига · Англия»
const COUNTRIES: Record<string, string> = {
  eng: 'Англия',
  esp: 'Испания',
  ger: 'Германия',
  ita: 'Италия',
  fra: 'Франция',
}

// «Премьер-лига · Англия · 3-е место»; the phone's nav bar has no room for the country
export function leagueLine(league: TeamLeague | null, withCountry = true): string {
  if (!league) return ''
  const country = withCountry ? COUNTRIES[league.id.split('.')[0]!] : undefined
  return [league.name, country, `${league.rank}-е место`].filter(Boolean).join(' · ')
}

function played(matches: TeamMatch[]): TeamMatch[] {
  return matches.filter((m) => m.state === 'post' && m.result !== '')
}

// The last results, oldest first, as the board's «В В Н В П»
export function form(matches: TeamMatch[], n = 5): Result[] {
  return played(matches)
    .slice(-n)
    .map((m) => m.result as Result)
}

// The team's Flow averaged over the matches read; null before any is
export function averageFlow(matches: TeamMatch[], teamId: string): number | null {
  const flows = matches.map((m) => flowOf(m, teamId)).filter((f) => f !== null)
  if (flows.length === 0) return null
  return flows.reduce((a, b) => a + b, 0) / flows.length
}

export interface FlowBar {
  matchId: string
  opponent: Team
  flow: number
}

// "Flow по матчам": the last matches read, oldest first
export function lastFlows(matches: TeamMatch[], teamId: string, n = 10): FlowBar[] {
  const bars: FlowBar[] = []
  for (const m of matches) {
    const flow = flowOf(m, teamId)
    if (flow !== null) bars.push({ matchId: m.id, opponent: opponentOf(m, teamId), flow })
  }
  return bars.slice(-n)
}

// "Ближайшие матчи": the one under way first, then the ones to come
export function upcoming(matches: TeamMatch[], n = Infinity): TeamMatch[] {
  return matches.filter((m) => m.state !== 'post').slice(0, n)
}

// "Результаты": the newest first
export function results(matches: TeamMatch[], n = Infinity): TeamMatch[] {
  return played(matches).reverse().slice(0, n)
}

// The team's goals first, as on the board: «3:1» a win, «1:2» a loss
export function scoreFor(m: TeamMatch, teamId: string): string {
  if (!m.score) return '—'
  const { home, away } = m.score
  return isHome(m, teamId) ? `${home}:${away}` : `${away}:${home}`
}

export function venueOf(m: TeamMatch, teamId: string): string {
  return isHome(m, teamId) ? 'Дома' : 'В гостях'
}

// The opponent; a match of another competition names it: «Inter · Лига чемпионов»
export function opponentLine(m: TeamMatch, teamId: string, leagueId: string | null): string {
  const name = opponentOf(m, teamId).name
  return m.leagueId === leagueId ? name : `${name} · ${m.league}`
}

function dayKey(date: Date, timeZone?: string): string {
  return new Intl.DateTimeFormat('en-CA', { timeZone }).format(date)
}

// «20 сент»
export function dayMonth(date: Date, timeZone?: string): string {
  return new Intl.DateTimeFormat('ru', { day: 'numeric', month: 'short', timeZone })
    .format(date)
    .replace('.', '')
}

// «Сегодня, 21:00» or «28 сент, 19:30»
export function kickoffLine(m: TeamMatch, now: Date, timeZone?: string): string {
  const at = new Date(m.kickoffAt)
  const time = new Intl.DateTimeFormat('ru', {
    hour: '2-digit',
    minute: '2-digit',
    timeZone,
  }).format(at)
  const day = dayKey(at, timeZone) === dayKey(now, timeZone) ? 'Сегодня' : dayMonth(at, timeZone)
  return `${day}, ${time}`
}

// ESPN names only the line of a squad's player
export const POSITION_WORD: Record<string, string> = {
  goalkeeper: 'Вратарь',
  defender: 'Защитник',
  midfielder: 'Полузащитник',
  forward: 'Нападающий',
}
