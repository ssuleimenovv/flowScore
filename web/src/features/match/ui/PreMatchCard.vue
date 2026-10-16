<script setup lang="ts">
import { computed, ref } from 'vue'
import type { MatchPreMatch } from '@/shared/api/types'
import AppCard from '@/shared/ui/AppCard.vue'
import AppSegmented from '@/shared/ui/AppSegmented.vue'
import type { MatchLive } from '../matchState'
import { formRows, headToHead, meetingRows } from '../prematch'
import TeamBadge from './TeamBadge.vue'

// "Перед матчем" from the MatchExtended boards: each team's last five
// matches with its Flow in them, or how the two teams' meetings went.
const { live, pre } = defineProps<{ live: MatchLive; pre: MatchPreMatch }>()

type Tab = 'form' | 'h2h'
const tab = ref<Tab>('form')
const tabs: Array<{ value: Tab; label: string }> = [
  { value: 'form', label: 'Форма' },
  { value: 'h2h', label: 'Очные встречи' },
]

const form = computed(() => formRows(pre, live.match))
const h2h = computed(() => headToHead(pre, live.match))
const meetings = computed(() => meetingRows(pre))
</script>

<template>
  <AppCard as="section" class="prematch fs-in" style="--fs-i: 5" aria-labelledby="prematch-title">
    <div class="head">
      <h2 id="prematch-title" class="title">Перед матчем</h2>
      <AppSegmented v-model="tab" :options="tabs" label="Что показать перед матчем" />
    </div>

    <ul v-if="tab === 'form'" class="rows">
      <li v-for="f in form" :key="f.side" class="form-row">
        <span class="team">
          <TeamBadge :code="f.code" :side="f.side" :logo="f.logo" size="sm" />
          <span class="name">{{ f.name }}</span>
        </span>
        <span
          class="results"
          :aria-label="`Последние матчи: ${f.results.map((r) => r.letter).join(' ')}`"
        >
          <span
            v-for="r in f.results"
            :key="r.id"
            class="result"
            :class="r.outcome"
            :title="r.title"
            aria-hidden="true"
          >
            {{ r.letter }}
          </span>
        </span>
        <span class="form-flow">
          <span class="flow-value" :class="f.side">{{ f.flow ?? '—' }}</span>
          <span class="flow-note">{{ f.flowNote }}</span>
        </span>
      </li>
    </ul>

    <div v-else class="h2h">
      <div v-if="meetings.length" class="tally">
        <div class="tally-labels">
          <span class="home">{{ h2h.home }}</span>
          <span class="draw">{{ h2h.draws }}</span>
          <span class="away">{{ h2h.away }}</span>
        </div>
        <div class="tally-bar" aria-hidden="true">
          <span class="home" :style="{ width: h2h.widths.home }" />
          <span class="draw" :style="{ width: h2h.widths.draw }" />
          <span class="away" :style="{ width: h2h.widths.away }" />
        </div>
        <span class="note">{{ h2h.note }}</span>
      </div>
      <p v-else class="note">Команды ещё не встречались.</p>

      <ul class="rows">
        <li v-for="m in meetings" :key="m.id" class="meeting">
          <span class="date">{{ m.date }}</span>
          <span class="what">
            <span class="name">{{ m.title }}</span>
            <span class="competition"
              >{{ m.competition }}<span class="phone-date"> · {{ m.date }}</span></span
            >
          </span>
          <span class="score">{{ m.score }}</span>
          <span class="meeting-flow">{{ m.flow ? `Flow ${m.flow}` : '' }}</span>
        </li>
      </ul>
    </div>
  </AppCard>
</template>

<style scoped>
.prematch {
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

.rows {
  display: flex;
  flex-direction: column;
  padding: 0;
  list-style: none;
}

/* EVE Everton   Н Н В В Н   49 */
.form-row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto auto;
  gap: var(--fs-space-16);
  align-items: center;
  padding: 14px 0;
  border-top: 1px solid var(--fs-line);
}

.team {
  display: flex;
  gap: var(--fs-space-12);
  align-items: center;
  min-width: 0;
}

.name {
  overflow: hidden;
  font-size: var(--fs-text-subheadline);
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.results {
  display: flex;
  gap: 6px;
}

.result {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 8px;
  font-size: var(--fs-text-footnote);
  font-weight: 700;
}

.result.win {
  background: var(--fs-ok);
  color: var(--fs-on-accent);
}

.result.draw {
  background: var(--fs-fill);
  color: var(--fs-text);
}

.result.loss {
  background: var(--fs-home);
  color: var(--fs-on-accent);
}

.form-flow {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.flow-value {
  font-family: var(--fs-font-display);
  font-size: 24px;
  font-weight: 800;
  line-height: 1;
}

.flow-value.home {
  color: var(--fs-home);
}

.flow-value.away {
  color: var(--fs-away);
}

.flow-note,
.note,
.competition,
.date,
.meeting-flow {
  color: var(--fs-muted);
  font-size: var(--fs-text-caption);
}

.h2h {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.tally {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-8);
}

.tally-labels {
  display: flex;
  justify-content: space-between;
  gap: var(--fs-space-8);
  font-size: var(--fs-text-footnote);
  font-weight: 600;
}

.tally-labels .home {
  color: var(--fs-home);
}

.tally-labels .draw {
  color: var(--fs-muted);
}

.tally-labels .away {
  color: var(--fs-away);
}

.tally-bar {
  display: flex;
  gap: 2px;
  height: 8px;
  overflow: hidden;
  border-radius: var(--fs-radius-full);
}

.tally-bar .home {
  background: var(--fs-home);
}

.tally-bar .draw {
  background: var(--fs-track);
}

.tally-bar .away {
  background: var(--fs-away);
}

/* 21 мар 2026   Everton — Chelsea   3:0   Flow 58–61 */
.meeting {
  display: grid;
  grid-template-columns: 110px minmax(0, 1fr) 52px 90px;
  gap: var(--fs-space-12);
  align-items: center;
  padding: 11px 0;
  border-top: 1px solid var(--fs-line);
}

.date,
.meeting-flow,
.score {
  font-family: var(--fs-font-mono);
  font-variant-numeric: tabular-nums;
}

.what {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.phone-date {
  display: none;
}

.score {
  font-size: var(--fs-text-subheadline);
  font-weight: 700;
}

.meeting-flow {
  text-align: right;
}

/* Phone, as on MobileMatchExtended: no Flow column, the date under the
   meeting, smaller result squares */
@media (max-width: 767px) {
  .prematch {
    padding: var(--fs-space-16);
  }

  .head {
    flex-wrap: wrap;
  }

  .form-row {
    grid-template-columns: minmax(0, 1fr) auto;
  }

  .form-flow,
  .date,
  .meeting-flow {
    display: none;
  }

  .result {
    width: 24px;
    height: 24px;
  }

  .meeting {
    grid-template-columns: minmax(0, 1fr) 44px;
  }

  .phone-date {
    display: inline;
  }
}
</style>
