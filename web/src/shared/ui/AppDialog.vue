<script setup lang="ts">
import { useTemplateRef, watch } from 'vue'
import AppCard from './AppCard.vue'

// A modal on the native <dialog>. showModal() gives us the focus trap, Esc,
// the inert page behind and the ::backdrop for free.
const open = defineModel<boolean>('open', { required: true })
const { label } = defineProps<{ label: string }>()

const dialog = useTemplateRef('dialog')

watch(
  [open, dialog],
  ([isOpen, el]) => {
    if (!el) return
    if (isOpen && !el.open) el.showModal()
    if (!isOpen && el.open) el.close()
  },
  { immediate: true },
)
</script>

<template>
  <!-- The dialog has no padding, so a click on the dialog itself is a click on the backdrop -->
  <dialog
    ref="dialog"
    class="dialog"
    :aria-label="label"
    @close="open = false"
    @click.self="open = false"
  >
    <AppCard class="panel">
      <slot />
    </AppCard>
  </dialog>
</template>

<style scoped>
.dialog {
  width: min(380px, 100% - 2 * var(--fs-screen-padding));
  max-width: none;
  /* The browser centers a modal with margin: auto; our reset sets margin: 0 on everything */
  margin: auto;
  padding: 0;
  overflow: visible;
  border: 0;
  background: transparent;
  color: var(--fs-text);
}

.dialog[open] {
  animation: fs-in var(--fs-duration-lift) var(--fs-ease) both;
}

.dialog::backdrop {
  background: rgb(0 0 0 / 55%);
}

.dialog[open]::backdrop {
  animation: fade var(--fs-duration-lift) var(--fs-ease) both;
}

@keyframes fade {
  from {
    opacity: 0;
  }
}

.panel {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: 28px;
}
</style>
