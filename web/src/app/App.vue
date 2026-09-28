<script setup lang="ts">
import { computed } from 'vue'
import { RouterView, useRoute } from 'vue-router'
import DesktopHeader from './layout/DesktopHeader.vue'
import MobileNavBar from './layout/MobileNavBar.vue'
import TabBar from './layout/TabBar.vue'
import ToastHost from '@/shared/toast/ToastHost.vue'

const route = useRoute()
const active = computed(() => route.meta.tab)
const inner = computed(() => route.meta.inner === true)
</script>

<template>
  <!-- Both layouts are in the DOM; CSS shows one of them at the 768 px breakpoint -->
  <div class="desktop-only">
    <DesktopHeader :active />
  </div>
  <div class="mobile-only">
    <MobileNavBar :inner />
  </div>

  <main class="content">
    <!-- The key remounts the page when only the id changes (/matches/1 → /matches/2) -->
    <RouterView :key="route.path" />
  </main>

  <div class="mobile-only">
    <TabBar :active />
  </div>
  <ToastHost />
</template>

<style scoped>
.content {
  /* Room for the fixed tab bar, so the last row is not hidden under it */
  padding-bottom: calc(var(--fs-tab-height) + var(--fs-safe-bottom));
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
