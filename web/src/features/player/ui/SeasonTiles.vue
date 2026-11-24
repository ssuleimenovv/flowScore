<script setup lang="ts">
import type { Tile } from '../player'

// The season in six numbers, as on PlayerSeason: one row on the desktop,
// three by two on the phone
const { tiles } = defineProps<{ tiles: Tile[] }>()
</script>

<template>
  <ul class="tiles fs-in">
    <li v-for="t in tiles" :key="t.label" class="tile">
      <span class="label">{{ t.label }}</span>
      <span class="value">{{ t.value }}</span>
      <span class="note">{{ t.note }}</span>
    </li>
  </ul>
</template>

<style scoped>
.tiles {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: var(--fs-space-8);
  margin: 0;
  padding: 0;
  list-style: none;
}

.tile {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
  padding: var(--fs-space-12);
  border: 1px solid var(--fs-line);
  border-radius: var(--fs-radius-lg);
  background: var(--fs-surface);
}

.label {
  color: var(--fs-muted);
  font-size: 11px;
  font-weight: 600;
}

.value {
  font-family: var(--fs-font-display);
  font-size: 28px;
  font-weight: 800;
  line-height: 1;
}

.note {
  color: var(--fs-muted);
  font-family: var(--fs-font-mono);
  font-size: 10px;
  font-variant-numeric: tabular-nums;
}

@media (min-width: 768px) {
  .tiles {
    grid-template-columns: repeat(6, minmax(0, 1fr));
    gap: var(--fs-space-12);
  }

  .tile {
    gap: var(--fs-space-4);
    padding: var(--fs-space-16) 18px;
  }

  .label,
  .note {
    font-size: var(--fs-text-caption);
  }

  .value {
    font-size: 40px;
  }
}
</style>
