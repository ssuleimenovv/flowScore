<script setup lang="ts">
import type { Team } from '@/shared/api/types'
import TeamBadge from './TeamBadge.vue'

// "Everton — Chelsea" in a list, each name next to its crest, as in FotMob.
// A logo is null for a team without one; its letters stand in for it.
const { home, away, homeLogo, awayLogo } = defineProps<{
  home: Team
  away: Team
  homeLogo: string | null
  awayLogo: string | null
}>()
</script>

<template>
  <span class="pair">
    <TeamBadge :code="home.code" side="home" size="mini" :logo="homeLogo" />
    <span class="name">{{ home.name }}</span>
    <span class="dash" aria-hidden="true">—</span>
    <span class="name">{{ away.name }}</span>
    <TeamBadge :code="away.code" side="away" size="mini" :logo="awayLogo" />
  </span>
</template>

<style scoped>
/* Not called .teams, which the lists use for their own blocks: a parent's
   scoped rule for that class would reach this root as well */
.pair {
  display: flex;
  gap: 6px;
  align-items: center;
  min-width: 0;
}

/* A long pair is cut at the end of each name, so both crests stay */
.name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.dash {
  flex-shrink: 0;
}
</style>
