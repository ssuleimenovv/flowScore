<script setup lang="ts">
import { RESULT_LETTER, RESULT_WORD, type Result } from '../team'

// «В», «Н» or «П» on a square: a win in the team's color, a draw and a loss quiet
const { result } = defineProps<{ result: Result }>()
</script>

<template>
  <span class="chip" :class="result" :aria-label="RESULT_WORD[result]">
    {{ RESULT_LETTER[result] }}
  </span>
</template>

<style scoped>
.chip {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: var(--fs-radius-sm);
  font-size: var(--fs-text-footnote);
  font-weight: 700;
}

.win {
  background: var(--fs-home);
  color: var(--fs-on-home);
}

.draw {
  background: var(--fs-surface-2);
  color: var(--fs-text);
}

.loss {
  background: var(--fs-line);
  color: var(--fs-text);
}
</style>
