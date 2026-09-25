import { createRouter, createWebHistory } from 'vue-router'

// Development-only pages; they are not part of the producation build
const devRoutes = [
  {
    path: '/dev/stream/:matchId?',
    name: 'dev-stream',
    component: () => import('@/pages/DevStreamPage.vue'),
  },
]

export const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('@/pages/HomePage.vue'),
    },
    ...(import.meta.env.DEV ? devRoutes : []),
  ],
})
