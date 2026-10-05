<script setup lang="ts">
import { computed, ref } from 'vue'
import AppButton from '@/shared/ui/AppButton.vue'
import AppCard from '@/shared/ui/AppCard.vue'
import { chronicle } from '../chronicle'
import type { MatchLive } from '../matchState'

// "Хроника" from the Match boards: the moments that moved the flow, newest on top.
const { live } = defineProps<{ live: MatchLive }>()

// Enough to fill the sidebar next to the wave; the rest on demand
const LIMIT = 10
const expanded = ref(false)

const items = computed(() => chronicle(live.events, live.match))
const visible = computed(() => (expanded.value ? items.value : items.value.slice(0, LIMIT)))
</script>

<template>
  <AppCard as="section" class="chronicle fs-in" style="--fs-i: 3" aria-labelledby="chronicle-title">
    <div class="head">
      <h2 id="chronicle-title" class="title">Хроника</h2>
      <span class="hint">Flow-эффект события</span>
    </div>

    <!-- On the desktop the list is a window of fixed height: new events and rows
         of one or two lines do not move the cards under it -->
    <div class="viewport" :class="{ open: expanded }">
      <p v-if="items.length === 0" class="empty">Первые события появятся здесь</p>

      <!-- A new event slides in on top and pushes the others down -->
      <TransitionGroup v-else tag="ol" name="row" class="list">
        <li v-for="item in visible" :key="item.id" class="row" :class="{ goal: item.goal }">
          <span class="minute">{{ item.minute }}</span>
          <span class="code" :class="item.goal" aria-hidden="true">{{ item.code }}</span>
          <div class="body">
            <span class="event">{{ item.title }}</span>
            <span v-if="item.detail" class="detail">{{ item.detail }}</span>
          </div>
          <span v-if="item.impact" class="impact" :class="item.impact.side">
            <span class="fs-sr-only">Flow </span>{{ item.impact.text }}
          </span>
        </li>
      </TransitionGroup>
    </div>

    <!-- The button keeps its place while there is nothing more to show -->
    <AppButton
      variant="secondary"
      size="small"
      class="more"
      :class="{ idle: items.length <= LIMIT }"
      :disabled="items.length <= LIMIT"
      :aria-hidden="items.length <= LIMIT || undefined"
      :aria-expanded="expanded"
      @click="expanded = !expanded"
    >
      {{ expanded ? 'Свернуть' : `Показать все · ${items.length}` }}
    </AppButton>
  </AppCard>
</template>

<style scoped>
.chronicle {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-12);
  padding: 22px 20px;
}

.head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 var(--fs-space-4);
}

.title {
  font-size: var(--fs-text-title-3);
  font-weight: 700;
  line-height: 1.2;
  letter-spacing: -0.01em;
}

.hint {
  color: var(--fs-muted);
  font-size: var(--fs-text-caption);
}

.empty {
  padding: var(--fs-space-12) var(--fs-space-8);
  color: var(--fs-muted);
  font-size: var(--fs-text-subheadline);
}

.list {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-4);
  padding: 0;
  list-style: none;
}

.row {
  display: grid;
  grid-template-columns: 36px 42px minmax(0, 1fr) auto;
  gap: 10px;
  align-items: center;
  padding: 10px var(--fs-space-8);
  border-radius: var(--fs-radius-md);
}

.row.goal {
  background: var(--fs-surface-2);
}

.minute {
  color: var(--fs-muted);
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-footnote);
  font-variant-numeric: tabular-nums;
}

.code {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 24px;
  border-radius: 6px;
  background: var(--fs-surface-2);
  color: var(--fs-muted);
  font-size: var(--fs-text-caption);
  font-weight: 700;
  letter-spacing: 0.04em;
}

.code.home {
  background: var(--fs-home);
  color: var(--fs-on-home);
}

.code.away {
  background: var(--fs-away);
  color: var(--fs-on-away);
}

.body {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.event,
.detail {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.event {
  font-size: var(--fs-text-subheadline);
  font-weight: 600;
}

.detail {
  color: var(--fs-muted);
  font-size: var(--fs-text-caption);
}

.impact {
  padding: 3px var(--fs-space-8);
  border-radius: 6px;
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-caption);
  font-variant-numeric: tabular-nums;
  font-weight: 600;
}

.impact.home {
  background: var(--fs-home-soft);
  color: var(--fs-home);
}

.impact.away {
  background: var(--fs-away-soft);
  color: var(--fs-away);
}

.more {
  align-self: center;
}

.more.idle {
  visibility: hidden;
}

/* Desktop: ten rows fit, the eleventh fades out under the edge */
@media (min-width: 768px) {
  .viewport:not(.open) {
    height: 540px;
    overflow: hidden;
    mask-image: linear-gradient(to bottom, #000 88%, transparent);
  }
}

.row-enter-active,
.row-move {
  transition:
    opacity var(--fs-duration-lift) var(--fs-ease),
    transform var(--fs-duration-lift) var(--fs-ease);
}

.row-enter-from {
  opacity: 0;
  transform: translateY(-8px);
}

/* Phone: the list fills the "События" tab on its own, without a header */
@media (max-width: 767px) {
  .chronicle {
    padding: var(--fs-space-12) var(--fs-space-8);
  }

  .head {
    display: none;
  }

  /* The tab has nothing under the list, so nothing to keep in place */
  .more.idle {
    display: none;
  }

  .detail,
  .impact {
    font-size: var(--fs-text-footnote);
  }
}
</style>
