<script setup lang="ts">
// The Match screen (Match and MobileMatch boards).
import { computed, onUnmounted, watchEffect } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useNavTitle } from '@/app/layout/navigation'
import ChronicleCard from '@/features/match/ui/ChronicleCard.vue'
import ChronicleSkeleton from '@/features/match/ui/ChronicleSkeleton.vue'
import ExplainCard from '@/features/match/ui/ExplainCard.vue'
import FlowWaveCard from '@/features/match/ui/FlowWaveCard.vue'
import FlowWaveSkeleton from '@/features/match/ui/FlowWaveSkeleton.vue'
import MatchHero from '@/features/match/ui/MatchHero.vue'
import MatchHeroCompact from '@/features/match/ui/MatchHeroCompact.vue'
import MatchHeroSkeleton from '@/features/match/ui/MatchHeroSkeleton.vue'
import MatchTabs, { type MatchTab } from '@/features/match/ui/MatchTabs.vue'
import OutcomeCard from '@/features/match/ui/OutcomeCard.vue'
import SimulatorCta from '@/features/match/ui/SimulatorCta.vue'
import StatsCard from '@/features/match/ui/StatsCard.vue'
import { useMatchClock } from '@/features/match/useMatchClock'
import { useMatchLive } from '@/features/match/useMatchLive'
import { ApiError } from '@/shared/api/http'
import { resolveScreenState } from '@/shared/state/screenState'
import { useMediaQuery } from '@/shared/state/useMediaQuery'
import { useNetwork } from '@/shared/state/useNetwork'
import AppCard from '@/shared/ui/AppCard.vue'
import ErrorState from '@/shared/ui/ErrorState.vue'
import NotFoundState from '@/shared/ui/NotFoundState.vue'
import OfflineBanner from '@/shared/ui/OfflineBanner.vue'
import ReconnectBanner from '@/shared/ui/ReconnectBanner.vue'
import StateMessage from '@/shared/ui/StateMessage.vue'
import WakingNote from '@/shared/ui/WakingNote.vue'

const route = useRoute()
const {
  state,
  request,
  error,
  updatedAt,
  socketStatus,
  retryAt,
  secondsToRetry,
  reconnectNow,
  retry,
} = useMatchLive(String(route.params.matchId))
const { online } = useNetwork()

const match = computed(() => state.value?.match)
const { seconds } = useMatchClock(match)

const screen = computed(() =>
  resolveScreenState({
    status: request.value,
    hasData: state.value !== null,
    isEmpty: false,
    online: online.value,
    socket: socketStatus.value,
  }),
)

// The phone shows one section at a time. The open tab lives in the URL
// (?tab=events), so "Назад" and a shared link come back to it.
const router = useRouter()
const TABS: MatchTab[] = ['flow', 'events', 'stats', 'ai']
const tab = computed<MatchTab>({
  get: () => TABS.find((t) => t === route.query.tab) ?? 'flow',
  set: (next) =>
    void router.replace({ query: { ...route.query, tab: next === 'flow' ? undefined : next } }),
})

// On the phone a section is a tab panel; on the desktop everything is on screen
// at once and the tab roles would only confuse a screen reader.
const phone = useMediaQuery('(max-width: 767px)')
function panel(id: MatchTab) {
  return {
    id: `panel-${id}`,
    class: { active: tab.value === id },
    ...(phone.value ? { role: 'tabpanel', 'aria-labelledby': `tab-${id}` } : {}),
  }
}

const notFound = computed(() => error.value instanceof ApiError && error.value.status === 404)

// Mobile nav bar: "Premier League" / "37-й тур · Etihad Stadium", as on MobileMatch
const { setTitle } = useNavTitle()
watchEffect(() => {
  const m = match.value
  setTitle(
    m
      ? {
          title: m.competition.name,
          subtitle: [m.competition.round && `${m.competition.round}-й тур`, m.venue?.name]
            .filter(Boolean)
            .join(' · '),
        }
      : null,
  )
})
onUnmounted(() => setTitle(null))
</script>

