<script setup lang="ts">
import { computed, ref } from 'vue'
import type { MatchLineups, Side } from '@/shared/api/types'
import AppCard from '@/shared/ui/AppCard.vue'
import AppSegmented from '@/shared/ui/AppSegmented.vue'
import StateMessage from '@/shared/ui/StateMessage.vue'
import { benchRows, formationsLine, listRows, pitchPlayers } from '../lineups'
import type { MatchLive } from '../matchState'

// "Составы" from the MatchExtended boards: both elevens on one pitch or as a
// list with what each player did and their part in the Flow, then the bench.
// lineups is null before the source has them.
const { live, lineups } = defineProps<{ live: MatchLive; lineups: MatchLineups | null }>()

type View = 'pitch' | 'list'
const view = ref<View>('pitch')
const views: Array<{ value: View; label: string }> = [
  { value: 'pitch', label: 'Поле' },
  { value: 'list', label: 'Список' },
]

const SIDES: Side[] = ['home', 'away']
const teams = computed(() => {
  const l = lineups
  if (!l) return []
  return SIDES.map((side) => ({
    side,
    name: live.match[side].name,
    formation: l[side].formation,
    pitch: pitchPlayers(l[side], side, live.events),
    list: listRows(l[side], side, live.events),
    bench: benchRows(l[side]),
  }))
})
const subtitle = computed(() =>
  lineups
    ? `${formationsLine(live.match.home.name, live.match.away.name, lineups.home, lineups.away)} · справа у игрока — его вклад во Flow`
    : '',
)
</script>

<template>
  <AppCard as="section" class="lineups fs-in" style="--fs-i: 4" aria-labelledby="lineups-title">
    <div class="head">
      <div class="heading">
        <h2 id="lineups-title" class="title">Составы</h2>
        <span v-if="lineups" class="subtitle">{{ subtitle }}</span>
      </div>
      <AppSegmented v-if="lineups" v-model="view" :options="views" label="Вид составов" />
    </div>

    <StateMessage
      v-if="!lineups"
      icon="person"
      title="Составов пока нет"
      text="Появятся примерно за час до начала матча."
    />

    <template v-else>
      <div
        v-if="view === 'pitch'"
        class="pitch"
        role="img"
        :aria-label="`Расстановка на поле: ${live.match.home.name} слева, ${live.match.away.name} справа`"
      >
        <span class="halfway" />
        <span class="centre" />
        <span class="box left" />
        <span class="box right" />
        <template v-for="t in teams" :key="t.side">
          <div
            v-for="p in t.pitch"
            :key="p.id"
            class="player"
            :style="{ left: p.left, top: p.top }"
          >
            <span class="shirt" :class="t.side">
              {{ p.number }}
              <span v-if="p.goal" class="scored" />
              <span v-if="p.card" class="booking" :class="p.card" />
            </span>
            <span class="name">{{ p.name }}</span>
            <span v-if="p.off" class="off">↓ {{ p.off }}</span>
          </div>
        </template>
      </div>

      <div v-else class="lists">
        <div v-for="t in teams" :key="t.side" class="list">
          <span class="team">{{ t.formation ? `${t.name} · ${t.formation}` : t.name }}</span>
          <ul class="rows">
            <li v-for="r in t.list" :key="r.id" class="row">
              <span class="number">{{ r.number }}</span>
              <span class="player-name">{{ r.name }}</span>
              <span class="done">{{ r.events }}</span>
              <span class="flow">{{ r.flow }}</span>
            </li>
          </ul>
        </div>
      </div>

      <div class="bench">
        <span class="team">Запасные</span>
        <div class="bench-lists">
          <ul v-for="t in teams" :key="t.side" class="rows" :aria-label="`Запасные: ${t.name}`">
            <li v-for="b in t.bench" :key="b.id" class="bench-row">
              <span class="number">{{ b.number }}</span>
              <span class="bench-name">{{ b.name }}</span>
              <span v-if="b.on" class="on" :class="t.side">{{ b.on }}</span>
            </li>
          </ul>
        </div>
      </div>
    </template>
  </AppCard>
</template>

<style scoped>
.lineups {
  display: flex;
  flex-direction: column;
  gap: 18px;
  padding: var(--fs-space-24) 28px;
}

.head {
  display: flex;
  gap: var(--fs-space-16);
  align-items: center;
  justify-content: space-between;
}

.heading {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-4);
  min-width: 0;
}

.title {
  font-size: 20px;
  font-weight: 700;
  line-height: 1.2;
  letter-spacing: -0.01em;
}

.subtitle {
  color: var(--fs-muted);
  font-size: var(--fs-text-footnote);
}

/* Both elevens face each other: the home team attacks to the right */
.pitch {
  position: relative;
  height: 440px;
  overflow: hidden;
  border: 1px solid var(--fs-line);
  border-radius: 14px;
  background: var(--fs-surface-2);
}

