<script setup lang="ts">
import { computed } from 'vue'
import { RouterView, useRoute } from 'vue-router'
import DesktopHeader from './layout/DesktopHeader.vue'
import MobileNavBar from './layout/MobileNavBar.vue'
import TabBar from './layout/TabBar.vue'
import { useBackdrop, useScrolledPast } from './layout/navigation'
import CompetitionBackdrop from '@/features/league/ui/CompetitionBackdrop.vue'
import ToastHost from '@/shared/toast/ToastHost.vue'
import { useSessionStore } from '@/stores/session'

const route = useRoute()
const active = computed(() => route.meta.tab)
const inner = computed(() => route.meta.inner === true)
const bare = computed(() => route.meta.bare === true)

// A competition's page stands on its artwork: the header lies clear on it
// until the page is scrolled off the art, then it is the usual bar again
const { art } = useBackdrop()
const scrolled = useScrolledPast(240)
const onArt = computed(() => art.value !== null && !scrolled.value)

// A kept token may have expired since the last visit: the server says whose it is
void useSessionStore().check()
</script>

<template>
  <!-- Вход is a screen of its own, as on the Auth boards -->
  <RouterView v-if="bare" />
  <template v-else>
    <CompetitionBackdrop v-if="art" :competition="art" />

    <!-- Both layouts are in the DOM; CSS shows one of them at the 768 px breakpoint -->
    <div class="desktop-only">
      <DesktopHeader :active :on-art />
    </div>
    <div class="mobile-only">
      <MobileNavBar :inner :on-art />
    </div>

    <main class="content">
      <!-- The key remounts the page when only the id changes (/matches/1 → /matches/2) -->
      <RouterView :key="route.path" />

      <!-- The real matches come from ESPN; the StatsBomb Open Data licence asks
           for this credit wherever its data, the demo matches, is shown -->
      <p class="credit">
        Данные матчей: ESPN ·
        <a href="https://github.com/statsbomb/open-data" target="_blank" rel="noopener">
          StatsBomb Open Data
        </a>
      </p>
    </main>

    <div class="mobile-only">
      <TabBar :active />
    </div>
  </template>
  <ToastHost />
</template>

<style scoped>
.content {
  /* Above a competition's artwork, which lies under the whole top of the screen */
  position: relative;
  z-index: 1;
  /* Room for the fixed tab bar, so the last row is not hidden under it */
  padding-bottom: calc(var(--fs-tab-height) + var(--fs-safe-bottom));
}

.credit {
  padding: var(--fs-space-16) var(--fs-screen-padding) var(--fs-space-24);
  color: var(--fs-faint);
  font-size: var(--fs-text-caption);
  text-align: center;
}

.credit a {
  color: inherit;
}

.desktop-only {
  display: none;
}

/* display: contents keeps the sticky header a direct child of the page flow */
.mobile-only {
  display: contents;
}

@media (min-width: 768px) {
  .desktop-only {
    display: contents;
  }

  .mobile-only {
    display: none;
  }

  .content {
    padding-bottom: 0;
  }
}
</style>
