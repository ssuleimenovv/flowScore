<script setup lang="ts">
import { RouterLink, useRouter } from 'vue-router'
import AppIcon from '@/shared/ui/AppIcon.vue'
import AppLogo from '@/shared/ui/AppLogo.vue'
import ThemeToggle from './ThemeToggle.vue'
import { useNavTitle } from './navigation'

// "root" is a tab screen (logo, bell, theme); "inner" is a pushed screen
// (Назад, title, theme), as on MobileMatch.
defineProps<{ inner: boolean }>()

const router = useRouter()
const { title } = useNavTitle()

function back() {
  // Opened by a direct link, there is no page to go back to
  if (window.history.state?.back) router.back()
  else router.push('/')
}
</script>

<template>
  <header class="bar" :class="{ inner }">
    <div v-if="!inner" class="row root">
      <RouterLink to="/" class="brand fs-link" aria-label="FlowScore, на главную">
        <AppLogo variant="filled" :size="26" />
        <span class="name">FlowScore</span>
      </RouterLink>
      <div class="actions">
        <RouterLink to="/notifications" class="icon-button fs-link" aria-label="Уведомления">
          <AppIcon name="bell" :size="22" />
        </RouterLink>
        <ThemeToggle variant="plain" />
      </div>
    </div>

    <div v-else class="row pushed">
      <button type="button" class="back fs-link" @click="back">
        <AppIcon name="back" />
        Назад
      </button>
      <div class="title">
        <span class="title-main">{{ title?.title }}</span>
        <span v-if="title?.subtitle" class="title-sub">{{ title.subtitle }}</span>
      </div>
      <div class="end">
        <ThemeToggle variant="plain" />
      </div>
    </div>
  </header>
</template>

<style scoped>
.bar {
  position: sticky;
  top: 0;
  z-index: 10;
  padding-top: var(--fs-safe-top);
  background: var(--fs-bar-bg);
  -webkit-backdrop-filter: var(--fs-material);
  backdrop-filter: var(--fs-material);
}

.inner {
  border-bottom: var(--fs-hairline) solid var(--fs-separator);
}

.row {
  height: var(--fs-nav-height);
}

.root {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 var(--fs-space-8) 0 var(--fs-inset);
}

.brand {
  display: flex;
  gap: var(--fs-space-8);
  align-items: center;
  color: var(--fs-text);
  text-decoration: none;
}

.name {
  font-size: var(--fs-text-headline);
  font-weight: 700;
  letter-spacing: -0.01em;
}

.actions {
  display: flex;
}

.icon-button {
  display: flex;
  align-items: center;
  justify-content: center;
  width: var(--fs-touch);
  height: var(--fs-touch);
  color: var(--fs-home);
}

.pushed {
  display: grid;
  grid-template-columns: 110px minmax(0, 1fr) 110px;
  align-items: center;
  padding: 0 var(--fs-space-8);
}

.back {
  display: flex;
  gap: 2px;
  align-items: center;
  height: var(--fs-touch);
  padding: 0;
  border: 0;
  background: none;
  color: var(--fs-home);
  font-size: var(--fs-text-body);
  cursor: pointer;
}

.title {
  display: flex;
  flex-direction: column;
  align-items: center;
  min-width: 0;
}

.title-main {
  max-width: 100%;
  overflow: hidden;
  font-size: var(--fs-text-headline);
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.title-sub {
  color: var(--fs-muted);
  font-size: var(--fs-text-caption);
  white-space: nowrap;
}

.end {
  display: flex;
  justify-content: flex-end;
}
</style>
