<script setup lang="ts">
import PlayerPhoto from '@/features/team/ui/PlayerPhoto.vue'
import TeamLink from '@/features/team/ui/TeamLink.vue'
import type { PlayerInfo } from '@/shared/api/types'
import { aboutLine } from '../player'

// The player's card from the Player board: photo with the number, the line,
// the name. On the phone (MobilePlayer) it lies right on the page, and the
// club goes up into the nav bar.
const { player } = defineProps<{ player: PlayerInfo }>()
</script>

<template>
  <section class="card fs-in">
    <PlayerPhoto :number="player.number" size="hero" />
    <div class="head">
      <h1 class="name">{{ player.name }}</h1>
      <span class="about">
        <template v-for="(part, i) in aboutLine(player)" :key="part">
          <span v-if="i > 0"> · </span>
          <span>{{ part }}</span>
          <!-- The desktop names the club after the position -->
          <span v-if="i === 0 && player.team" class="club">
            · <TeamLink :team-id="player.team.id">{{ player.team.name }}</TeamLink>
          </span>
        </template>
      </span>
    </div>
  </section>
</template>

<style scoped>
/* Phone, as on MobilePlayer: the photo beside the name */
.card {
  display: flex;
  gap: var(--fs-space-16);
  align-items: flex-end;
}

.head {
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-width: 0;
}

.name {
  font-size: 34px;
  font-weight: 700;
  letter-spacing: -0.025em;
  line-height: 1.05;
  overflow-wrap: anywhere;
}

.about {
  color: var(--fs-muted);
  font-size: var(--fs-text-footnote);
}

.club {
  display: none;
}

/* Desktop, as on Player: one card, the line above the name */
@media (min-width: 768px) {
  .card {
    gap: var(--fs-space-32);
    align-items: center;
    padding: var(--fs-space-32);
    border: 1px solid var(--fs-line);
    border-radius: var(--fs-radius-lg);
    background: var(--fs-surface);
  }

  .head {
    flex-direction: column-reverse;
    gap: var(--fs-space-12);
  }

  .name {
    font-size: 56px;
  }

  .about {
    font-family: var(--fs-font-mono);
    font-variant-numeric: tabular-nums;
  }

  .club {
    display: inline;
  }
}
</style>
