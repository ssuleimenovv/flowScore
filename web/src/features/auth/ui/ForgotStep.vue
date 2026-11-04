<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useToast } from '@/shared/toast/useToast'
import AppButton from '@/shared/ui/AppButton.vue'
import AppIcon from '@/shared/ui/AppIcon.vue'
import { useSessionStore } from '@/stores/session'
import { authError } from '../auth'

// "Восстановление пароля" from the AuthCode board: the email a code goes to.
// Every email gets the same answer, so the screen never tells whether one is
// signed up
const emit = defineEmits<{ sent: [email: string] }>()

const session = useSessionStore()
const { show } = useToast()
const route = useRoute()
const router = useRouter()

const email = ref('')
const busy = ref(false)
const ready = computed(() => email.value.includes('@'))

async function submit() {
  if (!ready.value || busy.value) return
  busy.value = true
  try {
    emit('sent', await session.forgot(email.value))
  } catch (err) {
    show({ icon: 'warning', tone: 'accent', title: authError(err) })
  } finally {
    busy.value = false
  }
}

// "Ко входу" keeps where Вход was to lead on to
function back() {
  void router.push({ path: '/login', query: route.query })
}
</script>

<template>
  <form class="step" novalidate @submit.prevent="submit">
    <button type="button" class="back fs-link" @click="back">
      <AppIcon name="back" :size="18" />
      Ко входу
    </button>

    <div class="intro">
      <span class="tile"><AppIcon name="lock" :size="24" /></span>
      <h2 class="title">Восстановление пароля</h2>
      <p class="text">Пришлём код на почту, с ним задашь новый пароль.</p>
    </div>

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

    <AppButton class="submit" type="submit" :loading="busy" :disabled="!ready">
      {{ busy ? 'Отправляем…' : 'Отправить код' }}
    </AppButton>
    <p class="note">Если эта почта зарегистрирована, код придёт в течение минуты.</p>
  </form>
</template>

<style scoped src="./step.css"></style>
