<script setup lang="ts">
import type { LogRow } from '../player'

// "Матчи сезона" from PlayerSeason: a table on the desktop, two lines a
// match on the phone. The Flow column comes with D6. Old matches have no
// page of their own yet, so the rows lead nowhere.
const { rows } = defineProps<{ rows: LogRow[] }>()
</script>

<template>
  <section class="card fs-in">
    <div class="head">
      <h2 class="title">Матчи сезона</h2>
      <span class="hint">журнал</span>
    </div>
    <p v-if="rows.length === 0" class="empty">В этом сезоне игрок пока не выходил на поле.</p>
    <template v-else>
      <div class="row header" aria-hidden="true">
        <span>Дата</span>
        <span>Соперник</span>
        <span />
        <span>Счёт</span>
        <span>Мин</span>
        <span>Г</span>
        <span>А</span>
        <span>xG</span>
      </div>
      <ul class="list">
        <li v-for="r in rows" :key="r.id" class="row">
          <span class="date mono">{{ r.date }}</span>
          <span class="who">
            <span class="opponent"
              >{{ r.opponent }}<span class="phone"> · {{ r.score }}</span></span
            >
            <span class="phone sub mono">{{ r.date }} · {{ r.minutes }}′ · xG {{ r.xg }}</span>
          </span>
          <span class="venue">{{ r.venue }}</span>
          <span class="score mono">{{ r.score }}</span>
          <span class="minutes mono">{{ r.minutes }}</span>
          <span class="goals mono">{{ r.goals }}</span>
          <span class="assists mono">{{ r.assists }}</span>
          <span class="xg mono">{{ r.xg }}</span>
          <span class="phone ga mono" :aria-label="`голы ${r.goals}, ассисты ${r.assists}`">
            {{ r.goals }} + {{ r.assists }}
          </span>
        </li>
      </ul>
      <span class="legend">Д — дома, Г — в гостях.</span>
    </template>
  </section>
</template>

<style scoped>
/* Phone, as on MobilePlayerSeason */
.card {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: var(--fs-space-16);
  border: 1px solid var(--fs-line);
  border-radius: var(--fs-radius-lg);
  background: var(--fs-surface);
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
  letter-spacing: -0.01em;
  line-height: 1.2;
}

.list {
  margin: 0;
  padding: 0;
  list-style: none;
}

.row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 10px;
  align-items: center;
  min-height: 56px;
  border-top: 1px solid var(--fs-line);
}

.mono {
  font-family: var(--fs-font-mono);
  font-variant-numeric: tabular-nums;
}

.hint,
.header,
.date,
.venue,
.score,
.minutes,
.goals,
.assists,
.xg,
.legend {
  display: none;
}

.who {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.opponent {
  overflow: hidden;
  font-size: var(--fs-text-subheadline);
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.sub,
.hint,
.legend {
  color: var(--fs-muted);
  font-size: var(--fs-text-caption);
}

.ga {
  font-size: var(--fs-text-footnote);
}

.empty {
  color: var(--fs-muted);
  font-size: var(--fs-text-footnote);
}

/* Desktop, as on PlayerSeason: a table */
@media (min-width: 768px) {
  .card {
    padding: 22px 28px;
  }

  .row {
    grid-template-columns: 80px minmax(0, 1fr) 36px 52px 44px 32px 32px 56px;
    gap: var(--fs-space-12);
    min-height: 0;
    padding: var(--fs-space-12) 0;
  }

  .header {
    display: grid;
    padding: 0 0 10px;
    border-top: 0;
    color: var(--fs-muted);
    font-size: var(--fs-text-caption);
    font-weight: 600;
  }

  .hint,
  .date,
  .venue,
  .score,
  .minutes,
  .goals,
  .assists,
  .xg,
  .legend {
    display: block;
  }

  .phone {
    display: none;
  }

  .date,
  .venue {
    color: var(--fs-muted);
    font-size: var(--fs-text-caption);
  }

  .score,
  .minutes,
  .goals,
  .assists,
  .xg {
    font-size: var(--fs-text-subheadline);
  }

  .minutes {
    color: var(--fs-muted);
  }

  .goals {
    font-weight: 700;
  }
}
</style>
