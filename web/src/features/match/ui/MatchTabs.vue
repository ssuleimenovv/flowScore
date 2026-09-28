<script setup lang="ts">
import { nextTick } from 'vue'

export type MatchTab = 'flow' | 'events' | 'stats' | 'ai'

// "Поток · События · Стат. · AI" from MobileMatch. Keyboard works as in any
// tab list: the arrows move between tabs, Tab jumps into the open panel.
const model = defineModel<MatchTab>({ required: true })

const tabs: Array<{ id: MatchTab; label: string }> = [
  { id: 'flow', label: 'Поток' },
  { id: 'events', label: 'События' },
  { id: 'stats', label: 'Стат.' },
  { id: 'ai', label: 'AI' },
]

function onKey(e: KeyboardEvent) {
  const step = e.key === 'ArrowRight' ? 1 : e.key === 'ArrowLeft' ? -1 : 0
  if (step === 0) return
  e.preventDefault()
  const current = tabs.findIndex((t) => t.id === model.value)
  const next = tabs[(current + step + tabs.length) % tabs.length]!
  model.value = next.id
  void nextTick(() => document.getElementById(`tab-${next.id}`)?.focus())
}
</script>

<template>
  <div class="tabs" role="tablist" aria-label="Разделы матча" @keydown="onKey">
    <button
      v-for="t in tabs"
      :id="`tab-${t.id}`"
      :key="t.id"
      type="button"
      role="tab"
      class="tab fs-press"
      :aria-selected="t.id === model"
      :aria-controls="`panel-${t.id}`"
      :tabindex="t.id === model ? 0 : -1"
      @click="model = t.id"
    >
      {{ t.label }}
    </button>
  </div>
</template>

<style scoped>
.tabs {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 2px;
  padding: 2px;
  border-radius: 11px;
  background: var(--fs-fill);
}

.tab {
  height: 40px;
  border: 0;
  border-radius: 10px;
  background: transparent;
  color: var(--fs-muted);
  font-size: var(--fs-text-footnote);
  font-weight: 600;
  cursor: pointer;
  transition:
    transform var(--fs-duration-press) ease,
    background-color var(--fs-duration-toggle) var(--fs-ease-switch),
    color var(--fs-duration-toggle) var(--fs-ease-switch);
}

.tab[aria-selected='true'] {
  background: var(--fs-seg-selected);
  color: var(--fs-text);
}
</style>
