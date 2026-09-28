import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

const stub = () => import('@/pages/StubPage.vue')

// Development-only pages; they are not part of the production build.
const devRoutes: RouteRecordRaw[] = [
  {
    path: '/dev/stream/:matchId?',
    name: 'dev-stream',
    component: () => import('@/pages/DevStreamPage.vue'),
    meta: { tab: 'home', inner: true },
  },
  {
    path: '/dev/ui',
    name: 'dev-ui',
    component: () => import('@/pages/DevUiPage.vue'),
  },
]

export const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('@/pages/HomePage.vue'),
      meta: { tab: 'home' },
    },
    // Screens from the mockup that are not built yet
    { path: '/leagues', component: stub, meta: { tab: 'leagues', stubTitle: 'Лиги' } },
    { path: '/simulator', component: stub, meta: { tab: 'simulator', stubTitle: 'Симулятор' } },
    { path: '/search', component: stub, meta: { tab: 'search', stubTitle: 'Поиск' } },
    { path: '/settings', component: stub, meta: { tab: 'profile', stubTitle: 'Настройки' } },
    { path: '/favorites', component: stub, meta: { tab: 'favorites', stubTitle: 'Избранное' } },
    { path: '/notifications', component: stub, meta: { inner: true, stubTitle: 'Уведомления' } },
    ...(import.meta.env.DEV ? devRoutes : []),
  ],
})
