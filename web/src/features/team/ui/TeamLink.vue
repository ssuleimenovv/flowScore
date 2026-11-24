<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { teamPath } from '../team'

// A team's name, or its crest and name, that lead to its page. A demo
// match's teams are StatsBomb's and have no page: without an id it stays text.
const { teamId = null } = defineProps<{ teamId?: string | null }>()
</script>

<template>
  <RouterLink v-if="teamId" :to="teamPath(teamId)" class="link fs-link"><slot /></RouterLink>
  <span v-else><slot /></span>
</template>

<style scoped>
.link {
  color: inherit;
  text-decoration: none;
}

/* A mouse finds the link by the underline; a finger needs no hint */
@media (hover: hover) {
  .link:hover {
    text-decoration: underline;
    text-underline-offset: 3px;
  }
}
</style>
