import { createPinia, setActivePinia } from 'pinia'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { useSessionStore } from '../session'

const session = { token: 't1', user: { id: 'u1', email: 'ayan@mail.com', name: 'Аян' } }

function respond(status: number, body: unknown) {
  const fetch = vi.fn<typeof globalThis.fetch>(
    async () =>
      new Response(JSON.stringify(body), {
        status,
        headers: {
          'Content-Type': status < 400 ? 'application/json' : 'application/problem+json',
        },
      }),
  )
  vi.stubGlobal('fetch', fetch)
  return fetch
}

describe('session', () => {
  beforeEach(() => {
    localStorage.clear()
    setActivePinia(createPinia())
  })
  afterEach(() => vi.unstubAllGlobals())

  it('signs in and keeps the session for the next visit', async () => {
    respond(200, session)
    const store = useSessionStore()
    await store.signIn('ayan@mail.com', 'flowscore1')

    expect(store.signedIn).toBe(true)
    expect(store.user?.name).toBe('Аян')
    // A new visit: a fresh store reads it back
    setActivePinia(createPinia())
    expect(useSessionStore().token).toBe('t1')
  })

  it('stays signed out on a wrong password', async () => {
    respond(401, { type: 'about:blank', title: 'Wrong email or password', status: 401 })
    const store = useSessionStore()

    await expect(store.signIn('ayan@mail.com', 'nope')).rejects.toMatchObject({ status: 401 })
    expect(store.signedIn).toBe(false)
  })

  it('sends the token to /me and signs out when it is turned away', async () => {
    localStorage.setItem('fs-session', JSON.stringify(session))
    const fetch = respond(401, { type: 'about:blank', title: 'Session expired', status: 401 })
    const store = useSessionStore()
    await store.check()

    const init = fetch.mock.calls[0]![1]!
    expect((init.headers as Record<string, string>).Authorization).toBe('Bearer t1')
    expect(store.signedIn).toBe(false)
    expect(localStorage.getItem('fs-session')).toBeNull()
  })

  it('stays signed in when the server cannot be reached', async () => {
    localStorage.setItem('fs-session', JSON.stringify(session))
    vi.stubGlobal(
      'fetch',
      vi.fn(async () => {
        throw new TypeError('Failed to fetch')
      }),
    )
    const store = useSessionStore()
    await store.check()

    expect(store.signedIn).toBe(true)
  })
})
