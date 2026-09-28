<script setup lang="ts">
import { RouterLink } from 'vue-router'
import AppIcon from '@/shared/ui/AppIcon.vue'
import { tabs, type TabId } from './navigation'

defineProps<{ active: TabId | undefined }>()
</script>

<template>
  <nav aria-label="Панель вкладок" class="bar">
    <RouterLink
      v-for="tab in tabs"
      :key="tab.id"
      :to="tab.to"
      class="tab fs-link"
      :aria-current="tab.id === active ? 'page' : undefined"
    >
      <AppIcon :name="tab.icon" :size="25" />
      {{ tab.label }}
    </RouterLink>
  </nav>
</template>

<style scoped>
.bar {
  position: fixed;
  right: 0;
  bottom: 0;
  left: 0;
  z-index: 10;
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  align-items: start;
  padding: 0 var(--fs-space-4) var(--fs-safe-bottom);
  border-top: var(--fs-hairline) solid var(--fs-separator);
  background: var(--fs-bar-bg);
  -webkit-backdrop-filter: var(--fs-material);
  backdrop-filter: var(--fs-material);
}

.tab {
  display: flex;
  flex-direction: column;
  gap: 2px;
  align-items: center;
  justify-content: center;
  height: var(--fs-tab-height);
  color: var(--fs-muted);
  font-size: var(--fs-text-tab);
  font-weight: 500;
  letter-spacing: 0.01em;
  text-decoration: none;
}

.tab[aria-current='page'] {
  color: var(--fs-home);
}
</style>
