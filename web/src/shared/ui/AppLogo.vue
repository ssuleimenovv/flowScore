<script setup lang="ts">
// The FlowScore mark: a pulse line in a rounded square.
// "outline" is the desktop header, "filled" the mobile nav bar.
const { variant = 'outline', size = 30 } = defineProps<{
  variant?: 'outline' | 'filled'
  size?: number
}>()
</script>

<template>
  <svg :width="size" :height="size" viewBox="0 0 30 30" fill="none" aria-hidden="true">
    <rect
      x="1"
      y="1"
      width="28"
      height="28"
      rx="8"
      :fill="variant === 'filled' ? 'var(--fs-accent)' : 'none'"
      :stroke="variant === 'filled' ? 'none' : 'var(--fs-home)'"
      stroke-width="2"
    />
    <path
      d="M6 16h4l2.5-6.5 4 12 2.5-5.5H24"
      :stroke="variant === 'filled' ? 'var(--fs-on-accent)' : 'var(--fs-home)'"
      :stroke-width="variant === 'filled' ? 2.4 : 2.2"
      stroke-linecap="round"
      stroke-linejoin="round"
    />
  </svg>
</template>
