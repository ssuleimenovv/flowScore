<script setup lang="ts">
import { computed, useId } from 'vue'
import type { Prediction } from '@/shared/api/types'
import AppCard from '@/shared/ui/AppCard.vue'
import AppIcon from '@/shared/ui/AppIcon.vue'
import type { MatchLive } from '../matchState'
import { finalColumns, finalVerdict, formatPreMatch, outcomeColumns, resultOf } from '../outcome'

// "Вероятность исхода" from the Match boards. full: the desktop card next to
// the explanation; brief: the phone's flow tab, without the chip and footer.
// After the final whistle it is "Итог матча": the chances before kick-off,
// the result marked, and whether the model expected it.

const {
  live,
  prediction,
  variant = 'full',
} = defineProps<{
  live: MatchLive
  prediction: Prediction
  variant?: 'full' | 'brief'
}>()

const titleId = useId()
const finished = computed(() => live.match.status === 'finished')
const result = computed(() => (finished.value ? resultOf(live.score) : null))
const columns = computed(() =>
  finished.value
    ? finalColumns(live.match, prediction, live.score)
    : outcomeColumns(live.match, prediction),
)
// The bar draws the same chances as the columns
const shown = computed(() => (finished.value ? prediction.preMatch : prediction.current))
</script>

<template>
  <AppCard
    as="section"
    class="outcome fs-in"
    :class="`is-${variant}`"
    style="--fs-i: 4"
    :aria-labelledby="titleId"
  >
    <div class="head">
      <h3 :id="titleId" class="title">{{ finished ? 'Итог матча' : 'Вероятность исхода' }}</h3>
      <span v-if="variant === 'full'" class="chip">
        <AppIcon name="sparkle" :size="14" class="spark" />
        Prediction
      </span>
    </div>

    <!-- The three chances as one bar: home, draw, away -->
    <div class="bar" aria-hidden="true">
      <span class="part home" :style="{ width: `${shown.home}%` }" />
      <span class="part draw" :style="{ width: `${shown.draw}%` }" />

      <span class="part away" />
    </div>

    <dl class="columns">
      <div
        v-for="column in columns"
        :key="column.key"
        class="column"
        :class="[column.key, { missed: result && column.key !== result }]"
      >
        <dt class="label">{{ column.label }}</dt>
        <dd class="value">{{ column.value }}%</dd>
        <dd class="change">{{ column.change }}</dd>
      </div>
    </dl>

    <p v-if="finished" class="verdict">{{ finalVerdict(live.match, prediction, live.score) }}</p>

    <div v-if="variant === 'full'" class="foot">
      <span>{{ finished ? 'Шансы до матча' : formatPreMatch(prediction) }}</span>
      <span class="model">{{ prediction.model }}</span>
    </div>
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
  gap: var(--fs-space-8);
  align-items: center;
  justify-content: space-between;
}

.title {
  font-size: var(--fs-text-title-3);
  font-weight: 600;
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

.bar {
  display: flex;
  gap: 3px;
  height: 14px;
}

/* A part slides when the chances change; the away part takes what is left */
.part {
  border-radius: 4px;
  transition: width var(--fs-duration-grow) var(--fs-ease);
}

.part.home {
  background: var(--fs-home);
  transform-origin: left;
  animation: grow var(--fs-duration-bar-h) var(--fs-ease) both;
}

.part.draw {
  background: var(--fs-faint);
}

.part.away {
  flex-grow: 1;
  background: var(--fs-away);
}

/* Label first in the DOM for screen readers ("Ничья: 30%"), number first on screen */
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
  font-size: 40px;
  font-weight: 800;
  line-height: 1;
  font-variant-numeric: tabular-nums;
}

.change {
  order: 3;
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-caption);
  font-variant-numeric: tabular-nums;
  white-space: nowrap; /* "▲ 14 с начала" on one line, even in a phone's third */
}

/* After the match the outcomes that did not happen step back */
.column.missed {
  opacity: 0.4;
}

.verdict {
  font-size: var(--fs-text-subheadline);
  line-height: 1.45;
}


.home .value,
.home .change {
  color: var(--fs-home);
}

.draw .value,
.draw .change {
  color: var(--fs-muted);
}

.away .value,
.away .change {
  color: var(--fs-away);
}

.foot {
  display: flex;
  justify-content: space-between;
  padding-top: 14px;
  border-top: 1px solid var(--fs-line);
  color: var(--fs-muted);
  font-size: var(--fs-text-footnote);
}

.model {
  font-family: var(--fs-font-mono);
  font-variant-numeric: tabular-nums;
}

@keyframes grow {
  from {
    transform: scaleX(0);
  }
}

/* Phone: the MobileMatch card is tighter, with smaller numbers */
.is-brief {
  gap: 14px;
  padding: 18px var(--fs-space-16);
}

.is-brief .title {
  font-size: var(--fs-text-headline);
}

.is-brief .bar {
  height: 12px;
}

.is-brief .value {
  font-size: 32px;
}
</style>
