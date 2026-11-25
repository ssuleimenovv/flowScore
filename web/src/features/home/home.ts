import { competitionGroup } from '@/features/league/league'
import type { MatchSummary, Probabilities, Side } from '@/shared/api/types'
import { explain } from '../match/explain'
import { formatMinute } from '../match/format'

// The home screen's text and shapes. Pure functions, so they are easy to test.

// One competition's matches of a day: a card of its own (HomeByCompetitionExtended)
export interface DayCompetition {
  id: string // "uefa.champions"
  name: string // "Лига чемпионов"
  matches: MatchSummary[] // by kick-off, then by the home team
}

// Real matches not started yet, under the day they are played on
export interface Day {
  key: string // "2026-10-10", the viewer's own calendar day
  title: string // "Сегодня", "Завтра", "Суббота, 10 октября"
  matches: MatchSummary[] // by kick-off, then by the home team
  competitions: DayCompetition[] // in the chips' order
}

export interface MatchDay {
  live: MatchSummary[] // real matches live or at the break, by kick-off
  days: Day[] // real matches not started yet, soonest first
  demo: { live: MatchSummary[]; upcoming: MatchSummary[] } // the replays of old matches
}

// The matches of one competition, by its id, or of all when competition is null. Real
// matches come first, by day; the demo replays go apart, so a replay of 2016
// is never taken for tonight's match. A finished match is not on the home
// screen: the mockup shows only what is on and what is next.
export function matchDay(
  items: MatchSummary[],
  competition: string | null,
  now = new Date(),
  timeZone?: string,
): MatchDay {
  const shown = items.filter((m) => competition === null || m.competition.id === competition)
  const playing = (m: MatchSummary) => m.status === 'live' || m.status === 'halftime'
  const real = shown.filter((m) => m.source !== 'replay')
  const demo = shown.filter((m) => m.source === 'replay')

  const days: Day[] = []
  for (const m of real.filter((m) => m.status === 'scheduled').sort(byKickoff)) {
    const key = dayKey(new Date(m.kickoffAt), timeZone)
    let day = days.find((d) => d.key === key)
    if (!day) {
      day = {
        key,
        title: dayName(new Date(m.kickoffAt), now, timeZone),
        matches: [],
        competitions: [],
      }
      days.push(day)
    }
    day.matches.push(m)
    let comp = day.competitions.find((c) => c.id === m.competition.id)
    if (!comp) {
      comp = { id: m.competition.id, name: m.competition.name, matches: [] }
      day.competitions.push(comp)
    }
    comp.matches.push(m)
  }
  for (const d of days)
    d.competitions.sort((a, b) => rank(a.id) - rank(b.id) || a.name.localeCompare(b.name))

  return {
    live: real.filter(playing),
    days,
    demo: { live: demo.filter(playing), upcoming: demo.filter((m) => m.status === 'scheduled') },
  }
}

// The same order at every poll: ESPN lists the matches of one kick-off in no
// fixed order, so they would change places
function byKickoff(a: MatchSummary, b: MatchSummary): number {
  return (
    Date.parse(a.kickoffAt) - Date.parse(b.kickoffAt) ||
    a.home.name.localeCompare(b.home.name) ||
    a.id.localeCompare(b.id)
  )
}

// "2026-10-10" in the viewer's time zone: a match at 00:30 is on the next day
function dayKey(date: Date, timeZone?: string): string {
  return new Intl.DateTimeFormat('en-CA', { timeZone }).format(date)
}

// "Сегодня", "Завтра", or the date for a day further on
function dayName(date: Date, now: Date, timeZone?: string): string {
  const key = dayKey(date, timeZone)
  if (key === dayKey(now, timeZone)) return 'Сегодня'
  if (key === dayKey(new Date(now.getTime() + 24 * 3600 * 1000), timeZone)) return 'Завтра'
  return dayTitle(date, timeZone)
}

// The order of the competitions on the screens: the leagues, the European
// cups, the national cups, as the gateway follows them
const ORDER = [
  'eng.1',
  'esp.1',
  'ger.1',
  'ita.1',
  'fra.1',
  'uefa.champions',
  'uefa.europa',
  'uefa.europa.conf',
  'uefa.super_cup',
  'eng.fa',
  'eng.league_cup',
  'ger.dfb_pokal',
  'esp.copa_del_rey',
  'ita.coppa_italia',
  'fra.coupe_de_france',
]

function rank(id: string): number {
  return ORDER.includes(id) ? ORDER.indexOf(id) : ORDER.length
}

export interface Competition {
  id: string // "esp.1"
  name: string // "Ла Лига"
}

// The home screen's picker, as on the league page with «Все» before the groups
export type Group = 'Все' | 'Лиги' | 'Еврокубки' | 'Кубки'

