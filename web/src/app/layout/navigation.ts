import { readonly, ref } from 'vue'
import type { IconName } from '@/shared/up/AppIcon.vue'

export type TabId = 'home' | 'leagues' | 'simulator' | 'search' | 'profile' | 'favorites'

export interface NavItem {
  id: TabId
  label: string
  to: string
  icon: IconName
}

// Mobile tab bar, as on the Mobile * iOS boards
export const tabs: NavItem[] = [
  { id: 'home', label: 'Главная', to: '/', icon: 'home' },
  { id: 'leagues', label: 'Лиги', to: '/leagues', icon: 'trophy' },
  { id: 'simulator', label: 'Симулятор', to: '/simulator', icon: 'sliders' },
  { id: 'search', label: 'Поиск', to: '/search', icon: 'search' },
  { id: 'profile', label: 'Профиль', to: '/settings', icon: 'person' },
]

// Dekstop header sections, as on the Desktop boards
export const sections: NavItem[] = [
  { id: 'home', label: 'Матчи', to: '/', icon: 'home' },
  { id: 'leagues', label: 'Лиги', to: '/leagues', icon: 'trophy' },
  { id: 'simulator', label: 'Симулятор', to: '/simulator', icon: 'sliders' },
  { id: 'favorites', label: 'Избранное', to: '/favorites', icon: 'person' },
]

// Title of an inner screen for the mobile nav bar ("Премьер-лига",
// "8-й тур · Этихад"). The page sets it; the shell shows it.
const title = ref<{ title: string; subtitle?: string } | null>(null)

export function useNavTitle() {
  return {
    title: readonly(title),
    setTitle(next: { title: string; subtitle?: string } | null) {
      title.value = next
    },
  }
}
