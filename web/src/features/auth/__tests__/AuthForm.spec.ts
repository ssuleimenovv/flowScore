import { flushPromises, mount } from '@vue/test-utils'
import { createPinia } from 'pinia'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { createMemoryHistory, createRouter } from 'vue-router'
import { useToast } from '@/shared/toast/useToast'
import AuthForm from '../ui/AuthForm.vue'

const session = { token: 't1', user: { id: 'u1', email: 'ayan@mail.com', name: 'Аян' } }

function respond(status: number, body: unknown) {
  const fetch = vi.fn<typeof globalThis.fetch>(
    async () =>
      new Response(JSON.stringify(body), {
        status,
        headers: {
          'Content-Type': status < 400 ? 'application/json' : 'application/problem+json',
        },
      }),
  )
  vi.stubGlobal('fetch', fetch)
  return fetch
}

async function mountAt(path: string) {
  const router = createRouter({
    history: createMemoryHistory(),
    routes: [{ path: '/:any(.*)*', component: { render: () => null } }],
  })
  await router.push(path)
  const wrapper = mount(AuthForm, { global: { plugins: [router, createPinia()] } })
  return { wrapper, router }
}

describe('AuthForm', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
    localStorage.clear()
  })

  it('signs in on Enter and goes back where the user came from', async () => {
    const fetch = respond(200, session)
    const { wrapper, router } = await mountAt('/login?redirect=/matches/1')

    await wrapper.find('input[type=email]').setValue('ayan@mail.com')
    await wrapper.find('input[type=password]').setValue('flowscore1')
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    expect(String(fetch.mock.calls[0]![0])).toMatch(/\/api\/v1\/auth\/login$/)
    expect(router.currentRoute.value.fullPath).toBe('/matches/1')
  })

  it('asks a sign-up for a name and the password rule', async () => {
    const { wrapper } = await mountAt('/login')
    await wrapper.findAll('[role=tab]')[1]!.trigger('click')

    expect(wrapper.find('input[autocomplete=name]').exists()).toBe(true)
    expect(wrapper.text()).toContain('Не короче 8 символов')
    expect(wrapper.text()).not.toContain('Забыли пароль?')

    await wrapper.find('input[autocomplete=name]').setValue('Аян')
    await wrapper.find('input[type=email]').setValue('ayan@mail.com')
    await wrapper.find('input[type=password]').setValue('short')
    expect(wrapper.find('button[type=submit]').attributes('disabled')).toBeDefined()
  })

  it('tells a wrong password in a toast and stays on the page', async () => {
    respond(401, { type: 'about:blank', title: 'Wrong email or password', status: 401 })
    const { wrapper, router } = await mountAt('/login')

    await wrapper.find('input[type=email]').setValue('ayan@mail.com')
    await wrapper.find('input[type=password]').setValue('nope')
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    const { toasts, dismiss } = useToast()
    expect(toasts.value.map((t) => t.title)).toContain('Неверная почта или пароль')
    expect(router.currentRoute.value.path).toBe('/login')
    for (const t of toasts.value) dismiss(t.id)
  })

  it('does not follow a redirect to another site', async () => {
    respond(200, session)
    const { wrapper, router } = await mountAt('/login?redirect=//evil.example')

    await wrapper.find('input[type=email]').setValue('ayan@mail.com')
    await wrapper.find('input[type=password]').setValue('flowscore1')
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    expect(router.currentRoute.value.fullPath).toBe('/')
  })
})
