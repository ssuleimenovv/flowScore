<script setup lang="ts">
import AppToast from './AppToast.vue'
import { useToast } from './useToast'

const { toasts, dismiss } = useToast()
</script>

<template>
  <!-- The live region stays in the DOM, so a screen reader hears every new toast -->
  <TransitionGroup tag="div" name="toast" class="host" role="status">
    <AppToast
      v-for="t in toasts"
      :key="t.id"
      :icon="t.icon"
      :tone="t.tone"
      :title="t.title"
      :text="t.text"
      class="item"
      @close="dismiss(t.id)"
    />
  </TransitionGroup>
</template>

<style scoped>
/* Phone: above the tab bar, full width. Desktop: bottom right corner */
.host {
  position: fixed;
  right: var(--fs-screen-padding);
  bottom: calc(var(--fs-tab-height) + var(--fs-safe-bottom) + var(--fs-space-12));
  left: var(--fs-screen-padding);
  z-index: 20;
  display: grid;
  gap: 10px;
  pointer-events: none;
}

@media (min-width: 768px) {
  .host {
    bottom: var(--fs-space-24);
    left: auto;
    width: 380px;
  }
}

/* The host spans the width, but only the toasts catch clicks */
.item {
  pointer-events: auto;
  box-shadow: 0 8px 24px rgb(0 0 0 / 18%);
}

.toast-enter-active,
.toast-leave-active,
.toast-move {
  transition:
    opacity var(--fs-duration-lift) var(--fs-ease),
    transform var(--fs-duration-lift) var(--fs-ease);
}

.toast-enter-from {
  opacity: 0;
  transform: translateY(16px);
}

.toast-leave-to {
  opacity: 0;
  transform: scale(0.96);
}

/* A leaving toast steps out of the flow, so the others slide into its place */
.toast-leave-active {
  position: absolute;
  width: 100%;
}
</style>
