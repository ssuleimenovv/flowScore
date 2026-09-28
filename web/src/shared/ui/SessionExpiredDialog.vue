<script setup lang="ts">
import { useRoute } from 'vue-router'
import AppButton from './AppButton.vue'
import AppDialog from './AppDialog.vue'
import AppIcon from './AppIcon.vue'

// 401 from the States board. After sign-in the user comes back to this page,
// which is what "открытый матч сохранится" promises.
const open = defineModel<boolean>('open', { required: true })
const route = useRoute()
</script>

<template>
  <AppDialog v-model:open="open" label="Сессия истекла">
    <span class="tile">
      <AppIcon name="lock" :size="22" />
    </span>
    <h2 class="title">Сессия истекла</h2>
    <p class="text">Мы не смогли обновить токен. Войди снова — открытый матч сохранится.</p>
    <div class="actions">
      <AppButton
        class="grow"
        :to="{ path: '/login', query: { redirect: route.fullPath } }"
        @click="open = false"
      >
        Войти снова
      </AppButton>
      <AppButton variant="secondary" @click="open = false">Позже</AppButton>
    </div>
  </AppDialog>
</template>

<style scoped>
.tile {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  border-radius: 14px;
  background: var(--fs-home-soft);
  color: var(--fs-home);
}

.title {
  font-size: var(--fs-text-title-3);
  font-weight: 600;
}

.text {
  color: var(--fs-muted);
  font-size: var(--fs-text-subheadline);
  line-height: 1.5;
}

.actions {
  display: flex;
  gap: 10px;
}

.grow {
  flex-grow: 1;
}
</style>
