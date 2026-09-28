<script setup lang="ts">
import type { Side } from '@/shared/api/types'

// The team tile with its three-letter code, tinted with the side's color.
// "lg" is the desktop scoreboard, "md" the phone one.
const { code, side, size = 'lg' } = defineProps<{ code: string; side: Side; size?: 'lg' | 'md' }>()
</script>

<template>
  <span class="badge" :class="[side, size]" aria-hidden="true">{{ code }}</span>
</template>

<style scoped>
.badge {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  font-family: var(--fs-font-display);
  font-weight: 800;
  letter-spacing: 0.04em;
}

.lg {
  width: 72px;
  height: 72px;
  border-radius: var(--fs-radius-lg);
  font-size: 23px;
}

.md {
  width: 52px;
  height: 52px;
  border-radius: 15px;
  font-size: var(--fs-text-headline);
}

.home {
  background: var(--fs-home-soft);
  color: var(--fs-home);
}

.away {
  background: var(--fs-away-soft);
  color: var(--fs-away);
}
</style>
