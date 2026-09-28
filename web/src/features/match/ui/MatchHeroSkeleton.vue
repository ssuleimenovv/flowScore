<script setup lang="ts">
import AppCard from '@/shared/ui/AppCard.vue'
import AppSkeleton from '@/shared/ui/AppSkeleton.vue'

// The scoreboard while it loads: the same blocks in the same places, so nothing
// jumps when the data arrives. "compact" follows MatchHeroCompact.
const { compact = false } = defineProps<{ compact?: boolean }>()
</script>

<template>
  <section v-if="compact" class="compact" aria-busy="true" aria-label="Загрузка матча">
    <AppSkeleton width="72px" :height="24" round />
    <div class="scoreboard">
      <div class="team">
        <AppSkeleton width="52px" :height="52" />
        <AppSkeleton width="80px" :height="14" />
      </div>
      <AppSkeleton width="110px" :height="54" />
      <div class="team">
        <AppSkeleton width="52px" :height="52" />
        <AppSkeleton width="80px" :height="14" />
      </div>
    </div>
    <AppSkeleton :height="8" round />
  </section>

  <AppCard v-else class="card" aria-busy="true" aria-label="Загрузка матча">
    <div class="row">
      <AppSkeleton width="72px" :height="24" round />
      <AppSkeleton width="40%" :height="14" />
    </div>
    <div class="scoreboard wide">
      <div class="side">
        <AppSkeleton width="72px" :height="72" />
        <div class="lines">
          <AppSkeleton width="70%" :height="20" />
          <AppSkeleton width="50%" :height="40" />
        </div>
      </div>
      <AppSkeleton width="160px" :height="76" />
      <div class="side reverse">
        <AppSkeleton width="72px" :height="72" />
        <div class="lines end">
          <AppSkeleton width="70%" :height="20" />
          <AppSkeleton width="50%" :height="40" />
        </div>
      </div>
    </div>
    <AppSkeleton :height="10" round />
  </AppCard>
</template>

<style scoped>
.compact {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-16);
  padding: 28px 0 var(--fs-space-16);
}

.card {
  display: flex;
  flex-direction: column;
  gap: 22px;
  padding: var(--fs-space-24) 28px;
}

.row {
  display: flex;
  gap: var(--fs-space-12);
  align-items: center;
}

.scoreboard {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto minmax(0, 1fr);
  gap: var(--fs-space-8);
  align-items: center;
}

.scoreboard.wide {
  gap: var(--fs-space-24);
}

.team {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-8);
  align-items: center;
}

.side {
  display: flex;
  gap: 18px;
  align-items: center;
}

.reverse {
  flex-direction: row-reverse;
}

.lines {
  display: flex;
  flex-grow: 1;
  flex-direction: column;
  gap: 10px;
}

.lines.end {
  align-items: flex-end;
}
</style>
