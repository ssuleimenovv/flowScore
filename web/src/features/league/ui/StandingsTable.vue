<script setup lang="ts">
import TeamBadge from '@/features/match/ui/TeamBadge.vue'
import type { Standing } from '@/shared/api/types'
import { teamLogo } from '@/shared/media/logos'
import AppCard from '@/shared/ui/AppCard.vue'
import { flowText, flowWidth, goals, zoneMark } from '../league'

// "Турнирная таблица" from the League boards. The desktop has every column;
// the phone keeps the place, the team, games, points and Flow, as on MobileLeague.
const { rows } = defineProps<{ rows: Standing[] }>()
</script>

<template>
  <AppCard as="section" class="table fs-in" style="--fs-i: 1" aria-label="Турнирная таблица">
    <div class="grid head" aria-hidden="true">
      <span>#</span>
      <span>Команда</span>
      <span>И</span>
      <span class="wide">В</span>
      <span class="wide">Н</span>
      <span class="wide">П</span>
      <span class="wide">Мячи</span>
      <span>О</span>
      <span><span class="wide">Средний </span>Flow</span>
    </div>
    <div v-for="r in rows" :key="r.team.id" class="grid row">
      <span class="place">
        <span class="zone" :class="zoneMark(r.zone)" />
        <span class="num">{{ r.rank }}</span>
      </span>
      <span class="team">
        <TeamBadge
          class="code wide"
          :code="r.team.code"
          :logo="teamLogo(r.team.id, 30)"
          size="xs"
        />
        <span class="name">{{ r.team.name }}</span>
      </span>
      <span class="num muted" :aria-label="`игр: ${r.played}`">{{ r.played }}</span>
      <span class="num wide">{{ r.won }}</span>
      <span class="num wide">{{ r.drawn }}</span>
      <span class="num wide">{{ r.lost }}</span>
      <span class="num muted wide">{{ goals(r) }}</span>
      <span class="num points" :aria-label="`очков: ${r.points}`">{{ r.points }}</span>
      <span class="flow" :aria-label="`средний Flow: ${flowText(r.flow)}`">
        <span class="track wide" aria-hidden="true">
          <span class="bar" :style="{ width: flowWidth(r.flow) }" />
        </span>
        <span class="num value">{{ flowText(r.flow) }}</span>
      </span>
    </div>
    <div class="legend">
      <span class="key"><span class="mark champions" />Лига чемпионов</span>
      <span class="key"><span class="mark europa" />Лига Европы</span>
    </div>
  </AppCard>
</template>

<style scoped>
/* Phone, as on MobileLeague: place, team, games, points, Flow */
.table {
  padding: var(--fs-space-8) 14px;
}

.grid {
  display: grid;
  grid-template-columns: 30px minmax(0, 1fr) 26px 34px 40px;
  gap: var(--fs-space-8);
  align-items: center;
}

.wide {
  display: none;
}

.head {
  padding: var(--fs-space-8) 0;
  color: var(--fs-muted);
  font-size: var(--fs-text-caption);
  font-weight: 600;
}

.row {
  min-height: 48px;
  border-top: 1px solid var(--fs-line);
  color: var(--fs-text);
  transition: background-color var(--fs-duration-color);
}

@media (hover: hover) {
  .row:hover {
    background: var(--fs-surface-2);
  }
}

.place {
  display: flex;
  gap: 6px;
  align-items: center;
}

.zone {
  width: 3px;
  height: 18px;
  border-radius: 2px;
  background: var(--fs-track);
}

.champions {
  background: var(--fs-home);
}

.europa {
  background: var(--fs-away);
}

.num {
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-footnote);
  font-variant-numeric: tabular-nums;
}

.place .num {
  font-weight: 600;
}

.muted {
  color: var(--fs-muted);
}

.points {
  font-size: var(--fs-text-subheadline);
  font-weight: 700;
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
  white-space: nowrap;
  text-overflow: ellipsis;
}

/* The crest is a desktop column, as on League; MobileLeague has none. Two
   classes, so it outweighs TeamBadge's own display */
.code.wide {
  display: none;
}

.flow .value {
  color: var(--fs-home);
}

.legend {
  display: none;
}

/* Desktop, as on League: every column, the Flow bar, the legend */
@media (min-width: 768px) {
  .table {
    padding: 20px 28px;
  }

  .grid {
    grid-template-columns: 44px minmax(0, 1fr) repeat(4, 40px) 64px 48px 170px;
    gap: var(--fs-space-12);
  }

  .wide {
    display: inline;
  }

  .code.wide {
    display: flex;
  }

  .head {
    padding: 0 0 var(--fs-space-12);
  }

  .row {
    min-height: 0;
    padding: var(--fs-space-12) 0;
  }

  .place {
    gap: var(--fs-space-8);
  }

  .zone {
    width: 4px;
    height: 22px;
  }

  .num {
    font-size: var(--fs-text-subheadline);
  }

  .flow {
    display: flex;
    gap: 10px;
    align-items: center;
  }

  .flow .value {
    width: 24px;
    color: var(--fs-text);
    font-size: var(--fs-text-footnote);
  }

  .track {
    display: block;
    flex-grow: 1;
    height: 6px;
    overflow: hidden;
    border-radius: var(--fs-radius-full);
    background: var(--fs-track);
  }

  .bar {
    display: block;
    height: 100%;
    border-radius: var(--fs-radius-full);
    background: var(--fs-home);
    transform-origin: left;
    animation: grow var(--fs-duration-bar-h) var(--fs-ease) both;
  }

  .legend {
    display: flex;
    gap: 20px;
    padding-top: var(--fs-space-16);
    color: var(--fs-muted);
    font-size: var(--fs-text-caption);
  }

  .key {
    display: flex;
    gap: var(--fs-space-8);
    align-items: center;
  }

  .mark {
    width: 4px;
    height: 14px;
    border-radius: 2px;
  }
}

@keyframes grow {
  from {
    transform: scaleX(0);
  }
}
</style>
