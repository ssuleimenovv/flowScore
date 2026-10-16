import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import TeamBadge from '../ui/TeamBadge.vue'

describe('TeamBadge', () => {
  it('shows the crest without the colored tile', () => {
    const badge = mount(TeamBadge, {
      props: { code: 'EVE', side: 'home', logo: 'https://x/368.png' },
    })
    expect(badge.find('img').attributes('src')).toBe('https://x/368.png')
    expect(badge.classes()).toContain('crest')
    expect(badge.text()).toBe('')
  })

  it('falls back to the three letters when there is no crest', () => {
    const badge = mount(TeamBadge, { props: { code: 'MCI', side: 'home' } })
    expect(badge.find('img').exists()).toBe(false)
    expect(badge.text()).toBe('MCI')
    expect(badge.classes()).toEqual(expect.arrayContaining(['home', 'lg']))
  })

  it('falls back when the crest fails to load, and tries again for a new one', async () => {
    const badge = mount(TeamBadge, {
      props: { code: 'EVE', side: 'away', logo: 'https://x/broken.png' },
    })
    await badge.find('img').trigger('error')
    expect(badge.find('img').exists()).toBe(false)
    expect(badge.text()).toBe('EVE')

    await badge.setProps({ logo: 'https://x/368.png' })
    expect(badge.find('img').exists()).toBe(true)
  })

  it('has a neutral tile without a side, as in the league table', () => {
    const badge = mount(TeamBadge, { props: { code: 'ARS', size: 'xs' } })
    expect(badge.classes()).toEqual(expect.arrayContaining(['neutral', 'xs']))
  })
})
