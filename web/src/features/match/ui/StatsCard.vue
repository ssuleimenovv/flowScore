<script setup lang="ts">
import { computed } from 'vue'
import AppCard from '@/shared/ui/AppCard.vue'
import type { MatchLive } from '../matchState'
import { statLines } from '../stats'

// "02 Статистика матча" from the Match boards. Each side's bar grows toward
// the label in the middle, so the longer bar reads as "this team has more".
const { live } = defineProps<{ live: MatchLive }>()

const lines = computed(() => statLines(live.match.stats))
</script>

<template>
  <AppCard as="section" class="stats fs-in" style="--fs-i: 3" aria-labelledby="stats-title">
    <div class="head">
      <h2 id="stats-title" class="title"><span class="index">02</span>Статистика матча</h2>
      <div class="legend">
        <span class="team"><span class="swatch home" />{{ live.match.home.name }}</span>
        <span class="team"><span class="swatch away" />{{ live.match.away.name }}</span>
      </div>
    </div>

    <ul class="rows">
      <li v-for="line in lines" :key="line.key" class="row">
        <span class="label">{{ line.label }}</span>
        <span class="value home">
          <span class="fs-sr-only">{{ live.match.home.name }}: </span>{{ line.home }}
        </span>
        <span class="value away">
          <span class="fs-sr-only">{{ live.match.away.name }}: </span>{{ line.away }}
        </span>
        <div class="track home" aria-hidden="true">
          <div class="bar" :style="{ width: line.homeWidth }" />
        </div>
        <div class="track away" aria-hidden="true">
          <div class="bar" :style="{ width: line.awayWidth }" />
        </div>
      </li>
    </ul>
  </AppCard>
</template>

<style scoped>
.stats {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-16);
  padding: var(--fs-space-24) 28px;
}

.head {
  display: flex;
  gap: var(--fs-space-16);
  align-items: center;
  justify-content: space-between;
}

.title {
  font-size: 24px;
  font-weight: 700;
  line-height: 1.2;
  letter-spacing: -0.02em;
}

.index {
  margin-right: var(--fs-space-12);
  color: var(--fs-home);
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-footnote);
  vertical-align: middle;
}

.legend {
  display: flex;
  gap: var(--fs-space-16);
  color: var(--fs-muted);
  font-size: var(--fs-text-footnote);
}

.team {
  display: flex;
  gap: 6px;
  align-items: center;
}

.swatch {
  width: 10px;
  height: 10px;
  border-radius: 3px;
}

.swatch.home,
.track.home .bar {
  background: var(--fs-home);
}

.swatch.away,
.track.away .bar {
  background: var(--fs-away);
}

.rows {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-16);
  padding: 0;
  list-style: none;
}

/* One line on the board: 14 ━━━━ Удары ━━ 6 */
.row {
  display: grid;
  grid-template-columns: 56px minmax(0, 1fr) 170px minmax(0, 1fr) 56px;
  gap: var(--fs-space-16);
  align-items: center;
}

.value {
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-subheadline);
  font-variant-numeric: tabular-nums;
  font-weight: 600;
}

.value.home {
  grid-area: 1 / 1;
}

.value.away {
  grid-area: 1 / 5;
  text-align: right;
}

.label {
  grid-area: 1 / 3;
  color: var(--fs-muted);
  font-size: var(--fs-text-subheadline);
  text-align: center;
}

.track {
  display: flex;
  height: 8px;
  overflow: hidden;
  border-radius: var(--fs-radius-full);
  background: var(--fs-track);
}

.track.home {
  grid-area: 1 / 2;
  justify-content: flex-end;
}

.track.away {
  grid-area: 1 / 4;
}

/* The bars grow out of the label when the card appears, and slide to the new
   share when a number changes */
.bar {
  border-radius: var(--fs-radius-full);
  animation: grow var(--fs-duration-bar-h) var(--fs-ease) both;
  transition: width var(--fs-duration-grow) var(--fs-ease);
}

.track.home .bar {
  transform-origin: right;
}

.track.away .bar {
  transform-origin: left;
}

@keyframes grow {
  from {
    transform: scaleX(0);
  }
}

/* Phone: the numbers above, both bars under them, as on MobileMatch */
@media (max-width: 767px) {
  .stats {
    gap: 18px;
    padding: 18px var(--fs-space-16);
  }

  .head {
    display: none;
  }

  .rows {
    gap: 18px;
  }

  .row {
    grid-template-columns: 1fr 1fr;
    gap: var(--fs-space-8) var(--fs-space-4);
  }

  .value {
    font-size: var(--fs-text-body);
  }

  .value.home {
    grid-area: 1 / 1;
  }

  .value.away {
    grid-area: 1 / 2;
  }

  /* The label spans both columns and sits in the middle, between the numbers */
  .label {
    grid-area: 1 / 1 / 2 / 3;
    font-size: var(--fs-text-footnote);
  }

  .track {
    height: 6px;
  }

  .track.home {
    grid-area: 2 / 1;
  }

  .track.away {
    grid-area: 2 / 2;
  }
}
</style>
