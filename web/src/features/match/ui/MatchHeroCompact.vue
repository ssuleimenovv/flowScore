<script setup lang="ts">
import { computed } from 'vue'
import { useTweenedNumber } from '@/shared/motion/useTweenedNumber'
import { flowShare, formatMinute } from '../format'
import type { MatchLive } from '../matchState'
import FlowShareBar from './FlowShareBar.vue'
import LiveBadge from './LiveBadge.vue'
import TeamBadge from './TeamBadge.vue'

// The scoreboard from the MobileMatch board: no card, teams stacked over their
// names, and the two Flow figures at the ends of one bar.
const { live, seconds } = defineProps<{ live: MatchLive; seconds: number }>()

const match = computed(() => live.match)
const flowHome = useTweenedNumber(computed(() => live.flow.home))
const flowAway = useTweenedNumber(computed(() => live.flow.away))
const share = computed(() => flowShare(live.flow))
const minute = computed(() => formatMinute(seconds, match.value.clock.period))
</script>

<template>
  <section class="hero" aria-label="Счёт и поток">
    <div class="top">
      <LiveBadge :status="match.status" :minute />
    </div>

    <div class="scoreboard">
      <div class="team">
        <TeamBadge :code="match.home.code" side="home" size="md" />
        <span class="name">{{ match.home.name }}</span>
      </div>
      <div class="score" :aria-label="`Счёт ${live.score.home}:${live.score.away}`">
        {{ live.score.home }}<span class="colon" aria-hidden="true">:</span>{{ live.score.away }}
      </div>
      <div class="team">
        <TeamBadge :code="match.away.code" side="away" size="md" />
        <span class="name">{{ match.away.name }}</span>
      </div>
    </div>

    <div class="flow-row">
      <span class="flow home">{{ flowHome }}</span>
      <div class="flow-middle">
        <span class="flow-label">FLOW MOMENTUM</span>
        <FlowShareBar :home="share.home" :height="8" />
      </div>
      <span class="flow away">{{ flowAway }}</span>
    </div>
  </section>
</template>

<style scoped>
.hero {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-16);
  padding-top: 18px;
}

.top {
  display: flex;
  align-items: center;
  min-height: var(--fs-touch);
}

.scoreboard {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto minmax(0, 1fr);
  gap: var(--fs-space-8);
  align-items: center;
}

.team {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-8);
  align-items: center;
  min-width: 0;
}

.name {
  max-width: 100%;
  overflow: hidden;
  font-size: var(--fs-text-subheadline);
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.score {
  font-family: var(--fs-font-display);
  font-size: 64px;
  font-weight: 800;
  line-height: 0.85;
  white-space: nowrap;
}

.colon {
  padding: 0 var(--fs-space-8);
  color: var(--fs-faint);
}

.flow-row {
  display: flex;
  gap: var(--fs-space-12);
  align-items: center;
}

.flow {
  font-family: var(--fs-font-display);
  font-size: 36px;
  font-weight: 800;
  line-height: 1;
}

.home {
  color: var(--fs-home);
}

.away {
  color: var(--fs-away);
}

.flow-middle {
  display: flex;
  flex-grow: 1;
  flex-direction: column;
  gap: 6px;
}

.flow-label {
  color: var(--fs-muted);
  font-size: var(--fs-text-tab);
  font-weight: 700;
  letter-spacing: 0.14em;
  text-align: center;
}
</style>
