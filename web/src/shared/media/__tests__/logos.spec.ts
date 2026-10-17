import { describe, expect, it } from 'vitest'
import { leagueLogo, matchTeamLogo, teamLogo } from '../logos'

describe('logos', () => {
  it("asks ESPN for a team's crest at twice the size on screen", () => {
    expect(teamLogo('368', 32)).toBe(
      'https://a.espncdn.com/combiner/i?img=/i/teamlogos/soccer/500/368.png&w=64&h=64',
    )
  })

  it("knows each league's logo by its own number", () => {
    expect(leagueLogo('eng.1', 18)).toContain('/i/leaguelogos/soccer/500/23.png&w=36')
    expect(leagueLogo('eng.1', 18, 'dark')).toContain('/i/leaguelogos/soccer/500-dark/23.png')
    expect(leagueLogo('uefa.champions', 18)).toBeNull()
  })

  it("finds a real team's crest by its ESPN id and a demo team's by its name", () => {
    const teams = {
      home: { id: '36', code: 'MCI', name: 'Manchester City' },
      away: { id: '363', code: 'CHE', name: 'Chelsea' },
    }
    expect(matchTeamLogo({ source: 'espn', ...teams }, 'away', 52)).toContain('/363.png&w=104')
    // A StatsBomb id means nothing to ESPN: Manchester City is 382 there
    expect(matchTeamLogo({ source: 'replay', ...teams }, 'home', 52)).toContain('/382.png')
    // Chelsea is not one of the demo teams
    expect(matchTeamLogo({ source: 'replay', ...teams }, 'away', 52)).toBeNull()
  })
})
