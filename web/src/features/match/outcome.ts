import type { Match, Prediction, Probabilities } from '@/shared/api/types'

export interface OutcomeColumn {
  key: keyof Probabilities
  value: number // percent, the three add up to 100
  label: string // "Победа Сити"
  change: string // "▲ 14 с начала", "▼ 5"
}

// The three columns of the "Вероятность исхода" card: the chance now and how
// it moved since the pre-match one. Only the first says "с начала": the row
// reads left to right, and saying it once is enough.
export function outcomeColumns(match: Match, prediction: Prediction): OutcomeColumn[] {
  const { current, preMatch } = prediction
  const columns = [
    { key: 'home', label: `Победа ${match.home.name}` },
    { key: 'draw', label: 'Ничья' },
    { key: 'away', label: `Победа ${match.away.name}` },
  ] as const
  return columns.map(({ key, label }, i) => ({
    key,
    label,
    value: current[key],
    change: change(current[key] - preMatch[key]) + (i === 0 ? ' с начала' : ''),
  }))
}

// "До матча: 38 · 30 · 32"
export function formatPreMatch({ preMatch }: Prediction): string {
  return `До матча: ${preMatch.home} · ${preMatch.draw} · ${preMatch.away}`
}

function change(delta: number): string {
  if (delta > 0) return `▲ ${delta}`
  if (delta < 0) return `▼ ${-delta}`
  return '= 0'
}
