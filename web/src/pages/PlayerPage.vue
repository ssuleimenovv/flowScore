<script setup lang="ts">
// The player page (Player and PlayerSeason boards): their card, then their
// season in our matches. Load and the impact on Flow come with D6.
import { computed, onUnmounted, watchEffect } from 'vue'
import { useRoute } from 'vue-router'
import { useNavTitle } from '@/app/layout/navigation'
import { finishing, logRow, navLine, seasonTiles, totals } from '@/features/player/player'
import FinishingCard from '@/features/player/ui/FinishingCard.vue'
import MatchLogCard from '@/features/player/ui/MatchLogCard.vue'
import PlayerHero from '@/features/player/ui/PlayerHero.vue'
import SeasonTiles from '@/features/player/ui/SeasonTiles.vue'
import { usePlayer } from '@/features/player/usePlayer'
import { ApiError } from '@/shared/api/http'
import { resolveScreenState } from '@/shared/state/screenState'
import { useNetwork } from '@/shared/state/useNetwork'
import AppCard from '@/shared/ui/AppCard.vue'
import AppSkeleton from '@/shared/ui/AppSkeleton.vue'
import ErrorState from '@/shared/ui/ErrorState.vue'
import NotFoundState from '@/shared/ui/NotFoundState.vue'
import OfflineBanner from '@/shared/ui/OfflineBanner.vue'
import StateMessage from '@/shared/ui/StateMessage.vue'
import WakingNote from '@/shared/ui/WakingNote.vue'

const route = useRoute()
const playerId = computed(() => String(route.params.playerId))
const { page, request, error, updatedAt, retry } = usePlayer(playerId)

const { online } = useNetwork()
const screen = computed(() =>
  resolveScreenState({
    status: request.value,
    hasData: page.value !== null,
    isEmpty: false,
    online: online.value,
  }),
)
const notFound = computed(() => error.value instanceof ApiError && error.value.status === 404)

const season = computed(() => (page.value ? totals(page.value.matches) : null))
const rows = computed(() => page.value?.matches.map((m) => logRow(m)) ?? [])

// Mobile nav bar: "Erling Haaland" / "Manchester City · №9", as on MobilePlayer
const { setTitle } = useNavTitle()
watchEffect(() => {
  const p = page.value?.player
  setTitle(p ? { title: p.name, subtitle: navLine(p) } : null)
})
onUnmounted(() => setTitle(null))
</script>

<template>
  <div class="page">
    <OfflineBanner v-if="screen.banner === 'offline'" :updated-at="updatedAt" />

    <template v-if="screen.view === 'content' && page">
      <PlayerHero :player="page.player" />
      <SeasonTiles v-if="season" :tiles="seasonTiles(season)" />
      <div class="season">
        <FinishingCard v-if="season" :f="finishing(season)" />
        <MatchLogCard :rows />
      </div>
    </template>

    <div
      v-else-if="screen.view === 'loading'"
      class="loading"
      aria-busy="true"
      aria-label="Загрузка игрока"
    >
      <WakingNote />
      <AppCard class="skeleton">
        <AppSkeleton width="104px" :height="104" />
        <AppSkeleton width="60%" :height="40" />
        <AppSkeleton width="40%" :height="16" />
      </AppCard>
    </div>

    <AppCard v-else-if="notFound" class="message">
      <NotFoundState title="Такого игрока нет" text="Возможно, ссылка устарела." />
    </AppCard>

    <AppCard v-else-if="screen.view === 'error'" class="message">
      <ErrorState title="Не удалось загрузить игрока" :error @retry="retry" />
    </AppCard>

    <AppCard v-else-if="screen.view === 'offline'" class="message">
      <StateMessage
        icon="offline"
        title="Нет подключения"
        text="Покажем игрока, когда сеть вернётся."
      />
    </AppCard>
  </div>
</template>

<style scoped>
/* Phone, as on MobilePlayer */
.page {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-16);
  padding: 18px var(--fs-screen-padding) var(--fs-space-24);
}

.loading {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-16);
}

.skeleton {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-12);
  padding: 20px;
}

.season {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-16);
}

/* Desktop, as on Player */
@media (min-width: 768px) {
  .page {
    gap: 22px;
    padding: var(--fs-space-24) var(--fs-page-padding) var(--fs-space-40);
  }

  /* PlayerSeason: "Голы против xG" in a 400 px column beside the log */
  .season {
    display: grid;
    grid-template-columns: 400px minmax(0, 1fr);
    gap: var(--fs-space-24);
    align-items: start;
  }
}
</style>
