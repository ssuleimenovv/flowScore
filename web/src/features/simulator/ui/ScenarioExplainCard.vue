<script setup lang="ts">
import AppCard from '@/shared/ui/AppCard.vue'
import AppIcon from '@/shared/ui/AppIcon.vue'
import type { ScenarioStory } from '../scenario'

// The Explainability card of the Simulator boards: what the scenario changes
// and by how much, in the model's numbers.
const { story } = defineProps<{ story: ScenarioStory }>()
</script>

<template>
  <AppCard as="section" class="explain fs-in" style="--fs-i: 3" aria-labelledby="why-title">
    <span class="chip">
      <AppIcon name="sparkle" :size="14" class="spark" />
      Explainability Agent
    </span>
    <h3 id="why-title" class="title">{{ story.headline }}</h3>
    <p v-if="!story.reasons.length" class="hint">
      Назначь удаление — здесь появится объяснение, как меняется исход.
    </p>
    <ul v-else class="reasons">
      <li v-for="reason in story.reasons" :key="reason.title" class="reason">
        <div class="text">
          <span class="what">{{ reason.title }}</span>
          <span class="detail">{{ reason.detail }}</span>
        </div>
        <span class="value" :class="reason.side">{{ reason.value }}</span>
      </li>
    </ul>
  </AppCard>
</template>

<style scoped>
.explain {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: var(--fs-space-24) 28px;
}

.chip {
  display: inline-flex;
  gap: 6px;
  align-items: center;
  align-self: flex-start;
  height: 26px;
  padding: 0 10px;
  border-radius: var(--fs-radius-sm);
  background: var(--fs-surface-2);
  font-size: var(--fs-text-caption);
  font-weight: 600;
  white-space: nowrap;
}

.spark {
  color: var(--fs-home);
}

.title {
  margin-bottom: 4px;
  font-size: var(--fs-text-title-3);
  font-weight: 600;
}

.hint {
  color: var(--fs-muted);
  font-size: var(--fs-text-subheadline);
  line-height: 1.5;
}

.reasons {
  padding: 0;
  list-style: none;
}

.reason {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: var(--fs-space-12);
  align-items: center;
  padding: var(--fs-space-12) 0;
  border-top: 1px solid var(--fs-line);
}

.text {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.what {
  font-size: var(--fs-text-subheadline);
  font-weight: 600;
}

.detail {
  color: var(--fs-muted);
  font-size: var(--fs-text-footnote);
}

.value {
  font-family: var(--fs-font-mono);
  font-size: var(--fs-text-footnote);
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  white-space: nowrap;
}

.value.home {
  color: var(--fs-home);
}

.value.away {
  color: var(--fs-away);
}

@media (max-width: 767px) {
  .explain {
    padding: 18px var(--fs-space-16);
  }
}
</style>
