<script setup lang="ts">
import { ref, watch } from 'vue'
import type { Side } from '@/shared/api/types'

// A team's crest, or its three letters when there is none: a demo match's
// teams have no crest, and a crest that fails to load falls back as well.
// "lg" is the desktop scoreboard, "md" the phone one, "sm" a card on the home
// screen, "xs" a row of the league table, "mini" a team's name in a list.
// Without a side the tile is neutral.
// "hero" the team page's card
const {
  code,
  side,
  logo = null,
  size = 'lg',
} = defineProps<{
  code: string
  side?: Side
  logo?: string | null
  size?: 'hero' | 'lg' | 'md' | 'sm' | 'xs' | 'mini'
}>()

const failed = ref(false)
watch(
  () => logo,
  () => (failed.value = false),
)
</script>

<template>
  <span
    class="badge"
    :class="[side ?? 'neutral', size, { crest: logo && !failed }]"
    aria-hidden="true"
  >
    <img
      v-if="logo && !failed"
      class="logo"
      :src="logo"
      alt=""
      loading="lazy"
      decoding="async"
      @error="failed = true"
    />
    <template v-else>{{ code }}</template>
  </span>
</template>

<style scoped>
.badge {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  font-family: var(--fs-font-display);
  font-weight: 800;
  letter-spacing: 0.04em;
}

/* The team page: 64px on MobileTeam, 112px on Team */
.hero {
  width: 64px;
  height: 64px;
  border-radius: 18px;
  font-size: 20px;
}

@media (min-width: 768px) {
  .hero {
    width: 112px;
    height: 112px;
    border-radius: 31px;
    font-size: 36px;
  }
}

.lg {
  width: 72px;
  height: 72px;
  border-radius: var(--fs-radius-lg);
  font-size: 23px;
}

.md {
  width: 52px;
  height: 52px;
  border-radius: 15px;
  font-size: var(--fs-text-headline);
}

.sm {
  width: 32px;
  height: 32px;
  border-radius: 9px;
  font-size: 10px;
}

.xs {
  width: 30px;
  height: 30px;
  border-radius: var(--fs-radius-sm);
  font-size: 11px;
}

.mini {
  width: 18px;
  height: 18px;
  border-radius: 5px;
  font-size: 7px;
}

.home {
  background: var(--fs-home-soft);
  color: var(--fs-home);
}

.away {
  background: var(--fs-away-soft);
  color: var(--fs-away);
}

.neutral {
  background: var(--fs-surface-2);
  color: var(--fs-text);
}

/* A crest stands on its own, as in FotMob: no colored tile behind it */
.badge.crest {
  background: transparent;
}

.logo {
  width: 88%;
  height: 88%;
  object-fit: contain;
}
</style>
