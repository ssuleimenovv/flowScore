<script setup lang="ts">
import type { MatchSummary } from '@/shared/api/types'
import AppButton from '@/shared/ui/AppButton.vue'
import AppCard from '@/shared/ui/AppCard.vue'
import AppIcon from '@/shared/ui/AppIcon.vue'
import { kickoffTime } from '../home'

// "AI-прогноз дня" from the Main board: the next match the outcome model is
// surest about, and the way into its simulator.
const { match, text } = defineProps<{ match: MatchSummary; text: string }>()
</script>

<template>
  <AppCard as="section" class="forecast fs-in" style="--fs-i: 5" aria-labelledby="forecast-title">
    <div class="head">
      <span class="chip">
        <AppIcon name="sparkle" :size="14" class="spark" />
        AI-прогноз дня
      </span>
      <span class="time">{{ kickoffTime(match) }}</span>
    </div>
    <div class="teams">
      <h2 id="forecast-title" class="names">{{ match.home.name }} — {{ match.away.name }}</h2>
      <span class="league">
        {{ match.competition.name }}
        <template v-if="match.competition.round">· {{ match.competition.round }}-й тур</template>
      </span>
    </div>
    <template v-if="match.prediction">
      <div class="bar" aria-hidden="true">
        <span class="part home" :style="{ width: `${match.prediction.current.home}%` }" />
        <span class="part draw" :style="{ width: `${match.prediction.current.draw}%` }" />
        <span class="part away" />
      </div>
      <div class="numbers">
        <span>П1 {{ match.prediction.current.home }}%</span>
        <span>X {{ match.prediction.current.draw }}%</span>
        <span>П2 {{ match.prediction.current.away }}%</span>
      </div>
    </template>
    <p class="text">{{ text }}</p>
    <AppButton
      variant="secondary"
      :to="{ name: 'simulator', params: { matchId: match.id } }"
      class="go"
    >
      Смоделировать сценарий
      <AppIcon name="arrow" :size="18" />
    </AppButton>
  </AppCard>
</template>

<style scoped>
.forecast {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-16);
  padding: var(--fs-space-24);
}

.head {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.chip {
  display: inline-flex;
  gap: 6px;
  align-items: center;
  height: 26px;
  padding: 0 10px;
  border-radius: var(--fs-radius-sm);
  background: var(--fs-surface-2);
  font-size: var(--fs-text-caption);
  font-weight: 600;
  white-space: nowrap;
}

.spark {
  color: var(--fs-home);
}

.time {
  color: var(--fs-muted);
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-footnote);
  font-variant-numeric: tabular-nums;
}

.teams {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.names {
  font-size: var(--fs-text-title-3);
  font-weight: 600;
}

.league {
  color: var(--fs-muted);
  font-size: var(--fs-text-footnote);
}

.bar {
  display: flex;
  gap: 3px;
  height: 10px;
}

.part {
  border-radius: 4px;
}

.part.home {
  background: var(--fs-home);
}

.part.draw {
  background: var(--fs-faint);
}

.part.away {
  flex-grow: 1;
  background: var(--fs-away);
}

.numbers {
  display: flex;
  justify-content: space-between;
  color: var(--fs-muted);
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-caption);
  font-variant-numeric: tabular-nums;
}

.text {
  color: var(--fs-muted);
  font-size: var(--fs-text-subheadline);
  line-height: 1.5;
}

.go {
  width: 100%;
}
</style>
