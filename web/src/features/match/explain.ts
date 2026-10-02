import type { FlowFactor, Side } from '@/shared/api/types'
import type { MatchLive } from './matchState'

export interface FactorLine {
  key: string
  text: string // "4 удара за последние 7 минут"
  value: string // "+14", "−5"
  side: Side // whose color: the leader's when it pushes the leader up
  width: string // the bar, as a share of the largest factor
}

export interface Explanation {
  title: string // "Сити забирает инициативу"
  text: string // "Поток вырос с 64 до 82 за 10 минут."
  lines: FactorLine[]
}

// Below this gap neither team is ahead in a way worth a headline
const EVEN = 10
// A change over 10 minutes small enough to call Flow steady
const STEADY = 5
// The bar of a factor this size is full; a larger one fills it and sets the scale
const FULL_BAR = 15

const other = (side: Side): Side => (side === 'home' ? 'away' : 'home')

// What Flow says right now, for the "Explainability" card. Every number comes
// from the factors the engine sent: a factor is a real part of the team's Flow
// (docs/FLOW.md, section 6), so the card explains it rather than guesses.
export function explain(live: MatchLive, limit: number): Explanation {
  const { flow, delta10, match } = live
  const leader: Side = flow.home >= flow.away ? 'home' : 'away'
  const name = match[leader].name
  const now = Math.round(flow[leader])
  const change = Math.round(delta10[leader])

  let title = `${name} контролирует игру`
  if (Math.abs(flow.home - flow.away) < EVEN) title = 'Равная игра'
  else if (change >= EVEN) title = `${name} забирает инициативу`

  let text = `Поток ${name} держится около ${now}.`
  if (change >= STEADY) text = `Поток ${name} вырос с ${now - change} до ${now} за 10 минут.`
  if (change <= -STEADY) text = `Поток ${name} снизился с ${now - change} до ${now} за 10 минут.`

  // Against the leader is negative: its own setbacks and what the other team
  // builds. What holds the other team back is left out: possession at 43% is
  // the leader's 57% seen from the other side, and would say the same twice
  const signed = live.factors
    .filter((f) => f.side === leader || f.value > 0)
    .map((f) => ({ f, v: f.side === leader ? f.value : -f.value }))
    .sort((a, b) => Math.abs(b.v) - Math.abs(a.v))
    .slice(0, limit)
  const scale = Math.max(FULL_BAR, ...signed.map(({ v }) => Math.abs(v)))

  const lines = signed.map(({ f, v }) => ({
    key: `${f.side}-${f.key}`,
    text: f.side === leader ? label(f) : `${match[f.side].name}: ${lowerFirst(label(f))}`,
    value: v > 0 ? `+${Math.round(v)}` : `−${Math.abs(Math.round(v))}`,
    side: v > 0 ? leader : other(leader),
    width: `${((Math.abs(v) / scale) * 100).toFixed(1)}%`,
  }))
  return { title, text, lines }
}

// The line of one factor: the events of the last minutes when there are any,
// else just what kind of events still holds Flow up
function label(f: FlowFactor): string {
  const n = f.count
  // The window is 10 minutes, so of the numbers ending in 1 only 1 itself comes up
  const recent =
    f.minutes === 1
      ? 'за последнюю минуту'
      : `за последние ${f.minutes} ${plural(f.minutes, 'минуту', 'минуты', 'минут')}`
  switch (f.key) {
    case 'shots':
      return n ? `${n} ${plural(n, 'удар', 'удара', 'ударов')} ${recent}` : 'Удары по воротам'
    case 'goals':
      return n > 1 ? `${n} ${plural(n, 'гол', 'гола', 'голов')} ${recent}` : 'Забитый гол'
    case 'key_passes':
      return n
        ? `${n} ${plural(n, 'ключевая передача', 'ключевые передачи', 'ключевых передач')} ${recent}`
        : 'Ключевые передачи'
    case 'possession': {
      // The value is the whole match decayed; the share is the last minutes. Show
      // the share only when both point the same way, or the line contradicts itself
      const share = f.share === undefined ? null : Math.round(f.share * 100)
      const agrees = share !== null && share !== 50 && share > 50 === f.value > 0
      return agrees ? `Владение ${share}% ${recent}` : 'Владение мячом'
    }
    case 'corners':
      return n ? `${n} ${plural(n, 'угловой', 'угловых', 'угловых')} ${recent}` : 'Угловые'
    case 'cards':
      return n > 1 ? `${n} жёлтые у соперника` : 'Жёлтая у соперника'
    case 'red_cards':
      return 'Удаление'
    case 'substitutions':
      return n ? `${n} ${plural(n, 'замена', 'замены', 'замен')} ${recent}` : 'Замены'
    default:
      return 'Другие события'
  }
}

// Russian plural: 1 удар, 2 удара, 5 ударов, 21 удар
function plural(n: number, one: string, few: string, many: string): string {
  const tens = n % 100
  if (tens >= 11 && tens <= 14) return many
  if (n % 10 === 1) return one
  if (n % 10 >= 2 && n % 10 <= 4) return few
  return many
}

function lowerFirst(s: string): string {
  return s.charAt(0).toLowerCase() + s.slice(1)
}
