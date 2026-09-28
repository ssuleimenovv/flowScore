<script setup lang="ts">
import { RouterLink, type RouteLocationRaw } from 'vue-router'
import AppIcon, { type IconName } from './AppIcon.vue'
import AppSpinner from './AppSpinner.vue'

// Buttons from the States board. With `to` it renders a link that looks like a button
// ("Смотреть расписание"), otherwise a real <button>.
const {
  variant = 'primary',
  size = 'regular',
  icon,
  to,
  loading = false,
  disabled = false,
} = defineProps<{
  variant?: 'primary' | 'secondary'
  size?: 'regular' | 'small'
  icon?: IconName
  to?: RouteLocationRaw
  loading?: boolean
  disabled?: boolean
}>()
</script>

<template>
  <RouterLink
    v-if="to"
    :to
    class="button fs-press"
    :class="[variant, size, { square: !$slots.default }]"
  >
    <AppIcon v-if="icon" :name="icon" :size="18" />
    <slot />
  </RouterLink>
  <button
    v-else
    type="button"
    class="button fs-press"
    :class="[variant, size, { loading, square: !$slots.default }]"
    :disabled="disabled || loading"
    :aria-busy="loading || undefined"
  >
    <!-- While loading the spinner takes the icon's place, as on "Повторить" → "Загружаем…" -->
    <AppSpinner v-if="loading" />
    <AppIcon v-else-if="icon" :name="icon" :size="18" />
    <slot />
  </button>
</template>

<style scoped>
.button {
  display: inline-flex;
  flex-shrink: 0;
  gap: var(--fs-space-8);
  align-items: center;
  justify-content: center;
  height: var(--fs-touch);
  padding: 0 18px;
  border-radius: var(--fs-radius-md);
  font-size: var(--fs-text-subheadline);
  font-weight: 600;
  white-space: nowrap;
  text-decoration: none;
  cursor: pointer;
}

.small {
  height: 36px;
  padding: 0 var(--fs-space-12);
  font-size: var(--fs-text-footnote);
}

/* Icon only (no text in the slot): a square. The caller gives it an aria-label */
.square {
  aspect-ratio: 1;
  padding: 0;
}

.primary {
  border: 0;
  background: var(--fs-accent);
  color: var(--fs-on-accent);
}

.secondary {
  border: 1px solid var(--fs-line);
  background: transparent;
  color: var(--fs-text);
}

/* The board shows hover as the team color, not a brightness change */
@media (hover: hover) {
  .primary:not(:disabled):hover {
    background: var(--fs-home);
    color: var(--fs-on-home);
    filter: none;
  }
}

.loading {
  cursor: progress;
}

.button:disabled:not(.loading) {
  border-color: transparent;
  background: var(--fs-surface-2);
  color: var(--fs-faint);
  cursor: not-allowed;
  filter: none;
  transform: none;
}
</style>
