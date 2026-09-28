import type { Match, MatchEvent, Side } from '@/shared/api/types'
import { surname } from './wave'

export interface ChronicleItem {
  id: string
  minute: string // "45+2′"
  code: string // "ГОЛ"
  title: string // "Гол! Foden"
  detail: string | null // "Ассист De Bruyne · 1:1"
  goal: Side | null // a goal row is highlighted in the scorer's color
  impact: { text: string; side: Side } | null // "+22" for the side it helped
}

// Fouls, offsides and take-ons weigh 0 until calibration (docs/FLOW.md), so they
// would only bury the moments that moved the flow.
const CODES: Partial<Record<MatchEvent['type'], string>> = {
  goal: 'ГОЛ',
  shot_on_target: 'УДР',
  shot_off_target: 'УДР',
  shot_blocked: 'УДР',
  corner: 'УГЛ',
  yellow_card: 'ЖК',
  red_card: 'КК',
  substitution: 'ЗАМ',
  halftime: 'ПЕР',
  fulltime: 'КОН',
}

// The chronicle, newest first, as MatchLive keeps the events.
export function chronicle(events: MatchEvent[], match: Match): ChronicleItem[] {
  // The score after each event, counted from the oldest one
  const scoreAfter = new Map<string, string>()
  let home = 0
  let away = 0
  for (let i = events.length - 1; i >= 0; i--) {
    const e = events[i]!
    if (e.type === 'goal' && e.side === 'home') home++
    if (e.type === 'goal' && e.side === 'away') away++
    scoreAfter.set(e.id, `${home}:${away}`)
  }

  const items: ChronicleItem[] = []
  for (const e of events) {
    const code = CODES[e.type]
    if (!code) continue
    items.push({
      id: e.id,
      minute: e.addedTime ? `${e.minute}+${e.addedTime}′` : `${e.minute}′`,
      code,
      ...describe(e, match, scoreAfter.get(e.id) ?? ''),
      goal: e.type === 'goal' ? e.side : null,
      impact: impact(e),
    })
  }
  return items
}

function describe(e: MatchEvent, match: Match, score: string) {
  const who = e.player ? surname(e.player.name) : null
  const team = e.side ? match[e.side].name : null
  const xg = typeof e.xG === 'number' ? `xG ${e.xG.toFixed(2)}` : null
  const labelled = (title: string, name: string | null) => (name ? `${title} · ${name}` : title)

  switch (e.type) {
    case 'goal':
      return {
        title: who ? `Гол! ${who}` : 'Гол!',
        detail: [e.assist && `Ассист ${surname(e.assist.name)}`, score].filter(Boolean).join(' · '),
      }
    case 'shot_on_target':
      return { title: labelled('Удар в створ', who), detail: xg }
    case 'shot_off_target':
      return { title: labelled('Удар мимо', who), detail: xg }
    case 'shot_blocked':
      return { title: labelled('Удар заблокирован', who), detail: xg }
    case 'corner':
      return { title: labelled('Угловой', team), detail: null }
    case 'yellow_card':
      return { title: labelled('Жёлтая', who), detail: team }
    case 'red_card':
      return { title: labelled('Красная', who), detail: team }
    case 'substitution':
      return { title: labelled('Замена', team), detail: who && `Уходит ${who}` }
    case 'halftime':
      return { title: 'Перерыв', detail: `Счёт ${score}` }
    default:
      return { title: 'Конец матча', detail: `Счёт ${score}` }
  }
}

// The event's Flow weight, credited to the side it helped: a card helps the
// opponent of the player who got it (docs/FLOW.md, section 4).
function impact(e: MatchEvent): ChronicleItem['impact'] {
  const n = Math.round(e.flowImpact ?? 0)
  if (n === 0 || !e.side) return null
  const card = e.type === 'yellow_card' || e.type === 'red_card'
  const side: Side = card ? (e.side === 'home' ? 'away' : 'home') : e.side
  return { text: n > 0 ? `+${n}` : `−${Math.abs(n)}`, side }
}
