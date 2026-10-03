<script setup lang="ts">
import { computed } from 'vue'
import type { Match, Probabilities, Simulation } from '@/shared/api/types'
import AppCard from '@/shared/ui/AppCard.vue'
import AppIcon from '@/shared/ui/AppIcon.vue'
import { formatChances, scenarioColumns } from '../scenario'

// "Исход матча" from the Simulator boards: the chances as the match stands,
// the scenario's under them, and how far each result moved.
const { match, simulation } = defineProps<{ match: Match; simulation: Simulation }>()

const columns = computed(() => scenarioColumns(match, simulation))
const width = (p: Probabilities, key: 'home' | 'draw') => `${p[key]}%`
</script>

<template>
  <AppCard as="section" class="outcome fs-in" style="--fs-i: 1" aria-labelledby="outcome-title">
    <div class="head">
      <h2 id="outcome-title" class="title">Исход матча</h2>
      <span class="chip">
        <AppIcon name="sparkle" :size="14" class="spark" />
        Prediction Agent
      </span>
    </div>

    <div class="row">
      <div class="caption">
        <span>Сейчас</span>
        <span class="mono">{{ formatChances(simulation.current) }}</span>
      </div>
      <div class="bar now" aria-hidden="true">
        <span class="part home" :style="{ width: width(simulation.current, 'home') }" />
        <span class="part draw" :style="{ width: width(simulation.current, 'draw') }" />
        <span class="part away" />
      </div>
    </div>

    <div class="row">
      <div class="caption scenario">
        <span>Сценарий</span>
        <span class="mono">{{ formatChances(simulation.scenario) }}</span>
      </div>
      <div class="bar" aria-hidden="true">
        <span class="part home" :style="{ width: width(simulation.scenario, 'home') }" />
        <span class="part draw" :style="{ width: width(simulation.scenario, 'draw') }" />
        <span class="part away" />
      </div>
    </div>

    <!-- Label first for screen readers ("Ничья: 34%"), number first on screen -->
    <dl class="columns" aria-live="polite">
      <div v-for="column in columns" :key="column.key" class="column" :class="column.key">
        <dt class="label">{{ column.label }}</dt>
        <dd class="value">{{ column.value }}%</dd>
        <dd class="change">{{ column.change }}</dd>
      </div>
    </dl>
  </AppCard>
</template>

<style scoped>
.outcome {
  display: flex;
  flex-direction: column;
  gap: 18px;
  padding: var(--fs-space-24) 28px;
}

.head {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.title {
  font-size: 24px;
  font-weight: 700;
  letter-spacing: -0.02em;
  line-height: 1.2;
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

.row {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-8);
}

.caption {
  display: flex;
  justify-content: space-between;
  color: var(--fs-muted);
  font-size: var(--fs-text-footnote);
}

.caption.scenario {
  color: var(--fs-text);
}

.caption.scenario span:first-child {
  font-weight: 600;
}

.mono {
  font-family: var(--fs-font-mono);
  font-variant-numeric: tabular-nums;
}

.bar {
  display: flex;
  gap: 3px;
  height: 14px;
}

.bar.now {
  height: 10px;
}

/* The scenario bar slides to the new chances; the away part takes the rest */
.part {
  border-radius: 4px;
  transition: width var(--fs-duration-grow) var(--fs-ease);
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

.columns {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: var(--fs-space-12);
}

.column {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.label {
  order: 2;
  color: var(--fs-muted);
  font-size: var(--fs-text-footnote);
}

.value {
  order: 1;
  font-family: var(--fs-font-display);
  font-size: 44px;
  font-weight: 800;
  line-height: 1;
  font-variant-numeric: tabular-nums;
}

.change {
  order: 3;
  color: var(--fs-muted);
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-caption);
  font-variant-numeric: tabular-nums;
}

.home .value {
  color: var(--fs-home);
}

.away .value {
  color: var(--fs-away);
}

/* Phone: the MobileSimulator card opens with the bars, no heading row */
@media (max-width: 767px) {
  .outcome {
    gap: 14px;
    padding: 18px var(--fs-space-16);
  }

  .head {
    display: none;
  }

  .value {
    font-size: 34px;
  }

  .change {
    font-size: var(--fs-text-footnote);
  }
}
</style>
