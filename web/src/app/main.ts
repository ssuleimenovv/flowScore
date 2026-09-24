import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import { router } from './router'

import '@/styles/tokens.css'
import '@/styles/base.css'

import { motion } from '@/shared/motion'

createApp(App).use(createPinia()).use(router).use(motion).mount('#app')
