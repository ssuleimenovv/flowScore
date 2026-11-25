// @vitest-environment node
import { describe, expect, it } from 'vitest'
import { competitionAccent, competitionArt, tintStyle } from '../accents'

describe('accents', () => {
  it('gives a cup its league’s accent and a strange competition none', () => {
    expect(competitionAccent('ger.dfb_pokal')).toEqual(competitionAccent('ger.1'))
    expect(competitionAccent('uefa.super_cup')).toEqual(competitionAccent('uefa.champions'))
    expect(competitionAccent('usa.1')).toBeNull()
    expect(competitionAccent(null)).toBeNull()
  })

  it('puts both themes’ colours on the page', () => {
    expect(tintStyle(competitionAccent('esp.1'))).toEqual({
      '--fs-tint-dark': '#ff6b00',
      '--fs-tint-light': '#b24b00',
    })
    expect(tintStyle(null)).toEqual({})
  })

  it('puts a cup on its league’s artwork, and leaves a competition without one plain', () => {
    expect(competitionArt('esp.copa_del_rey')?.key).toBe('esp.1')
    expect(competitionArt('uefa.super_cup')?.key).toBe('uefa.champions')
    expect(competitionArt('uefa.europa')).toBeNull()
    expect(competitionArt(null)).toBeNull()
  })
})
