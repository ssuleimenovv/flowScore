<script setup lang="ts">
import type { RoundMatch } from '@/shared/api/types'
import AppCard from '@/shared/ui/AppCard.vue'
import AppSegmented from '@/shared/ui/AppSegmented.vue'
import AppSkeleton from '@/shared/ui/AppSkeleton.vue'
import { pairing, roundStatus, scoreText } from '../league'

// "Матчи тура" from the League boards: the round tabs and the round's matches.
// A match under way opens its screen; the gateway follows only those live.
const round = defineModel<number>({ required: true })
const { tabs, matches, loading, now } = defineProps<{
  tabs: number[]
  matches: RoundMatch[]
  loading: boolean
  now: Date
}>()

defineEmits<{ retry: [] }>()
</script>

<template>
  <AppCard as="section" class="round fs-in" style="--fs-i: 2" aria-labelledby="round-title">
    <h2 id="round-title" class="title">Матчи тура</h2>
    <AppSegmented
      v-model="round"
      class="tabs"
      label="Тур"
      :options="tabs.map((r) => ({ value: r, label: `Тур ${r}` }))"
    />
    <div v-if="loading" class="list" aria-busy="true" aria-label="Загрузка матчей тура">
      <div v-for="i in 3" :key="i" class="row">
        <div class="what">
          <AppSkeleton width="70%" :height="15" />
          <AppSkeleton width="30%" :height="12" />
        </div>
      </div>
    </div>
    <div v-else class="list">
      <component
        :is="m.state === 'in' ? 'RouterLink' : 'div'"
        v-for="m in matches"
        :key="m.id"
        v-bind="m.state === 'in' ? { to: { name: 'match', params: { matchId: m.id } } } : {}"
        class="row"
        :class="{ 'fs-link': m.state === 'in' }"
      >
        <div class="what">
          <span class="match">{{ pairing(m) }}</span>
          <span class="score">{{ scoreText(m) }}</span>
        </div>
        <span v-if="m.state === 'in'" class="live">
          <span class="dot" aria-hidden="true" />LIVE {{ roundStatus(m, now) }}
        </span>
        <span v-else class="status">{{ roundStatus(m, now) }}</span>
      </component>
      <p v-if="matches.length === 0" class="none">В этом туре матчей нет</p>
    </div>
  </AppCard>
</template>

<style scoped>
.round {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: var(--fs-space-16);
}

.title {
  font-size: var(--fs-text-title-3);
  font-weight: 700;
  letter-spacing: -0.01em;
  line-height: 1.2;
}

.tabs :deep(.option) {
  flex: 1 1 0;
}

.list {
  display: flex;
  flex-direction: column;
}

.row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: var(--fs-space-12);
  align-items: center;
  padding: var(--fs-space-12) 0;
  border-top: 1px solid var(--fs-line);
  color: var(--fs-text);
  text-decoration: none;
}

.what {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.match {
  font-size: var(--fs-text-headline);
  font-weight: 600;
}

.score,
.status {
  color: var(--fs-muted);
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-footnote);
  font-variant-numeric: tabular-nums;
}

.live {
  display: inline-flex;
  gap: 6px;
  align-items: center;
  height: 24px;
  padding: 0 10px;
  border-radius: var(--fs-radius-full);
  background: var(--fs-live-soft);
  color: var(--fs-live);
  font-size: var(--fs-text-caption);
  font-weight: 700;
  letter-spacing: 0.08em;
  white-space: nowrap;
}

.dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--fs-live);
  animation: pulse var(--fs-duration-pulse) ease-in-out infinite;
}

.none {
  padding: var(--fs-space-12) 0;
  border-top: 1px solid var(--fs-line);
  color: var(--fs-muted);
  font-size: var(--fs-text-subheadline);
}

@keyframes pulse {
  50% {
    opacity: 0.4;
  }
}

/* Desktop, as on League: a little more room and the smaller names */
@media (min-width: 768px) {
  .round {
    gap: var(--fs-space-12);
    padding: 20px var(--fs-space-24);
  }

  .tabs :deep(.option) {
    height: 36px;
  }

  .match {
    font-size: var(--fs-text-subheadline);
  }

  .score,
  .status {
    font-size: var(--fs-text-caption);
  }
}
</style>
