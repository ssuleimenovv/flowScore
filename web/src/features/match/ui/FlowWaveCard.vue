<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, useTemplateRef, watch } from 'vue'
import AppCard from '@/shared/ui/AppCard.vue'
import AppSegmented from '@/shared/ui/AppSegmented.vue'
import type { MatchLive } from '../matchState'
import { buildWave, goalMarks } from '../wave'

// "01 Волна потока" from the Match boards: who controls the game minute by
// minute. Up is the home side, down the away side.
const { live, seconds } = defineProps<{ live: MatchLive; seconds: number }>()

const windows = [
  { value: 1, label: '1 мин' },
  { value: 5, label: '5 мин' },
  { value: 15, label: '15 мин' },
]
const smoothing = ref(5)

const wave = computed(() => buildWave(live.points, smoothing.value))
const total = computed(() => wave.value.length)
const marks = computed(() => goalMarks(live.events))
const home = computed(() => live.match.home.name)
const away = computed(() => live.match.away.name)

// Where a minute's column sits, as a share of the chart width
const center = (minute: number) => `${((minute - 0.5) / total.value) * 100}%`

// The clock ticks on the client and can lag a fast replay; the wave never
// shows "now" behind a minute it has already drawn.
const nowMinute = computed(() => {
  if (live.match.status !== 'live') return null
  const played = live.points[live.points.length - 1]?.minute ?? 0
  return Math.min(Math.max(Math.floor(seconds / 60), played), total.value)
})

// Bar height inside its half: a percentage of the half, never thinner than 1 px.
// Unplayed minutes get a 2 px stub, as on the board.
function height(value: number | null): string {
  if (value === null) return '2px'
  return value > 0 ? `max(1px, ${value}%)` : '1px'
}

function tone(value: number | null, side: 'up' | 'down'): string {
  if (value === null) return 'future'
  if (side === 'up') return value > 0 ? 'home' : 'neutral'
  return value < 0 ? 'away' : 'neutral'
}

// Hover (or drag a finger) to read one minute
const hovered = ref<number | null>(null)
const columns = useTemplateRef('columns')

function pick(e: PointerEvent) {
  const el = columns.value
  if (!el) return
  const box = el.getBoundingClientRect()
  const index = Math.floor(((e.clientX - box.left) / box.width) * total.value)
  hovered.value = Math.min(total.value - 1, Math.max(0, index))
}

const tipLeft = computed(() => (hovered.value === null ? '0' : center(hovered.value + 1)))

const tip = computed(() => {
  if (hovered.value === null) return null
  const value = wave.value[hovered.value] ?? null
  const minute = `${hovered.value + 1}′`
  if (value === null) return `${minute} · ещё не сыграно`
  const n = Math.round(value)
  if (n === 0) return `${minute} · поровну`
  return n > 0 ? `${minute} · ${home.value} +${n}` : `${minute} · ${away.value} +${-n}`
})

// What a screen reader gets instead of 90 bars
const summary = computed(() => {
  const played = wave.value.filter((v) => v !== null)
  const last = Math.round(played[played.length - 1] ?? 0)
  const side = last > 0 ? home.value : last < 0 ? away.value : null
  return side
    ? `Волна потока по минутам. Сейчас перевес у ${side}: +${Math.abs(last)}`
    : 'Волна потока по минутам. Сейчас поровну'
})

// Goal pills can collide (two goals a few minutes apart) or stick out of the
// card at the edges. After each render we measure them, keep them inside and
// move a colliding pill one row down.
const MARK_ROW = 30
const marksBox = useTemplateRef('marksBox')
const placed = ref<Array<{ left: number; row: number }>>([])

function placeMarks() {
  const box = marksBox.value
  if (!box) return
  const width = box.clientWidth
  const rowEnds: number[] = []
  // DOM children keep the v-for order, so child i is marks[i]
  placed.value = Array.from(box.children, (pill, i) => {
    const w = (pill as HTMLElement).offsetWidth
    const mid = ((marks.value[i]!.minute - 0.5) / total.value) * width
    const left = Math.min(Math.max(0, mid - w / 2), width - w)
    let row = rowEnds.findIndex((end) => end + 6 <= left)
    if (row === -1) row = rowEnds.length
    rowEnds[row] = left + w
    return { left, row }
  })
}

