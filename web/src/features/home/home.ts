import type { MatchSummary, Probabilities, Side } from '@/shared/api/types'
import { explain } from '../match/explain'
import { formatMinute } from '../match/format'

// The home screen's text and shapes. Pure functions, so they are easy to test.

// Real matches not started yet, under the day they are played on
export interface Day {
  key: string // "2026-10-10", the viewer's own calendar day
  title: string // "Сегодня", "Завтра", "Суббота, 10 октября"
  matches: MatchSummary[]
}

export interface MatchDay {
  live: MatchSummary[] // real matches live or at the break, by kick-off
  days: Day[] // real matches not started yet, soonest first
  demo: { live: MatchSummary[]; upcoming: MatchSummary[] } // the replays of old matches
}

// The matches of one competition, or of all when competition is null. Real
// matches come first, by day; the demo replays go apart, so a replay of 2016
// is never taken for tonight's match. A finished match is not on the home
// screen: the mockup shows only what is on and what is next.
export function matchDay(
  items: MatchSummary[],
  competition: string | null,
  now = new Date(),
  timeZone?: string,
): MatchDay {
  const shown = items.filter((m) => competition === null || m.competition.name === competition)
  const playing = (m: MatchSummary) => m.status === 'live' || m.status === 'halftime'
  const real = shown.filter((m) => m.source !== 'replay')
  const demo = shown.filter((m) => m.source === 'replay')

  const days: Day[] = []
  for (const m of real.filter((m) => m.status === 'scheduled')) {
    const key = dayKey(new Date(m.kickoffAt), timeZone)
    let day = days.find((d) => d.key === key)
    if (!day) {
      day = { key, title: dayName(new Date(m.kickoffAt), now, timeZone), matches: [] }
      days.push(day)
    }
    day.matches.push(m)
  }

  return {
    live: real.filter(playing),
    days,
    demo: { live: demo.filter(playing), upcoming: demo.filter((m) => m.status === 'scheduled') },
  }
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

// The competitions the list has, for the chips: "Premier League"
export function competitions(items: MatchSummary[]): string[] {
  return [...new Set(items.map((m) => m.competition.name))].sort()
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
  match: string // "Manchester City — Arsenal"
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
      match: `${m.home.name} — ${m.away.name}`,
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
