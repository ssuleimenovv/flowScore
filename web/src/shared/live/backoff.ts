export interface BackoffOptions {
  baseMs: number
  maxMs: number
}

export const DEFAULT_BACKOFF: BackoffOptions = { baseMs: 1000, maxMs: 30_000 }

// Exponential backoff with "equal jitter": the delay doubles with each attempt
// up to maxMs, and a random half of it is added. Without the random part, every
// client dropped by a server restart would come back at the same moment.
export function backoffDelay(
  attempt: number,
  options: BackoffOptions = DEFAULT_BACKOFF,
  random: () => number = Math.random,
): number {
  const cap = Math.min(options.maxMs, options.baseMs * 2 ** attempt)
  return Math.round(cap / 2 + (random() * cap) / 2)
}
