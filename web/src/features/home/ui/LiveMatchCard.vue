<script setup lang="ts">
import { computed } from 'vue'
import type { MatchSummary } from '@/shared/api/types'
import AppIcon from '@/shared/ui/AppIcon.vue'
import LiveBadge from '../../match/ui/LiveBadge.vue'
import TeamBadge from '../../match/ui/TeamBadge.vue'
import { flowShare, insightLine, liveMinute, sparkline } from '../home'

// A live match on the Main boards: score, the flow wave so far, who leads the
// flow and why. The whole card is the way into the match.
const { match, now } = defineProps<{ match: MatchSummary; now: number }>()

const minute = computed(() => liveMinute(match, now))
const share = computed(() => `${flowShare(match).toFixed(1)}%`)
const homeLine = computed(() => sparkline(match, 'home'))
const awayLine = computed(() => sparkline(match, 'away'))
const insight = computed(() => insightLine(match))
</script>

<template>
  <RouterLink :to="{ name: 'match', params: { matchId: match.id } }" class="card fs-card">
    <div class="top">
      <span class="league">{{ match.competition.name }}</span>
      <LiveBadge :status="match.status" :minute />
    </div>

    <div class="teams">
      <div class="team">
        <TeamBadge :code="match.home.code" side="home" size="sm" />
        <span class="name">{{ match.home.name }}</span>
        <span class="goals">{{ match.score.home }}</span>
      </div>
      <div class="team">
        <TeamBadge :code="match.away.code" side="away" size="sm" />
        <span class="name">{{ match.away.name }}</span>
        <span class="goals">{{ match.score.away }}</span>
      </div>
    </div>

    <svg class="wave" viewBox="0 0 300 60" preserveAspectRatio="none" aria-hidden="true">
      <polyline :points="awayLine" class="line away" />
      <polyline :points="homeLine" class="line home" />
    </svg>

    <div
      class="share"
      :aria-label="`Flow ${Math.round(match.flow.home)} : ${Math.round(match.flow.away)}`"
    >
      <span class="value home">{{ Math.round(match.flow.home) }}</span>
      <div class="bar" aria-hidden="true">
        <span class="part home" :style="{ width: share }" />
        <span class="part away" />
      </div>
      <span class="value away">{{ Math.round(match.flow.away) }}</span>
    </div>

    <div class="insight">
      <AppIcon name="sparkle" :size="14" class="spark" />
      <span>{{ insight }}</span>
    </div>
  </RouterLink>
</template>

<style scoped>
.card {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: 20px;
  border: 1px solid var(--fs-line);
  border-radius: var(--fs-radius-lg);
  background: var(--fs-surface);
  color: var(--fs-text);
  text-decoration: none;
}

.top {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.league {
  color: var(--fs-muted);
  font-size: var(--fs-text-footnote);
}

.teams {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.team {
  display: flex;
  gap: var(--fs-space-12);
  align-items: center;
}

.name {
  overflow: hidden;
  flex-grow: 1;
  font-size: var(--fs-text-subheadline);
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.goals {
  font-family: var(--fs-font-display);
  font-size: 30px;
  font-weight: 800;
  line-height: 1;
  font-variant-numeric: tabular-nums;
}

.wave {
  display: block;
  width: 100%;
  height: 56px;
}

.line {
  fill: none;
  stroke-width: 2;
  stroke-linecap: round;
  stroke-linejoin: round;
  vector-effect: non-scaling-stroke;
}

.line.home {
  stroke: var(--fs-home);
}

.line.away {
  stroke: var(--fs-away);
}

.share {
  display: flex;
  gap: 10px;
  align-items: center;
}

.value {
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-footnote);
  font-weight: 600;
  font-variant-numeric: tabular-nums;
}

.value.home {
  color: var(--fs-home);
}

.value.away {
  color: var(--fs-away);
}

.bar {
  display: flex;
  flex-grow: 1;
  gap: 2px;
  height: 6px;
}

.part {
  border-radius: var(--fs-radius-full);
}

.part.home {
  background: var(--fs-home);
  transition: width var(--fs-duration-grow) var(--fs-ease);
}

.part.away {
  flex-grow: 1;
  background: var(--fs-away);
}

.insight {
  display: flex;
  gap: var(--fs-space-8);
  align-items: center;
  padding-top: var(--fs-space-12);
  border-top: 1px solid var(--fs-line);
  color: var(--fs-muted);
  font-size: var(--fs-text-footnote);
}

.spark {
  flex-shrink: 0;
  color: var(--fs-home);
}

/* Phone: the MobileHome card is tighter, with larger names */
@media (max-width: 767px) {
  .card {
    padding: var(--fs-space-16);
  }

  .name {
    font-size: var(--fs-text-headline);
  }
}
</style>
