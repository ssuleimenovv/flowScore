<script setup lang="ts">
import { useTheme } from '@/shared/theme/useTheme'
import AppIcon from '@/shared/ui/AppIcon.vue'

// "boxed" is the bordered desktop button, "plain" the mobile nav bar one.
const { variant } = defineProps<{ variant: 'boxed' | 'plain' }>()
const { theme, toggle } = useTheme()
</script>

<template>
  <button
    type="button"
    class="toggle fs-press"
    :class="variant"
    aria-label="Переключить тему"
    @click="toggle"
  >
    <!-- The icon shows the theme a tap switches to -->
    <AppIcon :name="theme === 'dark' ? 'sun' : 'moon'" :size="variant === 'boxed' ? 20 : 22" />
  </button>
</template>

<style scoped>
.toggle {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  width: var(--fs-touch);
  height: var(--fs-touch);
  padding: 0;
  border-radius: var(--fs-radius-md);
  cursor: pointer;
}

.boxed {
  border: 1px solid var(--fs-line);
  background: var(--fs-surface);
  color: var(--fs-text);
}

.plain {
  border: 0;
  background: transparent;
  color: var(--fs-home);
}
</style>
