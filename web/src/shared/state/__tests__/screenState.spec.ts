// @vitest-environment node
import { describe, expect, it } from 'vitest'
import { resolveScreenState, type ScreenInput } from '../screenState'

const base: ScreenInput = {
  status: 'success',
  hasData: true,
  isEmpty: false,
  online: true,
}

describe('resolveScreenState', () => {
  it.each([
    ['content when data is loaded', {}, 'content', null],
    ['empty when response has no items', { isEmpty: true }, 'empty', null],
    ['loading on first request', { status: 'pending', hasData: false }, 'loading', null],
    ['error when request failed without data', { status: 'error', hasData: false }, 'error', null],
    [
      'offline instead of error without data',
      { status: 'error', hasData: false, online: false },
      'offline',
      null,
    ],
    ['keeps content when offline with data', { online: false }, 'content', 'offline'],
    ['keeps content when request failed with data', { status: 'error' }, 'content', null],
    ['reconnecting when socket dropped', { socket: 'connecting' }, 'content', 'reconnecting'],
    [
      'offline banner wins over reconnecting',
      { online: false, socket: 'closed' },
      'content',
      'offline',
    ],
  ] as const)('%s', (_, patch, view, banner) => {
    expect(resolveScreenState({ ...base, ...patch })).toEqual({ view, banner })
  })
})
