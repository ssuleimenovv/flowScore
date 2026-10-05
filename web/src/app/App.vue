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

    <!-- The StatsBomb Open Data licence asks for this credit wherever the data is shown -->
    <p class="credit">
      Данные матчей:
      <a href="https://github.com/statsbomb/open-data" target="_blank" rel="noopener">
        StatsBomb Open Data
      </a>
    </p>


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
