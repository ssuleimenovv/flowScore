<script setup lang="ts">
import AppCard from '@/shared/ui/AppCard.vue'
import NotFoundState from '@/shared/ui/NotFoundState.vue'
</script>

<template>
  <div class="page">
    <AppCard class="card fs-in">
      <NotFoundState
        title="Такой страницы нет"
        text="Возможно, ссылка устарела или её перенесли."
      />
    </AppCard>
  </div>
</template>

<style scoped>
.page {
  padding: var(--fs-space-40) var(--fs-screen-padding);
}

.card {
  max-width: 440px;
  margin: 0 auto;
}
</style>
