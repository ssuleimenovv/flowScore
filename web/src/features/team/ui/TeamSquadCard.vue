<script setup lang="ts">
import type { SquadPlayer } from '@/shared/api/types'
import type { RequestStatus } from '@/shared/state/screenState'
import AppSkeleton from '@/shared/ui/AppSkeleton.vue'
import ErrorState from '@/shared/ui/ErrorState.vue'
import { POSITION_WORD } from '../team'
import PlayerPhoto from './PlayerPhoto.vue'

// "Состав" from the Team board. ESPN gives no minutes, load or injuries, so
// the board's "За 7 дней", "Нагрузка" and "Статус" wait for their data
// (D6, D4). On the phone the position stands where the load bar will be.
const { players, request, error } = defineProps<{
  players: SquadPlayer[]
  request: RequestStatus
  error: unknown
}>()
defineEmits<{ retry: [] }>()
</script>

<template>
  <section class="card fs-in">
    <ErrorState
      v-if="request === 'error'"
      title="Не удалось загрузить состав"
      :error
      @retry="$emit('retry')"
    />
    <div
      v-else-if="request !== 'success'"
      class="loading"
      aria-busy="true"
      aria-label="Загрузка состава"
    >
      <AppSkeleton v-for="i in 7" :key="i" :height="48" />
    </div>
    <p v-else-if="players.length === 0" class="empty">ESPN пока не дал состав этой команды.</p>
    <template v-else>
      <div class="row header" aria-hidden="true">
        <span>№</span>
        <span>Игрок</span>
        <span>Позиция</span>
      </div>
      <ul class="list">
        <li v-for="p in players" :key="p.id" class="row">
          <span class="number">{{ p.number || '—' }}</span>
          <span class="who">
            <PlayerPhoto :number="p.number" />
            <span class="text">
              <span class="name">{{ p.name }}</span>
              <span class="line">{{ POSITION_WORD[p.position] }}</span>
            </span>
          </span>
          <span class="position">{{ POSITION_WORD[p.position] }}</span>
        </li>
      </ul>
    </template>
  </section>
</template>

<style scoped>
/* Phone, as on MobileTeam: photo with the number, the name, the position */
.card {
  display: flex;
  flex-direction: column;
  padding: 6px var(--fs-space-16);
  border: 1px solid var(--fs-line);
  border-radius: var(--fs-radius-lg);
  background: var(--fs-surface);
}

.list {
  margin: 0;
  padding: 0;
  list-style: none;
}

.row {
  display: grid;
  grid-template-columns: minmax(0, 1fr);
  gap: var(--fs-space-12);
  align-items: center;
  padding: var(--fs-space-12) 0;
  border-top: 1px solid var(--fs-line);
}

.header,
.number,
.position {
  display: none;
}

.who {
  display: flex;
  gap: var(--fs-space-12);
  align-items: center;
  min-width: 0;
}

.text {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.name {
  overflow: hidden;
  font-size: var(--fs-text-body);
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.line {
  color: var(--fs-muted);
  font-size: var(--fs-text-footnote);
}

.loading {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-12);
  padding: var(--fs-space-12) 0;
}

.empty {
  padding: var(--fs-space-12) 0;
  color: var(--fs-muted);
  font-size: var(--fs-text-footnote);
}

/* Desktop, as on Team: a table of №, Игрок, Позиция */
@media (min-width: 768px) {
  .card {
    padding: 20px 28px;
  }

  .row {
    grid-template-columns: 48px minmax(0, 1fr) 240px;
    gap: var(--fs-space-16);
    padding: 14px 0;
  }

  .header {
    display: grid;
    padding: 0 0 var(--fs-space-12);
    border-top: 0;
    color: var(--fs-muted);
    font-size: var(--fs-text-caption);
    font-weight: 600;
  }

  .number {
    display: block;
    color: var(--fs-faint);
    font-family: var(--fs-font-display);
    font-size: 24px;
    font-weight: 800;
  }

  .position {
    display: block;
    color: var(--fs-muted);
    font-size: var(--fs-text-subheadline);
  }

  .who {
    gap: 14px;
  }

  .name {
    font-size: var(--fs-text-subheadline);
  }

  .line {
    display: none;
  }
}
</style>
