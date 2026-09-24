import { readonly, ref } from 'vue'

const online = ref(navigator.onLine)

window.addEventListener('online', () => (online.value = true))
window.addEventListener('offline', () => (online.value = false))

export function useNetwork() {
  return { online: readonly(online) }
}
