import { computed, onScopeDispose, ref, type Ref } from 'vue'
import type { Match } from '@/shared/api/types'

// The server sends the clock once, with the moment it was measured
// (observedAt). While the match is live the client keeps it ticking.
export function useMatchClock(match: Ref<Match | undefined>) {
  const now = ref(Date.now())
  const timer = setInterval(() => (now.value = Date.now()), 1000)
  onScopeDispose(() => clearInterval(timer))

  const seconds = computed(() => {
    const m = match.value
    if (!m) return 0
    const { elapsedSeconds, observedAt } = m.clock
    if (m.status !== 'live') return elapsedSeconds
    return elapsedSeconds + Math.max(0, (now.value - Date.parse(observedAt)) / 1000)
  })

  return { seconds }
}
