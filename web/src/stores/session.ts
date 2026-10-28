import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { ApiError, apiGet, apiPost } from '@/shared/api/http'
import type { Session, User } from '@/shared/api/types'

const KEY = 'fs-session'

// Who is signed in, kept between visits. A token the server turns away (it
// counts for 15 minutes until refresh tokens come) signs the user out.
export const useSessionStore = defineStore('session', () => {
  const saved = read()
  const token = ref<string | null>(saved?.token ?? null)
  const user = ref<User | null>(saved?.user ?? null)
  const signedIn = computed(() => token.value !== null)

  function keep(s: Session) {
    token.value = s.token
    user.value = s.user
    write(s)
  }

  async function signIn(email: string, password: string) {
    keep(await apiPost<Session>('/auth/login', { email, password }))
  }

  async function signUp(name: string, email: string, password: string) {
    keep(await apiPost<Session>('/auth/register', { name, email, password }))
  }

  function signOut() {
    token.value = null
    user.value = null
    write(null)
  }

  // Asks the server whose the kept token is: on a new visit it may have
  // expired. Only a 401 signs out; offline, the user stays as they were
  async function check() {
    if (!token.value) return
    try {
      keep({ token: token.value, user: await apiGet<User>('/me', undefined, token.value) })
    } catch (err) {
      if (err instanceof ApiError && err.status === 401) signOut()
    }
  }

  return { token, user, signedIn, signIn, signUp, signOut, check }
})

// localStorage can throw (private mode, blocked storage): then the session
// lasts this visit only
function read(): Session | null {
  try {
    const raw = localStorage.getItem(KEY)
    return raw ? (JSON.parse(raw) as Session) : null
  } catch {
    return null
  }
}

function write(s: Session | null) {
  try {
    if (s) localStorage.setItem(KEY, JSON.stringify(s))
    else localStorage.removeItem(KEY)
  } catch {
    // kept in memory only
  }
}
