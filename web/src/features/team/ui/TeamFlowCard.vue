<script setup lang="ts">
import { computed } from 'vue'
import type { TeamMatch } from '@/shared/api/types'
import { lastFlows } from '../team'

// "Flow по матчам" from the Team board: the team's Flow in its last ten
// matches read, oldest first. A match of 70 or more stands out in its color.
const { matches, teamId } = defineProps<{ matches: TeamMatch[]; teamId: string }>()

const bars = computed(() =>
  lastFlows(matches, teamId).map((b) => ({ ...b, value: Math.round(b.flow) })),
)
</script>

<template>
  <section class="card fs-in">
    <div class="head">
      <h2 class="title">Flow по матчам</h2>
      <span class="note">последние {{ bars.length }} · красным ≥ 70</span>
    </div>
    <ol v-if="bars.length > 0" class="bars">
      <li
        v-for="(b, i) in bars"
        :key="b.matchId"
        class="bar"
        :style="{ '--v': b.value, '--i': i }"
        :aria-label="`${b.opponent.name}: Flow ${b.value}`"
      >
        <span class="label" aria-hidden="true">{{ b.value }}</span>
        <span class="fill" :class="{ high: b.value >= 70 }" />
        <span class="label" aria-hidden="true">{{ b.opponent.code }}</span>
      </li>
    </ol>
    <p v-else class="empty">Flow появится после первого сыгранного матча.</p>
  </section>
</template>

<style scoped>
/* Phone, as on MobileTeam */
.card {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: 18px var(--fs-space-16);
  border: 1px solid var(--fs-line);
  border-radius: var(--fs-radius-lg);
  background: var(--fs-surface);
}

.head {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.title {
  font-size: 20px;
  font-weight: 700;
  letter-spacing: -0.01em;
  line-height: 1.2;
}

.note {
  display: none;
  color: var(--fs-muted);
  font-size: var(--fs-text-footnote);
}

.bars {
  display: flex;
  gap: var(--fs-space-4);
  align-items: flex-end;
  height: 130px;
  margin: 0;
  padding: 0;
  list-style: none;
}

.bar {
  display: flex;
  flex: 1 1 0;
  flex-direction: column;
  gap: 6px;
  align-items: center;
  justify-content: flex-end;
  min-width: 0;
  height: 100%;
}

.label {
  color: var(--fs-muted);
  font-family: var(--fs-font-mono);
  font-size: 10px;
  font-variant-numeric: tabular-nums;
}

/* The bars grow out of the floor one after another, as on the board */
.fill {
  width: 100%;
  height: calc(var(--v) * 0.9px);
  border-radius: 6px 6px 2px 2px;
  background: var(--fs-faint);
  transform-origin: bottom;
  animation: grow var(--fs-duration-bar) var(--fs-ease) both;
  animation-delay: calc(var(--i) * var(--fs-stagger));
}

.high {
  background: var(--fs-home);
}

@keyframes grow {
  from {
    transform: scaleY(0);
  }
}

.empty {
  color: var(--fs-muted);
  font-size: var(--fs-text-footnote);
}

/* Desktop, as on Team */
@media (min-width: 768px) {
  .card {
    gap: 18px;
    padding: var(--fs-space-24) 28px;
  }

  .title {
    font-size: 24px;
    letter-spacing: -0.02em;
  }

  .note {
    display: inline;
  }

  .bars {
    gap: 10px;
    height: 200px;
  }

  .label {
    font-size: var(--fs-text-caption);
  }

  .fill {
    height: calc(var(--v) * 1.5px);
  }
}
</style>
