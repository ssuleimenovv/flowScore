<script setup lang="ts">
import AppCard from '@/shared/ui/AppCard.vue'
import type { Peak } from '../home'

// "Пики потока" from the Main boards: the live matches where a team's Flow
// rose the most over the last five minutes.
const { peaks } = defineProps<{ peaks: Peak[] }>()
</script>

<template>
  <AppCard as="section" class="peaks fs-in" style="--fs-i: 4" aria-labelledby="peaks-title">
    <div class="head">
      <h2 id="peaks-title" class="title">Пики потока</h2>
      <span class="hint">за 5 минут</span>
    </div>
    <RouterLink
      v-for="(peak, i) in peaks"
      :key="peak.id"
      :to="{ name: 'match', params: { matchId: peak.id } }"
      class="row fs-link"
    >
      <span class="rank">{{ String(i + 1).padStart(2, '0') }}</span>
      <div class="what">
        <span class="match">{{ peak.match }}</span>
        <span class="rise" :class="peak.side">+{{ peak.rise }} за 5′</span>
      </div>
      <span class="flow" :class="peak.side">{{ peak.flow }}</span>
    </RouterLink>
    <p v-if="peaks.length === 0" class="none">Скачков потока сейчас нет</p>
  </AppCard>
</template>

<style scoped>
.peaks {
  display: flex;
  flex-direction: column;
  padding: 22px var(--fs-space-24);
}

.head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: var(--fs-space-12);
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
  grid-template-columns: 40px minmax(0, 1fr) auto;
  gap: var(--fs-space-12);
  align-items: center;
  padding: var(--fs-space-12) 0;
  border-top: 1px solid var(--fs-line);
  color: var(--fs-text);
  text-decoration: none;
}

.rank {
  color: var(--fs-faint);
  font-family: var(--fs-font-display);
  font-size: 28px;
  font-weight: 800;
}

.what {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.match {
  overflow: hidden;
  font-size: var(--fs-text-subheadline);
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.rise {
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-caption);
  font-variant-numeric: tabular-nums;
}

.flow {
  font-family: var(--fs-font-display);
  font-size: 34px;
  font-weight: 800;
  line-height: 1;
}

.home {
  color: var(--fs-home);
}

.away {
  color: var(--fs-away);
}

.none {
  padding: var(--fs-space-12) 0;
  border-top: 1px solid var(--fs-line);
  color: var(--fs-muted);
  font-size: var(--fs-text-subheadline);
}

@media (max-width: 767px) {
  .peaks {
    padding: 18px var(--fs-space-16);
  }

  .match {
    font-size: var(--fs-text-headline);
  }

  .rise {
    font-size: var(--fs-text-footnote);
  }
}
</style>
