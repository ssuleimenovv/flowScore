<script setup lang="ts">
import AppIcon, { type IconName } from './AppIcon.vue'

// The live event banner from the States board ("ГОЛ! Фоден · 58′"): the loudest
// surface in the app, so it is kept for goals and big flow swings.
const { icon, title, text, time } = defineProps<{
  icon: IconName
  title: string
  text: string
  time: string
}>()
</script>

<template>
  <div class="banner" role="status">
    <span class="tile">
      <AppIcon :name="icon" :size="22" />
    </span>
    <div class="body">
      <span class="title">{{ title }}</span>
      <span class="text">{{ text }}</span>
    </div>
    <span class="time">{{ time }}</span>
  </div>
</template>

<style scoped>
.banner {
  display: grid;
  grid-template-columns: 44px minmax(0, 1fr) auto;
  gap: var(--fs-space-12);
  align-items: center;
  padding: var(--fs-space-12) 14px;
  border-radius: 18px;
  background: var(--fs-accent);
  color: var(--fs-on-accent);
}

.tile {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  border-radius: var(--fs-radius-md);
  background: color-mix(in srgb, var(--fs-on-accent) 18%, transparent);
}

.body {
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.title {
  font-size: var(--fs-text-subheadline);
  font-weight: 700;
}

.text {
  font-size: var(--fs-text-footnote);
}

.time {
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-caption);
  font-variant-numeric: tabular-nums;
}
</style>
