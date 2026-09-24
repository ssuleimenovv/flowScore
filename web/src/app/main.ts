import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import { router } from './router'

import '@fontsource-variable/onest'
import '@fontsource-variable/jetbrains-mono'
import '@fontsource-variable/big-shoulders-display'

import '@/styles/tokens.css'
import '@/styles/base.css'

import { motion } from '@/shared/motion'

createApp(App).use(createPinia()).use(router).use(motion).mount('#app')
