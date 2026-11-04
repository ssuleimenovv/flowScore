import { flushPromises, mount } from '@vue/test-utils'
import { afterEach, describe, expect, it, vi } from 'vitest'
import type { Stage, Tie, TieList } from '@/shared/api/types'
import PlayoffCard from '../ui/PlayoffCard.vue'

const stages: Stage[] = [
  { id: 'first-round', name: '1-й раунд', ties: 16 },
  { id: 'second-round', name: '2-й раунд', ties: 8 },
  { id: 'round-of-16', name: '1/8', ties: 8 },
]

const team = (id: string, code: string, name: string) => ({ id, code, name })

const shootout: Tie = {
  home: { team: team('1', 'SAS', 'Sassuolo'), goals: [1], total: 1, shootout: 4, through: true },
  away: { team: team('2', 'FRO', 'Frosinone'), goals: [1], total: 1, shootout: 3, through: false },
  legs: [{ matchId: '9', kickoffAt: '2026-09-02T13:00:00Z', state: 'post', clock: '' }],
  finish: 'pens',
}

function respond(items: Tie[], status = 200) {
  const fetch = vi.fn<typeof globalThis.fetch>(async (url) => {
    const stage = String(url).split('/stages/')[1]?.split('/')[0] ?? ''
    const body: TieList = { leagueId: 'ita.coppa_italia', stage, items }
    return new Response(JSON.stringify(body), {
      status,
      headers: { 'Content-Type': 'application/json' },
    })
  })
  vi.stubGlobal('fetch', fetch)
  return fetch
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('PlayoffCard', () => {
  it('opens on the stage being played and draws its ties', async () => {
    const fetch = respond([shootout])
    const card = mount(PlayoffCard, {
      props: { leagueId: 'ita.coppa_italia', stages, current: 'second-round', main: true },
    })
    expect(card.find('[aria-busy="true"]').exists()).toBe(true)
    await flushPromises()

    expect(String(fetch.mock.calls[0]?.[0])).toContain(
      '/leagues/ita.coppa_italia/stages/second-round/ties',
    )
    expect(card.find('h2').text()).toBe('2-й раунд')
    expect(card.find('.side').text()).toBe('1 матч')
    expect(card.find('[aria-pressed="true"]').text()).toBe('2-й раунд')
    expect(card.findAll('.line')).toHaveLength(2)
    expect(card.find('.columns').text()).toBe('Счёт')
    expect(card.find('.note').text()).toBe('пен. 4:3')
    expect(card.find('.mark.through').exists()).toBe(true)
    expect(card.find('.name.muted').text()).toBe('Frosinone')
    expect(card.find('.footnote').text()).toBe(
      'Один матч. При ничьей — дополнительное время и пенальти.',
    )
  })

  it('is headed "Плей-офф" beside a league phase, and loads the stage picked', async () => {
    const fetch = respond([shootout])
    const card = mount(PlayoffCard, {
      props: { leagueId: 'uefa.champions', stages, current: '' },
    })
    await flushPromises()
    expect(card.find('h2').text()).toBe('Плей-офф')
    expect(card.find('.side').text()).toBe('один матч')
    expect(card.find('[aria-pressed="true"]').text()).toBe('1/8')

    await card.findAll('button')[0]!.trigger('click')
    await flushPromises()
    expect(String(fetch.mock.lastCall?.[0])).toContain('/stages/first-round/ties')
    expect(card.find('[aria-pressed="true"]').text()).toBe('1-й раунд')
  })

  it('says when the ties failed to load, and tries again', async () => {
    const fetch = respond([], 500)
    const card = mount(PlayoffCard, {
      props: { leagueId: 'ita.coppa_italia', stages, current: 'second-round' },
    })
    await flushPromises()
    expect(card.text()).toContain('Не удалось загрузить пары')

    respond([shootout])
    await card.find('.error button').trigger('click')
    await flushPromises()
    expect(fetch).toHaveBeenCalledTimes(1)
    expect(card.findAll('.line')).toHaveLength(2)
  })
})
