<script setup lang="ts">
import type { MatchSummary } from '@/shared/api/types'
import { kickoffTime } from '../home'

// "Позже сегодня" from the Main boards: kick-off time, the teams and the
// model's chances before the match. The phone shows time and teams only.
const { matches } = defineProps<{ matches: MatchSummary[] }>()
</script>

<template>
  <div class="list fs-in" style="--fs-i: 3">
    <RouterLink
      v-for="m in matches"
      :key="m.id"
      :to="{ name: 'match', params: { matchId: m.id } }"
      class="row fs-link"
    >
      <span class="time">{{ kickoffTime(m) }}</span>
      <div class="teams">
        <span class="names">{{ m.home.name }} — {{ m.away.name }}</span>
        <span class="league">{{ m.competition.name }}</span>
      </div>
      <div v-if="m.prediction" class="chances">
        <div class="bar" aria-hidden="true">
          <span class="part home" :style="{ width: `${m.prediction.current.home}%` }" />
          <span class="part draw" :style="{ width: `${m.prediction.current.draw}%` }" />
          <span class="part away" />
        </div>
        <div class="numbers">
          <span>П1 {{ m.prediction.current.home }}%</span>
          <span>X {{ m.prediction.current.draw }}%</span>
          <span>П2 {{ m.prediction.current.away }}%</span>
        </div>
      </div>
    </RouterLink>
    <p v-if="matches.length === 0" class="none">Больше матчей сегодня нет</p>
  </div>
</template>

<style scoped>
.list {
  display: flex;
  flex-direction: column;
  padding: 4px var(--fs-space-24);
  border: 1px solid var(--fs-line);
  border-radius: var(--fs-radius-lg);
  background: var(--fs-surface);
}

.row,
.none {
  border-top: 1px solid var(--fs-line);
}

/* The first line of a card has nothing above it to separate from */
.row:first-child,
.none:first-child {
  border-top: 0;
}

.row {
  display: grid;
  grid-template-columns: 64px minmax(0, 1fr) 260px;
  gap: var(--fs-space-16);
  align-items: center;
  padding: 14px 0;
  color: var(--fs-text);
  text-decoration: none;
}

.time {
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-subheadline);
  font-weight: 600;
  font-variant-numeric: tabular-nums;
}

.teams {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.names {
  overflow: hidden;
  font-size: var(--fs-text-subheadline);
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.league {
  color: var(--fs-muted);
  font-size: var(--fs-text-caption);
}

.chances {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.bar {
  display: flex;
  gap: 2px;
  height: 8px;
}

.part {
  border-radius: 3px;
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

.none {
  padding: 20px 0;
  color: var(--fs-muted);
  font-size: var(--fs-text-subheadline);
}

/* Phone: time and teams, as on MobileHome */
@media (max-width: 767px) {
  .list {
    padding: 0 var(--fs-space-16);
  }

  .row {
    grid-template-columns: 48px minmax(0, 1fr);
    gap: var(--fs-space-12);
  }

  .time,
  .names {
    font-size: var(--fs-text-headline);
  }

  .league {
    font-size: var(--fs-text-footnote);
  }

  .chances {
    display: none;
  }
}
</style>
