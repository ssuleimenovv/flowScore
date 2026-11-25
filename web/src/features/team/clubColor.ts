import { competitionAccent, type Accent } from '@/features/league/accents'

// A club's accent (AccentsExtended, «Цвет клуба»): its colour from ESPN at
// the lightness that reads, 4.5:1 or more, on the page and both cards of a
// theme. Grey, black and white are no colour: the club takes its league's.

type Theme = keyof Accent

// Each theme's page, card and second card (tokens.css)
const SURFACES: Record<Theme, string[]> = {
  dark: ['#0b0909', '#151112', '#1f1a1b'],
  light: ['#f5f3f2', '#ffffff', '#efebea'],
}

// WCAG AA for normal text
const READABLE = 4.5

type Rgb = [number, number, number] // 0–1 each

function rgbOf(hex: string): Rgb {
  const n = parseInt(hex.replace('#', ''), 16)
  return [(n >> 16) & 255, (n >> 8) & 255, n & 255].map((c) => c / 255) as Rgb
}

function hexOf(rgb: Rgb): string {
  return (
    '#' +
    rgb
      .map((c) =>
        Math.round(c * 255)
          .toString(16)
          .padStart(2, '0'),
      )
      .join('')
  )
}

// WCAG's relative luminance
function luminance(rgb: Rgb): number {
  const [r, g, b] = rgb.map((c) => (c <= 0.04045 ? c / 12.92 : ((c + 0.055) / 1.055) ** 2.4)) as Rgb
  return 0.2126 * r + 0.7152 * g + 0.0722 * b
}

// How well one colour reads on the other: 1 to 21
export function contrast(a: string, b: string): number {
  const [hi, lo] = [luminance(rgbOf(a)), luminance(rgbOf(b))].sort((x, y) => y - x) as [
    number,
    number,
  ]
  return (hi + 0.05) / (lo + 0.05)
}

// Hue and saturation stay; only the lightness moves
type Hsl = [number, number, number] // hue 0–360, saturation and lightness 0–1

function hslOf([r, g, b]: Rgb): Hsl {
  const max = Math.max(r, g, b)
  const min = Math.min(r, g, b)
  const l = (max + min) / 2
  if (max === min) return [0, 0, l]
  const d = max - min
  const s = l > 0.5 ? d / (2 - max - min) : d / (max + min)
  const h =
    max === r ? (g - b) / d + (g < b ? 6 : 0) : max === g ? (b - r) / d + 2 : (r - g) / d + 4
  return [h * 60, s, l]
}

function rgbFrom([h, s, l]: Hsl): Rgb {
  const k = (n: number) => (n + h / 30) % 12
  const a = s * Math.min(l, 1 - l)
  const f = (n: number) => l - a * Math.max(-1, Math.min(k(n) - 3, Math.min(9 - k(n), 1)))
  return [f(0), f(8), f(4)]
}

// Grey, black or white: too little of a hue to tell a club by
function colourless(hex: string): boolean {
  const [, s, l] = hslOf(rgbOf(hex))
  return s < 0.2 || l < 0.08 || l > 0.94
}

// The colour at the lightness nearest its own that reads on all of a theme's
// surfaces: lighter on the dark theme, darker on the light one
function readable(hex: string, theme: Theme): string {
  const [h, s, l] = hslOf(rgbOf(hex))
  const step = theme === 'dark' ? 0.01 : -0.01
  for (let x = l; x >= 0 && x <= 1; x += step) {
    const candidate = hexOf(rgbFrom([h, s, x]))
    if (SURFACES[theme].every((bg) => contrast(candidate, bg) >= READABLE)) return candidate
  }
  return theme === 'dark' ? '#ffffff' : '#000000'
}

// The club's main colour, fitted to both themes; null when it has none or
// it is grey, black or white. ESPN's second colour is often an away kit's,
// Juventus's pink, so it is not the club's.
export function clubAccent(colors: string[]): Accent | null {
  const main = colors[0] ? `#${colors[0]}` : null
  if (!main || !/^#[0-9a-f]{6}$/i.test(main) || colourless(main)) return null
  return { dark: readable(main, 'dark'), light: readable(main, 'light') }
}

// The club's page: its own colour, else its league's accent, else none
// (the page keeps --fs-home)
export function teamAccent(colors: string[], leagueId: string | null): Accent | null {
  return clubAccent(colors) ?? competitionAccent(leagueId)
}
