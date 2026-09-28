import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { useToast } from '../useToast'

const { toasts, show, dismiss } = useToast()

describe('useToast', () => {
  beforeEach(() => vi.useFakeTimers())

  afterEach(() => {
    for (const t of toasts.value) dismiss(t.id)
    vi.useRealTimers()
  })

  it('hides a toast after its duration', () => {
    show({ icon: 'check', tone: 'ok', title: 'Сценарий сохранён', duration: 3000 })
    expect(toasts.value).toHaveLength(1)

    vi.advanceTimersByTime(2999)
    expect(toasts.value).toHaveLength(1)
    vi.advanceTimersByTime(1)
    expect(toasts.value).toHaveLength(0)
  })

  it('keeps at most three, dropping the oldest', () => {
    for (const title of ['1', '2', '3', '4']) show({ icon: 'ball', tone: 'accent', title })
    expect(toasts.value.map((t) => t.title)).toEqual(['2', '3', '4'])
  })

  it('closes a toast on demand', () => {
    const id = show({ icon: 'ball', tone: 'accent', title: 'Гол!' })
    dismiss(id)
    expect(toasts.value).toHaveLength(0)
  })
})
