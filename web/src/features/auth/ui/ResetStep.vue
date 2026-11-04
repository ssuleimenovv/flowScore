<script setup lang="ts">
import { computed, onScopeDispose, ref, useTemplateRef, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ApiError } from '@/shared/api/http'
import { useToast } from '@/shared/toast/useToast'
import AppButton from '@/shared/ui/AppButton.vue'
import AppIcon from '@/shared/ui/AppIcon.vue'
import { useSessionStore } from '@/stores/session'
import { attemptsLeft, authError, countdown, nextPath, resetReady, triesLeft } from '../auth'
import CodeInput from './CodeInput.vue'

// "Новый пароль" from the AuthCode board: the code from the email and the
// new password. They set it and sign in
const { email } = defineProps<{ email: string }>()
const emit = defineEmits<{ back: [] }>()

const session = useSessionStore()
const { show } = useToast()
const route = useRoute()
const router = useRouter()

const code = ref('')
const password = ref('')
const reveal = ref(false)
const busy = ref(false)
const error = ref<string | null>(null)
const ready = computed(() => resetReady(code.value, password.value))
// A new digit is a new try: the old error goes
watch(code, () => (error.value = null))

// A code went out just before this screen; the next one after a minute
const RESEND_AFTER = 60
const wait = ref(RESEND_AFTER)
const timer = setInterval(() => {
  if (wait.value > 0) wait.value--
}, 1000)
onScopeDispose(() => clearInterval(timer))

// The sixth digit sends the form once the password is there too; until
// then it moves on to the password
const passwordField = useTemplateRef<HTMLInputElement>('passwordField')
function codeDone() {
  if (ready.value) void submit()
  else passwordField.value?.focus()
}

async function submit() {
  if (!ready.value || busy.value) return
  busy.value = true
  try {
    await session.reset(email, code.value, password.value)
    await router.replace(nextPath(route.query.redirect))
  } catch (err) {
    const left = attemptsLeft(err)
    if (left !== null) {
      error.value = triesLeft(left)
      if (left === 0) wait.value = 0
    } else {
      // 410: no code to check any more, a new one can go out at once
      if (err instanceof ApiError && err.status === 410) wait.value = 0
      show({ icon: 'warning', tone: 'accent', title: authError(err) })
    }
  } finally {
    busy.value = false
  }
}

async function resend() {
  try {
    await session.forgot(email)
    wait.value = RESEND_AFTER
    code.value = ''
    show({ icon: 'check', tone: 'ok', title: 'Новый код отправлен' })
  } catch (err) {
    show({ icon: 'warning', tone: 'accent', title: authError(err) })
  }
}
</script>

<template>
  <form class="step" novalidate @submit.prevent="submit">
    <button type="button" class="back fs-link" @click="emit('back')">
      <AppIcon name="back" :size="18" />
      Изменить почту
    </button>

    <div class="intro">
      <span class="tile"><AppIcon name="lock" :size="24" /></span>
      <h2 class="title">Новый пароль</h2>
      <p class="text">
        Код отправлен на <span class="email">{{ email }}</span
        >.
      </p>
    </div>

    <div class="entry">
      <span class="field">Код из письма</span>
      <CodeInput v-model="code" :error="error !== null" @complete="codeDone" />
      <span v-if="error" class="error" role="alert">{{ error }}</span>
    </div>

    <label class="field">
      Новый пароль
      <span class="password">
        <input
          ref="passwordField"
          v-model="password"
          class="input"
          :type="reveal ? 'text' : 'password'"
          autocomplete="new-password"
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
    <span class="rule">Не короче 8 символов и хотя бы одна цифра.</span>

    <AppButton class="submit" type="submit" :loading="busy" :disabled="!ready">
      {{ busy ? 'Сохраняем…' : 'Сохранить и войти' }}
    </AppButton>

    <p class="resend">
      Не пришло письмо?
      <template v-if="wait > 0">
        Отправить ещё раз через <span class="timer">{{ countdown(wait) }}</span>
      </template>
      <button v-else type="button" class="again fs-link" @click="resend">Отправить ещё раз</button>
    </p>
  </form>
</template>

<style scoped src="./step.css"></style>
