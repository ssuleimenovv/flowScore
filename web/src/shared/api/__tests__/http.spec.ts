// @vitest-environment node
import { afterEach, describe, expect, it, vi } from 'vitest'
import { ApiError, apiGet } from '../http'

function respond(status: number, body: unknown, contentType = 'application/json') {
  vi.stubGlobal(
    'fetch',
    vi.fn<typeof fetch>(
      async () =>
        new Response(JSON.stringify(body), { status, headers: { 'Content-Type': contentType } }),
    ),
  )
}

describe('apiGet', () => {
  afterEach(() => vi.unstubAllGlobals())

  it('returns the JSON body', async () => {
    respond(200, { id: 'm1' })
    await expect(apiGet('/matches/m1')).resolves.toEqual({ id: 'm1' })
  })

  it('turns a problem response into ApiError', async () => {
    respond(
      404,
      { type: 'about:blank', title: 'Match not found', status: 404 },
      'application/problem+json',
    )

    const err = await apiGet('/matches/nope').catch((e: unknown) => e)

    expect(err).toBeInstanceOf(ApiError)
    expect(err).toMatchObject({ status: 404, message: 'Match not found' })
  })

  it('reports a network failure as status 0', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn<typeof fetch>(async () => {
        throw new TypeError('Failed to fetch')
      }),
    )

    await expect(apiGet('/matches/m1')).rejects.toMatchObject({ status: 0 })
  })
})
