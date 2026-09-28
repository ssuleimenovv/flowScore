<script setup lang="ts">
import AppCard from '@/shared/ui/AppCard.vue'
import AppSkeleton from '@/shared/ui/AppSkeleton.vue'
</script>

<template>
  <AppCard class="skeleton" aria-busy="true" aria-label="Загрузка хроники">
    <AppSkeleton width="110px" :height="20" />
    <div v-for="i in 6" :key="i" class="row">
      <AppSkeleton width="28px" :height="14" />
      <AppSkeleton width="42px" :height="24" />
      <div class="lines">
        <AppSkeleton width="80%" :height="14" />
        <AppSkeleton width="50%" :height="10" />
      </div>
    </div>
  </AppCard>
</template>

<style scoped>
.skeleton {
  display: flex;
  flex-direction: column;
  gap: 18px;
  padding: 22px 24px;
}

.row {
  display: flex;
  gap: 10px;
  align-items: center;
}

.lines {
  display: flex;
  flex-grow: 1;
  flex-direction: column;
  gap: 6px;
}
</style>
