// @vitest-environment node
import { describe, expect, it } from 'vitest'
import { backoffDelay } from '../backoff'

describe('backoffDelay', () => {
  it.each([
    [0, 0, 500],
    [0, 1, 1000],
    [2, 0, 2000],
    [2, 1, 4000],
    [10, 1, 30_000],
  ])('attempt %i with random %f waits %i ms', (attempt, random, want) => {
    expect(backoffDelay(attempt, undefined, () => random)).toBe(want)
  })
})
