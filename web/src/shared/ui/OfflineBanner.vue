<script setup lang="ts">
import { computed } from 'vue'
import AppIcon from './AppIcon.vue'

// No network, but the screen still has data (States board). The colors are
// inverted so the banner reads as a system message, not part of the content.
const { updatedAt } = defineProps<{ updatedAt: Date | null }>()

const time = computed(() =>
  updatedAt?.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' }),
)
</script>

<template>
  <div class="banner" role="status">
    <AppIcon name="offline" :size="18" />
    <p class="text">
      <b>Нет подключения.</b>
      {{
        time
          ? `Показаны данные на ${time} — обновим, когда сеть вернётся.`
          : 'Обновим, когда сеть вернётся.'
      }}
    </p>
  </div>
</template>

<style scoped>
.banner {
  display: flex;
  gap: var(--fs-space-12);
  align-items: center;
  padding: var(--fs-space-12) var(--fs-space-16);
  border-radius: 14px;
  background: var(--fs-text);
  color: var(--fs-bg);
  font-size: var(--fs-text-subheadline);
}

.banner > svg {
  flex-shrink: 0;
}

.text {
  flex-grow: 1;
}
</style>
