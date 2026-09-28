import 'vue-router'
import type { TabId } from './layout/navigation'

// Typed route meta: the shell reads it to highlight the tab and pick the nav bar
declare module 'vue-router' {
  interface RouteMeta {
    tab?: TabId // which tab or section is active
    inner?: boolean // pushed screen: "Назад" and a title instead of the logo
    stubTitle?: string // temporary, for screens that are not built yet
  }
}
