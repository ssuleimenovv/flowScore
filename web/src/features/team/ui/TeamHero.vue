<script setup lang="ts">
import { computed } from 'vue'
import { flowText } from '@/features/league/league'
import TeamBadge from '@/features/match/ui/TeamBadge.vue'
import type { TeamPage } from '@/shared/api/types'
import { teamLogo } from '@/shared/media/logos'
import { averageFlow, form, leagueLine } from '../team'
import ResultChip from './ResultChip.vue'

// The team's card from the Team board: crest, league and place, name, form,
// and its Flow averaged over the season. On the phone (MobileTeam) it lies
// right on the page, and the league and place go up into the nav bar.
const { page } = defineProps<{ page: TeamPage }>()

const results = computed(() => form(page.matches))
const flow = computed(() => averageFlow(page.matches, page.team.id))
</script>

<template>
  <section class="card fs-in">
    <TeamBadge
      class="badge"
      :code="page.team.code"
      side="home"
      :logo="teamLogo(page.team.id, 112)"
      size="hero"
    />
    <div class="head">
      <span v-if="page.league" class="league">{{ leagueLine(page.league) }}</span>
      <h1 class="name">{{ page.team.name }}</h1>
    </div>
    <div v-if="results.length > 0" class="form">
      <span class="label">Форма</span>
      <ol class="results" aria-label="Последние матчи, от старых к новым">
        <li v-for="(r, i) in results" :key="i"><ResultChip :result="r" /></li>
      </ol>
    </div>
    <div class="flow">
      <span class="value">{{ flowText(flow) }}</span>
      <span class="caption">
        <span class="word">FLOW</span>
        <span class="note">средний за сезон</span>
      </span>
    </div>
  </section>
</template>

<style scoped>
/* Phone, as on MobileTeam: crest and name, then the form and the Flow */
.card {
  display: grid;
  grid-template-areas:
    'badge head head'
    'form form flow';
  grid-template-columns: auto minmax(0, 1fr) auto;
  gap: 14px;
  align-items: center;
}

.badge {
  grid-area: badge;
}

.head {
  display: flex;
  grid-area: head;
  flex-direction: column;
  gap: var(--fs-space-12);
  min-width: 0;
}

.league {
  display: none;
  color: var(--fs-muted);
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-footnote);
  font-variant-numeric: tabular-nums;
}

.name {
  font-size: 34px;
  font-weight: 700;
  letter-spacing: -0.025em;
  line-height: 1.05;
  overflow-wrap: anywhere;
}

.form {
  display: flex;
  grid-area: form;
  gap: var(--fs-space-16);
  align-items: center;
}

/* The board's phone shows the chips without the word */
.label {
  display: none;
  color: var(--fs-muted);
  font-size: var(--fs-text-footnote);
}

.results {
  display: flex;
  gap: 6px;
  margin: 0;
  padding: 0;
  list-style: none;
}

.flow {
  display: flex;
  grid-area: flow;
  gap: 6px;
  align-items: baseline;
  justify-self: end;
}

.value {
  color: var(--fs-tint);
  font-family: var(--fs-font-display);
  font-size: 44px;
  font-weight: 800;
  line-height: 0.8;
}

.caption {
  display: flex;
  flex-direction: column;
}

.word {
  color: var(--fs-muted);
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 0.12em;
}

.note {
  display: none;
  color: var(--fs-muted);
  font-size: var(--fs-text-caption);
}

/* Desktop, as on Team: one card, the Flow at its right edge */
@media (min-width: 768px) {
  .card {
    grid-template-areas:
      'badge head flow'
      'badge form flow';
    column-gap: 28px;
    row-gap: var(--fs-space-12);
    padding: var(--fs-space-32);
    border: 1px solid var(--fs-line);
    border-radius: var(--fs-radius-lg);
    background: var(--fs-surface);
  }

  .league,
  .label {
    display: inline;
  }

  .name {
    font-size: 56px;
  }

  .flow {
    gap: 10px;
    align-items: flex-end;
  }

  .value {
    font-size: 88px;
  }

  .word {
    font-size: var(--fs-text-caption);
  }

  .note {
    display: inline;
  }
}
</style>
