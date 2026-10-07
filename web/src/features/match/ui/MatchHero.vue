<script setup lang="ts">
import { computed } from 'vue'
import { useTweenedNumber } from '@/shared/motion/useTweenedNumber'
import AppCard from '@/shared/ui/AppCard.vue'
import { flowShare, formatClock, formatDelta, formatMeta, formatMinute } from '../format'
import type { MatchLive } from '../matchState'
import FlowShareBar from './FlowShareBar.vue'
import LiveBadge from './LiveBadge.vue'
import TeamBadge from './TeamBadge.vue'

// The scoreboard card from the desktop Match board.
const { live, seconds } = defineProps<{ live: MatchLive; seconds: number }>()

const match = computed(() => live.match)
const flowHome = useTweenedNumber(computed(() => live.flow.home))
const flowAway = useTweenedNumber(computed(() => live.flow.away))
const share = computed(() => flowShare(live.flow))
const minute = computed(() => formatMinute(seconds, match.value.clock.period))

// "72:14 · 1-й тайм 0:1": the running clock and the first-half score
const clockLine = computed(() => {
  const { status, halftimeScore } = match.value
  if (status !== 'live' && status !== 'halftime') return null
  const half = halftimeScore ? `1-й тайм ${halftimeScore.home}:${halftimeScore.away}` : null
  return [formatClock(seconds), half].filter(Boolean).join(' · ')
})
</script>

<template>
  <AppCard as="section" class="hero fs-in" style="--fs-i: 1" aria-label="Счёт и поток">
    <div class="meta">
      <LiveBadge :status="match.status" :minute />
      <span>{{ formatMeta(match) }}</span>
    </div>

    <div class="scoreboard">
      <div class="team">
        <TeamBadge :code="match.home.code" side="home" />
        <div class="team-body">
          <span class="name">{{ match.home.name }}</span>
          <div class="flow-line">
            <span class="flow home">{{ flowHome }}</span>
            <span class="flow-label">FLOW</span>
            <span class="delta home">{{ formatDelta(live.delta10.home) }} за 10′</span>
          </div>
        </div>
      </div>

      <div class="center">
        <div class="score" :aria-label="`Счёт ${live.score.home}:${live.score.away}`">
          {{ live.score.home }}<span class="colon" aria-hidden="true">:</span>{{ live.score.away }}
        </div>
        <span v-if="clockLine" class="clock">{{ clockLine }}</span>
      </div>

      <div class="team away-team">
        <TeamBadge :code="match.away.code" side="away" />
        <div class="team-body end">
          <span class="name">{{ match.away.name }}</span>
          <div class="flow-line away-line">
            <span class="flow away end">{{ flowAway }}</span>
            <span class="flow-label">FLOW</span>
            <span class="delta away">{{ formatDelta(live.delta10.away) }} за 10′</span>
          </div>
        </div>
      </div>
    </div>

    <div class="share">
      <div class="share-head">
        <span>Доля потока за последние 10 минут</span>
        <span class="mono">{{ share.home }} / {{ share.away }}</span>
      </div>
      <FlowShareBar :home="share.home" />
    </div>
  </AppCard>
</template>

<style scoped>
.hero {
  display: flex;
  flex-direction: column;
  gap: 22px;
  padding: var(--fs-space-24) 28px;
}

.meta {
  display: flex;
  gap: var(--fs-space-12);
  align-items: center;
  color: var(--fs-muted);
  font-size: var(--fs-text-subheadline);
}

.scoreboard {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto minmax(0, 1fr);
  gap: var(--fs-space-24);
  align-items: center;
}

.team {
  display: flex;
  gap: 18px;
  align-items: center;
}

.away-team {
  flex-direction: row-reverse;
}

.team-body {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-8);
  min-width: 0;
}

.team-body.end {
  align-items: flex-end;
}

.name {
  overflow: hidden;
  font-size: 24px;
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* A long delta ("+47 за 10′" next to 3:3) goes under the number rather
   than into the score in the middle */
.flow-line {
  display: flex;
  flex-wrap: wrap;
  gap: 6px 10px;
  align-items: baseline;
}

/* The away side reads from the edge of the card inwards, as on the board */
.away-line {
  flex-direction: row-reverse;
}

.flow {
  min-width: 64px;
  font-family: var(--fs-font-display);
  font-size: 56px;
  font-weight: 800;
  line-height: 0.85;
}

.flow.end {
  text-align: right;
}

.flow-label {
  color: var(--fs-muted);
  font-size: var(--fs-text-caption);
  font-weight: 700;
  letter-spacing: 0.12em;
}

.delta,
.clock,
.mono {
  font-family: var(--fs-font-mono);
  font-variant-numeric: tabular-nums;
}

.delta {
  font-size: var(--fs-text-footnote);
  white-space: nowrap;
}

.home {
  color: var(--fs-home);
}

.away {
  color: var(--fs-away);
}

.center {
  display: flex;
  flex-direction: column;
  gap: 10px;
  align-items: center;
}

.score {
  font-family: var(--fs-font-display);
  font-size: 96px;
  font-weight: 800;
  line-height: 0.8;
  white-space: nowrap;
}

.colon {
  padding: 0 14px;
  color: var(--fs-faint);
}

.clock {
  color: var(--fs-muted);
  font-size: var(--fs-text-footnote);
}

.share {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-8);
}

.share-head {
  display: flex;
  justify-content: space-between;
  color: var(--fs-muted);
  font-size: var(--fs-text-footnote);
}
</style>
