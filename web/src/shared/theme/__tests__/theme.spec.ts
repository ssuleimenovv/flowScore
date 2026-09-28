// @vitest-environment node
import { afterEach, describe, expect, it, vi } from 'vitest'
import { DEFAULT_PREFERENCE, readPreference, resolveTheme } from '../theme'

describe('resolveTheme', () => {
  it.each([
    ['dark', false, 'dark'],
    ['light', true, 'light'],
    ['system', true, 'dark'],
    ['system', false, 'light'],
  ] as const)('%s with system dark = %s gives %s', (preference, systemDark, want) => {
    expect(resolveTheme(preference, systemDark)).toBe(want)
  })
})

describe('readPreference', () => {
  afterEach(() => vi.unstubAllGlobals())

  it('falls back to the default for an unknown value', () => {
    vi.stubGlobal('localStorage', { getItem: () => 'purple' })
    expect(readPreference()).toBe(DEFAULT_PREFERENCE)
  })

  it('falls back to the default when storage is blocked', () => {
    vi.stubGlobal('localStorage', {
      getItem: () => {
        throw new Error('SecurityError')
      },
    })
    expect(readPreference()).toBe(DEFAULT_PREFERENCE)
  })
})
