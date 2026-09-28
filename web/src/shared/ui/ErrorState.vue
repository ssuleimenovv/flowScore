<script setup lang="ts">
import { computed } from 'vue'
import { ApiError } from '@/shared/api/http'
import AppButton from './AppButton.vue'
import StateMessage from './StateMessage.vue'

// The load error from the States board: what happened, "Повторить",
// and a code line to quote to support ("код 503 · req 7f3a91").
const {
  error,
  title,
  retrying = false,
} = defineProps<{
  error: unknown
  title: string
  retrying?: boolean
}>()
defineEmits<{ retry: [] }>()

const status = computed(() => (error instanceof ApiError ? error.status : null))

const text = computed(() =>
  status.value === 0
    ? 'Сервер не отвечает. Данные не потеряны — попробуй ещё раз.'
    : 'Сервер ответил ошибкой. Данные не потеряны — попробуй ещё раз.',
)

const details = computed(() => {
  if (!(error instanceof ApiError) || error.status === 0) return null
  const requestId = error.problem?.requestId
  return requestId ? `код ${error.status} · req ${requestId}` : `код ${error.status}`
})
</script>

<template>
  <StateMessage icon="warning" tone="error" :title :text>
    <AppButton icon="retry" :loading="retrying" @click="$emit('retry')">
      {{ retrying ? 'Загружаем…' : 'Повторить' }}
    </AppButton>
    <span v-if="details" class="details">{{ details }}</span>
  </StateMessage>
</template>

<style scoped>
.details {
  color: var(--fs-faint);
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-caption);
  font-variant-numeric: tabular-nums;
}
</style>
