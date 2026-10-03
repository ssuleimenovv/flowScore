import type { Problem } from './types'

const BASE = '/api/v1'

// ApiError carries the RFC 9457 problem the server sent. status is 0 when
// the request never reached the server (offline, DNS, server down).
export class ApiError extends Error {
  readonly status: number
  readonly problem: Problem | null

  constructor(status: number, problem: Problem | null) {
    super(problem?.title ?? (status === 0 ? 'Network error' : `HTTP ${status}`))
    this.name = 'ApiError'
    this.status = status
    this.problem = problem
  }
}

export function apiGet<T>(path: string, signal?: AbortSignal): Promise<T> {
  return request<T>(path, { signal, headers: { Accept: 'application/json' } })
}

// POST with a JSON body, for questions too big for a URL ("what if…")
export function apiPost<T>(path: string, body: unknown, signal?: AbortSignal): Promise<T> {
  return request<T>(path, {
    method: 'POST',
    signal,
    headers: { Accept: 'application/json', 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
}

async function request<T>(path: string, init: RequestInit): Promise<T> {
  const { signal } = init
  let res: Response
  try {
    res = await fetch(BASE + path, init)
    } catch (err) {
    if (signal?.aborted) throw err // cancelled on purpose, not a network failure
    throw new ApiError(0, null)
  }

  if (!res.ok) {
    throw new ApiError(res.status, await readProblem(res))
  }
  return (await res.json()) as T
}


async function readProblem(res: Response): Promise<Problem | null> {
  if (!res.headers.get('Content-Type')?.includes('application/problem+json')) {
    return null
  }
  try {
    return (await res.json()) as Problem
  } catch {
    return null
  }
}
