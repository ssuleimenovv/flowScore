<script setup lang="ts">
// Development page: shows the raw match data. Not a screen from the mockup.
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useMatchLive } from '@/features/match/useMatchLive'
import { useNetwork } from '@/shared/state/useNetwork'
import { resolveScreenState } from '@/shared/state/screenState'

const route = useRoute()
const matchId = String(route.params.matchId ?? '3754314')

const { state, request, error, socketStatus, secondsToRetry, reconnectNow, retry } =
  useMatchLive(matchId)
const { online } = useNetwork()

const screen = computed(() =>
  resolveScreenState({
    status: request.value,
    hasData: state.value !== null,
    isEmpty: false,
    online: online.value,
    socket: socketStatus.value,
  }),
)

function signed(n: number): string {
  return `${n >= 0 ? '+' : '−'}${Math.abs(Math.round(n))}`
}
</script>

<template>
  <main class="dev">
    <p v-if="screen.banner === 'reconnecting'" class="banner">
      Live-поток прервался.
      <span v-if="secondsToRetry !== null">Переподключение через {{ secondsToRetry }} с…</span>
      <button type="button" @click="reconnectNow">Сейчас</button>
    </p>

    <p class="meta">
      socket {{ socketStatus }} · request {{ request }} · view {{ screen.view }}
      <span v-if="state">· seq {{ Object.values(state.seq).join('/') }}</span>
    </p>

    <p v-if="screen.view === 'loading'">Загрузка…</p>
    <p v-else-if="screen.view === 'offline'">Нет подключения.</p>
    <p v-else-if="screen.view === 'error'">
      Не удалось загрузить матч: {{ error }}
      <button type="button" @click="retry">Повторить</button>
    </p>

    <template v-if="state">
      <header class="score">
        <span>{{ state.match.home.code }}</span>
        <span class="num">{{ state.score.home }} : {{ state.score.away }}</span>
        <span>{{ state.match.away.code }}</span>
        <span class="meta">
          {{ state.match.competition.name }} · {{ state.match.competition.round }}-й тур ·
          {{ state.match.venue?.name }}
        </span>
      </header>

      <section class="flow">
        <div>
          <span class="num home">{{ Math.round(state.flow.home) }}</span>
          <span class="meta">{{ signed(state.delta10.home) }} за 10′</span>
        </div>
        <div>
          <span class="num away">{{ Math.round(state.flow.away) }}</span>
          <span class="meta">{{ signed(state.delta10.away) }} за 10′</span>
        </div>
      </section>

      <div class="wave" aria-label="Волна потока">
        <div v-for="p in state.points" :key="p.minute" class="col">
          <span class="up" :style="{ height: `${Math.max(0, p.home - p.away) / 2}%` }" />
          <span class="down" :style="{ height: `${Math.max(0, p.away - p.home) / 2}%` }" />
        </div>
      </div>

      <ol class="events">
        <li v-for="e in state.events.slice(0, 12)" :key="e.id">
          {{ e.minute }}{{ e.addedTime ? `+${e.addedTime}` : '' }}′ · {{ e.type }} · {{ e.side }} ·
          {{ e.player?.name }}
        </li>
      </ol>
    </template>
  </main>
</template>

<style scoped>
.dev {
  display: grid;
  gap: var(--fs-space-16);
  padding: var(--fs-space-24);
  font-family: var(--fs-font-mono);
}

.banner {
  display: flex;
  gap: var(--fs-space-12);
  align-items: center;
  padding: var(--fs-space-12) var(--fs-space-16);
  border-radius: var(--fs-radius-md);
  background: var(--fs-live-soft);
}

.meta {
  color: var(--fs-muted);
}

.score {
  display: flex;
  gap: var(--fs-space-16);
  align-items: baseline;
}

.flow {
  display: flex;
  gap: var(--fs-space-40);
}

.num {
  font-family: var(--fs-font-display);
  font-size: 4rem;
  font-weight: 800;
  font-variant-numeric: tabular-nums;
}

.home {
  color: var(--fs-home);
}

.away {
  color: var(--fs-away);
}

.wave {
  display: flex;
  gap: 2px;
  height: 200px;
}

.col {
  display: flex;
  flex: 1;
  flex-direction: column;
  justify-content: center;
}

.up,
.down {
  transition: height var(--fs-duration-grow) var(--fs-ease);
}

.up {
  background: var(--fs-home);
  border-radius: 3px 3px 0 0;
}

.down {
  background: var(--fs-away);
  border-radius: 0 0 3px 3px;
}
</style>
