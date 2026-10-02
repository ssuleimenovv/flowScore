<script setup lang="ts">
import { computed, onScopeDispose, ref, useId } from 'vue'
import AppCard from '@/shared/ui/AppCard.vue'
import AppIcon from '@/shared/ui/AppIcon.vue'
import { explain } from '../explain'
import type { MatchLive } from '../matchState'

// "Explainability Agent" from the Match boards: why Flow is what it is.
// full: the desktop card under the wave; brief: the same on the phone's flow
// tab, top three; factors: the phone's AI tab, "Что двигает поток".
const {
  live,
  updatedAt,
  variant = 'full',
} = defineProps<{
  live: MatchLive
  updatedAt: Date | null
  variant?: 'full' | 'brief' | 'factors'
}>()

// Each copy of the card on the page needs its own heading id
const titleId = useId()

const explanation = computed(() => explain(live, variant === 'brief' ? 3 : 5))

// "обновлено 12 с назад", kept current between updates
const now = ref(Date.now())
const timer = setInterval(() => (now.value = Date.now()), 1000)
onScopeDispose(() => clearInterval(timer))

const ago = computed(() => {
  if (!updatedAt) return null
  const seconds = Math.max(0, Math.round((now.value - updatedAt.getTime()) / 1000))
  const text = seconds < 3 ? 'только что' : `${seconds} с назад`
  return variant === 'full' ? `обновлено ${text}` : text
})
</script>

<template>
  <AppCard
    as="section"
    class="explain fs-in"
    :class="`is-${variant}`"
    style="--fs-i: 3"
    :aria-labelledby="titleId"
  >
    <div class="head">
      <span class="chip">
        <AppIcon name="sparkle" :size="14" class="spark" />
        {{ variant === 'brief' ? 'AI-разбор' : 'Explainability Agent' }}
      </span>
      <span v-if="ago && variant !== 'factors'" class="ago">{{ ago }}</span>
    </div>

    <h3 v-if="variant === 'factors'" :id="titleId" class="title">Что двигает поток</h3>
    <div v-else class="summary">
      <h3 :id="titleId" class="title">{{ explanation.title }}</h3>
      <p class="text">{{ explanation.text }}</p>
    </div>

    <ul v-if="explanation.lines.length" class="factors">
      <li v-for="line in explanation.lines" :key="line.key" class="factor">
        <span class="label">{{ line.text }}</span>
        <span class="track" aria-hidden="true">
          <span class="bar" :class="line.side" :style="{ width: line.width }" />
        </span>
        <span class="value" :class="line.side">{{ line.value }}</span>
      </li>
    </ul>
    <p v-else class="empty">Факторы появятся с первыми событиями матча.</p>
  </AppCard>
</template>

<style scoped>
.explain {
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

.ago {
  color: var(--fs-muted);
  font-size: var(--fs-text-caption);
}

.summary {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-8);
}

.title {
  font-size: var(--fs-text-title-3);
  font-weight: 600;
  line-height: 1.25;
}

.text {
  color: var(--fs-muted);
  font-size: var(--fs-text-subheadline);
  line-height: 1.55;
}

.factors {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-12);
  padding: 0;
  list-style: none;
}

/* "4 удара за последние 7 минут ━━━━ +14" */
.factor {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 110px 36px;
  gap: var(--fs-space-12);
  align-items: center;
  font-size: var(--fs-text-subheadline);
  line-height: 1.35;
}

.track {
  display: flex;
  height: 6px;
  overflow: hidden;
  border-radius: var(--fs-radius-full);
  background: var(--fs-track);
}

/* The bar grows in when the card appears and slides when the factor changes */
.bar {
  border-radius: var(--fs-radius-full);
  transform-origin: left;
  animation: grow var(--fs-duration-bar-h) var(--fs-ease) both;
  transition: width var(--fs-duration-grow) var(--fs-ease);
}

.value {
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-footnote);
  font-variant-numeric: tabular-nums;
  text-align: right;
}

.bar.home {
  background: var(--fs-home);
}

.bar.away {
  background: var(--fs-away);
}

.value.home {
  color: var(--fs-home);
}

.value.away {
  color: var(--fs-away);
}

.empty {
  color: var(--fs-muted);
  font-size: var(--fs-text-subheadline);
}

@keyframes grow {
  from {
    transform: scaleX(0);
  }
}

/* Phone: the MobileMatch cards are tighter, with a shorter bar */
.is-brief,
.is-factors {
  gap: 14px;
  padding: 18px var(--fs-space-16);
}

.is-brief .chip,
.is-brief .ago,
.is-factors .chip {
  font-size: var(--fs-text-footnote);
}

.is-brief .summary {
  gap: 6px;
}

.is-brief .text {
  line-height: 1.5;
}

.is-brief .factor,
.is-factors .factor {
  grid-template-columns: minmax(0, 1fr) 64px 36px;
}

.is-brief .factors {
  gap: 10px;
}
</style>
