<script setup lang="ts">
import { onScopeDispose, ref } from 'vue'
import AppSpinner from './AppSpinner.vue'

// Shown under a loading skeleton that takes too long. The demo backend runs on
// a free host that sleeps when nobody visits, and the first request wakes it:
// up to a minute of skeleton looks broken unless it says why.
const { after = 4000 } = defineProps<{ after?: number }>()

const visible = ref(false)
const timer = setTimeout(() => (visible.value = true), after)
onScopeDispose(() => clearTimeout(timer))
</script>

<template>
  <p v-if="visible" class="note fs-in" role="status">
    <AppSpinner />
    <span>
      Сервер просыпается — бесплатный хостинг засыпает без посетителей. Это займёт до
      минуты, страница откроется сама.
    </span>
  </p>
</template>

<style scoped>
.note {
  display: flex;
  gap: var(--fs-space-12);
  align-items: center;
  padding: var(--fs-space-16) var(--fs-space-24);
  border: 1px solid var(--fs-line);
  border-radius: var(--fs-radius-lg);
  background: var(--fs-surface);
  color: var(--fs-muted);
  font-size: var(--fs-text-subheadline);
  line-height: 1.45;
}

.note :deep(.spinner) {
  flex-shrink: 0;
}
</style>
