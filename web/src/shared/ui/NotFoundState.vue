<script setup lang="ts">
import AppButton from './AppButton.vue'

// 404 from the States board. The texts differ per screen: a match, a team, a page.
const { title, text } = defineProps<{ title: string; text: string }>()
</script>

<template>
  <div class="not-found">
    <span class="code" aria-hidden="true">404</span>
    <h1 class="title">{{ title }}</h1>
    <p class="text">{{ text }}</p>
    <AppButton variant="secondary" to="/">На главную</AppButton>
  </div>
</template>

<style scoped>
.not-found {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-12);
  align-items: center;
  padding: var(--fs-space-40) 28px;
  text-align: center;
}

.code {
  color: var(--fs-home);
  font-family: var(--fs-font-display);
  font-size: 120px;
  font-weight: 800;
  line-height: 0.8;
}

.title {
  font-size: var(--fs-text-title-3);
  font-weight: 600;
}

.text {
  color: var(--fs-muted);
  font-size: var(--fs-text-subheadline);
}
</style>
