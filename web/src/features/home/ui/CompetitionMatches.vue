<script setup lang="ts">
import { computed } from 'vue'
import { competitionAccent, competitionArt, tintStyle } from '@/features/league/accents'
import { competitionKind } from '@/features/league/league'
import CompetitionArt from '@/features/league/ui/CompetitionArt.vue'
import { leagueLogo, needsPlate } from '@/shared/media/logos'
import type { DayCompetition } from '../home'
import UpcomingList from './UpcomingList.vue'

// One competition's matches of a day (HomeByCompetitionExtended): a head on
// the competition's artwork, the same as its page's, and the rows below on
// the theme's card. A cup stands on its league's artwork; one without
// artwork gets its accent, or the theme's plain fill.
const { competition } = defineProps<{ competition: DayCompetition }>()

const art = computed(() => competitionArt(competition.id))
const accent = computed(() => competitionAccent(competition.id))
const logo = computed(() => leagueLogo(competition.id, 28, 'dark'))
const count = computed(() => {
  const n = competition.matches.length
  const ones = n % 10
  const tens = n % 100
  if (ones === 1 && tens !== 11) return `${n} матч`
  if (ones >= 2 && ones <= 4 && (tens < 12 || tens > 14)) return `${n} матча`
  return `${n} матчей`
})
// A cup has its knockout stages where a league has its table
const to = computed(() => (competitionKind(competition.id) === 'cup' ? 'Сетка' : 'Таблица'))
</script>

<template>
  <section class="card fs-in" :aria-label="competition.name">
    <div
      class="head fs-tinted"
      :class="{ painted: art || accent }"
      :style="[
        tintStyle(accent),
        art ? { '--fs-art-deep': art.deep, '--fs-art-word': art.word } : {},
      ]"
    >
      <CompetitionArt v-if="art" :art="art.key" />
      <div class="line">
        <span v-if="logo" class="logo">
          <img :class="{ plate: needsPlate(competition.id) }" :src="logo" alt="" />
        </span>
        <span class="titles">
          <span class="name">{{ competition.name }}</span>
          <span class="count">{{ count }}</span>
        </span>
        <RouterLink
          class="go fs-press"
          :to="{ name: 'leagues', query: { league: competition.id } }"
          :aria-label="`${to}: ${competition.name}`"
        >
          <span class="go-word">{{ to }}</span> →
        </RouterLink>
      </div>
    </div>
    <UpcomingList :matches="competition.matches" bare />
  </section>
</template>

<style scoped>
.card {
  overflow: hidden;
  border: 1px solid var(--fs-line);
  border-radius: var(--fs-radius-lg);
  background: var(--fs-surface);
}

/* The head: white on the artwork in both themes, as on the competition's
   page; the artwork's left side, under the text, is held dark */
.head {
  position: relative;
  height: 76px;
  overflow: hidden;
  background: var(--fs-surface-2);
  color: var(--fs-text);
}

.head.painted {
  background: var(--fs-art-deep, var(--fs-tint));
  color: var(--fs-on-art);
}

.line {
  position: relative;
  display: flex;
  gap: var(--fs-space-12);
  align-items: center;
  height: 100%;
  padding: 0 var(--fs-space-16);
}

.logo {
  display: inline-flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: 1px solid var(--fs-art-glass-line);
  border-radius: 10px;
  background: var(--fs-art-glass);
  backdrop-filter: blur(10px);
}

.logo img {
  width: 22px;
  height: 22px;
  object-fit: contain;
}

.logo img.plate {
  padding: 2px;
  border-radius: var(--fs-radius-full);
  background: var(--fs-on-accent);
}

.titles {
  display: flex;
  flex: 1;
  flex-direction: column;
  gap: 3px;
  min-width: 0;
}

.name {
  overflow: hidden;
  font-size: var(--fs-text-headline);
  font-weight: 700;
  letter-spacing: -0.01em;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.count {
  color: var(--fs-art-word, var(--fs-on-art-muted));
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-caption);
  font-variant-numeric: tabular-nums;
}

.head:not(.painted) .count {
  color: var(--fs-muted);
}

.go {
  display: inline-flex;
  flex-shrink: 0;
  align-items: center;
  gap: 4px;
  height: 32px;
  padding: 0 var(--fs-space-12);
  border: 1px solid var(--fs-art-glass-line);
  border-radius: var(--fs-radius-full);
  background: var(--fs-art-glass);
  color: inherit;
  font-size: var(--fs-text-footnote);
  font-weight: 600;
  text-decoration: none;
  backdrop-filter: blur(10px);
}

.go:hover {
  background: var(--fs-art-glass-hover);
}

.go-word {
  display: none;
}

.head:not(.painted) .logo,
.head:not(.painted) .go {
  border-color: var(--fs-line);
  background: var(--fs-surface);
}

@media (min-width: 768px) {
  .head {
    height: 84px;
  }

  .line {
    padding: 0 22px;
  }

  .name {
    font-size: 19px;
  }

  .go-word {
    display: inline;
  }
}

@media (prefers-reduced-transparency: reduce) {
  .logo,
  .go {
    background: var(--fs-art-glass-solid);
    backdrop-filter: none;
  }
}
</style>
