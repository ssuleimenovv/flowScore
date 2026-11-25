<script setup lang="ts">
import { computed } from 'vue'
import { competitionAccent, competitionArt, tintStyle } from '../accents'
import CompetitionArt from './CompetitionArt.vue'

// What a competition's page stands on (CompetitionBackdropsExtended): its
// artwork under the site's header, fading out, and a wash of its accent that
// runs on below and fades into the theme. The cards on top keep the theme's.
const { competition } = defineProps<{ competition: string }>()

const art = computed(() => competitionArt(competition))
</script>

<template>
  <div
    v-if="art"
    class="backdrop fs-tinted"
    :style="tintStyle(competitionAccent(competition))"
    aria-hidden="true"
  >
    <div class="art">
      <CompetitionArt :art="art.key" />
    </div>
  </div>
</template>

<style scoped>
/* Phone, as on the board's phones: the art behind the bar and the heading */
.backdrop {
  position: absolute;
  inset: 0 0 auto;
  z-index: 0;
  height: calc(760px + var(--fs-safe-top));
  background: linear-gradient(
    180deg,
    color-mix(in srgb, var(--fs-tint) 18%, transparent),
    transparent
  );
  pointer-events: none;
  animation: appear var(--fs-duration-enter) var(--fs-ease) both;
}

/* Only a fade: a slide would open a gap at the top of the screen */
@keyframes appear {
  from {
    opacity: 0;
  }
}

:root[data-theme='light'] .backdrop {
  background: linear-gradient(
    180deg,
    color-mix(in srgb, var(--fs-tint) 12%, transparent),
    transparent
  );
}

.art {
  position: absolute;
  inset: 0 0 auto;
  height: calc(330px + var(--fs-safe-top));
  overflow: hidden;
  -webkit-mask-image: linear-gradient(180deg, #000 62%, transparent);
  mask-image: linear-gradient(180deg, #000 62%, transparent);
}

/* Desktop, as on the board: the art under the 72 px header and the heading */
@media (min-width: 768px) {
  .backdrop {
    height: 900px;
  }

  .art {
    height: 430px;
  }
}
</style>
