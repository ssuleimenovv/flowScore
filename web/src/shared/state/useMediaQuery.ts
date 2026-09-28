import { onScopeDispose, readonly, ref } from 'vue'

// Whether a media query matches right now, kept up to date. For layout use CSS;
// this is for what CSS cannot switch, like ARIA roles.
export function useMediaQuery(query: string) {
  const list = matchMedia(query)
  const matches = ref(list.matches)

  const onChange = (e: MediaQueryListEvent) => (matches.value = e.matches)
  list.addEventListener('change', onChange)
  onScopeDispose(() => list.removeEventListener('change', onChange))

  return readonly(matches)
}
