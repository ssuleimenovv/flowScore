// @vitest-environment node
import { describe, expect, it } from 'vitest'
import type { Match, Simulation } from '@/shared/api/types'
import { clampMinute, formatChances, redMinutes, scenarioColumns, scenarioStory } from '../scenario'

const match = {
  home: { id: '36', code: 'MCI', name: 'Сити' },
  away: { id: '1', code: 'ARS', name: 'Арсенал' },
} as Match

const sim: Simulation = {
  minute: 38,
  current: { home: 32, draw: 40, away: 28 },
  scenario: { home: 16, draw: 34, away: 50 },
}

describe('redMinutes', () => {
  it('runs from the next minute to the 89th', () => {
    expect(redMinutes(38.7)).toEqual({ min: 39, max: 89 })
    expect(redMinutes(0)).toEqual({ min: 1, max: 89 })
  })

  it('has nothing left at the end of the match', () => {
    expect(redMinutes(89)).toBeNull()
  })

  it('moves a passed minute along with the clock', () => {
    expect(clampMinute(60, { min: 72, max: 89 })).toBe(72)
    expect(clampMinute(80, { min: 72, max: 89 })).toBe(80)
  })
})

describe('scenarioColumns', () => {
  it('shows the scenario and how far it moved', () => {
    expect(scenarioColumns(match, sim)).toEqual([
      { key: 'home', label: 'Победа Сити', value: 16, change: '−16' },
      { key: 'draw', label: 'Ничья', value: 34, change: '−6' },
      { key: 'away', label: 'Победа Арсенал', value: 50, change: '+22' },
    ])
    expect(formatChances(sim.current)).toBe('32 · 40 · 28')
  })
})

describe('scenarioStory', () => {
  it('waits for a scenario', () => {
    expect(scenarioStory(match, sim, null)).toEqual({
      headline: 'Сценарий пока не задан',
      reasons: [],
    })
  })

  it('explains a red in the model numbers', () => {
    expect(scenarioStory(match, sim, { side: 'home', minute: 60 })).toEqual({
      headline: 'Шансы Сити: 32% → 16%',
      reasons: [
        {
          title: 'Красная Сити на 60′',
          detail: 'вдесятером 30 минут до конца',
          value: '−16% Сити',
          side: 'away',
        },
      ],
    })
  })

  it('counts the minutes left in Russian', () => {
    const at = (minute: number) =>
      scenarioStory(match, sim, { side: 'away', minute }).reasons[0]?.detail
    expect(at(89)).toBe('вдесятером 1 минуту до конца')
    expect(at(87)).toBe('вдесятером 3 минуты до конца')
  })
})
