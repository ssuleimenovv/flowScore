<script setup lang="ts">
import type { Side } from '@/shared/api/types'

// "lg" is the desktop scoreboard, "md" the phone one, "sm" a card on the home screen.
const {
  code,
  side,
  size = 'lg',
} = defineProps<{ code: string; side: Side; size?: 'lg' | 'md' | 'sm' }>()

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

.sm {
  width: 32px;
  height: 32px;
  border-radius: 9px;
  font-size: 10px;
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
