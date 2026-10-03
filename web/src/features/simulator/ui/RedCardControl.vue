<script setup lang="ts">
import { useId } from 'vue'
import type { Match, Side } from '@/shared/api/types'

// "Красная карточка" from the Simulator boards: whose player is sent off and
// when. The minute slider starts at the next minute of the match.
const { match, range } = defineProps<{
  match: Match
  range: { min: number; max: number } | null // null when no minute is left
}>()
const side = defineModel<Side | null>('side', { required: true })
const minute = defineModel<number>('minute', { required: true })

const titleId = useId()
const options = [
  { value: null, label: 'Нет' },
  { value: 'home', label: match.home.name },
  { value: 'away', label: match.away.name },
] as const

function onMinute(event: Event) {
  minute.value = Number((event.target as HTMLInputElement).value)
}
</script>

<template>
  <div class="red" role="group" :aria-labelledby="titleId">
    <h3 :id="titleId" class="title">Красная карточка</h3>
    <div class="options">
      <button
        v-for="option in options"
        :key="String(option.value)"
        type="button"
        class="option fs-press"
        :aria-pressed="side === option.value"
        :disabled="!range && option.value !== null"
        @click="side = option.value"
      >
        {{ option.label }}
      </button>
    </div>
    <label class="minute">
      <span class="caption">
        <span>Минута удаления</span>
        <span class="value">{{ minute }}′</span>
      </span>
      <input
        type="range"
        :min="range?.min ?? 1"
        :max="range?.max ?? 89"
        step="1"
        :value="minute"
        :disabled="!side || !range"
        @input="onMinute"
      />
    </label>
  </div>
</template>

<style scoped>
.red {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.title {
  font-size: var(--fs-text-headline);
  font-weight: 600;
}

.options {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 6px;
}

/* The picked option is filled with the text color, as on the board */
.option {
  overflow: hidden;
  height: var(--fs-touch);
  padding: 0 var(--fs-space-8);
  border: 1px solid var(--fs-line);
  border-radius: 10px;
  background: transparent;
  color: var(--fs-text);
  font-size: var(--fs-text-subheadline);
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
  cursor: pointer;
  transition:
    background-color var(--fs-duration-color) var(--fs-ease),
    color var(--fs-duration-color) var(--fs-ease);
}

.option[aria-pressed='true'] {
  border-color: var(--fs-text);
  background: var(--fs-text);
  color: var(--fs-bg);
}

.option:disabled {
  color: var(--fs-faint);
  cursor: not-allowed;
}

.minute {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-8);
  color: var(--fs-muted);
  font-size: var(--fs-text-footnote);
}

.caption {
  display: flex;
  justify-content: space-between;
}

.value {
  color: var(--fs-text);
  font-family: var(--fs-font-mono);
  font-variant-numeric: tabular-nums;
}

input {
  width: 100%;
  accent-color: var(--fs-home);
}

input:disabled {
  opacity: 0.4;
}
</style>
