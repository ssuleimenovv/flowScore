<script setup lang="ts">
import { computed } from 'vue'
import AppCard from '@/shared/ui/AppCard.vue'
import AppIcon from '@/shared/ui/AppIcon.vue'
import type { MatchLive } from '../matchState'
import { shots, shotsWord, shotTotals } from '../shots'

// "Карта ударов" from the MatchPhases boards: where each team shot from, the
// size of a circle its xG, filled for a goal and empty for a miss.
const { live } = defineProps<{ live: MatchLive }>()

const items = computed(() => shots(live.events))
const totals = computed(() => shotTotals(live.events))
const label = computed(
  () => `Карта ударов: ${live.match.home.name} атакует вправо, ${live.match.away.name} влево`,
)
</script>

<template>
  <AppCard as="section" class="shot-map fs-in" style="--fs-i: 1" aria-labelledby="shots-title">
    <div class="head">
      <h2 id="shots-title" class="title">Карта ударов</h2>
      <span class="note">наша xG</span>
    </div>

    <div class="pitch" role="img" :aria-label="label">
      <span class="halfway" />
      <span class="box left" />
      <span class="box right" />
      <span class="mouth left" />
      <span class="mouth right" />
      <span
        v-for="s in items"
        :key="s.id"
        class="dot"
        :class="[s.side, s.outcome]"
        :title="s.label"
        :style="{ left: s.left, top: s.top, width: s.size, height: s.size }"
      />
      <div v-if="!items.length" class="empty">
        <AppIcon name="target" :size="28" />
        <span class="empty-title">Ударов пока нет</span>
        <span>Карта заполнится после первого удара</span>
      </div>
    </div>

    <div class="foot">
      <span class="legend" aria-hidden="true">
        <span class="key"><span class="mark home goal" />Гол</span>
        <span class="key"><span class="mark home on" />В створ</span>
        <span class="key"><span class="mark home block" />Блок</span>
        <span class="key"><span class="mark home off" />Мимо</span>
        <span>размер — наша xG</span>
      </span>
      <span class="totals">
        <span class="home-total">xG {{ totals.home.xg }}</span> ·
        {{ shotsWord(totals.home.shots) }} —
        <span class="away-total">xG {{ totals.away.xg }}</span> · {{ totals.away.shots }}
      </span>
    </div>
  </AppCard>
</template>

<style scoped>
.shot-map {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: var(--fs-space-24) 28px;
}

.head {
  display: flex;
  gap: var(--fs-space-12);
  align-items: center;
  justify-content: space-between;
}

.title {
  font-size: 20px;
  font-weight: 700;
  line-height: 1.2;
  letter-spacing: -0.01em;
}

.note {
  color: var(--fs-muted);
  font-size: var(--fs-text-caption);
}

/* The pitch is drawn in lines of the card's own color: the halfway line,
   both boxes and both goal mouths. A mouth is not called .goal, which is the
   class of a goal's circle */
.pitch {
  position: relative;
  height: 380px;
  overflow: hidden;
  border: 1px solid var(--fs-line);
  border-radius: 14px;
  background: var(--fs-surface-2);
}

.halfway {
  position: absolute;
  top: 0;
  bottom: 0;
  left: 50%;
  width: 1px;
  background: var(--fs-line);
}

.box {
  position: absolute;
  top: 22%;
  width: 16%;
  height: 56%;
  border: 1px solid var(--fs-line);
}

.box.left {
  left: 0;
  border-left: 0;
}

.box.right {
  right: 0;
  border-right: 0;
}

.mouth {
  position: absolute;
  top: 42%;
  width: 6px;
  height: 16%;
  background: var(--fs-text);
}

.mouth.left {
  left: 0;
  border-radius: 0 3px 3px 0;
}

.mouth.right {
  right: 0;
  border-radius: 3px 0 0 3px;
}

/* A shot: centred on where it was taken. It pops in when it happens */
.dot {
  position: absolute;
  box-sizing: border-box;
  border: 2px solid;
  border-radius: 50%;
  translate: -50% -50%;
  animation: pop var(--fs-duration-lift) var(--fs-ease) both;
}

.dot.home {
  border-color: var(--fs-home);
}

.dot.away {
  border-color: var(--fs-away);
}

/* On target: filled with the text color, so it stands out from a miss */
.dot.on {
  background: var(--fs-text);
}

.dot.block.home {
  background: var(--fs-home-soft);
}

.dot.block.away {
  background: var(--fs-away-soft);
}

/* A goal is filled in and ringed, so it stands out among the misses */
.dot.goal {
  z-index: 1;
  box-shadow: 0 0 0 2px var(--fs-text);
}

.dot.goal.home {
  background: var(--fs-home);
}

.dot.goal.away {
  background: var(--fs-away);
}

@keyframes pop {
  from {
    opacity: 0;
    scale: 0.4;
  }
}

.empty {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  gap: 6px;
  align-items: center;
  justify-content: center;
  color: var(--fs-muted);
  font-size: var(--fs-text-footnote);
}

.empty-title {
  color: var(--fs-text);
  font-size: var(--fs-text-subheadline);
  font-weight: 600;
}

.foot {
  display: flex;
  flex-wrap: wrap;
  gap: var(--fs-space-12);
  align-items: center;
  justify-content: space-between;
  color: var(--fs-muted);
  font-size: var(--fs-text-caption);
}

.legend {
  display: flex;
  gap: 14px;
  align-items: center;
}

.key {
  display: flex;
  gap: 6px;
  align-items: center;
}

.mark {
  box-sizing: border-box;
  width: 12px;
  height: 12px;
  border: 2px solid var(--fs-home);
  border-radius: 50%;
}

.mark.goal {
  border: 0;
  background: var(--fs-home);
  box-shadow: 0 0 0 2px var(--fs-text);
}

.mark.on {
  background: var(--fs-text);
}

.mark.block {
  background: var(--fs-home-soft);
}

.totals {
  color: var(--fs-text);
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-footnote);
  font-variant-numeric: tabular-nums;
}

.home-total,
.away-total {
  font-weight: 600;
}

.home-total {
  color: var(--fs-home);
}

.away-total {
  color: var(--fs-away);
}

/* Phone, as on MobileMatchPhases: a lower pitch on a tighter card */
@media (max-width: 767px) {
  .shot-map {
    padding: var(--fs-space-16);
  }

  .pitch {
    height: 220px;
  }

  .legend {
    flex-wrap: wrap;
  }
}
</style>
