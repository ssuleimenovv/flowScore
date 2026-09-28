<script setup lang="ts">
// The open arc from the States board. It keeps spinning under reduced motion
// a frozen spinner would read as "stuck", and a rotation this small is not a trigger
const { size = 18 } = defineProps<{ size?: number }>()
</script>

<template>
  <svg
    class="spinner"
    :width="size"
    :height="size"
    view-box="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="1.8"
    stroke-linecap="round"
    aria-hidden="true"
  >
    <path d="M12 3a9 9 0 1 1-9 9" />
  </svg>
</template>

<style scoped>
.spinner {
  flex-shrink: 0;
  animation: spin var(--fs-duration-spin) linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
