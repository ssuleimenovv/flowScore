import type { Match, MatchPreMatch, PastMatch, Side } from '@/shared/api/types'

// "Перед матчем" and "Новости" from the MatchExtended boards.

export type Outcome = 'win' | 'draw' | 'loss'

export interface FormRow {
  side: Side
  code: string // "EVE"
  name: string
  results: Array<{ id: string; letter: string; outcome: Outcome; title: string }> // oldest first
  flow: string | null // "49": the team's average Flow over the matches the season read
  flowNote: string // "ср. Flow, 4 матча"
}

export interface HeadToHead {
  home: string // "Everton · 1"
  draws: string // "Ничьи · 1"
  away: string // "Chelsea · 3"
  widths: { home: string; draw: string; away: string } // shares of the bar
  note: string // "последние 5 встреч"
}

export interface MeetingRow {
  id: string
  date: string // "21 мар 2026"
  title: string // "Everton — Chelsea"
  competition: string
  score: string // "3:0"
  flow: string | null // "58–61"
}

export interface NewsRow {
  headline: string
  url: string
  meta: string // "ESPN · 12 мин назад"
}

const LETTERS: Record<Outcome, string> = { win: 'В', draw: 'Н', loss: 'П' }
const MONTHS = ['янв', 'фев', 'мар', 'апр', 'мая', 'июн', 'июл', 'авг', 'сен', 'окт', 'ноя', 'дек']

// "21 мар 2026", in the viewer's time zone
export function shortDate(iso: string): string {
  const d = new Date(iso)
  return `${d.getDate()} ${MONTHS[d.getMonth()]} ${d.getFullYear()}`
}

function score(m: PastMatch): string {
  return `${m.score.home}:${m.score.away}`
}

// Each team's last matches as В / Н / П, oldest first, and its Flow in them
export function formRows(pre: MatchPreMatch, match: Match): FormRow[] {
  return (['home', 'away'] as const).map((side) => {
    const form = pre[side]
    return {
      side,
      code: match[side].code,
      name: match[side].name,
      results: form.matches.map((m) => ({
        id: m.id,
        letter: LETTERS[m.result],
        outcome: m.result,
        title: `${m.home.name} ${score(m)} ${m.away.name}, ${shortDate(m.date)} · ${m.competition}`,
      })),
      flow: form.flow === null ? null : String(Math.round(form.flow)),
      flowNote: `ср. Flow, ${form.flowGames} ${plural(form.flowGames, 'матч', 'матча', 'матчей')}`,
    }
  })
}

// "Everton · 1 — Ничьи · 1 — Chelsea · 3" over the last meetings
export function headToHead(pre: MatchPreMatch, match: Match): HeadToHead {
  const { homeWins, draws, awayWins, matches } = pre.meetings
  const total = Math.max(1, homeWins + draws + awayWins)
  const share = (n: number) => `${((n / total) * 100).toFixed(1)}%`
  return {
    home: `${match.home.name} · ${homeWins}`,
    draws: `Ничьи · ${draws}`,
    away: `${match.away.name} · ${awayWins}`,
    widths: { home: share(homeWins), draw: share(draws), away: share(awayWins) },
    note: `последние ${matches.length} ${plural(matches.length, 'встреча', 'встречи', 'встреч')}`,
  }
}

export function meetingRows(pre: MatchPreMatch): MeetingRow[] {
  return pre.meetings.matches.map((m) => ({
    id: m.id,
    date: shortDate(m.date),
    title: `${m.home.name} — ${m.away.name}`,
    competition: m.competition,
    score: score(m),
    flow: m.flow ? `${Math.round(m.flow.home)}–${Math.round(m.flow.away)}` : null,
  }))
}

// "ESPN · 12 мин назад"
export function newsRows(pre: MatchPreMatch, now: Date): NewsRow[] {
  return pre.news.map((n) => ({
    headline: n.headline,
    url: n.url,
    meta: `ESPN · ${ago(new Date(n.publishedAt), now)}`,
  }))
}

function ago(at: Date, now: Date): string {
  const minutes = Math.max(0, Math.floor((now.getTime() - at.getTime()) / 60_000))
  if (minutes < 1) return 'только что'
  if (minutes < 60) return `${minutes} мин назад`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours} ч назад`
  const days = Math.floor(hours / 24)
  return `${days} ${plural(days, 'день', 'дня', 'дней')} назад`
}

function plural(n: number, one: string, few: string, many: string): string {
  const tens = n % 100
  if (tens >= 11 && tens <= 14) return many
  if (n % 10 === 1) return one
  if (n % 10 >= 2 && n % 10 <= 4) return few
  return many
}
