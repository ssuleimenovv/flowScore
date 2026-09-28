<script setup lang="ts">
// Development-only: base components in every state, to compare with the States board.
import { ref } from 'vue'
import AppButton from '@/shared/ui/AppButton.vue'

const retrying = ref(false)

function retry() {
  retrying.value = true
  setTimeout(() => (retrying.value = false), 1600)
}
</script>

<template>
  <div class="page">
    <section class="group">
      <h2 class="label">КНОПКИ</h2>
      <div class="row">
        <AppButton>Основная</AppButton>
        <AppButton disabled>Недоступна</AppButton>
        <AppButton loading>Загрузка</AppButton>
        <AppButton variant="secondary">Вторичная</AppButton>
        <AppButton variant="secondary" icon="bell" aria-label="Уведомления" />
        <AppButton variant="secondary" size="small">Сейчас</AppButton>
        <AppButton variant="secondary" to="/">На главную</AppButton>
      </div>
    </section>

    <section class="group">
      <h2 class="label">ПОВТОР</h2>
      <div class="row">
        <AppButton icon="retry" :loading="retrying" @click="retry">
          {{ retrying ? 'Загружаем…' : 'Повторить' }}
        </AppButton>
      </div>
    </section>
  </div>
</template>

<style scoped>
.page {
  display: grid;
  gap: var(--fs-space-32);
  padding: var(--fs-space-24) var(--fs-screen-padding);
}

.group {
  display: grid;
  gap: var(--fs-space-12);
}

.label {
  color: var(--fs-muted);
  font-size: var(--fs-text-caption);
  font-weight: 700;
  letter-spacing: 0.1em;
}

.row {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  align-items: center;
}
</style>
