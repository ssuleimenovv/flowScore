// @vitest-environment node
import { describe, expect, it } from 'vitest'
import { competitionAccent } from '@/features/league/accents'
import { clubAccent, contrast, teamAccent } from '../clubColor'

describe('clubColor', () => {
  it('keeps a colour that reads and moves only the lightness of one that does not', () => {
    // Manchester City: sky blue reads on the dark theme, not on the light one
    const city = clubAccent(['99c5ea'])!
    expect(city.dark).toBe('#99c5ea')
    expect(city.light).toBe('#246ead')
    for (const bg of ['#0b0909', '#151112', '#1f1a1b'])
      expect(contrast(city.dark, bg)).toBeGreaterThanOrEqual(4.5)
    for (const bg of ['#f5f3f2', '#ffffff', '#efebea'])
      expect(contrast(city.light, bg)).toBeGreaterThanOrEqual(4.5)
  })

  it('takes the main colour only, and no grey, black or white for one', () => {
    // Arsenal: the red, not the white after it
    expect(clubAccent(['db0007', 'ffffff'])).toEqual({ dark: '#ff232a', light: '#d60007' })
    // Juventus, as ESPN gives it: black, then the pink of an away kit
    expect(clubAccent(['000000', 'e8a2b0'])).toBeNull()
    expect(clubAccent([])).toBeNull()
  })

  it('gives a club without a colour its league’s accent', () => {
    expect(teamAccent(['000000'], 'ita.1')).toEqual(competitionAccent('ita.1'))
    expect(teamAccent(['99c5ea'], 'eng.1')).toEqual(clubAccent(['99c5ea']))
    expect(teamAccent([], null)).toBeNull()
  })
})
