import { computed, readonly, ref, watch } from 'vue'
import {
  applyTheme,
  readPreference,
  resolveTheme,
  savePreference,
  type Theme,
  type ThemePreference,
} from './theme'

// One theme for the whole app, shared by every component that asks for it.
const preference = ref<ThemePreference>(readPreference())

const systemQuery = matchMedia('(prefers-color-scheme: dark)')
const systemDark = ref(systemQuery.matches)
systemQuery.addEventListener('change', (e) => (systemDark.value = e.matches))

const theme = computed<Theme>(() => resolveTheme(preference.value, systemDark.value))

applyTheme(theme.value)
watch(theme, (next) => switchTheme(next))

// switchTheme cross-fades the whole page instead of flipping it at once:
// a sudden jump from dark to light is hard on the eyes.
function switchTheme(next: Theme) {
  const reduce = matchMedia('(prefers-reduced-motion: reduce)').matches
  if (!document.startViewTransition || reduce) {
    applyTheme(next)
    return
  }
  document.startViewTransition(() => applyTheme(next))
}

export function useTheme() {
  return {
    preference: readonly(preference),
    theme,
    setPreference(next: ThemePreference) {
      preference.value = next
      savePreference(next)
    },
    toggle() {
      const next: Theme = theme.value === 'dark' ? 'light' : 'dark'
      preference.value = next
      savePreference(next)
    },
  }
}
