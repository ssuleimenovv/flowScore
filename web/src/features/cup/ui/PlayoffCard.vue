<script setup lang="ts">
import {
  computed,
  nextTick,
  onMounted,
  onScopeDispose,
  ref,
  useId,
  useTemplateRef,
  watch,
  type ComponentPublicInstance,
} from 'vue'
import TeamBadge from '@/features/match/ui/TeamBadge.vue'
import type { Stage } from '@/shared/api/types'
import { teamLogo } from '@/shared/media/logos'
import AppCard from '@/shared/ui/AppCard.vue'
import AppSegmented from '@/shared/ui/AppSegmented.vue'
import AppSkeleton from '@/shared/ui/AppSkeleton.vue'
import ErrorState from '@/shared/ui/ErrorState.vue'
import {
  isTwoLegged,
  legsNote,
  matchesNote,
  openingStage,
  tieCells,
  tieColumns,
  tieFootnote,
  tieMark,
  tieNote,
} from '../cup'
import { useTies } from '../useTies'

// "Плей-офф" from the CupsExtended boards: the stages, and a stage's ties, a
// line a team with its goals on the right. On a European cup's page it is the
// card beside the league phase; "main" is the domestic cup page's own column,
// headed by the stage and the number of its matches.
const {
  leagueId,
  stages,
  current,
  main = false,
} = defineProps<{ leagueId: string; stages: Stage[]; current: string; main?: boolean }>()

// The stage picked, else the one being played; a new cup starts there again
const picked = ref<string | null>(null)
watch(
  () => leagueId,
  () => (picked.value = null),
)
const stageId = computed<string>({
  get: () => {
    const id = picked.value
    return id && stages.some((s) => s.id === id) ? id : openingStage(stages, current)
  },
  set: (id) => (picked.value = id),
})
const stage = computed(() => stages.find((s) => s.id === stageId.value) ?? null)

const ties = useTies(
  computed(() => leagueId || null),
  computed(() => stageId.value || null),
)
const items = computed(() => ties.items.value)
const twoLegs = computed(() => isTwoLegged(items.value))
const loading = computed(() => ties.request.value === 'idle' || ties.request.value === 'pending')

const title = computed(() => (main ? (stage.value?.name ?? 'Плей-офф') : 'Плей-офф'))
const side = computed(() => {
  if (main) return stage.value ? matchesNote(stage.value, items.value) : ''
  return loading.value ? '' : legsNote(twoLegs.value)
})

// "ответный 11 мар" and "19:00" follow the clock
const now = ref(new Date())
const timer = setInterval(() => (now.value = new Date()), 60_000)
onScopeDispose(() => clearInterval(timer))

// Many stages scroll sideways, as on the cup page's board. They start from
// the first, moved only as far as the stage picked needs to be seen.
const tabs = useTemplateRef<ComponentPublicInstance>('tabs')
function reveal() {
  const row = tabs.value?.$el as HTMLElement | undefined
  const on = row?.querySelector<HTMLElement>('[aria-pressed="true"]')
  if (!row || !on) return
  row.scrollLeft = Math.max(0, on.offsetLeft + on.offsetWidth + 2 - row.clientWidth)
}
onMounted(reveal)
watch(
  () => [leagueId, stages],
  () => void nextTick(reveal),
)

// A tie as the board draws it: two lines, then the note under them
const rows = computed(() =>
  items.value.map((t) => ({
    key: `${t.home.team.id}-${t.away.team.id}`,
    lines: [t.home, t.away].map((s) => ({
      team: s.team,
      mark: tieMark(s),
      cells: tieCells(s, twoLegs.value),
    })),
    note: tieNote(t, now.value),
  })),
)

const titleId = useId()
</script>

<template>
  <AppCard
    as="section"
    class="playoff fs-in"
    :class="{ lead: main }"
    style="--fs-i: 2"
    :aria-labelledby="titleId"
  >
    <div class="head">
      <h2 :id="titleId" class="title">{{ title }}</h2>
      <span class="side">{{ side }}</span>
    </div>
    <AppSegmented
      v-if="stages.length > 1"
      ref="tabs"
      v-model="stageId"
      class="tabs"
      label="Этап"
      :options="stages.map((s) => ({ value: s.id, label: s.name }))"
    />

    <div v-if="loading" class="ties" aria-busy="true" aria-label="Загрузка пар">
      <div v-for="i in 3" :key="i" class="tie">
        <div v-for="j in 2" :key="j" class="line">
          <span class="mark open" />
          <AppSkeleton :height="22" width="22px" round />
          <AppSkeleton :width="j === 1 ? '60%' : '45%'" :height="15" />
          <AppSkeleton width="30px" :height="13" />
        </div>
      </div>
    </div>

    <ErrorState
      v-else-if="ties.request.value === 'error'"
      class="error"
      title="Не удалось загрузить пары"
      :error="ties.error.value"
      @retry="ties.retry"
    />

    <div v-else :key="stageId" class="ties fs-in">
      <div v-if="items.length" class="columns" aria-hidden="true">
        <span v-for="c in tieColumns(twoLegs)" :key="c" class="cell">{{ c }}</span>
      </div>
      <div v-for="r in rows" :key="r.key" class="tie">
        <div
          v-for="l in r.lines"
          :key="l.team.id"
          class="line"
          :aria-label="`${l.team.name}: ${l.cells.join(', ')}`"
        >
          <span class="mark" :class="l.mark" />
          <TeamBadge
            class="crest"
            :code="l.team.code"
            :logo="teamLogo(l.team.id, 22)"
            size="mini"
          />
          <span class="name" :class="{ muted: l.mark === 'out' }">{{ l.team.name }}</span>
          <span class="cells" aria-hidden="true">
            <span
              v-for="(c, i) in l.cells"
              :key="i"
              class="cell"
              :class="{ last: i === l.cells.length - 1 }"
              >{{ c }}</span
            >
          </span>
        </div>
        <div v-if="r.note" class="note">
          <span v-if="r.note.live" class="live">
            <span class="dot" aria-hidden="true" />{{ r.note.text }}
          </span>
          <span v-else class="when">{{ r.note.text }}</span>
        </div>
      </div>
      <p v-if="items.length === 0" class="none">Пар этого этапа пока нет</p>
    </div>

    <p v-if="!loading && items.length" class="footnote">{{ tieFootnote(twoLegs) }}</p>
  </AppCard>
