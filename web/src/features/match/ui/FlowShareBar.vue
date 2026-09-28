<script setup lang="ts">
// Two bars that split the width by each side's share of the flow.
const { home, height = 10 } = defineProps<{ home: number; height?: number }>()
</script>

<template>
  <div
    class="bar"
    role="img"
    :aria-label="`Доля потока: хозяева ${home}%, гости ${100 - home}%`"
    :style="{ height: `${height}px` }"
  >
    <div class="part home" :style="{ width: `${home}%` }" />
    <div class="part away" />
  </div>
</template>

<style scoped>
.bar {
  display: flex;
  gap: 3px;
}

.part {
  border-radius: var(--fs-radius-full);
}

.home {
  background: var(--fs-home);
  transition: width var(--fs-duration-grow) var(--fs-ease);
}

.away {
  flex-grow: 1;
  background: var(--fs-away);
}
</style>