// The real competitions the list has, for the chips, in the screens' order.
// The demo has a block of its own and no chip.
export function competitions(items: MatchSummary[]): Competition[] {
  const seen = new Map<string, Competition>()
  for (const m of items) {
    if (m.source === 'replay') continue
    seen.set(m.competition.id, { id: m.competition.id, name: m.competition.name })
  }
  return [...seen.values()].sort((a, b) => rank(a.id) - rank(b.id) || a.name.localeCompare(b.name))
}

// The groups that have a competition today, «Все» first
export function groups(list: Competition[]): Group[] {
  const present = new Set(list.map((c) => competitionGroup(c.id)))
  return ['Все', ...(['Лиги', 'Еврокубки', 'Кубки'] as const).filter((g) => present.has(g))]
}

// The matches a pick shows: a group's real ones; all of them, the demo too, for «Все»
export function inGroup(items: MatchSummary[], group: Group): MatchSummary[] {
  if (group === 'Все') return items
  return items.filter((m) => m.source !== 'replay' && competitionGroup(m.competition.id) === group)
}

// The match clock ticks on between polls, from the moment it was measured
export function liveMinute(m: MatchSummary, now: number): string {
  const { elapsedSeconds, observedAt, period } = m.clock
  const running = m.status === 'live' ? Math.max(0, (now - Date.parse(observedAt)) / 1000) : 0
  return formatMinute(elapsedSeconds + running, period)
}

// The card's one line of "AI": the agent's headline once it has written one,
// until then the headline the match screen builds from Flow
export function insightLine(m: MatchSummary): string {
  if (m.explanation) return m.explanation.title
  return explain({ match: m, flow: m.flow, delta10: m.delta10, factors: m.factors }, 0).title
}

// Each side's share of the flow bar in percent; even before anything happened
export function flowShare(m: MatchSummary): number {
  const total = m.flow.home + m.flow.away
  return total > 0 ? (m.flow.home / total) * 100 : 50
}

// The mini wave: one line per side over the match, in a 300 × 60 box, the
// way the card draws it. Time runs to the 95th minute, Flow from 0 to 100.
export function sparkline(m: MatchSummary, side: Side): string {
  return m.points
    .map((p) => `${((p.minute / 95) * 300).toFixed(1)},${(58 - p[side] * 0.52).toFixed(1)}`)
    .join(' ')
}

export interface Peak {
  id: string
  summary: MatchSummary // the teams and their crests
  side: Side // the team that rose
  rise: number // Flow points gained over the window
  flow: number // its Flow now
}

// The live matches where a team's Flow rose the most over the last minutes:
// "Пики потока за 5 минут". A match where nobody rose is not a peak.
export function flowPeaks(live: MatchSummary[], minutes = 5, limit = 3): Peak[] {
  const peaks: Peak[] = []
  for (const m of live) {
    const last = m.points[m.points.length - 1]
    if (!last) continue
    // The latest point at least `minutes` before the last one
    const earlier = m.points.filter((p) => p.minute <= last.minute - minutes)
    const before = earlier[earlier.length - 1]
    if (!before) continue
    const rises = { home: last.home - before.home, away: last.away - before.away }
    const side: Side = rises.home >= rises.away ? 'home' : 'away'
    if (rises[side] <= 0) continue
    peaks.push({
      id: m.id,
      summary: m,
      side,
      rise: Math.round(rises[side]),
      flow: Math.round(last[side]),
    })
  }
  return peaks.sort((a, b) => b.rise - a.rise).slice(0, limit)
}

// "19:40", in the viewer's time zone
export function kickoffTime(m: MatchSummary, timeZone?: string): string {
  return new Intl.DateTimeFormat('ru', { hour: '2-digit', minute: '2-digit', timeZone }).format(
    new Date(m.kickoffAt),
  )
}

// "Понедельник, 5 октября"
export function dayTitle(date: Date, timeZone?: string): string {
  const text = new Intl.DateTimeFormat('ru', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    timeZone,
  }).format(date)
  return text.charAt(0).toUpperCase() + text.slice(1)
}

// The day's forecast: the next match the model is surest about, with a line
// that says where the edge comes from. Only what the model uses is named.
export function dayForecast(
  upcoming: MatchSummary[],
): { match: MatchSummary; text: string } | null {
  const sure = (p: Probabilities) => Math.max(p.home, p.away)
  const withChances = upcoming.filter((m) => m.prediction)
  const match = withChances.sort(
    (a, b) => sure(b.prediction!.current) - sure(a.prediction!.current),
  )[0]
  if (!match?.prediction) return null

  const { home, away } = match.prediction.current
  const favourite = home >= away ? match.home.name : match.away.name
  const edge = Math.abs(home - away)
  const text =
    edge < 10
      ? `Равный матч: модель не видит явного фаворита, шансы команд различаются на ${edge} п.`
      : `${favourite} — фаворит: модель учитывает силу команд по xG прошлых матчей и преимущество своего поля.`
  return { match, text }
}
