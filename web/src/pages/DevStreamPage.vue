<script setup lang="ts">
// Development page: shows the raw live stream. Not a screen from the mockup.
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useMatchStream } from '@/shared/live/useMatchStream'
import { useNetwork } from '@/shared/state/useNetwork'
import { resolveScreenState } from '@/shared/state/screenState'

const route = useRoute()
const matchId = String(route.params.matchId ?? '3754314')

const { status, secondsToRetry, flow, delta10, wave, events, gaps, reconnectNow } =
  useMatchStream(matchId)
const { online } = useNetwork()

const screen = computed(() =>
  resolveScreenState({
    status: flow.value ? 'success' : 'pending',
    hasData: flow.value !== null,
    isEmpty: false,
    online: online.value,
    socket: status.value,
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
      match {{ matchId }} · socket {{ status }} · view {{ screen.view }} · gaps {{ gaps }}
    </p>

    <section v-if="flow" class="flow">
      <div>
        <span class="num home">{{ Math.round(flow.home) }}</span>
        <span v-if="delta10" class="delta">{{ signed(delta10.home) }} за 10′</span>
      </div>
      <div>
        <span class="num away">{{ Math.round(flow.away) }}</span>
        <span v-if="delta10" class="delta">{{ signed(delta10.away) }} за 10′</span>
      </div>
    </section>

    <div class="wave" aria-label="Волна потока">
      <div v-for="p in wave" :key="p.minute" class="col">
        <span class="up" :style="{ height: `${Math.max(0, p.home - p.away) / 2}%` }" />
        <span class="down" :style="{ height: `${Math.max(0, p.away - p.home) / 2}%` }" />
      </div>
    </div>

    <ol class="events">
      <li v-for="e in events.slice(0, 12)" :key="e.seq">
        {{ e.minute }}{{ e.addedTime ? `+${e.addedTime}` : '' }}′ · {{ e.type }} · {{ e.side }} ·
        {{ e.player?.name }}
      </li>
    </ol>
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

.meta,
.delta {
  color: var(--fs-muted);
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
