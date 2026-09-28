<script setup lang="ts">
// The card surface from the mockup. Padding and layout come from the caller's class.
const { as = 'div' } = defineProps<{ as?: string }>()
</script>

<template>
  <component :is="as" class="card">
    <slot />
  </component>
</template>

<style scoped>
.card {
  border: 1px solid var(--fs-line);
  border-radius: var(--fs-radius-lg);
  background: var(--fs-surface);
}
</style>
