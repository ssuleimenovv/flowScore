<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useToast } from '@/shared/toast/useToast'
import AppButton from '@/shared/ui/AppButton.vue'
import AppIcon from '@/shared/ui/AppIcon.vue'
import { useSessionStore } from '@/stores/session'
import { authError, formReady, type AuthMode } from '../auth'

// The form of the Auth board: "Вход", or "Регистрация" with a name and the
// password rule added. The button stays off until the form can be sent.
const session = useSessionStore()
const { show } = useToast()
const route = useRoute()
const router = useRouter()

const mode = ref<AuthMode>('login')
const name = ref('')
const email = ref('')
const password = ref('')
const reveal = ref(false)
const busy = ref(false)

const tabs: Array<{ value: AuthMode; label: string }> = [
  { value: 'login', label: 'Вход' },
  { value: 'reg', label: 'Регистрация' },
]
const ready = computed(() => formReady(mode.value, name.value, email.value, password.value))

// Back where the user came from (the session-expired dialog sends
// ?redirect=), else home. Only a path of this site: "//other.site" is not one
function next(): string {
  const r = route.query.redirect
  return typeof r === 'string' && r.startsWith('/') && !r.startsWith('//') ? r : '/'
}

async function submit() {
  if (!ready.value || busy.value) return
  busy.value = true
  try {
    if (mode.value === 'login') await session.signIn(email.value, password.value)
    else await session.signUp(name.value, email.value, password.value)
    await router.replace(next())
  } catch (err) {
    show({ icon: 'warning', tone: 'accent', title: authError(err) })
  } finally {
    busy.value = false
  }
}
</script>

<template>
  <form class="form" novalidate @submit.prevent="submit">
    <div class="tabs" role="tablist" aria-label="Вход или регистрация">
      <button
        v-for="t in tabs"
        :key="t.value"
        type="button"
        role="tab"
        class="tab fs-press"
        :aria-selected="mode === t.value"
        @click="mode = t.value"
      >
        {{ t.label }}
      </button>
    </div>

    <div class="fields">
      <label v-if="mode === 'reg'" class="field">
        Имя
        <input
          v-model="name"
          class="input"
          type="text"
          autocomplete="name"
          placeholder="Как к тебе обращаться"
        />
      </label>
      <label class="field">
        Эл. почта
        <input
          v-model="email"
          class="input"
          type="email"
          inputmode="email"
          autocomplete="email"
          placeholder="you@mail.kz"
        />
      </label>
      <label class="field">
        <span class="field-head">
          Пароль
          <RouterLink v-if="mode === 'login'" to="/forgot" class="forgot fs-link">
            Забыли пароль?
          </RouterLink>
        </span>
        <span class="password">
          <input
            v-model="password"
            class="input"
            :type="reveal ? 'text' : 'password'"
            :autocomplete="mode === 'login' ? 'current-password' : 'new-password'"
            placeholder="Минимум 8 символов"
          />
          <button
            type="button"
            class="reveal fs-press"
            :aria-label="reveal ? 'Скрыть пароль' : 'Показать пароль'"
            :aria-pressed="reveal"
            @click="reveal = !reveal"
          >
            <AppIcon name="eye" :size="18" />
          </button>
        </span>
      </label>
      <span v-if="mode === 'reg'" class="rule">Не короче 8 символов и хотя бы одна цифра.</span>
    </div>

    <AppButton class="submit" type="submit" :loading="busy" :disabled="!ready">
      {{ busy ? 'Проверяем…' : mode === 'login' ? 'Войти' : 'Создать аккаунт' }}
    </AppButton>

    <p class="legal">
      Продолжая, ты соглашаешься с условиями сервиса и политикой конфиденциальности.
    </p>
  </form>
</template>

<style scoped>
.form {
  display: flex;
  flex-direction: column;
  gap: 22px;
}

.tabs {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 2px;
  padding: 2px;
  border-radius: 11px;
  background: var(--fs-fill);
}

.tab {
  height: var(--fs-touch);
  border: 0;
  border-radius: 10px;
  background: transparent;
  color: var(--fs-muted);
  font-size: var(--fs-text-subheadline);
  font-weight: 600;
  cursor: pointer;
  transition:
    transform var(--fs-duration-press) ease,
    background-color var(--fs-duration-toggle) var(--fs-ease-switch),
    color var(--fs-duration-toggle) var(--fs-ease-switch);
}

.tab[aria-selected='true'] {
  background: var(--fs-seg-selected);
  color: var(--fs-text);
}

.fields {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: var(--fs-space-8);
  color: var(--fs-muted);
  font-size: var(--fs-text-footnote);
  font-weight: 600;
}

.field-head {
  display: flex;
  justify-content: space-between;
}

.forgot {
  color: var(--fs-home);
  text-decoration: none;
}

.input {
  width: 100%;
  height: 48px;
  padding: 0 14px;
  border: 1px solid var(--fs-line);
  border-radius: var(--fs-radius-md);
  background: var(--fs-bg);
  color: var(--fs-text);
  font: inherit;
  font-size: var(--fs-text-subheadline);
}

.password {
  position: relative;
  display: block;
}

.password .input {
  padding-right: 52px;
}

.reveal {
  position: absolute;
  top: 2px;
  right: 2px;
  display: flex;
  align-items: center;
  justify-content: center;
  width: var(--fs-touch);
  height: var(--fs-touch);
  padding: 0;
  border: 0;
  background: transparent;
  color: var(--fs-muted);
  cursor: pointer;
}

.reveal[aria-pressed='true'] {
  color: var(--fs-text);
}

.rule {
  color: var(--fs-muted);
  font-size: var(--fs-text-caption);
}

/* AppButton's root: full width and 52 px, as the board's main button */
.submit {
  width: 100%;
  height: 52px;
}

.legal {
  color: var(--fs-muted);
  font-size: var(--fs-text-caption);
  line-height: 1.5;
  text-align: center;
}
</style>
