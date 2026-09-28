import type { FlowValues, Match } from '@/shared/api/types'

// Text for the match screen. Pure functions, so they are easy to test.

const HALF = 45 * 60

// "72′", and "45+2′" / "90+3′" in added time, as a broadcast shows it.
export function formatMinute(seconds: number, period: Match['clock']['period']): string {
  const minute = Math.floor(seconds / 60)
  const end = period === 'first_half' ? HALF : period === 'second_half' ? 2 * HALF : null
  if (end !== null && seconds >= end) {
    return `${end / 60}+${Math.floor((seconds - end) / 60) + 1}′`
  }
  return `${minute}′`
}

// "72:14"
export function formatClock(seconds: number): string {
  const s = Math.floor(seconds)
  return `${Math.floor(s / 60)}:${String(s % 60).padStart(2, '0')}`
}

// "Premier League · 37-й тур · Etihad Stadium, Manchester"
export function formatMeta(match: Match): string {
  const { competition, venue } = match
  const place = [venue?.name, venue?.city].filter(Boolean).join(', ')
  return [competition.name, competition.round && `${competition.round}-й тур`, place]
    .filter(Boolean)
    .join(' · ')
}

// "▲ +18" or "▼ −21": the change of Flow over the last 10 minutes
export function formatDelta(delta: number): string {
  const n = Math.round(delta)
  return n >= 0 ? `▲ +${n}` : `▼ −${Math.abs(n)}`
}

// Each side's share of the flow in percent; 50/50 before anything happened
export function flowShare(flow: FlowValues): { home: number; away: number } {
  const total = flow.home + flow.away
  const home = total > 0 ? Math.round((flow.home / total) * 100) : 50
  return { home, away: 100 - home }
}
