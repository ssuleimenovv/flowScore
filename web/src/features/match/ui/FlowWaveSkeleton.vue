<script setup lang="ts">
import AppCard from '@/shared/ui/AppCard.vue'
import AppSkeleton from '@/shared/ui/AppSkeleton.vue'

// The wave card while it loads, as on the States board
const bars = Array.from({ length: 30 }, (_, i) => 20 + ((i * 37) % 70))
</script>

<template>
  <AppCard class="loading" aria-busy="true" aria-label="Загрузка волны">
    <AppSkeleton width="140px" :height="18" />
    <div class="bars">
      <AppSkeleton v-for="(h, i) in bars" :key="i" class="bar" :height="h" />
    </div>
  </AppCard>
</template>

<style scoped>
.loading {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-12);
  padding: 20px;
}

.bars {
  display: flex;
  gap: 3px;
  align-items: center;
  height: 120px;
}

.bar {
  flex: 1 1 0;
  border-radius: 3px;
}
</style>
