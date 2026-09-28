<script setup lang="ts">
// Development-only: base components in every state, to compare with the States board.
import { onUnmounted, ref } from 'vue'
import { ApiError } from '@/shared/api/http'
import AppButton from '@/shared/ui/AppButton.vue'
import AppCard from '@/shared/ui/AppCard.vue'
import AppSkeleton from '@/shared/ui/AppSkeleton.vue'
import AppSpinner from '@/shared/ui/AppSpinner.vue'
import ErrorState from '@/shared/ui/ErrorState.vue'
import NotFoundState from '@/shared/ui/NotFoundState.vue'
import OfflineBanner from '@/shared/ui/OfflineBanner.vue'
import ReconnectBanner from '@/shared/ui/ReconnectBanner.vue'
import StateMessage from '@/shared/ui/StateMessage.vue'
import LiveBanner from '@/shared/ui/LiveBanner.vue'
import SessionExpiredDialog from '@/shared/ui/SessionExpiredDialog.vue'
import AppToast from '@/shared/toast/AppToast.vue'
import { useToast, type ToastInput } from '@/shared/toast/useToast'

const retrying = ref(false)

function retry() {
  retrying.value = true
  setTimeout(() => (retrying.value = false), 1600)
}

const serverError = new ApiError(503, {
  type: 'about:blank',
  title: 'Service Unavailable',
  status: 503,
  requestId: '7f3a91',
})

// A fake reconnect cycle: 5 s to wait, then it starts over
const retryAt = ref(Date.now() + 5000)
const seconds = ref(5)
const ticker = setInterval(() => {
  const left = Math.ceil((retryAt.value - Date.now()) / 1000)
  if (left <= 0) retryAt.value = Date.now() + 5000
  seconds.value = Math.max(1, left)
}, 250)
onUnmounted(() => clearInterval(ticker))

// The four toasts from the board; "Показать тост" sends them one by one
const samples: ToastInput[] = [
  {
    icon: 'check',
    tone: 'ok',
    title: 'Сценарий сохранён',
    text: 'Он появится в профиле → Симуляции',
  },
  { icon: 'ball', tone: 'accent', title: 'Гол! Фоден, 58′', text: 'Манчестер Сити 1:1 Арсенал' },
  { icon: 'pulse', tone: 'accent', title: 'Скачок Flow +18', text: 'Сити забирает инициативу' },
  {
    icon: 'warning',
    tone: 'accent',
    title: 'Слишком много запросов',
    text: 'Подожди пару секунд и попробуй снова · 429',
  },
]
const toast = useToast()
let sample = 0
function showToast() {
  toast.show(samples[sample++ % samples.length]!)
}

const sessionExpired = ref(false)

// Skeleton bars of the flow wave, the same heights as on the board
const bars = Array.from({ length: 30 }, (_, i) => 20 + ((i * 37) % 70))
</script>

