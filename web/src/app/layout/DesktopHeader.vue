<script setup lang="ts">
import { RouterLink } from 'vue-router'
import AppIcon from '@/shared/ui/AppIcon.vue'
import AppLogo from '@/shared/ui/AppLogo.vue'
import ThemeToggle from './ThemeToggle.vue'
import { sections, type TabId } from './navigation'

defineProps<{ active: TabId | undefined }>()
</script>

<template>
  <header class="bar">
    <div class="start">
      <RouterLink to="/" class="brand fs-link">
        <AppLogo variant="outline" :size="30" />
        <span class="wordmark">FLOWSCORE</span>
      </RouterLink>

      <nav aria-label="Основная навигация" class="sections">
        <RouterLink
          v-for="item in sections"
          :key="item.id"
          :to="item.to"
          class="section fs-link"
          :aria-current="item.id === active ? 'page' : undefined"
        >
          {{ item.label }}
        </RouterLink>
      </nav>
    </div>

    <div class="end">
      <RouterLink to="/search" class="search fs-link">
        <AppIcon name="search" :size="18" />
        <span class="search-text">Команда, лига, игрок</span>
        <kbd class="search-key">/</kbd>
      </RouterLink>
      <RouterLink to="/notifications" class="square fs-link" aria-label="Уведомления">
        <AppIcon name="bell" :size="20" />
      </RouterLink>
      <ThemeToggle variant="boxed" />
      <RouterLink to="/settings" class="avatar fs-link" aria-label="Профиль">
        <AppIcon name="person" :size="20" />
      </RouterLink>
    </div>
  </header>
</template>

<style scoped>
.bar {
  position: sticky;
  top: 0;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 72px;
  padding: 0 var(--fs-page-padding);
  border-bottom: var(--fs-hairline) solid var(--fs-separator);
  background: var(--fs-bar-bg);
  -webkit-backdrop-filter: var(--fs-material);
  backdrop-filter: var(--fs-material);
}

.start {
  display: flex;
  gap: 36px;
  align-items: center;
}

.brand {
  display: flex;
  gap: 10px;
  align-items: center;
  color: var(--fs-text);
  text-decoration: none;
}

.wordmark {
  font-family: var(--fs-font-display);
  font-size: 26px;
  font-weight: 800;
  letter-spacing: 0.05em;
}

.sections {
  display: flex;
  gap: var(--fs-space-4);
}

.section {
  padding: 10px 14px;
  border-radius: 10px;
  color: var(--fs-muted);
  font-size: var(--fs-text-subheadline);
  font-weight: 500;
  text-decoration: none;
}

.section[aria-current='page'] {
  background: var(--fs-surface-2);
  color: var(--fs-text);
}

.end {
  display: flex;
  gap: 10px;
  align-items: center;
}

.search {
  display: flex;
  gap: 10px;
  align-items: center;
  width: 260px;
  height: var(--fs-touch);
  padding: 0 14px;
  border: 1px solid var(--fs-line);
  border-radius: var(--fs-radius-md);
  background: var(--fs-surface);
  color: var(--fs-muted);
  font-size: var(--fs-text-subheadline);
  text-decoration: none;
}

.search-text {
  flex-grow: 1;
}

.search-key {
  padding: 2px 6px;
  border: 1px solid var(--fs-line);
  border-radius: 6px;
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-caption);
}

.square,
.avatar {
  display: flex;
  align-items: center;
  justify-content: center;
  width: var(--fs-touch);
  height: var(--fs-touch);
  color: var(--fs-text);
}

.square {
  border: 1px solid var(--fs-line);
  border-radius: var(--fs-radius-md);
  background: var(--fs-surface);
}

.avatar {
  border-radius: var(--fs-radius-full);
  background: var(--fs-surface-2);
}

/* 768–1199: no room for the full search field, it becomes a square button */
@media (max-width: 1199px) {
  .search {
    justify-content: center;
    width: var(--fs-touch);
    padding: 0;
  }

  .search-text,
  .search-key {
    display: none;
  }
}
</style>
