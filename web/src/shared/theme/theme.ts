export type ThemePreference = 'dark' | 'light' | 'system'
export type Theme = 'dark' | 'light'

// Keep in sync with the inline script in index.html.
export const STORAGE_KEY = 'fs-theme'
export const DEFAULT_PREFERENCE: ThemePreference = 'dark' // as in the mockup

export function readPreference(): ThemePreference {
  try {
    const value = localStorage.getItem(STORAGE_KEY)
    return value === 'dark' || value === 'light' || value === 'system' ? value : DEFAULT_PREFERENCE
  } catch {
    return DEFAULT_PREFERENCE // storage blocked, e.g. in private mode
  }
}

export function savePreference(preference: ThemePreference): void {
  try {
    localStorage.setItem(STORAGE_KEY, preference)
  } catch {
    // the choice lasts until reload; nothing else to do
  }
}

export function resolveTheme(preference: ThemePreference, systemDark: boolean): Theme {
  if (preference === 'system') return systemDark ? 'dark' : 'light'
  return preference
}

// applyTheme switches the tokens in tokens.css and repaints the browser chrome
// (the iOS status bar area, the Android address bar) in the new background.
export function applyTheme(theme: Theme): void {
  const root = document.documentElement
  root.dataset.theme = theme

  const background = getComputedStyle(root).getPropertyValue('--fs-bg').trim()
  document.querySelector('meta[name="theme-color"]')?.setAttribute('content', background)
}
