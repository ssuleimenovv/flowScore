<script setup lang="ts">
// Temporary: marks a screen that is not built yet, so the navigation can be tried
import { useRoute } from 'vue-router'

const route = useRoute()
</script>

<template>
  <section class="stub">
    <h1>{{ route.meta.stubTitle }}</h1>
    <p>Экран еще не сделан</p>
  </section>
</template>

<style scoped>
.stub {
  display: grid;
  gap: var(--fs-space-8);
  padding: var(--fs-space-24) var(--fs-screen-padding);
}

p {
  color: var(--fs-muted);
}
</style>
