<script setup lang="ts">
// The Simulator screen (Simulator and MobileSimulator boards): "what if" for a
// live match, from the moment it is at.
import { computed, onUnmounted, ref, watchEffect } from 'vue'
import { useRoute } from 'vue-router'
import { useNavTitle } from '@/app/layout/navigation'
import { formatMinute } from '@/features/match/format'
import { useMatchClock } from '@/features/match/useMatchClock'
import { useMatchLive } from '@/features/match/useMatchLive'
import { clampMinute, redMinutes, scenarioStory } from '@/features/simulator/scenario'
import RedCardControl from '@/features/simulator/ui/RedCardControl.vue'
import ScenarioExplainCard from '@/features/simulator/ui/ScenarioExplainCard.vue'
import ScenarioOutcomeCard from '@/features/simulator/ui/ScenarioOutcomeCard.vue'
import { useSimulation } from '@/features/simulator/useSimulation'
import { ApiError } from '@/shared/api/http'
import type { Side } from '@/shared/api/types'
import AppButton from '@/shared/ui/AppButton.vue'
import AppCard from '@/shared/ui/AppCard.vue'
import AppIcon from '@/shared/ui/AppIcon.vue'
import AppSkeleton from '@/shared/ui/AppSkeleton.vue'
import ErrorState from '@/shared/ui/ErrorState.vue'
import NotFoundState from '@/shared/ui/NotFoundState.vue'
import StateMessage from '@/shared/ui/StateMessage.vue'

const route = useRoute()
const matchId = String(route.params.matchId)
const { state, error, retry } = useMatchLive(matchId)

const match = computed(() => state.value?.match)
const { seconds } = useMatchClock(match)
const minute = computed(() => Math.floor(seconds.value / 60))
const finished = computed(() => match.value?.status === 'finished')

// The scenario: one red card, whose and when. The minute follows the clock
// when the one picked has already passed
const side = ref<Side | null>(null)
const picked = ref(0)
const range = computed(() => redMinutes(minute.value))
const redMinute = computed(() => (range.value ? clampMinute(picked.value, range.value) : 0))
const red = computed(() =>
  side.value && range.value ? { side: side.value, minute: redMinute.value } : null,
)
const reds = computed(() => (red.value ? [red.value] : []))

// The model answers again when the scenario changes and every match minute
const moment = computed(
  () => `${minute.value}:${state.value?.score.home}:${state.value?.score.away}`,
)
const sim = useSimulation(matchId, reds, moment)

const story = computed(() =>
  match.value ? scenarioStory(match.value, sim.simulation.value, red.value) : null,
)

function pick(next: Side | null) {
  // A fresh red starts six minutes ahead, as the board does at 72′ → 78′
  if (next && !side.value) picked.value = minute.value + 6
  side.value = next
}

function reset() {
  side.value = null
}

// "Manchester City 1:1 Arsenal · 72′"
const line = computed(() => {
  const s = state.value
  if (!s) return ''
  const { home, away, status, clock } = s.match
  const when =
    status === 'scheduled'
      ? 'до матча'
      : status === 'finished'
        ? 'матч завершён'
        : formatMinute(seconds.value, clock.period)
  return `${home.name} ${s.score.home}:${s.score.away} ${away.name} · ${when}`
})

const notFound = computed(() => error.value instanceof ApiError && error.value.status === 404)
const noModel = computed(
  () => sim.error.value instanceof ApiError && sim.error.value.status === 404,
)

// Mobile nav bar: "Симулятор" / "Сити 1:1 Арсенал · 72′", as on MobileSimulator
const { setTitle } = useNavTitle()
watchEffect(() => setTitle(state.value ? { title: 'Симулятор', subtitle: line.value } : null))
onUnmounted(() => setTitle(null))
</script>

