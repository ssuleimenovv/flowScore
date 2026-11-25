<script setup lang="ts">
import type { Finishing } from '../player'

// "Голы против xG" from PlayerSeason: the two bars on the desktop; the
// phone (MobilePlayerSeason) says it in a line under the difference. The
// note under them names the matches without an xG.
const { f } = defineProps<{ f: Finishing }>()
</script>

<template>
  <section class="card fs-in">
    <div class="head">
      <h2 class="title">Голы против xG</h2>
      <span class="hint">реализация</span>
    </div>
    <div class="bars">
      <span class="label">Голы</span>
      <span class="track"><span class="fill goals" :style="{ width: f.goalsWidth }" /></span>
      <span class="num">{{ f.goals }}</span>
      <span class="label">xG</span>
      <span class="track"><span class="fill xg" :style="{ width: f.xgWidth }" /></span>
      <span class="num">{{ f.xg }}</span>
    </div>
    <div class="verdict">
      <span class="diff">{{ f.diff }}</span>
      <span class="text">
        <span class="phone">{{ f.goals }} при xG {{ f.xg }}. </span>{{ f.text }}
      </span>
    </div>
    <p v-if="f.note" class="note">{{ f.note }}</p>
  </section>
</template>

<style scoped>
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

.hint,
.bars {
  display: none;
}

.hint {
  color: var(--fs-muted);
  font-size: var(--fs-text-caption);
}

.label {
  font-size: var(--fs-text-footnote);
}

.track {
  height: 10px;
  overflow: hidden;
  border-radius: var(--fs-radius-full);
  background: var(--fs-track);
}

/* The bars grow from the left, as on the board */
.fill {
  display: block;
  height: 100%;
  transform-origin: left;
  animation: grow var(--fs-duration-bar-h) var(--fs-ease) both;
}

.goals {
  background: var(--fs-home);
}

.xg {
  background: var(--fs-faint);
}

@keyframes grow {
  from {
    transform: scaleX(0);
  }
}

.num {
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-subheadline);
  font-variant-numeric: tabular-nums;
  font-weight: 700;
  text-align: right;
}

.verdict {
  display: flex;
  gap: var(--fs-space-12);
  align-items: center;
}

.diff {
  color: var(--fs-home);
  font-family: var(--fs-font-display);
  font-size: 32px;
  font-weight: 800;
  line-height: 1;
}

.text {
  color: var(--fs-muted);
  font-size: var(--fs-text-footnote);
  line-height: 1.45;
}

/* The cup matches left out: ESPN gives their goals without a place */
.note {
  padding-top: var(--fs-space-12);
  border-top: 1px solid var(--fs-line);
  color: var(--fs-muted);
  font-size: var(--fs-text-caption);
  line-height: 1.45;
}

@media (min-width: 768px) {
  .card {
    padding: 20px var(--fs-space-24);
  }

  .hint {
    display: inline;
  }

  .bars {
    display: grid;
    grid-template-columns: 80px minmax(0, 1fr) 48px;
    gap: var(--fs-space-8) var(--fs-space-12);
    align-items: center;
  }

  .verdict {
    gap: 14px;
  }

  .diff {
    font-size: 36px;
  }

  .text {
    font-size: 14px;
  }

  .phone {
    display: none;
  }
}
</style>