<template>
  <div class="page">
    <OfflineBanner v-if="screen.banner === 'offline'" :updated-at="updatedAt" />
    <ReconnectBanner
      v-else-if="screen.banner === 'reconnecting'"
      :seconds="secondsToRetry"
      :retry-at="retryAt"
      @reconnect="reconnectNow"
    />

    <!-- Main column and sidebar, as on the Match board. On the phone the same
         cards become tab panels, and CSS shows only the open one -->
    <div v-if="screen.view === 'content' && state" class="layout">
      <div class="main">
        <!-- Both scoreboards are in the DOM; CSS shows one of them at 768 px -->
        <div class="wide"><MatchHero :live="state" :seconds /></div>
        <div class="narrow"><MatchHeroCompact :live="state" :seconds /></div>
        <div class="narrow"><MatchTabs v-model="tab" /></div>
        <div class="panel" v-bind="panel('flow')"><FlowWaveCard :live="state" :seconds /></div>
        <!-- Why Flow is what it is and who wins from here: side by side on the
             desktop, one under the other on the phone's flow tab -->
        <div class="wide">
          <div class="insights" :class="{ pair: state.prediction }">
            <ExplainCard :live="state" :updated-at="updatedAt" />
            <OutcomeCard v-if="state.prediction" :live="state" :prediction="state.prediction" />
          </div>
        </div>
        <div class="panel narrow" :class="{ active: tab === 'flow' }">
          <ExplainCard :live="state" :updated-at="updatedAt" variant="brief" />
          <OutcomeCard
            v-if="state.prediction"
            :live="state"
            :prediction="state.prediction"
            variant="brief"
          />
        </div>
        <div class="panel" v-bind="panel('stats')">
          <StatsCard v-if="state.match.stats.length" :live="state" />
          <AppCard v-else class="fs-in">
            <StateMessage
              icon="pulse"
              title="Статистики пока нет"
              text="Удары, xG и владение появятся здесь с первых минут матча."
            />
          </AppCard>
        </div>
      </div>
      <aside class="side">
        <div class="panel" v-bind="panel('events')"><ChronicleCard :live="state" /></div>
        <div class="panel" :class="{ active: tab === 'flow' }"><SimulatorCta :match-id="state.match.id"/></div>
        <!-- The phone's AI tab: every factor; the desktop shows them under the wave -->
        <div class="panel narrow-only" v-bind="panel('ai')">
          <ExplainCard :live="state" :updated-at="updatedAt" variant="factors" />
        </div>
      </aside>
    </div>

    <div v-else-if="screen.view === 'loading'" class="layout">
      <div class="main">
        <WakingNote />
        <div class="wide"><MatchHeroSkeleton /></div>
        <div class="narrow"><MatchHeroSkeleton compact /></div>
        <FlowWaveSkeleton />
      </div>
      <aside class="side">
        <ChronicleSkeleton />
      </aside>
    </div>

    <AppCard v-else-if="notFound" class="message">
      <NotFoundState
        title="Такого матча нет"
        text="Возможно, ссылка устарела или матч перенесли."
      />
    </AppCard>

    <AppCard v-else-if="screen.view === 'error'" class="message">
      <ErrorState title="Не удалось загрузить матч" :error @retry="retry" />
    </AppCard>

    <AppCard v-else-if="screen.view === 'offline'" class="message">
      <StateMessage
        icon="offline"
        title="Нет подключения"
        text="Покажем матч, когда сеть вернётся."
      />
    </AppCard>
  </div>
</template>

<style scoped>
/* Phone: the scoreboard sits right on the page, as on MobileMatch */
.page {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-16);
  padding: 0 var(--fs-screen-padding) var(--fs-space-16);
}

.layout,
.main,
.side {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-16);
}

.wide {
  display: none;
}

.narrow {
  display: contents;
}

.message {
  margin-top: var(--fs-space-16);
}

/* Phone: only the open tab's cards are shown. Showing a card again replays its
   fs-in entrance, so switching tabs animates for free */
@media (max-width: 767px) {
  .panel:not(.active) {
    display: none;
  }
}

/* Desktop: cards on a padded page, as on the Match board */
@media (min-width: 768px) {
  .page {
    gap: var(--fs-space-24);
    padding: 28px var(--fs-page-padding) var(--fs-space-40);
  }

  .layout,
  .main,
  .side {
    gap: var(--fs-space-24);
  }

  .wide {
    display: contents;
  }

  .narrow,
  .narrow-only {
    display: none;
  }

  .message {
    margin-top: 0;
  }

  .insights {
    display: flex;
    flex-direction: column;
    gap: var(--fs-space-24);
  }
}

/* The explanation and the outcome side by side, as on the Match board */
@media (min-width: 1024px) {
  .insights.pair {
    display: grid;
    grid-template-columns: minmax(0, 1.35fr) minmax(0, 1fr);
  }
}

/* Wide desktop: the sidebar stands next to the main column. Below 1200 it moves under it */
@media (min-width: 1200px) {
  .layout {
    display: grid;
    grid-template-columns: minmax(0, 1fr) var(--fs-sidebar-max);
    align-items: start;
  }
}
</style>