<template>
  <div class="page">
    <section class="group">
      <h2 class="label">КНОПКИ</h2>
      <div class="row">
        <AppButton>Основная</AppButton>
        <AppButton disabled>Недоступна</AppButton>
        <AppButton loading>Загрузка</AppButton>
        <AppButton variant="secondary">Вторичная</AppButton>
        <AppButton variant="secondary" icon="bell" aria-label="Уведомления" />
        <AppButton variant="secondary" size="small">Сейчас</AppButton>
        <AppButton variant="secondary" to="/">На главную</AppButton>
      </div>
    </section>

    <div class="columns">
      <section class="group">
        <h2 class="label">СКЕЛЕТОН</h2>
        <AppCard
          class="skeleton-card fs-in"
          style="--fs-i: 1"
          aria-busy="true"
          aria-label="Загрузка матча"
        >
          <div class="between">
            <AppSkeleton width="90px" :height="14" />
            <AppSkeleton width="64px" :height="22" round />
          </div>
          <div class="team">
            <AppSkeleton width="32px" :height="32" />
            <AppSkeleton class="grow" />
            <AppSkeleton width="20px" :height="24" />
          </div>
          <div class="team">
            <AppSkeleton width="32px" :height="32" />
            <AppSkeleton width="60%" />
          </div>
          <AppSkeleton :height="52" />
          <AppSkeleton :height="6" round />
        </AppCard>
        <AppCard
          class="skeleton-card fs-in"
          style="--fs-i: 2"
          aria-busy="true"
          aria-label="Загрузка волны"
        >
          <AppSkeleton width="140px" :height="18" />
          <div class="bars">
            <AppSkeleton v-for="(h, i) in bars" :key="i" class="bar" :height="h" />
          </div>
        </AppCard>
        <p class="connecting">
          <AppSpinner class="accent" />
          Подключаемся к live-потоку…
        </p>
      </section>

      <section class="group">
        <h2 class="label">ПУСТО</h2>
        <AppCard class="fs-in" style="--fs-i: 3">
          <StateMessage
            icon="ball"
            title="Сейчас нет live-матчей"
            text="Ближайший стартует в 19:00 — Тобол vs Ордабасы"
          >
            <AppButton variant="secondary" to="/">Смотреть расписание</AppButton>
          </StateMessage>
        </AppCard>
        <AppCard class="fs-in" style="--fs-i: 4">
          <StateMessage
            icon="search"
            title="Ничего не нашли по «аррсенал»"
            text="Проверь написание или поищи по лиге"
          />
        </AppCard>
      </section>

      <section class="group">
        <h2 class="label">ОШИБКИ · 5xx · offline · WS</h2>
        <AppCard class="fs-in" style="--fs-i: 5">
          <ErrorState
            title="Не удалось загрузить матч"
            :error="serverError"
            :retrying
            @retry="retry"
          />
        </AppCard>
        <OfflineBanner :updated-at="new Date(2026, 8, 28, 21, 14)" />
        <ReconnectBanner :seconds :retry-at="retryAt" @reconnect="retryAt = Date.now() + 5000" />
      </section>
    </div>

    <div class="columns">
      <section class="group">
        <h2 class="label">ТОСТЫ</h2>
        <AppToast v-for="t in samples" :key="t.title" v-bind="t" />
        <AppButton variant="secondary" @click="showToast">Показать тост</AppButton>
      </section>

      <section class="group">
        <h2 class="label">LIVE-БАННЕР</h2>
        <LiveBanner
          icon="ball"
          title="ГОЛ! Фоден · 58′"
          text="Манчестер Сити 1:1 Арсенал · Flow +22"
          time="сейчас"
        />
        <h2 class="label">401 · СЕССИЯ ИСТЕКЛА</h2>
        <AppButton variant="secondary" @click="sessionExpired = true">Открыть модалку</AppButton>
        <SessionExpiredDialog v-model:open="sessionExpired" />
      </section>

      <section class="group">
        <h2 class="label">404</h2>
        <AppCard class="fs-in" style="--fs-i: 6">
          <NotFoundState
            title="Такого матча нет"
            text="Возможно, ссылка устарела или матч перенесли."
          />
        </AppCard>
      </section>
    </div>
  </div>
</template>

<style scoped>
.page {
  display: grid;
  gap: var(--fs-space-32);
  padding: var(--fs-space-24) var(--fs-screen-padding);
}

.columns {
  display: grid;
  gap: var(--fs-space-24);
  align-items: start;
}

@media (min-width: 1200px) {
  .columns {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

.group {
  display: grid;
  gap: var(--fs-space-12);
}

.narrow {
  max-width: 440px;
}

.label {
  color: var(--fs-muted);
  font-size: var(--fs-text-caption);
  font-weight: 700;
  letter-spacing: 0.1em;
}

.row {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  align-items: center;
}

.skeleton-card {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: 20px;
}

.between {
  display: flex;
  justify-content: space-between;
}

.team {
  display: flex;
  gap: var(--fs-space-12);
  align-items: center;
}

.grow {
  flex-grow: 1;
}

.bars {
  display: flex;
  gap: 3px;
  align-items: center;
  height: 120px;
}

.bar {
  flex: 1 1 0;
  border-radius: 3px;
}

.connecting {
  display: flex;
  gap: 10px;
  align-items: center;
  color: var(--fs-muted);
  font-size: var(--fs-text-subheadline);
}

.accent {
  color: var(--fs-home);
}
</style>
