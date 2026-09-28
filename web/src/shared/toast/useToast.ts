import { readonly, ref } from 'vue'
import type { IconName } from '@/shared/ui/AppIcon.vue'

export interface Toast {
  id: number
  icon: IconName
  tone: 'ok' | 'accent'
  title: string
  text?: string
}

export type ToastInput = Omit<Toast, 'id'> & { duration?: number }

const DEFAULT_DURATION = 5000
// More than three at once is noise; the oldest goes first
const MAX_VISIBLE = 3

const toasts = ref<Toast[]>([])
const timers = new Map<number, ReturnType<typeof setTimeout>>()
let nextId = 1

function dismiss(id: number) {
  clearTimeout(timers.get(id))
  timers.delete(id)
  toasts.value = toasts.value.filter((t) => t.id !== id)
}

function show({ duration = DEFAULT_DURATION, ...toast }: ToastInput): number {
  const id = nextId++
  toasts.value = [...toasts.value, { id, ...toast }]
  while (toasts.value.length > MAX_VISIBLE) dismiss(toasts.value[0]!.id)
  timers.set(
    id,
    setTimeout(() => dismiss(id), duration),
  )
  return id
}

// One queue for the whole app: any code can show a toast, ToastHost in App.vue draws them.
export function useToast() {
  return { toasts: readonly(toasts), show, dismiss }
}
