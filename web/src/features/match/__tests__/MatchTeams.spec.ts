import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import MatchTeams from '../ui/MatchTeams.vue'

const everton = { id: '368', code: 'EVE', name: 'Everton' }
const chelsea = { id: '363', code: 'CHE', name: 'Chelsea' }

describe('MatchTeams', () => {
  it('puts each name next to its crest', () => {
    const teams = mount(MatchTeams, {
      props: {
        home: everton,
        away: chelsea,
        homeLogo: 'https://x/368.png',
        awayLogo: 'https://x/363.png',
      },
    })
    expect(teams.findAll('.name').map((n) => n.text())).toEqual(['Everton', 'Chelsea'])
    expect(teams.findAll('img').map((i) => i.attributes('src'))).toEqual([
      'https://x/368.png',
      'https://x/363.png',
    ])
  })

  it("shows a team's letters when it has no crest", () => {
    const teams = mount(MatchTeams, {
      props: { home: everton, away: chelsea, homeLogo: null, awayLogo: 'https://x/363.png' },
    })
    expect(teams.findAll('img')).toHaveLength(1)
    expect(teams.text()).toContain('EVE')
  })
})
