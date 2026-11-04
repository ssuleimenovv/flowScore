<script setup lang="ts">
import { computed, onScopeDispose, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ApiError } from '@/shared/api/http'
import { useToast } from '@/shared/toast/useToast'
import AppButton from '@/shared/ui/AppButton.vue'
import AppIcon from '@/shared/ui/AppIcon.vue'
import { useSessionStore } from '@/stores/session'
import { attemptsLeft, authError, countdown, nextPath, triesLeft } from '../auth'
import CodeInput from './CodeInput.vue'

// "Проверь почту" from the AuthCode board: the code from the email confirms
// it and signs in
const { email } = defineProps<{ email: string }>()
const emit = defineEmits<{ back: [] }>()

const session = useSessionStore()
const { show } = useToast()
const route = useRoute()
const router = useRouter()

const code = ref('')
const busy = ref(false)
const error = ref<string | null>(null)
const ready = computed(() => /^\d{6}$/.test(code.value))
// A new digit is a new try: the old error goes
watch(code, () => (error.value = null))

// A code went out just before this screen; the next one after a minute
const RESEND_AFTER = 60
const wait = ref(RESEND_AFTER)
const timer = setInterval(() => {
  if (wait.value > 0) wait.value--
}, 1000)
onScopeDispose(() => clearInterval(timer))

async function submit() {
  if (!ready.value || busy.value) return
  busy.value = true
  try {
    await session.verify(email, code.value)
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
    await session.resend(email)
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
      <span class="tile"><AppIcon name="mail" :size="24" /></span>
      <h2 class="title">Проверь почту</h2>
      <p class="text">
        Мы отправили 6-значный код на <span class="email">{{ email }}</span
        >. Он действует 10&nbsp;минут.
      </p>
    </div>

    <div class="entry">
      <CodeInput v-model="code" :error="error !== null" @complete="submit" />
      <span v-if="error" class="error" role="alert">{{ error }}</span>
    </div>

    <AppButton class="submit" type="submit" :loading="busy" :disabled="!ready">
      {{ busy ? 'Проверяем…' : 'Подтвердить' }}
    </AppButton>

    <p class="resend">
      Не пришло письмо?
      <template v-if="wait > 0">
        Отправить ещё раз через <span class="timer">{{ countdown(wait) }}</span>
      </template>
      <button v-else type="button" class="again fs-link" @click="resend">Отправить ещё раз</button>
    </p>
    <p class="note">Проверь папку «Спам»: письмо приходит от FlowScore.</p>
  </form>
</template>
<style scoped src="./step.css"></style>