</template>

<style scoped>
/* Phone, as on the CupsExtended phone frames. The lines are set tighter than
   the page's text, as the board's are */
.playoff {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-12);
  min-width: 0;
  padding: var(--fs-space-16);
  line-height: 1.25;
}

.head {
  display: flex;
  gap: var(--fs-space-12);
  justify-content: space-between;
  align-items: center;
}

.title {
  font-size: var(--fs-text-title-3);
  font-weight: 700;
  letter-spacing: -0.01em;
  line-height: 1.2;
}

.side,
.footnote {
  color: var(--fs-muted);
  font-size: var(--fs-text-caption);
}

.side {
  flex-shrink: 0;
}

/* Few stages share the width; many keep their size and scroll */
.tabs {
  position: relative;
  overflow-x: auto;
  scrollbar-width: none;
}

.tabs::-webkit-scrollbar {
  display: none;
}

.tabs :deep(.option) {
  flex: 1 1 0;
  height: 32px;
  padding: 0 14px;
}

.lead .tabs :deep(.option) {
  flex: 0 0 auto;
}

.ties {
  display: flex;
  flex-direction: column;
}

.columns {
  display: flex;
  gap: 6px;
  justify-content: flex-end;
  padding-top: 2px;
  color: var(--fs-muted);
  font-size: 11px;
  font-weight: 600;
}

.tie {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-8);
  padding: var(--fs-space-12) 0;
  border-top: 1px solid var(--fs-line);
}

.line {
  display: grid;
  grid-template-columns: 4px 22px minmax(0, 1fr) auto;
  gap: 10px;
  align-items: center;
}

.mark {
  width: 4px;
  height: 20px;
  border-radius: 2px;
}

.mark.through {
  background: var(--fs-ok);
}

.mark.out {
  background: var(--fs-faint);
}

.mark.open {
  background: var(--fs-track);
}

/* A crest is 22 px on the board, a three-letter tile round */
.line .crest {
  width: 22px;
  height: 22px;
  border-radius: 50%;
}

.name {
  overflow: hidden;
  color: var(--fs-text);
  font-size: var(--fs-text-subheadline);
  font-weight: 600;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.name.muted {
  color: var(--fs-muted);
}

.cells {
  display: flex;
  gap: 6px;
}

.cell {
  width: 30px;
  text-align: right;
}

.cells .cell {
  color: var(--fs-muted);
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-footnote);
  font-variant-numeric: tabular-nums;
  font-weight: 500;
}

.cells .last {
  color: var(--fs-text);
  font-weight: 700;
}

.note {
  display: flex;
  gap: var(--fs-space-8);
  align-items: center;
  padding-left: 46px;
}

.when {
  color: var(--fs-muted);
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-caption);
}

.live {
  display: inline-flex;
  gap: 6px;
  align-items: center;
  height: 22px;
  padding: 0 9px;
  border-radius: var(--fs-radius-full);
  background: var(--fs-live-soft);
  color: var(--fs-live);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.04em;
  white-space: nowrap;
}

.dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--fs-live);
  animation: pulse var(--fs-duration-pulse) ease-in-out infinite;
}

.none {
  padding: var(--fs-space-12) 0;
  border-top: 1px solid var(--fs-line);
  color: var(--fs-muted);
  font-size: var(--fs-text-subheadline);
}

@keyframes pulse {
  50% {
    opacity: 0.4;
  }
}

@media (prefers-reduced-motion: reduce) {
  .dot {
    animation: none;
  }
}

/* Desktop: the European card beside the league phase, the cup's own column as wide as the table */
@media (min-width: 768px) {
  .playoff {
    padding: 20px var(--fs-space-24);
  }

  .playoff.lead {
    padding: 20px 28px;
  }
}
</style>