function markStyle(i: number) {
  const place = placed.value[i]
  const mark = marks.value[i]!
  // Before the first measurement: centered on its minute
  if (!place) return { left: center(mark.minute), transform: 'translateX(-50%)' }
  return { left: `${place.left}px`, top: `${place.row * MARK_ROW}px` }
}

const marksHeight = computed(() => {
  const rows = Math.max(1, ...placed.value.map((p) => p.row + 1))
  return `${(rows - 1) * MARK_ROW + 26}px`
})

watch(marks, placeMarks, { flush: 'post', immediate: true })

let resize: ResizeObserver | undefined
onMounted(() => {
  resize = new ResizeObserver(placeMarks)
  if (marksBox.value) resize.observe(marksBox.value)
  // The pill width changes once the web font arrives
  void document.fonts.ready.then(placeMarks)
})
onUnmounted(() => resize?.disconnect())

const ticks = [0, 15, 30, 45, 60, 75, 90]
</script>

<template>
  <AppCard as="section" class="wave fs-in" style="--fs-i: 2" aria-labelledby="wave-title">
    <div class="head">
      <div class="titles">
        <h2 id="wave-title" class="title"><span class="index">01</span>Волна потока</h2>
        <p class="lead">Кто контролирует игру по минутам. Вверх — {{ home }}, вниз — {{ away }}.</p>
      </div>
      <AppSegmented v-model="smoothing" :options="windows" label="Окно сглаживания" />
    </div>

    <div
      class="chart"
      @pointermove="pick"
      @pointerdown="pick"
      @pointerleave="hovered = null"
      @pointercancel="hovered = null"
    >
      <span v-if="tip" class="tip" :style="{ left: tipLeft }">{{ tip }}</span>

      <div
        ref="columns"
        class="columns"
        :class="{ dimmed: hovered !== null }"
        role="img"
        :aria-label="summary"
      >
        <div
          v-for="(value, i) in wave"
          :key="i"
          class="column"
          :class="{ active: i === hovered }"
          :style="{ '--i': i }"
        >
          <div class="half top">
            <div class="bar up" :class="tone(value, 'up')" :style="{ height: height(value) }" />
          </div>
          <div class="axis" />
          <div class="half bottom">
            <div
              class="bar down"
              :class="tone(value, 'down')"
              :style="{ height: height(value === null ? null : -value) }"
            />
          </div>
        </div>
      </div>

      <template v-if="nowMinute !== null">
        <div class="now-line" :style="{ left: `${(nowMinute / total) * 100}%` }" />
        <span class="now-label" :style="{ left: `${(nowMinute / total) * 100}%` }">
          <span class="now-word">сейчас </span>{{ nowMinute }}′
        </span>
      </template>
    </div>

    <div ref="marksBox" class="marks" :style="{ height: marksHeight }">
      <span
        v-for="(mark, i) in marks"
        :key="`${mark.minute}-${mark.label}`"
        class="mark"
        :class="mark.side"
        :style="markStyle(i)"
      >
        {{ mark.label }}
      </span>
    </div>

    <div class="ticks" aria-hidden="true">
      <span v-for="t in ticks" :key="t" :class="{ minor: t % 45 !== 0 }">{{ t }}′</span>
    </div>
  </AppCard>
</template>

<style scoped>
.wave {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: var(--fs-space-24) 28px;
}

.head {
  display: flex;
  gap: var(--fs-space-16);
  align-items: flex-start;
  justify-content: space-between;
}

.titles {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-8);
}

.title {
  font-size: 24px;
  font-weight: 700;
  line-height: 1.2;
  letter-spacing: -0.02em;
}

.index {
  margin-right: var(--fs-space-12);
  color: var(--fs-home);
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-footnote);
  vertical-align: middle;
}

.lead {
  color: var(--fs-muted);
  font-size: var(--fs-text-subheadline);
}

