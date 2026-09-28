<script setup lang="ts">
// A pulsing placeholder block. The container around skeletons carries aria-busy
// and a label; the blocks themselves mean nothing to a screen reader.
const {
  width,
  height = 16,
  round = false,
} = defineProps<{
  width?: string
  height?: number
  round?: boolean
}>()
</script>

<template>
  <span
    class="skeleton"
    :class="{ round }"
    :style="{ width, height: `${height}px` }"
    aria-hidden="true"
  />
</template>

<style scoped>
.skeleton {
  display: block;
  flex-shrink: 0;
  border-radius: var(--fs-radius-sm);
  background: var(--fs-surface-2);
  animation: pulse var(--fs-duration-skeleton) ease-in-out infinite;
}

.round {
  border-radius: var(--fs-radius-full);
}

@keyframes pulse {
  50% {
    opacity: 0.4;
  }
}
</style>
