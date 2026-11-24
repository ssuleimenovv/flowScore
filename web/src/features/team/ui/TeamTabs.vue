<script setup lang="ts">
import { nextTick } from 'vue'

export type TeamTab = 'overview' | 'squad' | 'fixtures'

// "Обзор · Состав · Расписание" from the Team board, underlined in the
// team's color. Keyboard works as in any tab list: the arrows move between tabs.
const model = defineModel<TeamTab>({ required: true })

const tabs: Array<{ id: TeamTab; label: string }> = [
  { id: 'overview', label: 'Обзор' },
  { id: 'squad', label: 'Состав' },
  { id: 'fixtures', label: 'Расписание' },
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
  <div class="tabs" role="tablist" aria-label="Разделы команды" @keydown="onKey">
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
/* Phone: from edge to edge of the screen, three equal tabs */
.tabs {
  display: flex;
  margin: 0 calc(-1 * var(--fs-screen-padding));
  padding: 0 var(--fs-space-8);
  border-bottom: 1px solid var(--fs-line);
}

.tab {
  flex: 1 1 0;
  height: 48px;
  margin-bottom: -1px;
  border: 0;
  border-bottom: 3px solid transparent;
  background: transparent;
  color: var(--fs-muted);
  font-size: var(--fs-text-subheadline);
  font-weight: 600;
  cursor: pointer;
  transition:
    transform var(--fs-duration-press) ease,
    color var(--fs-duration-toggle) var(--fs-ease-switch),
    border-color var(--fs-duration-toggle) var(--fs-ease-switch);
}

.tab[aria-selected='true'] {
  border-bottom-color: var(--fs-home);
  color: var(--fs-text);
}

/* Desktop: the tabs keep their width, 28 px apart */
@media (min-width: 768px) {
  .tabs {
    gap: 28px;
    margin: 0;
    padding: 0;
  }

  .tab {
    flex: none;
    padding: 0 var(--fs-space-4);
  }
}
</style>
