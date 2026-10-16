<script setup lang="ts">
// The Leagues screen (League and MobileLeague boards): the table with each
// team's average Flow, the matches of a round and the season's Flow records.
// The league and the round live in the URL (?league=ger.1&round=6), so "Назад"
// and a shared link come back to them.
import { computed, onScopeDispose, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { roundTabs, seasonLine } from '@/features/league/league'
import FlowRecords from '@/features/league/ui/FlowRecords.vue'
import RoundMatches from '@/features/league/ui/RoundMatches.vue'
import StandingsTable from '@/features/league/ui/StandingsTable.vue'
import { useLeague, useLeagues, useRound } from '@/features/league/useLeague'
import { leagueLogo } from '@/shared/media/logos'
import { resolveScreenState } from '@/shared/state/screenState'
import { useNetwork } from '@/shared/state/useNetwork'
import AppCard from '@/shared/ui/AppCard.vue'
import AppSkeleton from '@/shared/ui/AppSkeleton.vue'
import AppSpinner from '@/shared/ui/AppSpinner.vue'
import ErrorState from '@/shared/ui/ErrorState.vue'
import OfflineBanner from '@/shared/ui/OfflineBanner.vue'
import StateMessage from '@/shared/ui/StateMessage.vue'
import WakingNote from '@/shared/ui/WakingNote.vue'

const route = useRoute()
const router = useRouter()

const leagues = useLeagues()
// The league in the URL, if the gateway follows it; else the first one
const leagueId = computed<string | null>(() => {
  const ids = leagues.items.value.map((l) => l.id)
  const q = route.query.league
  return typeof q === 'string' && ids.includes(q) ? q : (ids[0] ?? null)
})
function pick(id: string) {
  const first = leagues.items.value[0]?.id
  void router.replace({ query: { league: id === first ? undefined : id } })
}

const { league, standings, request, preparing, error, updatedAt, retry } = useLeague(leagueId)

// The round in the URL, if the season has it; else the one being played
const round = computed<number>({
  get: () => {
    const l = league.value
    const q = Number(route.query.round)
    if (!l) return 0
    return Number.isInteger(q) && q >= 1 && q <= l.rounds ? q : l.round
  },
  set: (n) =>
    void router.replace({
      query: { ...route.query, round: n === league.value?.round ? undefined : String(n) },
    }),
})
const roundId = computed(() => (round.value > 0 ? round.value : null))
const tabs = computed(() =>
  league.value ? roundTabs(league.value.round, league.value.rounds) : [],
)
const matches = useRound(leagueId, roundId)

// The statuses of the round ("22:00", "Завтра") follow the clock
const now = ref(new Date())
const timer = setInterval(() => (now.value = new Date()), 60_000)
onScopeDispose(() => clearInterval(timer))

const { online } = useNetwork()
const screen = computed(() =>
  resolveScreenState({
    status: leagues.request.value === 'error' ? 'error' : request.value,
    hasData: league.value !== null,
    isEmpty: standings.value.length === 0,
    online: online.value,
  }),
)
const listError = computed(() =>
  leagues.request.value === 'error' ? leagues.error.value : error.value,
)
function retryAll() {
  if (leagues.request.value === 'error') void leagues.retry()
  else void retry()
}
</script>

<template>
  <div class="page">
    <header class="top">
      <div class="heading">
        <span v-if="league" class="season">{{ seasonLine(league) }}</span>
        <AppSkeleton v-else width="180px" :height="14" />
        <h1 class="h1">
          <span class="desktop-word">Турнирная <span class="accent">таблица</span></span>
          <span class="phone-word">Таблица</span>
        </h1>
      </div>
      <div v-if="leagues.items.value.length > 1" class="chips" role="group" aria-label="Лиги">
        <button
          v-for="l in leagues.items.value"
          :key="l.id"
          type="button"
          class="chip fs-press"
          :aria-pressed="l.id === leagueId"
          @click="pick(l.id)"
        >
          <template v-if="leagueLogo(l.id, 18)">
            <img class="chip-logo on-light" :src="leagueLogo(l.id, 18)!" alt="" />
            <img class="chip-logo on-dark" :src="leagueLogo(l.id, 18, 'dark')!" alt="" />
          </template>
          {{ l.name }}
        </button>
      </div>
    </header>

    <OfflineBanner v-if="screen.banner === 'offline'" class="banner" :updated-at="updatedAt" />

    <template v-if="screen.view === 'content' && league">
      <div class="main">
        <StandingsTable :rows="standings" />
      </div>
      <aside class="side">
        <RoundMatches
          v-model="round"
          :tabs
          :matches="matches.items.value"
          :loading="matches.request.value !== 'success'"
          :now
        />
        <FlowRecords :records="league.records" class="records" />
      </aside>
    </template>

    <AppCard v-else-if="screen.view === 'empty'" class="main notice">
      <StateMessage
        icon="trophy"
        title="Таблицы пока нет"
        text="Сезон ещё не начался. Таблица появится после первого тура."
      />
    </AppCard>

    <AppCard v-else-if="screen.view === 'error' && preparing" class="main notice">
      <StateMessage
        icon="pulse"
        title="Считаем статистику сезона"
        text="Прогоняем сыгранные матчи через Flow. Обычно пара минут — таблица откроется сама."
      >
        <AppSpinner />
      </StateMessage>
    </AppCard>

    <AppCard v-else-if="screen.view === 'error'" class="main notice">
      <ErrorState title="Не удалось загрузить лигу" :error="listError" @retry="retryAll" />
    </AppCard>

    <AppCard v-else-if="screen.view === 'offline'" class="main notice">
      <StateMessage
        icon="offline"
        title="Нет подключения"
        text="Покажем таблицу, когда сеть вернётся."
      />
    </AppCard>

    <div v-else class="main" aria-busy="true" aria-label="Загрузка таблицы">
      <WakingNote />
      <AppCard class="loading">
        <AppSkeleton v-for="i in 8" :key="i" :height="28" />
      </AppCard>
    </div>
  </div>
</template>

<style scoped>
/* Phone, as on MobileLeague: one column, the table and then the round */
.page {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-16);
  padding: 18px var(--fs-screen-padding) var(--fs-space-24);
}

