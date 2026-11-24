<script setup lang="ts">
// The team page (Team and MobileTeam boards): its card, then "Обзор",
// "Состав" and "Расписание". The open tab lives in the URL (?tab=squad),
// so "Назад" and a shared link come back to it.
import { computed, onScopeDispose, onUnmounted, ref, watchEffect } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useNavTitle } from '@/app/layout/navigation'
import { leagueLine, results, upcoming } from '@/features/team/team'
import TeamFlowCard from '@/features/team/ui/TeamFlowCard.vue'
import TeamHero from '@/features/team/ui/TeamHero.vue'
import TeamMatchesCard from '@/features/team/ui/TeamMatchesCard.vue'
import TeamSquadCard from '@/features/team/ui/TeamSquadCard.vue'
import TeamTabs, { type TeamTab } from '@/features/team/ui/TeamTabs.vue'
import { useSquad, useTeam } from '@/features/team/useTeam'
import { ApiError } from '@/shared/api/http'
import { resolveScreenState } from '@/shared/state/screenState'
import { useNetwork } from '@/shared/state/useNetwork'
import AppCard from '@/shared/ui/AppCard.vue'
import AppSkeleton from '@/shared/ui/AppSkeleton.vue'
import AppSpinner from '@/shared/ui/AppSpinner.vue'
import ErrorState from '@/shared/ui/ErrorState.vue'
import NotFoundState from '@/shared/ui/NotFoundState.vue'
import OfflineBanner from '@/shared/ui/OfflineBanner.vue'
import StateMessage from '@/shared/ui/StateMessage.vue'
import WakingNote from '@/shared/ui/WakingNote.vue'

const route = useRoute()
const router = useRouter()

const teamId = computed(() => String(route.params.teamId))
const { page, request, preparing, error, updatedAt, retry } = useTeam(teamId)

// "Сегодня, 21:00" follows the clock
const now = ref(new Date())
const timer = setInterval(() => (now.value = new Date()), 60_000)
onScopeDispose(() => clearInterval(timer))

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

const TABS: TeamTab[] = ['overview', 'squad', 'fixtures']
const tab = computed<TeamTab>({
  get: () => TABS.find((t) => t === route.query.tab) ?? 'overview',
  set: (next) =>
    void router.replace({ query: { ...route.query, tab: next === 'overview' ? undefined : next } }),
})

// The squad is asked for once its tab is opened
const {
  items: players,
  request: squadRequest,
  error: squadError,
  retry: retrySquad,
} = useSquad(
  teamId,
  computed(() => tab.value === 'squad'),
)

// Mobile nav bar: "Manchester City" / "Премьер-лига · 3-е место", as on MobileTeam
const { setTitle } = useNavTitle()
watchEffect(() => {
  const p = page.value
  setTitle(p ? { title: p.team.name, subtitle: leagueLine(p.league, false) } : null)
})
onUnmounted(() => setTitle(null))
</script>

<template>
  <div class="page">
    <OfflineBanner v-if="screen.banner === 'offline'" :updated-at="updatedAt" />

    <template v-if="screen.view === 'content' && page">
      <TeamHero :page />
      <TeamTabs v-model="tab" />
      <div
        v-if="tab === 'overview'"
        id="panel-overview"
        class="overview"
        role="tabpanel"
        aria-labelledby="tab-overview"
      >
        <div class="main">
          <TeamFlowCard :matches="page.matches" :team-id="page.team.id" />
          <TeamMatchesCard
            kind="upcoming"
            title="Ближайшие матчи"
            phone-title="Ближайшие"
            :matches="upcoming(page.matches, 4)"
            :team-id="page.team.id"
            :league-id="page.league?.id ?? null"
            :now
            empty="Ближайших матчей пока нет."
          />
        </div>
        <aside class="side">
          <TeamMatchesCard
            kind="results"
            title="Результаты"
            :matches="results(page.matches, 4)"
            :team-id="page.team.id"
            :league-id="page.league?.id ?? null"
            :now
            empty="Сыгранных матчей пока нет."
          />
        </aside>
      </div>
      <TeamSquadCard
        v-else-if="tab === 'squad'"
        id="panel-squad"
        role="tabpanel"
        aria-labelledby="tab-squad"
        :players
        :request="squadRequest"
        :error="squadError"
        @retry="retrySquad"
      />
      <div
        v-else
        id="panel-fixtures"
        class="fixtures"
        role="tabpanel"
        aria-labelledby="tab-fixtures"
      >
        <TeamMatchesCard
          kind="upcoming"
          title="Расписание"
          :matches="upcoming(page.matches)"
          :team-id="page.team.id"
          :league-id="page.league?.id ?? null"
          :now
          empty="Ближайших матчей пока нет."
        />
        <TeamMatchesCard
          kind="results"
          title="Сыгранные"
          :matches="results(page.matches)"
          :team-id="page.team.id"
          :league-id="page.league?.id ?? null"
          :now
          empty="Сыгранных матчей пока нет."
        />
      </div>
    </template>

    <div
      v-else-if="screen.view === 'loading'"
      class="loading"
      aria-busy="true"
      aria-label="Загрузка команды"
    >
      <WakingNote />
      <AppCard class="skeleton">
        <AppSkeleton width="112px" :height="112" />
        <AppSkeleton width="60%" :height="40" />
        <AppSkeleton width="40%" :height="28" />
      </AppCard>
    </div>

    <AppCard v-else-if="notFound" class="message">
      <NotFoundState
        title="Такой команды нет"
        text="Возможно, ссылка устарела или мы не следим за её лигой."
      />
    </AppCard>

    <AppCard v-else-if="screen.view === 'error' && preparing" class="message">
      <StateMessage
        icon="pulse"
        title="Считаем статистику сезона"
        text="Прогоняем сыгранные матчи через Flow. Обычно пара минут — страница откроется сама."
      >
        <AppSpinner />
      </StateMessage>
    </AppCard>

    <AppCard v-else-if="screen.view === 'error'" class="message">
      <ErrorState title="Не удалось загрузить команду" :error @retry="retry" />
    </AppCard>

    <AppCard v-else-if="screen.view === 'offline'" class="message">
      <StateMessage
        icon="offline"
        title="Нет подключения"
        text="Покажем команду, когда сеть вернётся."
      />
    </AppCard>
  </div>
</template>

<style scoped>
/* Phone, as on MobileTeam */
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

.overview,
.main,
.fixtures {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-16);
}

/* MobileTeam shows no results on the overview */
.side {
  display: none;
}

/* Desktop, as on Team */
@media (min-width: 768px) {
  /* Team: the main column beside a 400 px one */
  .overview {
    display: grid;
    grid-template-columns: minmax(0, 1fr) 400px;
    gap: var(--fs-space-24);
    align-items: start;
  }

  .main,
  .side {
    display: flex;
    flex-direction: column;
    gap: var(--fs-space-24);
    min-width: 0;
  }

  /* Team: the schedule and the results side by side */
  .fixtures {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: var(--fs-space-24);
    align-items: start;
  }

  .page {
    gap: 22px;
    padding: var(--fs-space-24) var(--fs-page-padding) var(--fs-space-40);
  }
}
</style>
