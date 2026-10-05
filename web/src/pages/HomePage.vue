<script setup lang="ts">
// The home screen (Main and MobileHome boards): what is on now and what is next.
import { computed, onScopeDispose, ref } from 'vue'
import { competitions, dayForecast, dayTitle, flowPeaks, matchDay } from '@/features/home/home'
import DayForecast from '@/features/home/ui/DayForecast.vue'
import FlowPeaks from '@/features/home/ui/FlowPeaks.vue'
import LiveMatchCard from '@/features/home/ui/LiveMatchCard.vue'
import UpcomingList from '@/features/home/ui/UpcomingList.vue'
import { useMatchList } from '@/features/home/useMatchList'
import AppCard from '@/shared/ui/AppCard.vue'
import AppSkeleton from '@/shared/ui/AppSkeleton.vue'
import ErrorState from '@/shared/ui/ErrorState.vue'
import WakingNote from '@/shared/ui/WakingNote.vue'

const { items, request, error, retry } = useMatchList()

// One clock for every card: the minutes tick on between polls
const now = ref(Date.now())
const timer = setInterval(() => (now.value = Date.now()), 1000)
onScopeDispose(() => clearInterval(timer))

// null is "Все"
const competition = ref<string | null>(null)
const chips = computed(() => [null, ...competitions(items.value)])

const day = computed(() => matchDay(items.value, competition.value))
const peaks = computed(() => flowPeaks(day.value.live))
const forecast = computed(() => dayForecast(day.value.upcoming))
const today = dayTitle(new Date())
</script>

<template>
  <div class="page">
    <header class="top">
      <div class="heading">
        <span class="date">{{ today }}</span>
        <h1 class="h1">Поток <span class="accent">сегодня</span></h1>
      </div>
      <div class="chips" role="group" aria-label="Лиги">
        <button
          v-for="c in chips"
          :key="c ?? 'all'"
          type="button"
          class="chip fs-press"
          :aria-pressed="competition === c"
          @click="competition = c"
        >
          {{ c ?? 'Все' }}
        </button>
      </div>
    </header>

    <template v-if="request === 'success'">
      <div class="main">
        <section class="live" aria-labelledby="live-title">
          <div class="section-head">
            <h2 id="live-title" class="h2"><span class="index">01</span>Сейчас в эфире</h2>
            <span class="count">live: {{ day.live.length }}</span>
          </div>
          <div v-if="day.live.length" class="cards">
            <LiveMatchCard
              v-for="(m, i) in day.live"
              :key="m.id"
              :match="m"
              :now
              class="fs-in"
              :style="{ '--fs-i': i + 1 }"
            />
          </div>
          <AppCard v-else class="nothing fs-in">Сейчас нет live-матчей</AppCard>
        </section>

        <section class="later" aria-labelledby="later-title">
          <h2 id="later-title" class="h2"><span class="index">02</span>Позже сегодня</h2>
          <UpcomingList :matches="day.upcoming" />
        </section>
      </div>

      <aside class="side">
        <FlowPeaks :peaks class="peaks" />
        <DayForecast
          v-if="forecast"
          :match="forecast.match"
          :text="forecast.text"
          class="forecast"
        />
      </aside>
    </template>

    <div v-else-if="request === 'error'" class="main">
      <AppCard>
        <ErrorState title="Не удалось загрузить матчи" :error @retry="retry" />
      </AppCard>
    </div>

    <div v-else class="main" aria-busy="true" aria-label="Загрузка матчей">
      <WakingNote />
      <div class="cards">
        <AppCard v-for="i in 2" :key="i" class="loading">
          <AppSkeleton width="40%" :height="14" />
          <AppSkeleton :height="32" />
          <AppSkeleton :height="32" />
          <AppSkeleton :height="56" />
        </AppCard>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Phone: one column, in the order of MobileHome: live, peaks, later */
.page {
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding: var(--fs-space-8) var(--fs-screen-padding) var(--fs-space-24);
}

.main,
.side {
  display: contents;
}

.top {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.heading {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-8);
}

.date {
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
  flex-shrink: 0;
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

.chip[aria-pressed='true'] {
  border-color: var(--fs-text);
  background: var(--fs-text);
  color: var(--fs-bg);
}

.live,
.later {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-12);
}

.live {
  order: 1;
}

.peaks {
  order: 2;
}

.later {
  order: 3;
}

/* The day's forecast is a desktop card, as on the boards */
.forecast {
  display: none;
}

.section-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.h2 {
  font-size: var(--fs-text-title-3);
  font-weight: 700;
  letter-spacing: -0.01em;
  line-height: 1.2;
}

.index {
  display: none;
  margin-right: var(--fs-space-12);
  color: var(--fs-home);
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-footnote);
  vertical-align: middle;
}

.count {
  color: var(--fs-muted);
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-footnote);
  font-variant-numeric: tabular-nums;
}

.cards {
  display: grid;
  gap: var(--fs-space-12);
}

.nothing {
  padding: 32px;
  color: var(--fs-muted);
  font-size: var(--fs-text-subheadline);
  text-align: center;
}

.loading {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-12);
  padding: 20px;
}

/* Desktop: the live cards two by two, the side cards in a column of their own */
@media (min-width: 768px) {
  .page {
    gap: 32px;
    padding: 32px var(--fs-page-padding) var(--fs-space-40);
  }

  .main,
  .side {
    display: flex;
    flex-direction: column;
    gap: 32px;
    min-width: 0;
  }

  .side {
    gap: var(--fs-space-24);
  }

  .h1 {
    font-size: 48px;
  }

  .h2 {
    font-size: 24px;
    letter-spacing: -0.02em;
  }

  .index {
    display: inline;
  }

  .live {
    gap: var(--fs-space-16);
  }

  .cards {
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: var(--fs-space-16);
  }

  .forecast {
    display: flex;
  }

  /* The phone's order does not apply here: peaks first, as on Main */
  .peaks {
    order: 0;
  }
}

/* Wide desktop: the side column runs next to the heading and the lists, as on Main */
@media (min-width: 1200px) {
  .page {
    display: grid;
    grid-template-columns: minmax(0, 1fr) 380px;
    align-items: start;
  }

  .top {
    grid-column: 1;
  }

  .side {
    grid-row: 1 / span 2;
    grid-column: 2;
    padding-top: 4px;
  }
}
</style>