.main,
.side {
  display: contents;
}

.top {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-16);
}

.heading {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.season {
  color: var(--fs-muted);
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-footnote);
  font-variant-numeric: tabular-nums;
}

.h1 {
  font-size: 34px;
  font-weight: 700;
  letter-spacing: -0.025em;
  line-height: 1.05;
}

.desktop-word {
  display: none;
}

.accent {
  color: var(--fs-home);
}

.chips {
  display: flex;
  gap: var(--fs-space-8);
  overflow-x: auto;
}

/* The picked chip is filled with the text color, as on the board */
.chip {
  display: inline-flex;
  flex-shrink: 0;
  gap: var(--fs-space-8);
  align-items: center;
  height: 40px;
  padding: 0 var(--fs-space-16);
  border: 1px solid var(--fs-line);
  border-radius: var(--fs-radius-full);
  background: var(--fs-surface);
  color: var(--fs-text);
  font-size: var(--fs-text-subheadline);
  font-weight: 600;
  white-space: nowrap;
  cursor: pointer;
}

.chip-logo {
  width: 18px;
  height: 18px;
  object-fit: contain;
}

/* The logo drawn for the chip's background: a chip is dark in the dark
   theme and light once picked, the other way round in the light theme */
.chip .on-light,
.chip[aria-pressed='true'] .on-dark,
:root[data-theme='light'] .chip[aria-pressed='true'] .on-light {
  display: none;
}

.chip[aria-pressed='true'] .on-light,
:root[data-theme='light'] .chip .on-light,
:root[data-theme='light'] .chip[aria-pressed='true'] .on-dark {
  display: block;
}

:root[data-theme='light'] .chip .on-dark {
  display: none;
}

.chip[aria-pressed='true'] {
  border-color: var(--fs-text);
  background: var(--fs-text);
  color: var(--fs-bg);
}

/* The season's records are a desktop card, as on the boards */
.records {
  display: none;
}

.notice {
  display: block;
}

.loading {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-12);
  padding: 20px;
}

/* Desktop, as on League: the table beside a 380 px column */
@media (min-width: 768px) {
  .page {
    display: grid;
    grid-template-columns: minmax(0, 1fr) 380px;
    gap: 20px 32px;
    align-items: start;
    padding: 32px var(--fs-page-padding) var(--fs-space-40);
  }

  .top,
  .banner,
  .main {
    grid-column: 1;
    min-width: 0;
  }

  .top {
    gap: 20px;
  }

  .heading {
    gap: 10px;
  }

  .h1 {
    font-size: 48px;
  }

  .desktop-word {
    display: inline;
  }

  .phone-word {
    display: none;
  }

  .chips {
    flex-wrap: wrap;
  }

  .main {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .side {
    display: flex;
    flex-direction: column;
    grid-row: 1 / span 3;
    grid-column: 2;
    gap: var(--fs-space-24);
    padding-top: var(--fs-space-8);
  }

  .records {
    display: flex;
  }
}
</style>
