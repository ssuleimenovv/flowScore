<script setup lang="ts">
import AppIcon, { type IconName } from './AppIcon.vue'

// Empty and error states from the States board: icon tile, title, text and
// actions in the slot. Put it inside an AppCard.
const {
  icon,
  title,
  text,
  tone = 'neutral',
} = defineProps<{
  icon: IconName
  title: string
  text?: string
  tone?: 'neutral' | 'error'
}>()
</script>

<template>
  <div class="message" :class="tone">
    <span class="tile">
      <AppIcon :name="icon" :size="26" />
    </span>
    <h2 class="title">{{ title }}</h2>
    <p v-if="text" class="text">{{ text }}</p>
    <div v-if="$slots.default" class="actions">
      <slot />
    </div>
  </div>
</template>

<style scoped>
.message {
  display: flex;
  flex-direction: column;
  gap: 10px;
  align-items: center;
  padding: 28px;
  text-align: center;
}

.tile {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 56px;
  height: 56px;
  border-radius: 16px;
  background: var(--fs-surface-2);
  color: var(--fs-muted);
}

.title {
  font-size: var(--fs-text-headline);
  font-weight: 600;
}

.text {
  color: var(--fs-muted);
  font-size: var(--fs-text-subheadline);
  line-height: 1.5;
}

/* Actions sit a little lower than the text, as on the board */
.actions {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-12);
  align-items: center;
  margin-top: var(--fs-space-4);
}

.error {
  gap: var(--fs-space-12);
}

.error .tile {
  background: var(--fs-home-soft);
  color: var(--fs-home);
}

.error .title {
  font-size: var(--fs-text-title-3);
}

.error .text {
  max-width: 320px;
}

.error .actions {
  margin-top: 0;
}
</style>
