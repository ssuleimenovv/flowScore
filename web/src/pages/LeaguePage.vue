<script setup lang="ts">
// The Leagues screen (League and MobileLeague boards): the table with each
// team's average Flow, the matches of a round and the season's Flow records.
// A European cup adds its knockouts beside the table, and a national cup has
// its knockouts in place of the table (CupsExtended board).
// The league and the round live in the URL (?league=ger.1&round=6), so "Назад"
// and a shared link come back to them.
import { computed, onScopeDispose, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import PlayoffCard from '@/features/cup/ui/PlayoffCard.vue'
import {
  competitionGroups,
  competitionKind,
  isEuroCup,
  roundTabs,
  seasonLine,
} from '@/features/league/league'
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

const groups = computed(() => competitionGroups(leagues.items.value))
const kind = computed(() => competitionKind(leagueId.value))

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

// Nothing to show yet: a league before its table, a cup before its first
// stage. The Super Cup has no stages, only the round of its one match.
const isEmpty = computed(() => {
  const l = league.value
  if (!l || l.stages.length > 0) return false
  return kind.value === 'cup' ? l.rounds === 0 : standings.value.length === 0
})

const { online } = useNetwork()
const screen = computed(() =>
  resolveScreenState({
    status: leagues.request.value === 'error' ? 'error' : request.value,
    hasData: league.value !== null,
    isEmpty: isEmpty.value,
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
        <h1 v-if="kind === 'cup'" class="h1">
          <span class="desktop-word">Сетка <span class="accent">кубка</span></span>
          <span class="phone-word">Кубок</span>
        </h1>
        <h1 v-else class="h1">
          <span class="desktop-word">Турнирная <span class="accent">таблица</span></span>
          <span class="phone-word">Таблица</span>
        </h1>
      </div>
      <div v-if="leagues.items.value.length > 1" class="groups">
        <div v-for="(g, i) in groups" :key="g.name" class="group">
          <span :id="`group-${i}`" class="group-name">{{ g.name }}</span>
          <div class="chips" role="group" :aria-labelledby="`group-${i}`">
            <button
              v-for="l in g.items"
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
        </div>
      </div>
    </header>

    <OfflineBanner v-if="screen.banner === 'offline'" class="banner" :updated-at="updatedAt" />

    <!-- A national cup: its knockouts in place of the table, or the Super
         Cup's one match; the records beside them -->
    <template v-if="screen.view === 'content' && league && kind === 'cup'">
      <div class="main">
        <PlayoffCard
          v-if="league.stages.length > 0"
          main
          :league-id="league.id"
          :stages="league.stages"
          :current="league.stage"
        />
        <RoundMatches
          v-else
          v-model="round"
          :tabs
          :matches="matches.items.value"
          :loading="matches.request.value !== 'success'"
          :now
        />
      </div>
      <aside class="side">
        <FlowRecords :records="league.records" :stages="league.stages" class="records" />
      </aside>
    </template>

    <!-- A league, or a European cup: once its knockouts begin, they take the
         place of the rounds of its league phase -->
    <template v-else-if="screen.view === 'content' && league">
      <div class="main">
        <StandingsTable :rows="standings" :cup="isEuroCup(leagueId)" />
      </div>
      <aside class="side">
        <PlayoffCard
          v-if="league.stages.length > 0"
          class="playoff"
          :league-id="league.id"
          :stages="league.stages"
          :current="league.stage"
        />
        <RoundMatches
          v-else
          v-model="round"
          :tabs
          :matches="matches.items.value"
          :loading="matches.request.value !== 'success'"
          :now
        />
        <FlowRecords :records="league.records" :stages="league.stages" class="records" />
      </aside>
    </template>

    <AppCard v-else-if="screen.view === 'empty' && kind === 'cup'" class="main notice">
      <StateMessage
        icon="trophy"
        title="Сетки пока нет"
        text="Кубок ещё не начался. Пары появятся после жеребьёвки."
      />
    </AppCard>

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

/* The chips in groups, each under its name (CompetitionsExtended). On the
   phone a group is one row that scrolls from edge to edge of the screen. */
.groups {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.group-name {
  color: var(--fs-muted);
  font-size: var(--fs-text-caption);
  font-weight: 600;
}

.chips {
  display: flex;
  gap: var(--fs-space-8);
  margin: 0 calc(-1 * var(--fs-screen-padding));
  padding: 0 var(--fs-screen-padding);
  overflow-x: auto;
  scrollbar-width: none;
}

/* On the phone, as on the cup board, the knockouts come before the table */
.top,
.banner {
  order: -2;
}

.playoff {
  order: -1;
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
    margin: 0;
    padding: 0;
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
