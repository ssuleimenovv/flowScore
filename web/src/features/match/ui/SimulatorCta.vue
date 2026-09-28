<script setup lang="ts">
import AppButton from '@/shared/ui/AppButton.vue'
import AppIcon from '@/shared/ui/AppIcon.vue'

// "Что, если…?" from the Match boards: the way from a live match into the simulator.
</script>

<template>
  <section class="cta" aria-labelledby="cta-title">
    <h2 id="cta-title" class="title">Что, если…?</h2>
    <p class="text">
      Убери игрока, добавь красную на 75′ или включи дождь — AI пересчитает исход и поток.
    </p>
    <AppButton to="/simulator" class="go">
      Открыть симулятор
      <AppIcon name="arrow" :size="18" class="arrow" />
    </AppButton>
  </section>
</template>

<style scoped>
/* Inverted colors, like the offline banner: it stands apart from the data cards */
.cta {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: var(--fs-space-24);
  border-radius: var(--fs-radius-lg);
  background: var(--fs-text);
  color: var(--fs-bg);
}

.title {
  font-size: 28px;
  font-weight: 700;
  line-height: 1.1;
  letter-spacing: -0.02em;
}

.text {
  font-size: var(--fs-text-subheadline);
  line-height: 1.5;
}

.go {
  align-self: flex-start;
}

.arrow {
  transition: transform var(--fs-duration-nudge) var(--fs-ease);
}

/* The arrow leans toward where the button leads */
@media (hover: hover) {
  .go:hover .arrow {
    transform: translateX(4px);
  }
}

/* Phone: a full-width button with a bigger touch target, as on MobileMatch */
@media (max-width: 767px) {
  .cta {
    gap: var(--fs-space-12);
    padding: 20px;
  }

  .title {
    font-size: 23px;
  }

  /* .cta in front outweighs the button's own size rules */
  .cta .go {
    align-self: stretch;
    height: 50px;
    padding: 0 20px;
    border-radius: 14px;
    font-size: var(--fs-text-body);
  }
}
</style>
