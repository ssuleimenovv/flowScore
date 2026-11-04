import { flushPromises, mount } from '@vue/test-utils'
import { createPinia } from 'pinia'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { createMemoryHistory, createRouter } from 'vue-router'
import { useToast } from '@/shared/toast/useToast'
import ResetStep from '../ui/ResetStep.vue'

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
  await router.push('/forgot?redirect=/leagues')
  const wrapper = mount(ResetStep, {
    props: { email: 'ayan@mail.com' },
    global: { plugins: [router, createPinia()] },
    attachTo: document.body,
  })
  return { wrapper, router }
}

describe('ResetStep', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
    localStorage.clear()
    document.body.innerHTML = ''
    const { toasts, dismiss } = useToast()
    for (const t of toasts.value) dismiss(t.id)
  })

  it('moves on to the password, tells a wrong code, then sets the password', async () => {
    const fetch = vi
      .fn<typeof globalThis.fetch>()
      .mockResolvedValueOnce(
        answer(422, { type: 'about:blank', title: 'Wrong code', status: 422, attemptsLeft: 4 }),
      )
      .mockResolvedValueOnce(answer(200, session))
    vi.stubGlobal('fetch', fetch)
    const { wrapper, router } = await mountStep()
    const [codeField, passwordField] = wrapper.findAll('input')

    // Six digits without a password: nothing goes out, the password is next
    await codeField!.setValue('482190')
    await flushPromises()
    expect(fetch).not.toHaveBeenCalled()
    expect(document.activeElement).toBe(passwordField!.element)
    expect(wrapper.find('button[type=submit]').attributes('disabled')).toBeDefined()

    await passwordField!.setValue('newpass22')
    await wrapper.find('form').trigger('submit')
    await flushPromises()
    expect(wrapper.find('[role=alert]').text()).toBe('Неверный код. Осталось 4 попытки.')

    await codeField!.setValue('482193')
    await flushPromises()
    expect(String(fetch.mock.calls[1]![0])).toMatch(/\/api\/v1\/auth\/reset$/)
    expect(JSON.parse(String(fetch.mock.calls[1]![1]!.body))).toEqual({
      email: 'ayan@mail.com',
      code: '482193',
      password: 'newpass22',
    })
    expect(router.currentRoute.value.fullPath).toBe('/leagues')
  })

  it('asks for a new code once the old one is gone', async () => {
    const fetch = vi
      .fn<typeof globalThis.fetch>()
      .mockResolvedValueOnce(
        answer(410, { type: 'about:blank', title: 'No code to check', status: 410 }),
      )
      .mockResolvedValueOnce(answer(202, { email: 'ayan@mail.com' }))
    vi.stubGlobal('fetch', fetch)
    const { wrapper } = await mountStep()

    const [codeField, passwordField] = wrapper.findAll('input')
    await passwordField!.setValue('newpass22')
    await codeField!.setValue('482193')
    await flushPromises()
    expect(useToast().toasts.value.map((t) => t.title)).toContain('Код истёк, отправь новый')

    await wrapper.find('button.again').trigger('click')
    await flushPromises()
    expect(String(fetch.mock.calls[1]![0])).toMatch(/\/api\/v1\/auth\/forgot$/)
    expect(useToast().toasts.value.map((t) => t.title)).toContain('Новый код отправлен')
  })
})
