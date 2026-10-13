<script setup lang="ts">
import { computed } from 'vue'
import type { FlowRecord } from '@/shared/api/types'
import AppCard from '@/shared/ui/AppCard.vue'
import { recordRows } from '../league'

// "Рекорды Flow" from the League board: the season's matches where a team's
// Flow was highest on average over the whole match.
const { records } = defineProps<{ records: FlowRecord[] }>()
const rows = computed(() => recordRows(records))
</script>

<template>
  <AppCard as="section" class="records fs-in" style="--fs-i: 3" aria-labelledby="records-title">
    <div class="head">
      <h2 id="records-title" class="title">Рекорды Flow</h2>
      <span class="hint">сезон</span>
    </div>
    <div v-for="r in rows" :key="r.id" class="row">
      <span class="rank" aria-hidden="true">{{ r.rank }}</span>
      <div class="what">
        <span class="match">{{ r.match }}</span>
        <span class="round">{{ r.round }}</span>
      </div>
      <span class="flow" :aria-label="`средний Flow ${r.flow}`">{{ r.flow }}</span>
    </div>
    <p v-if="rows.length === 0" class="none">Рекорды появятся после первых матчей</p>
  </AppCard>
</template>

<style scoped>
.records {
  display: flex;
  flex-direction: column;
  padding: 20px var(--fs-space-24);
}

.head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: 10px;
}

.title {
  font-size: var(--fs-text-title-3);
  font-weight: 700;
  letter-spacing: -0.01em;
  line-height: 1.2;
}

.hint {
  color: var(--fs-muted);
  font-size: var(--fs-text-caption);
}

.row {
  display: grid;
  grid-template-columns: 36px minmax(0, 1fr) auto;
  gap: var(--fs-space-12);
  align-items: center;
  padding: var(--fs-space-12) 0;
  border-top: 1px solid var(--fs-line);
}

.rank {
  color: var(--fs-faint);
  font-family: var(--fs-font-display);
  font-size: 26px;
  font-weight: 800;
}

.what {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.match {
  font-size: var(--fs-text-subheadline);
  font-weight: 600;
}

.round {
  color: var(--fs-muted);
  font-size: var(--fs-text-caption);
}

.flow {
  color: var(--fs-home);
  font-family: var(--fs-font-display);
  font-size: 32px;
  font-weight: 800;
}

.none {
  padding: var(--fs-space-12) 0;
  border-top: 1px solid var(--fs-line);
  color: var(--fs-muted);
  font-size: var(--fs-text-subheadline);
}
</style>
