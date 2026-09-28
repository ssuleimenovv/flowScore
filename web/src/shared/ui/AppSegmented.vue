<script setup lang="ts" generic="T extends string | number">
// The segmented control from the mockup ("1 мин · 5 мин · 15 мин").
const model = defineModel<T>({ required: true })
const { options, label } = defineProps<{
  options: Array<{ value: T; label: string }>
  label: string
}>()
</script>

<template>
  <div class="segmented" role="group" :aria-label="label">
    <button
      v-for="option in options"
      :key="option.value"
      type="button"
      class="option fs-press"
      :aria-pressed="option.value === model"
      @click="model = option.value"
    >
      {{ option.label }}
    </button>
  </div>
</template>

<style scoped>
.segmented {
  display: flex;
  flex-shrink: 0;
  gap: 2px;
  padding: 2px;
  border-radius: 10px;
  background: var(--fs-fill);
}

.option {
  height: 36px;
  padding: 0 var(--fs-space-12);
  border: 0;
  border-radius: 9px;
  background: transparent;
  color: var(--fs-muted);
  font-size: var(--fs-text-footnote);
  font-weight: 600;
  white-space: nowrap;
  cursor: pointer;
  transition:
    transform var(--fs-duration-press) ease,
    background-color var(--fs-duration-toggle) var(--fs-ease-switch),
    color var(--fs-duration-toggle) var(--fs-ease-switch);
}

.option[aria-pressed='true'] {
  background: var(--fs-seg-selected);
  color: var(--fs-text);
}
</style>