.halfway {
  position: absolute;
  top: 0;
  bottom: 0;
  left: 50%;
  width: 1px;
  background: var(--fs-line);
}

.centre {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 120px;
  height: 120px;
  border: 1px solid var(--fs-line);
  border-radius: 50%;
  translate: -50% -50%;
}

.box {
  position: absolute;
  top: 25%;
  width: 13%;
  height: 50%;
  border: 1px solid var(--fs-line);
}

.box.left {
  left: 0;
  border-left: 0;
}

.box.right {
  right: 0;
  border-right: 0;
}

.player {
  position: absolute;
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-4);
  align-items: center;
  width: 84px;
  translate: -50% -50%;
}

.shirt {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border-radius: 50%;
  box-shadow: 0 0 0 2px var(--fs-surface);
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-footnote);
  font-variant-numeric: tabular-nums;
  font-weight: 600;
}

.shirt.home {
  background: var(--fs-home);
  color: var(--fs-on-home);
}

.shirt.away {
  background: var(--fs-away);
  color: var(--fs-on-away);
}

/* A goal: a dot on the shirt's shoulder; a card at the other one. Not
   called .card, which is AppCard's root class and would get these rules */
.scored {
  position: absolute;
  top: -4px;
  right: -6px;
  width: 14px;
  height: 14px;
  box-sizing: border-box;
  border: 2px solid var(--fs-surface);
  border-radius: 50%;
  background: var(--fs-text);
}

.booking {
  position: absolute;
  top: -4px;
  left: -4px;
  width: 9px;
  height: 12px;
  border-radius: 2px;
}

.booking.yellow {
  background: var(--fs-card-yellow);
}

.booking.red {
  background: var(--fs-card-red);
}

.name {
  font-size: var(--fs-text-caption);
  font-weight: 600;
  white-space: nowrap;
}

.off {
  color: var(--fs-muted);
  font-family: var(--fs-font-mono);
  font-size: 11px;
  font-variant-numeric: tabular-nums;
}

.lists,
.bench-lists {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 28px;
}

.list,
.bench {
  display: flex;
  flex-direction: column;
}

.bench {
  gap: var(--fs-space-8);
}

.team {
  padding-bottom: var(--fs-space-8);
  color: var(--fs-muted);
  font-size: var(--fs-text-caption);
  font-weight: 600;
}

.bench .team {
  padding-bottom: 0;
}

.rows {
  display: flex;
  flex-direction: column;
  padding: 0;
  list-style: none;
}

.row {
  display: grid;
  grid-template-columns: 32px minmax(0, 1fr) auto 44px;
  gap: 10px;
  align-items: center;
  padding: 9px 0;
  border-top: 1px solid var(--fs-line);
}

.bench-row {
  display: grid;
  grid-template-columns: 28px minmax(0, 1fr) auto;
  gap: 10px;
  align-items: center;
  padding: var(--fs-space-8) 0;
  border-top: 1px solid var(--fs-line);
}

.number {
  color: var(--fs-muted);
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-footnote);
  font-variant-numeric: tabular-nums;
}

.bench-row .number {
  font-size: var(--fs-text-caption);
}

/* A long name is cut rather than wrapped, so every row keeps its height */
.player-name,
.bench-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.player-name {
  font-size: var(--fs-text-subheadline);
  font-weight: 600;
}

.bench-name {
  font-size: 14px;
  font-weight: 500;
}

.done {
  color: var(--fs-muted);
  font-size: var(--fs-text-caption);
}

.flow {
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-footnote);
  font-variant-numeric: tabular-nums;
  font-weight: 600;
  text-align: right;
}

.on {
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-caption);
  font-variant-numeric: tabular-nums;
  font-weight: 600;
}

.on.home {
  color: var(--fs-home);
}

.on.away {
  color: var(--fs-away);
}

/* Phone, as on MobileMatchExtended: a lower pitch with smaller shirts, one
   list under the other, the bench still side by side */
@media (max-width: 767px) {
  .lineups {
    gap: 14px;
    padding: var(--fs-space-16);
  }

  .subtitle {
    display: none;
  }

  .pitch {
    height: 300px;
  }

  /* The lines are about 40 px apart here: a longer name is cut short so
     it does not run into the next line's */
  .player {
    width: 42px;
  }

  .shirt {
    width: 26px;
    height: 26px;
    font-size: 11px;
  }

  .name {
    max-width: 100%;
    overflow: hidden;
    font-size: 10px;
    text-overflow: ellipsis;
  }

  .lists {
    grid-template-columns: 1fr;
    gap: var(--fs-space-16);
  }

  .bench-lists {
    gap: var(--fs-space-16);
  }
}
</style>
