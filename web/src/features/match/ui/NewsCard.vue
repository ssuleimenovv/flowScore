<script setup lang="ts">
import { computed, onScopeDispose, ref } from 'vue'
import type { MatchPreMatch } from '@/shared/api/types'
import AppCard from '@/shared/ui/AppCard.vue'
import { newsRows } from '../prematch'

// "Новости" from the MatchExtended boards: the league's latest from ESPN,
// each opening the story there.
const { pre } = defineProps<{ pre: MatchPreMatch }>()

// "12 мин назад" keeps up with the clock
const now = ref(new Date())
const timer = setInterval(() => (now.value = new Date()), 60_000)
onScopeDispose(() => clearInterval(timer))

const rows = computed(() => newsRows(pre, now.value))
</script>

<template>
  <AppCard as="section" class="news fs-in" style="--fs-i: 6" aria-labelledby="news-title">
    <h2 id="news-title" class="title">Новости</h2>
    <ul class="rows">
      <li v-for="n in rows" :key="n.url">
        <a class="story fs-press" :href="n.url" target="_blank" rel="noopener">
          <span class="headline">{{ n.headline }}</span>
          <span class="meta">{{ n.meta }}</span>
        </a>
      </li>
    </ul>
  </AppCard>
</template>

<style scoped>
.news {
  display: flex;
  flex-direction: column;
  padding: 20px var(--fs-space-24);
}

.title {
  padding-bottom: 10px;
  font-size: 20px;
  font-weight: 700;
  line-height: 1.2;
  letter-spacing: -0.01em;
}

.rows {
  padding: 0;
  list-style: none;
}

.story {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-4);
  padding: var(--fs-space-12) 0;
  border-top: 1px solid var(--fs-line);
  color: var(--fs-text);
  text-decoration: none;
}

.headline {
  font-size: var(--fs-text-subheadline);
  font-weight: 600;
  line-height: 1.3;
}

.story:hover .headline {
  text-decoration: underline;
}

.meta {
  color: var(--fs-muted);
  font-size: var(--fs-text-caption);
}

@media (max-width: 767px) {
  .news {
    padding: var(--fs-space-16);
  }
}
</style>
