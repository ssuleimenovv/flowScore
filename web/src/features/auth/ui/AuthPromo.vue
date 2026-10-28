<script setup lang="ts">
import AppIcon, { type IconName } from '@/shared/ui/AppIcon.vue'
import AppLogo from '@/shared/ui/AppLogo.vue'
import { flowBars } from '../auth'

// The left half of the Auth board: what FlowScore is, over a Flow wave
const bars = flowBars()
const features: Array<{ icon: IconName; text: string }> = [
  { icon: 'pulse', text: 'Flow в реальном времени' },
  { icon: 'sparkle', text: 'AI объясняет каждый скачок' },
  { icon: 'sliders', text: 'Симулятор «Что, если…?»' },
]
</script>

<template>
  <section class="promo">
    <RouterLink to="/" class="brand fs-link" aria-label="FlowScore, на главную">
      <AppLogo />
      <span class="wordmark" aria-hidden="true">FLOWSCORE</span>
    </RouterLink>

    <div class="pitch">
      <h1 class="title">Почувствуй, кто <span class="accent">владеет игрой</span></h1>
      <div class="wave" aria-hidden="true">
        <div v-for="(b, i) in bars" :key="i" class="column">
          <div class="half up">
            <span class="bar" :class="{ home: b.home }" :style="{ height: `${b.up}px` }" />
          </div>
          <div class="axis" />
          <div class="half down">
            <span class="bar" :class="{ away: b.away }" :style="{ height: `${b.down}px` }" />
          </div>
        </div>
      </div>
    </div>

    <ul class="features">
      <li v-for="f in features" :key="f.text" class="feature">
        <AppIcon :name="f.icon" :size="22" class="feature-icon" />
        {{ f.text }}
      </li>
    </ul>
  </section>
</template>

<style scoped>
.promo {
  display: flex;
  flex-direction: column;
  flex-grow: 1;
  justify-content: space-between;
  gap: var(--fs-space-32);
  padding: 48px 56px;
  border-right: 1px solid var(--fs-line);
  background: var(--fs-surface);
}

.brand {
  display: flex;
  gap: 10px;
  align-items: center;
  align-self: flex-start;
  color: var(--fs-text);
  text-decoration: none;
}

.wordmark {
  font-family: var(--fs-font-display);
  font-size: 28px;
  font-weight: 800;
  letter-spacing: 0.05em;
}

.pitch {
  display: flex;
  flex-direction: column;
  gap: 28px;
}

.title {
  font-size: 64px;
  font-weight: 700;
  line-height: 1.05;
  letter-spacing: -0.025em;
}

.accent {
  color: var(--fs-home);
}

/* 40 columns: a bar grows up from the line, another down from it */
.wave {
  display: flex;
  gap: 4px;
  height: 201px;
}

.column {
  display: flex;
  flex: 1 1 0;
  flex-direction: column;
}

.half {
  display: flex;
  height: 100px;
}

.up {
  align-items: flex-end;
}

.down {
  align-items: flex-start;
}

.axis {
  height: 1px;
  background: var(--fs-line);
}

.bar {
  width: 100%;
  background: var(--fs-line);
  animation: grow var(--fs-duration-bar) var(--fs-ease) both;
}

.up .bar {
  border-radius: 3px 3px 0 0;
  transform-origin: bottom;
}

.down .bar {
  border-radius: 0 0 3px 3px;
  transform-origin: top;
}

.bar.home {
  background: var(--fs-home);
}

.bar.away {
  background: var(--fs-away);
}

.features {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 20px;
  padding: 0;
  list-style: none;
}

.feature {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-8);
  font-size: var(--fs-text-subheadline);
  font-weight: 600;
}

.feature-icon {
  color: var(--fs-home);
}

@keyframes grow {
  from {
    transform: scaleY(0);
  }
}
</style>
