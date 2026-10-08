import type { Match, Prediction, Probabilities, Score } from '@/shared/api/types'

export interface OutcomeColumn {
  key: keyof Probabilities
  value: number // percent, the three add up to 100
  label: string // "Победа Сити"
  change: string // "▲ 14 с начала", "▼ 5"; after the match "✓ итог" or ""
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

// After the final whistle the chances are 100 · 0 · 0 and say nothing. The
// card shows instead what the model expected before kick-off and which of
// it came true: the columns are the pre-match chances, the result marked.
export function finalColumns(match: Match, prediction: Prediction, score: Score): OutcomeColumn[] {
  const result = resultOf(score)
  return outcomeColumns(match, prediction).map((column) => ({
    ...column,
    value: prediction.preMatch[column.key],
    change: column.key === result ? '✓ итог' : '',
  }))
}

// "Победа Liverpool 4:5. Неожиданно: до матча модель давала этому 20%"
export function finalVerdict(match: Match, prediction: Prediction, score: Score): string {
  const { preMatch } = prediction
  const result = resultOf(score)
  const outcome =
    result === 'draw' ? 'Ничья' : `Победа ${result === 'home' ? match.home.name : match.away.name}`
  const line = `${outcome} ${score.home}:${score.away}.`
  // The model's favourite is the outcome it gave the most
  const favourite = (['home', 'draw', 'away'] as const).reduce((a, b) =>
    preMatch[b] > preMatch[a] ? b : a,
  )
  return result === favourite
    ? `${line} Модель это ждала: ${preMatch[result]}% до матча.`
    : `${line} Неожиданно: до матча модель давала этому ${preMatch[result]}%.`
}

export function resultOf(score: Score): keyof Probabilities {
  if (score.home > score.away) return 'home'
  if (score.home < score.away) return 'away'
  return 'draw'
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
