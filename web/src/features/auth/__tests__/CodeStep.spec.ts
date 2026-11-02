import { flushPromises, mount } from '@vue/test-utils'
import { createPinia } from 'pinia'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { createMemoryHistory, createRouter } from 'vue-router'
import { useToast } from '@/shared/toast/useToast'
import CodeStep from '../ui/CodeStep.vue'

const session = { token: 't1', user: { id: 'u1', email: 'ayan@mail.com', name: 'Аян' } }

function answer(status: number, body: unknown) {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': status < 400 ? 'application/json' : 'application/problem+json' },
  })
}

async function mountStep() {
  const router = createRouter({
    history: createMemoryHistory(),
    routes: [{ path: '/:any(.*)*', component: { render: () => null } }],
  })
  await router.push('/login?redirect=/matches/1')
  const wrapper = mount(CodeStep, {
    props: { email: 'ayan@mail.com' },
    global: { plugins: [router, createPinia()] },
  })
  return { wrapper, router }
}

describe('CodeStep', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
    vi.useRealTimers()
    localStorage.clear()
    const { toasts, dismiss } = useToast()
    for (const t of toasts.value) dismiss(t.id)
  })

  it('sends the sixth digit, tells a wrong code, then signs in', async () => {
    const fetch = vi
      .fn<typeof globalThis.fetch>()
      .mockResolvedValueOnce(
        answer(422, { type: 'about:blank', title: 'Wrong code', status: 422, attemptsLeft: 4 }),
      )
      .mockResolvedValueOnce(answer(200, session))
    vi.stubGlobal('fetch', fetch)
    const { wrapper, router } = await mountStep()

    // Pasted with a space: the input keeps the digits and sends at six
    await wrapper.find('input').setValue('482 190')
    await flushPromises()
    expect(wrapper.find('[role=alert]').text()).toBe('Неверный код. Осталось 4 попытки.')

    await wrapper.find('input').setValue('482193')
    await flushPromises()
    expect(JSON.parse(String(fetch.mock.calls[1]![1]!.body))).toEqual({
      email: 'ayan@mail.com',
      code: '482193',
    })
    expect(router.currentRoute.value.fullPath).toBe('/matches/1')
  })

  it('lets a new code go out after a minute', async () => {
    vi.useFakeTimers()
    const fetch = vi.fn<typeof globalThis.fetch>(async () =>
      answer(202, { email: 'ayan@mail.com' }),
    )
    vi.stubGlobal('fetch', fetch)
    const { wrapper } = await mountStep()

    expect(wrapper.text()).toContain('через 1:00')
    await vi.advanceTimersByTimeAsync(18_000)
    expect(wrapper.text()).toContain('через 0:42')

    await vi.advanceTimersByTimeAsync(42_000)
    await wrapper.find('button.again').trigger('click')
    await flushPromises()
    expect(String(fetch.mock.calls[0]![0])).toMatch(/\/api\/v1\/auth\/resend$/)
    expect(useToast().toasts.value.map((t) => t.title)).toContain('Новый код отправлен')
    expect(wrapper.text()).toContain('через 1:00')
  })
})
