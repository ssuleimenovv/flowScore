import type { MatchSummary, Probabilities, Side } from '@/shared/api/types'
import { explain } from '../match/explain'
import { formatMinute } from '../match/format'

// The home screen's text and shapes. Pure functions, so they are easy to test.

export interface MatchDay {
  live: MatchSummary[] // live or at the break, by kick-off
  upcoming: MatchSummary[] // not started yet, soonest first
}

// The matches of one competition, or of all when competition is null. A
// finished match is not on the home screen: the mockup shows only what is on
// and what is next.
export function matchDay(items: MatchSummary[], competition: string | null): MatchDay {
  const shown = items.filter((m) => competition === null || m.competition.name === competition)
  return {
    live: shown.filter((m) => m.status === 'live' || m.status === 'halftime'),
    upcoming: shown.filter((m) => m.status === 'scheduled'),
  }
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
export function kickoffTime(m: MatchSummary): string {
  return new Intl.DateTimeFormat('ru', { hour: '2-digit', minute: '2-digit' }).format(
    new Date(m.kickoffAt),
  )
}

// "Понедельник, 5 октября"
export function dayTitle(date: Date): string {
  const text = new Intl.DateTimeFormat('ru', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
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
