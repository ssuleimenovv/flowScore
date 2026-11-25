// Each competition's accent (AccentsExtended): its brand colour at the
// lightness that reads, 4.5:1 or more, on the page and the cards of a theme
export interface Accent {
  dark: string
  light: string
}

const ACCENTS: Record<string, Accent> = {
  'eng.1': { dark: '#a468d6', light: '#3d195b' },
  'esp.1': { dark: '#ff6b00', light: '#b24b00' },
  'ger.1': { dark: '#fa2f3e', light: '#d20515' },
  'ita.1': { dark: '#497bff', light: '#1b59ff' },
  'fra.1': { dark: '#c8d400', light: '#686e00' },
  'uefa.champions': { dark: '#8076e4', light: '#3d2dd6' },
  'uefa.europa': { dark: '#f26f21', light: '#b44a0b' },
  'uefa.europa.conf': { dark: '#00b140', light: '#007b2d' },
}

// A domestic cup takes its country's league, the Super Cup the Champions League's
const SAME_AS: Record<string, string> = {
  'eng.fa': 'eng.1',
  'eng.league_cup': 'eng.1',
  'esp.copa_del_rey': 'esp.1',
  'ger.dfb_pokal': 'ger.1',
  'ita.coppa_italia': 'ita.1',
  'fra.coupe_de_france': 'fra.1',
  'uefa.super_cup': 'uefa.champions',
}

export function competitionAccent(id: string | null): Accent | null {
  if (!id) return null
  return ACCENTS[SAME_AS[id] ?? id] ?? null
}

// The style that puts an accent on a page's root, with the class fs-tinted
export function tintStyle(accent: Accent | null): Record<string, string> {
  return accent ? { '--fs-tint-dark': accent.dark, '--fs-tint-light': accent.light } : {}
}

// The competitions with an artwork of their own (CompetitionBackdropsExtended)
export type ArtKey =
  | 'uefa.champions'
  | 'uefa.europa'
  | 'uefa.europa.conf'
  | 'esp.1'
  | 'eng.1'
  | 'ger.1'
  | 'ita.1'
  | 'fra.1'

export interface Art {
  key: ArtKey
  word: string // the heading's highlighted word, bright on the artwork
  deep: string // the colour the white text stands on: the picked chip's text
}

const ARTS: Record<ArtKey, Omit<Art, 'key'>> = {
  'uefa.champions': { word: '#8fd8ff', deep: '#071042' },
  'uefa.europa': { word: '#ffb066', deep: '#2a0f02' },
  'uefa.europa.conf': { word: '#7dffb0', deep: '#03241a' },
  'esp.1': { word: '#ffd27a', deep: '#a3140f' },
  'eng.1': { word: '#f4b8ff', deep: '#24093f' },
  'ger.1': { word: '#ffc2c6', deep: '#7a0710' },
  'ita.1': { word: '#a9c8ff', deep: '#0a1f6b' },
  'fra.1': { word: '#e4f25a', deep: '#0f1400' },
}

// A competition's artwork; a cup stands on its league's. null for one
// without: its page keeps the theme's plain background.
export function competitionArt(id: string | null): Art | null {
  if (!id) return null
  const key = (SAME_AS[id] ?? id) as ArtKey
  return key in ARTS ? { key, ...ARTS[key] } : null
}
