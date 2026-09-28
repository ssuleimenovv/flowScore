<script setup lang="ts">
import { computed } from 'vue'
import AppButton from './AppButton.vue'
import AppSpinner from './AppSpinner.vue'

// The live stream dropped (States board). The bar fills up until the next attempt;
// "Сейчас" skips the wait.
const { seconds, retryAt } = defineProps<{
  seconds: number | null
  retryAt: number | null
}>()
defineEmits<{ reconnect: [] }>()

const countdown = computed(() =>
  seconds ? `Переподключение через ${seconds} с…` : 'Переподключаемся…',
)

// A CSS animation runs the bar, so there is no per-frame JS. It restarts
// whenever retryAt changes (the key below).
const duration = computed(() => (retryAt === null ? 0 : Math.max(0, retryAt - Date.now())))
</script>

<template>
  <div class="banner" role="status">
    <div class="row">
      <AppSpinner class="spinner" />
      <!-- The countdown changes every second; a screen reader would read each tick -->
      <p class="text">
        <b>Live-поток прервался.</b> <span aria-hidden="true">{{ countdown }}</span>
      </p>
      <AppButton variant="secondary" size="small" @click="$emit('reconnect')">Сейчас</AppButton>
    </div>
    <div v-if="retryAt !== null" class="track">
      <div :key="retryAt" class="fill" :style="{ animationDuration: `${duration}ms` }" />
    </div>
  </div>
</template>

<style scoped>
.banner {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: var(--fs-space-12) var(--fs-space-16);
  border: 1px solid var(--fs-home);
  border-radius: 14px;
  background: var(--fs-home-soft);
  font-size: var(--fs-text-subheadline);
}

.row {
  display: flex;
  gap: 10px;
  align-items: center;
}

.spinner {
  color: var(--fs-home);
}

.text {
  flex-grow: 1;
}

.track {
  height: 3px;
  overflow: hidden;
  border-radius: var(--fs-radius-full);
  background: var(--fs-track);
}

.fill {
  height: 100%;
  background: var(--fs-home);
  transform-origin: left;
  animation: fill linear both;
}

@keyframes fill {
  from {
    transform: scaleX(0);
  }
}
</style>