<template>
  <div class="page">
    <div v-if="state && match" class="layout">
      <header class="top">
        <div class="heading">
          <RouterLink class="back fs-link" :to="{ name: 'match', params: { matchId } }">
            <AppIcon name="back" :size="18" />
            {{ line }}
          </RouterLink>
          <h1 class="h1">Что, если<span class="accent">…?</span></h1>
        </div>
        <AppButton variant="secondary" icon="retry" class="reset wide" @click="reset">
          Сбросить
        </AppButton>
        <AppButton
          variant="secondary"
          icon="retry"
          class="reset narrow"
          aria-label="Сбросить сценарий"
          @click="reset"
        />
      </header>

      <AppCard v-if="finished" as="section" class="controls fs-in">
        <StateMessage
          icon="check"
          title="Матч завершён"
          text="Сценарии считаются для идущего матча: после финального свистка исход уже известен."
        />
      </AppCard>
      <AppCard v-else as="section" class="controls fs-in" style="--fs-i: 2" aria-label="Сценарий">
        <RedCardControl
          :match
          :range
          :side
          :minute="redMinute"
          @update:side="pick"
          @update:minute="picked = $event"
        />
      </AppCard>

      <div class="results">
        <ScenarioOutcomeCard
          v-if="sim.simulation.value && !finished"
          :match
          :simulation="sim.simulation.value"
          class="outcome"
        />
        <AppCard v-else-if="noModel" class="outcome">
          <StateMessage
            icon="warning"
            title="Модель исхода не загружена"
            text="Сервер запущен без model.json, поэтому посчитать сценарий нечем."
          />
        </AppCard>
        <ScenarioExplainCard v-if="story && !finished" :story class="why" />
      </div>
    </div>

    <div v-else-if="!error" class="layout" aria-busy="true" aria-label="Загрузка симулятора">
      <AppCard class="skeleton controls">
        <AppSkeleton width="40%" :height="20" />
        <AppSkeleton :height="44" />
        <AppSkeleton :height="20" />
      </AppCard>
      <div class="results">
        <AppCard class="skeleton outcome">
          <AppSkeleton width="30%" :height="24" />
          <AppSkeleton :height="14" />
          <AppSkeleton :height="44" />
        </AppCard>
      </div>
    </div>

    <AppCard v-else-if="notFound" class="message">
      <NotFoundState
        title="Такого матча нет"
        text="Возможно, ссылка устарела или матч перенесли."
      />
    </AppCard>

    <AppCard v-else class="message">
      <ErrorState title="Не удалось загрузить матч" :error @retry="retry" />
    </AppCard>
  </div>
</template>

<style scoped>
/* Phone: the outcome first, then the controls and the explanation, as on
   MobileSimulator */
.page {
  padding: var(--fs-space-16) var(--fs-screen-padding);
}

.layout,
.results {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-16);
}

.results {
  display: contents;
}

.top {
  display: flex;
  gap: var(--fs-space-24);
  align-items: flex-end;
  justify-content: space-between;
}

.heading {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.back {
  display: none;
  gap: 6px;
  align-items: center;
  color: var(--fs-muted);
  font-size: var(--fs-text-subheadline);
  text-decoration: none;
}

.h1 {
  font-size: 34px;
  font-weight: 700;
  letter-spacing: -0.025em;
  line-height: 1.05;
}

.accent {
  color: var(--fs-home);
}

.reset.wide {
  display: none;
}

.outcome {
  order: 1;
}

.controls {
  display: flex;
  flex-direction: column;
  gap: 28px;
  order: 2;
  padding: 18px var(--fs-space-16);
}

.why {
  order: 3;
}

.skeleton {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-16);
  padding: var(--fs-space-24);
}

/* Desktop: controls on the left, results on the right, as on the Simulator board */
@media (min-width: 768px) {
  .page {
    padding: 28px var(--fs-page-padding) var(--fs-space-40);
  }

  .layout {
    display: grid;
    grid-template-columns: 440px minmax(0, 1fr);
    gap: var(--fs-space-24);
    align-items: start;
  }

  .top {
    grid-column: 1 / -1;
  }

  .back {
    display: flex;
  }

  .h1 {
    font-size: 48px;
  }

  .reset.wide {
    display: inline-flex;
  }

  .reset.narrow {
    display: none;
  }

  /* The phone's order (outcome first) does not apply to the two columns */
  .controls {
    order: 0;
    padding: var(--fs-space-24);
  }

  .results {
    display: flex;
    min-width: 0;
  }
}

/* Below 1024 the two columns do not fit: everything in one */
@media (min-width: 768px) and (max-width: 1023px) {
  .layout {
    grid-template-columns: minmax(0, 1fr);
  }
}
</style>
