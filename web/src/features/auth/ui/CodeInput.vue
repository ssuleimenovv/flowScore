<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

// Six cells for the code from the email, as on the AuthCode board. Under
// them is one input, so pasting, the phone's "code from Mail" suggestion and
// a password manager all fill the whole code at once
const model = defineModel<string>({ required: true })
const { error = false } = defineProps<{ error?: boolean }>()
const emit = defineEmits<{ complete: [] }>()

const input = ref<HTMLInputElement | null>(null)
const focused = ref(false)
const cells = computed(() => Array.from({ length: 6 }, (_, i) => model.value[i] ?? ''))
// The cell the next digit goes to; none when all six are in
const active = computed(() => (focused.value && model.value.length < 6 ? model.value.length : -1))

function onInput(e: Event) {
  const el = e.target as HTMLInputElement
  // Digits only, and six at most: "482 193" pasted is "482193"
  const digits = el.value.replace(/\D/g, '').slice(0, 6)
  el.value = digits
  model.value = digits
  if (digits.length === 6) emit('complete')
}

onMounted(() => input.value?.focus())
</script>

<template>
  <label class="code" :class="{ error }">
    <span class="cells" aria-hidden="true">
      <span
        v-for="(d, i) in cells"
        :key="i"
        class="cell"
        :class="{ filled: d !== '', active: i === active }"
      >
        {{ d }}<span v-if="i === active" class="caret" />
      </span>
    </span>
    <input
      ref="input"
      class="input"
      :value="model"
      type="text"
      inputmode="numeric"
      autocomplete="one-time-code"
      aria-label="Код из письма"
      @input="onInput"
      @focus="focused = true"
      @blur="focused = false"
    />
  </label>
</template>

<style scoped>
.code {
  position: relative;
  display: block;
}

/* Phone, as on MobileAuthCode */
.cells {
  display: grid;
  grid-template-columns: repeat(6, minmax(0, 1fr));
  gap: var(--fs-space-8);
}

.cell {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 56px;
  border: 1px solid var(--fs-line);
  border-radius: var(--fs-radius-md);
  background: var(--fs-bg);
  color: var(--fs-text);
  font-family: var(--fs-font-mono);
  font-size: 24px;
  font-weight: 600;
  transition: border-color var(--fs-duration-color);
}

.cell.filled {
  border-color: var(--fs-muted);
}

.cell.active {
  border: 2px solid var(--fs-home);
}

.error .cell {
  border: 1.5px solid var(--fs-home);
}

.caret {
  width: 2px;
  height: 1em;
  border-radius: 1px;
  background: var(--fs-home);
  animation: blink var(--fs-duration-pulse) steps(1) infinite;
}

/* The input sits over the cells, unseen: a tap on any cell focuses it.
   16 px keeps iOS from zooming in on focus */
.input {
  position: absolute;
  inset: 0;
  width: 100%;
  border: 0;
  background: transparent;
  color: transparent;
  font-size: 16px;
  caret-color: transparent;
  opacity: 0;
}

/* Desktop, as on AuthCode */
@media (min-width: 768px) {
  .cells {
    gap: 10px;
  }

  .cell {
    height: 64px;
    font-size: 28px;
  }
}

@keyframes blink {
  50% {
    opacity: 0;
  }
}
</style>
