<script setup lang="ts">
import AppIcon, { type IconName } from '@/shared/ui/AppIcon.vue'

// A toast from the States board: icon tile, title, text and a close button.
const {
  icon,
  tone = 'accent',
  title,
  text,
} = defineProps<{
  icon: IconName
  tone?: 'ok' | 'accent'
  title: string
  text?: string
}>()
defineEmits<{ close: [] }>()
</script>

<template>
  <div class="toast">
    <span class="tile" :class="tone">
      <AppIcon :name="icon" :size="18" />
    </span>
    <div class="body">
      <span class="title">{{ title }}</span>
      <span v-if="text" class="text">{{ text }}</span>
    </div>
    <button type="button" class="close fs-press" aria-label="Закрыть" @click="$emit('close')">
      <AppIcon name="close" :size="16" />
    </button>
  </div>
</template>

<style scoped>
.toast {
  display: grid;
  grid-template-columns: 36px minmax(0, 1fr) 32px;
  gap: var(--fs-space-12);
  align-items: start;
  padding: 14px;
  border: 1px solid var(--fs-line);
  border-radius: 14px;
  background: var(--fs-surface);
}

.tile {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: var(--fs-surface-2);
}

.ok {
  color: var(--fs-ok);
}

.accent {
  color: var(--fs-home);
}

.body {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.title {
  font-size: var(--fs-text-subheadline);
  font-weight: 600;
}

.text {
  color: var(--fs-muted);
  font-size: var(--fs-text-footnote);
}

.close {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  padding: 0;
  border: 0;
  background: transparent;
  color: var(--fs-muted);
  cursor: pointer;
}

/* The button is drawn at 32 px, but the finger gets the full 44 */
.close::before {
  position: absolute;
  inset: -6px;
  content: '';
}
</style>
