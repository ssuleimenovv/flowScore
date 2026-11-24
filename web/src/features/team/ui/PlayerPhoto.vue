<script setup lang="ts">
// A player's photo square from the Team board. ESPN has photos of few
// players, so it holds the board's silhouette. "number" pins the shirt
// number to its corner, as on MobileTeam; the desktop has a column for it.
const { number = '' } = defineProps<{ number?: string }>()
</script>

<template>
  <span class="photo" aria-hidden="true">
    <svg viewBox="0 0 48 48" preserveAspectRatio="xMidYMax meet">
      <circle cx="24" cy="19" r="8.5" />
      <path d="M7 48c1.5-10 8.5-15.5 17-15.5S39.5 38 41 48z" />
    </svg>
    <span v-if="number" class="number">{{ number }}</span>
  </span>
</template>

<style scoped>
/* Phone, as on MobileTeam */
.photo {
  position: relative;
  display: block;
  flex-shrink: 0;
  width: 48px;
  height: 48px;
  overflow: hidden;
  border-radius: 10px;
  background: var(--fs-surface-2);
}

svg {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  fill: var(--fs-faint);
}

.number {
  position: absolute;
  right: 3px;
  bottom: 3px;
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 16px;
  height: 16px;
  padding: 0 var(--fs-space-4);
  border-radius: 5px;
  background: var(--fs-accent);
  color: var(--fs-on-accent);
  font-family: var(--fs-font-display);
  font-size: 11px;
  font-weight: 800;
}

/* Desktop, as on Team: smaller, the number has its own column */
@media (min-width: 768px) {
  .photo {
    width: 44px;
    height: 44px;
    border-radius: 9px;
  }

  .number {
    display: none;
  }
}
</style>