.chart {
  position: relative;
  padding-top: 22px;
  /* A horizontal drag reads the minutes; a vertical one still scrolls the page */
  touch-action: pan-y;
}

.columns {
  display: flex;
  gap: 2px;
}

.column {
  display: flex;
  flex: 1 1 0;
  flex-direction: column;
  min-width: 0;
  transition: opacity var(--fs-duration-fade);
}

.dimmed .column:not(.active) {
  opacity: 0.35;
}

.half {
  display: flex;
  height: 100px;
}

.top {
  align-items: flex-end;
}

.bottom {
  align-items: flex-start;
}

.axis {
  height: 1px;
  background: var(--fs-line);
}

/* Bars grow out of the axis one after another; a new minute grows on its own */
.bar {
  width: 100%;
  animation: grow var(--fs-duration-bar) var(--fs-ease) both;
  animation-delay: calc(var(--i) * var(--fs-stagger) / 8);
  transition:
    height var(--fs-duration-grow) var(--fs-ease),
    background-color var(--fs-duration-color);
}

.up {
  border-radius: 2px 2px 0 0;
  transform-origin: bottom;
}

.down {
  border-radius: 0 0 2px 2px;
  transform-origin: top;
}

@keyframes grow {
  from {
    transform: scaleY(0);
  }
}

.home {
  background: var(--fs-home);
}

.away {
  background: var(--fs-away);
}

.neutral {
  background: var(--fs-line);
}

.future {
  background: var(--fs-track);
}

.tip {
  position: absolute;
  top: -10px;
  z-index: 2;
  display: flex;
  align-items: center;
  height: 26px;
  padding: 0 10px;
  border-radius: var(--fs-radius-sm);
  background: var(--fs-text);
  color: var(--fs-bg);
  font-size: var(--fs-text-caption);
  font-weight: 600;
  white-space: nowrap;
  transform: translateX(-50%);
  pointer-events: none;
}

.now-line {
  position: absolute;
  top: 18px;
  bottom: 0;
  border-left: 1px dashed var(--fs-home);
  animation: pulse var(--fs-duration-pulse) ease-in-out infinite;
  pointer-events: none;
}

@keyframes pulse {
  50% {
    opacity: 0.4;
  }
}

.now-label {
  position: absolute;
  top: 0;
  color: var(--fs-muted);
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-caption);
  font-variant-numeric: tabular-nums;
  white-space: nowrap;
  transform: translateX(-50%);
}

.marks {
  position: relative;
  transition: height var(--fs-duration-nudge) var(--fs-ease);
}

.mark {
  position: absolute;
  top: 0;
  display: flex;
  align-items: center;
  height: 24px;
  padding: 0 10px;
  border-radius: var(--fs-radius-full);
  font-size: var(--fs-text-caption);
  font-weight: 600;
  white-space: nowrap;
  transition:
    left var(--fs-duration-nudge) var(--fs-ease),
    top var(--fs-duration-nudge) var(--fs-ease);
}

.mark.home {
  color: var(--fs-on-home);
}

.mark.away {
  color: var(--fs-on-away);
}

.ticks {
  display: flex;
  justify-content: space-between;
  color: var(--fs-muted);
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-caption);
  font-variant-numeric: tabular-nums;
}

@media (prefers-reduced-motion: reduce) {
  .now-line {
    animation: none;
  }
}

/* Phone: the MobileMatch card is lower and plainer */
@media (max-width: 767px) {
  .wave {
    gap: var(--fs-space-12);
    padding: 18px var(--fs-space-16);
  }

  .head {
    align-items: center;
  }

  .title {
    font-size: var(--fs-text-title-3);
    letter-spacing: -0.01em;
  }

  .index,
  .lead,
  .now-word,
  .minor {
    display: none;
  }

  .head :deep(.option) {
    height: 32px;
  }

  .columns {
    gap: 1px;
  }

  .half {
    height: 74px;
  }

  .tip,
  .mark {
    font-size: var(--fs-text-footnote);
  }
}
</style>
