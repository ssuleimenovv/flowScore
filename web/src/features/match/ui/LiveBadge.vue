<script setup lang="ts">
import { computed } from 'vue'
import type { Match } from '@/shared/api/types'

// "● LIVE 72′" from the Match boards. Other statuses get a quiet grey pill.
const { status, minute } = defineProps<{ status: Match['status']; minute: string }>()

const label = computed(() => {
  switch (status) {
    case 'live':
      return `LIVE ${minute}`
    case 'halftime':
      return 'ПЕРЕРЫВ'
    case 'finished':
      return 'ЗАВЕРШЁН'
    case 'postponed':
      return 'ПЕРЕНЕСЁН'
    default:
      return 'СКОРО'
  }
})

const isLive = computed(() => status === 'live' || status === 'halftime')
</script>

<template>
  <span class="badge" :class="{ live: isLive }">
    <span v-if="isLive" class="dot" aria-hidden="true" />
    {{ label }}
  </span>
</template>

<style scoped>
.badge {
  display: inline-flex;
  gap: 6px;
  align-items: center;
  height: 24px;
  padding: 0 10px;
  border-radius: var(--fs-radius-full);
  background: var(--fs-surface-2);
  color: var(--fs-muted);
  font-size: var(--fs-text-caption);
  font-weight: 700;
  letter-spacing: 0.08em;
  white-space: nowrap;
}

.live {
  background: var(--fs-live-soft);
  color: var(--fs-live);
}

.dot {
  position: relative;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: currentColor;
}

/* The ring that keeps spreading from the dot: "this is happening now" */
.dot::after {
  position: absolute;
  inset: 0;
  border-radius: 50%;
  background: inherit;
  content: '';
  animation: ping var(--fs-duration-pulse) ease-out infinite;
}

@keyframes ping {
  from {
    opacity: 0.75;
    transform: scale(1);
  }

  to {
    opacity: 0;
    transform: scale(3);
  }
}

/* An endless decoration: nothing to lose when it stops */
@media (prefers-reduced-motion: reduce) {
  .dot::after {
    animation: none;
  }
}
</style>
