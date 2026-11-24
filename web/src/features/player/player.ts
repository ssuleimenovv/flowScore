import { dayMonth, POSITION_WORD } from '@/features/team/team'
import type { PlayerInfo, PlayerMatch, Team } from '@/shared/api/types'

// The player's page reads their matches from their team's side, as the
// team's page does. Here: who the player is, in a line.

// «Нападающий · Norway» under the name; the desktop puts the club between them
export function aboutLine(p: PlayerInfo): string[] {
  return [POSITION_WORD[p.position], p.country].filter((s): s is string => Boolean(s))
}

// The phone's nav bar, as on MobilePlayer: «Manchester City · №9»
export function navLine(p: PlayerInfo): string {
  return [p.team?.name, p.number && `№${p.number}`].filter(Boolean).join(' · ')
}

export interface Totals {
  matches: number
  starts: number
  minutes: number
  goals: number
  assists: number
  shots: number
  onTarget: number
  xg: number
}

// The season so far, summed over the matches the player got on the pitch in
export function totals(matches: PlayerMatch[]): Totals {
  const t: Totals = {
    matches: 0,
    starts: 0,
    minutes: 0,
    goals: 0,
    assists: 0,
    shots: 0,
    onTarget: 0,
    xg: 0,
  }
  for (const m of matches) {
    t.matches++
    if (m.starter) t.starts++
    t.minutes += m.minutes
    t.goals += m.goals
    t.assists += m.assists
    t.shots += m.shots
    t.onTarget += m.onTarget
    t.xg += m.xg
  }
  return t
}

// «1.15»: how many per 90 minutes played; «—» before any
export function per90(n: number, minutes: number): string {
  return minutes > 0 ? ((n / minutes) * 90).toFixed(2) : '—'
}

const whole = new Intl.NumberFormat('ru')

export interface Tile {
  label: string
  value: string
  note: string
}

// The six tiles of PlayerSeason
export function seasonTiles(t: Totals): Tile[] {
  const perMatch = t.matches > 0 ? `${Math.round(t.minutes / t.matches)} за матч` : '—'
  return [
    { label: 'Матчи', value: String(t.matches), note: `${t.starts} в старте` },
    { label: 'Минуты', value: whole.format(t.minutes), note: perMatch },
    { label: 'Голы', value: String(t.goals), note: `${per90(t.goals, t.minutes)} на 90′` },
    { label: 'Ассисты', value: String(t.assists), note: `${per90(t.assists, t.minutes)} на 90′` },
    { label: 'xG · наша', value: t.xg.toFixed(1), note: `${per90(t.xg, t.minutes)} на 90′` },
    { label: 'Удары', value: String(t.shots), note: `${t.onTarget} в створ` },
  ]
}

export interface Finishing {
  goals: number
  xg: string
  diff: string // «+1.6», «−0.8»
  goalsWidth: string
  xgWidth: string
  text: string
}

// "Голы против xG": half a goal either way is within what luck gives
export function finishing(t: Totals): Finishing {
  const top = Math.max(t.goals, t.xg)
  const width = (n: number) => (top > 0 ? `${Math.round((n / top) * 100)}%` : '0%')
  const d = t.goals - t.xg
  let text = 'Почти ровно по xG: голы соответствуют качеству моментов.'
  if (t.shots === 0) text = 'Ударов в этом сезоне пока нет.'
  else if (d > 0.5) text = 'Забивает больше ожидаемого: реализует моменты лучше, чем обещает xG.'
  else if (d < -0.5) text = 'Забивает меньше ожидаемого: моменты есть, реализация отстаёт.'
  return {
    goals: t.goals,
    xg: t.xg.toFixed(1),
    diff: `${d < 0 ? '−' : '+'}${Math.abs(d).toFixed(1)}`,
    goalsWidth: width(t.goals),
    xgWidth: width(t.xg),
    text,
  }
}

export function opponentOf(m: PlayerMatch): Team {
  return m.side === 'home' ? m.away : m.home
}

// The player's team's goals first: «3:1» a win
export function scoreFor(m: PlayerMatch): string {
  if (!m.score) return '—'
  const { home, away } = m.score
  return m.side === 'home' ? `${home}:${away}` : `${away}:${home}`
}

// A line of "Матчи сезона"
export interface LogRow {
  id: string
  date: string
  opponent: string
  venue: 'Д' | 'Г'
  score: string
  minutes: number
  goals: number
  assists: number
  xg: string
}

export function logRow(m: PlayerMatch, timeZone?: string): LogRow {
  return {
    id: m.id,
    date: dayMonth(new Date(m.kickoffAt), timeZone),
    opponent: opponentOf(m).name,
    venue: m.side === 'home' ? 'Д' : 'Г',
    score: scoreFor(m),
    minutes: m.minutes,
    goals: m.goals,
    assists: m.assists,
    xg: m.xg.toFixed(2),
  }
}

// The player's page: /players/253989
export function playerPath(playerId: string): string {
  return `/players/${encodeURIComponent(playerId)}`
}
