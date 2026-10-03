import type { Match, Probabilities, Side, Simulation } from '@/shared/api/types'

// The regular end of the match on the clock, as "до конца" counts it
const FULL_TIME = 90

// The minutes a red can still be shown: from the next one to the 89th, the
// last with time left to play a man down
export function redMinutes(now: number): { min: number; max: number } | null {
  const min = Math.max(1, Math.floor(now) + 1)
  const max = FULL_TIME - 1
  return min > max ? null : { min, max }
}

// A minute picked earlier may have passed while the match went on: it moves
// along with the clock instead of pointing at the past
export function clampMinute(minute: number, range: { min: number; max: number }): number {
  return Math.min(range.max, Math.max(range.min, minute))
}

export interface ScenarioColumn {
  key: keyof Probabilities
  value: number
  label: string // "Победа Сити"
  change: string // "+14", "−5", "0"
}

// The three results of the scenario and how far each moved from the match as
// it stands
export function scenarioColumns(match: Match, sim: Simulation): ScenarioColumn[] {
  const columns = [
    { key: 'home', label: `Победа ${match.home.name}` },
    { key: 'draw', label: 'Ничья' },
    { key: 'away', label: `Победа ${match.away.name}` },
  ] as const
  return columns.map(({ key, label }) => ({
    key,
    label,
    value: sim.scenario[key],
    change: signed(sim.scenario[key] - sim.current[key]),
  }))
}

// "32 · 40 · 28"
export function formatChances(p: Probabilities): string {
  return `${p.home} · ${p.draw} · ${p.away}`
}

export interface ScenarioReason {
  title: string // "Красная Сити на 60′"
  detail: string // "вдесятером 30 мин до конца"
  value: string // "−16% Сити"
  side: Side // whose color the value takes
}

export interface ScenarioStory {
  headline: string
  reasons: ScenarioReason[]
}

// What the Explainability card says about the scenario. Every number is the
// model's: the home team's chance before and after the change.
export function scenarioStory(
  match: Match,
  sim: Simulation | null,
  red: { side: Side; minute: number } | null,
): ScenarioStory {
  if (!sim || !red) return { headline: 'Сценарий пока не задан', reasons: [] }

  const home = match.home.name
  const change = sim.scenario.home - sim.current.home
  const team = match[red.side].name
  const left = FULL_TIME - red.minute
  return {
    headline: `Шансы ${home}: ${sim.current.home}% → ${sim.scenario.home}%`,
    reasons: [
      {
        title: `Красная ${team} на ${red.minute}′`,
        detail: `вдесятером ${left} ${plural(left, 'минуту', 'минуты', 'минут')} до конца`,
        value: `${signed(change)}% ${home}`,
        side: change >= 0 ? 'home' : 'away',
      },
    ],
  }
}

function signed(n: number): string {
  if (n > 0) return `+${n}`
  if (n < 0) return `−${-n}`
  return '0'
}

function plural(n: number, one: string, few: string, many: string): string {
  const tens = n % 100
  if (tens >= 11 && tens <= 14) return many
  if (n % 10 === 1) return one
  if (n % 10 >= 2 && n % 10 <= 4) return few
  return many
}
