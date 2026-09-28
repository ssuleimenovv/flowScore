import type { StatRow } from '@/shared/api/types'

export interface StatLine {
  key: string
  label: string // "Удары"
  home: string // "14"
  away: string // "6"
  homeWidth: string // "70%": the side's share of the row, drawn as its bar
  awayWidth: string
}

// The rows of "Статистика матча", in the server's order. A key the client does
// not know yet is left out rather than shown without a label.
const LABELS: Record<string, string> = {
  possession: 'Владение',
  shots: 'Удары',
  shots_on_target: 'В створ',
  xg: 'xG',
  key_passes: 'Ключевые передачи',
  pressure_tackles: 'Отборы в прессинге',
}

export function statLines(rows: StatRow[]): StatLine[] {
  return rows
    .filter((row) => row.key in LABELS)
    .map((row) => {
      const total = row.home + row.away
      const share = (n: number) => `${total > 0 ? ((n / total) * 100).toFixed(1) : 0}%`
      return {
        key: row.key,
        label: LABELS[row.key]!,
        home: formatStat(row.key, row.home),
        away: formatStat(row.key, row.away),
        homeWidth: share(row.home),
        awayWidth: share(row.away),
      }
    })
}

// "64%", "1.84", "14", the way the Match board writes each row
function formatStat(key: string, value: number): string {
  if (key === 'possession') return `${Math.round(value)}%`
  if (key === 'xg') return value.toFixed(2)
  return String(Math.round(value))
}
