<script setup lang="ts">
// The Match screen (Match and MobileMatch boards).
import { computed, onUnmounted, watchEffect } from 'vue'
import { useRoute } from 'vue-router'
import { useNavTitle } from '@/app/layout/navigation'
import MatchHero from '@/features/match/ui/MatchHero.vue'
import FlowWaveCard from '@/features/match/ui/FlowWaveCard.vue'
import FlowWaveSkeleton from '@/features/match/ui/FlowWaveSkeleton.vue'
import MatchHeroCompact from '@/features/match/ui/MatchHeroCompact.vue'
import MatchHeroSkeleton from '@/features/match/ui/MatchHeroSkeleton.vue'
import { useMatchClock } from '@/features/match/useMatchClock'
import { useMatchLive } from '@/features/match/useMatchLive'
import { ApiError } from '@/shared/api/http'
import { resolveScreenState } from '@/shared/state/screenState'
import { useNetwork } from '@/shared/state/useNetwork'
import AppCard from '@/shared/ui/AppCard.vue'
import ErrorState from '@/shared/ui/ErrorState.vue'
import NotFoundState from '@/shared/ui/NotFoundState.vue'
import OfflineBanner from '@/shared/ui/OfflineBanner.vue'
import ReconnectBanner from '@/shared/ui/ReconnectBanner.vue'
import StateMessage from '@/shared/ui/StateMessage.vue'

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

    <!-- Both layouts are in the DOM; CSS shows one of them at the 768 px breakpoint -->
    <template v-if="screen.view === 'content' && state">
      <div class="wide"><MatchHero :live="state" :seconds /></div>
      <div class="narrow"><MatchHeroCompact :live="state" :seconds /></div>
      <FlowWaveCard :live="state" :seconds />
    </template>

    <template v-else-if="screen.view === 'loading'">
      <div class="wide"><MatchHeroSkeleton /></div>
      <div class="narrow"><MatchHeroSkeleton compact /></div>
      <FlowWaveSkeleton />
    </template>

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

.wide {
  display: none;
}

.narrow {
  display: contents;
}

.message {
  margin-top: var(--fs-space-16);
}

/* Desktop: cards on a padded page, as on the Match board */
@media (min-width: 768px) {
  .page {
    gap: var(--fs-space-24);
    padding: 28px var(--fs-page-padding) var(--fs-space-40);
  }

  .wide {
    display: contents;
  }

  .narrow {
    display: none;
  }

  .message {
    margin-top: 0;
  }
}
</style>
