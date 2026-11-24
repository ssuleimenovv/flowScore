<script setup lang="ts">
import { RouterLink } from 'vue-router'
import LiveBadge from '@/features/match/ui/LiveBadge.vue'
import type { TeamMatch } from '@/shared/api/types'
import { dayMonth, kickoffLine, opponentLine, scoreFor, venueOf, type Result } from '../team'
import ResultChip from './ResultChip.vue'

// "Ближайшие матчи" and "Результаты" from the Team board. The one under way
// is a live pill that leads to the match; the results begin with «В», «Н», «П».
const { kind, title, phoneTitle, matches, teamId, leagueId, now, empty } = defineProps<{
  kind: 'upcoming' | 'results'
  title: string
  phoneTitle?: string // MobileTeam says just «Ближайшие»
  matches: TeamMatch[]
  teamId: string
  leagueId: string | null
  now: Date
  empty: string
}>()

function live(m: TeamMatch): string {
  return `${scoreFor(m, teamId)} · ${m.clock.replace(/'/g, '′')}`
}
</script>

<template>
  <section class="card fs-in" :class="kind">
    <h2 class="title">
      <span :class="{ 'desktop-word': phoneTitle }">{{ title }}</span>
      <span v-if="phoneTitle" class="phone-word">{{ phoneTitle }}</span>
    </h2>
    <p v-if="matches.length === 0" class="row empty">{{ empty }}</p>
    <div v-for="m in matches" :key="m.id" class="row">
      <ResultChip v-if="kind === 'results'" :result="m.result as Result" />
      <div class="what">
        <span class="opponent">{{ opponentLine(m, teamId, leagueId) }}</span>
        <span class="when">
          <template v-if="kind === 'results'">{{ dayMonth(new Date(m.kickoffAt)) }}</template>
          <template v-else>{{ kickoffLine(m, now) }} · {{ venueOf(m, teamId) }}</template>
        </span>
      </div>
      <RouterLink v-if="m.state === 'in'" :to="`/matches/${m.id}`" class="fs-link">
        <LiveBadge status="live" :minute="live(m)" />
      </RouterLink>
      <span v-else class="score">{{ kind === 'results' ? scoreFor(m, teamId) : '—' }}</span>
    </div>
  </section>
</template>

<style scoped>
.card {
  display: flex;
  flex-direction: column;
  padding: 18px var(--fs-space-16);
  border: 1px solid var(--fs-line);
  border-radius: var(--fs-radius-lg);
  background: var(--fs-surface);
}

.title {
  padding-bottom: var(--fs-space-8);
  font-size: 20px;
  font-weight: 700;
  letter-spacing: -0.01em;
  line-height: 1.2;
}

.desktop-word {
  display: none;
}

.row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: var(--fs-space-12);
  align-items: center;
  padding: var(--fs-space-12) 0;
  border-top: 1px solid var(--fs-line);
}

.results .row {
  grid-template-columns: 28px minmax(0, 1fr) auto;
  padding: 10px 0;
}

.what {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.opponent {
  overflow: hidden;
  font-size: var(--fs-text-subheadline);
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.when {
  color: var(--fs-muted);
  font-size: var(--fs-text-caption);
}

.score {
  color: var(--fs-muted);
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-footnote);
  font-variant-numeric: tabular-nums;
}

.results .score {
  color: var(--fs-text);
  font-size: var(--fs-text-subheadline);
  font-weight: 600;
}

.empty {
  display: block;
  color: var(--fs-muted);
  font-size: var(--fs-text-footnote);
}

@media (min-width: 768px) {
  .card {
    padding: 20px 28px;
  }

  .results {
    padding: 20px var(--fs-space-24);
  }

  .desktop-word {
    display: inline;
  }

  .phone-word {
    display: none;
  }
}
</style>
