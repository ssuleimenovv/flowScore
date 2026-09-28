<script setup lang="ts">
import AppButton from '@/shared/ui/AppButton.vue'
</script>

<template>
  <!-- Temporary until the Main board is built: the way into the demo replay -->
  <div class="home">
    <h1>FlowScore</h1>
    <AppButton to="/matches/3754314">Демо-матч: Ман Сити — Арсенал</AppButton>
  </div>
</template>

<style scoped>
.home {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-16);
  align-items: flex-start;
  padding: var(--fs-space-24) var(--fs-screen-padding);
}
</style>
