import { ApiError } from '@/shared/api/http'
import type { WrongCode } from '@/shared/api/types'

export type AuthMode = 'login' | 'reg'

// The rule the form tells, the same the server checks: at least 8
// characters, one of them a digit
export function passwordOk(p: string): boolean {
  return [...p].length >= 8 && /\d/.test(p)
}

// The form is ready to send: the button stays off until it is
export function formReady(mode: AuthMode, name: string, email: string, password: string): boolean {
  const filled = email.includes('@') && password.length > 0
  return mode === 'login' ? filled : filled && name.trim() !== '' && passwordOk(password)
}

// What the toast says when the server says no
export function authError(err: unknown): string {
  if (!(err instanceof ApiError)) return 'Что-то пошло не так'
  switch (err.status) {
    case 0:
      return 'Нет связи с сервером'
    case 401:
      return 'Неверная почта или пароль'
    case 409:
      return 'Эта почта уже зарегистрирована'
    case 422:
      return 'Проверь имя, почту и пароль'
    case 410:
      return 'Код истёк, отправь новый'
    case 429:
      return 'Подожди минуту перед новым кодом'
    case 502:
      return 'Не получилось отправить письмо, попробуй ещё раз'
    default:
      return 'Сервер ответил ошибкой, попробуй ещё раз'
  }
}

export interface Bar {
  up: number // px above the line
  down: number // px below it
  home: boolean // the bar above is the home team's color
  away: boolean // the bar below is the away team's
}

// The Flow wave of the left panel, as the Auth board draws it: 40 bars,
// home above the line and away below, home pulling ahead at the end
export function flowBars(n = 40): Bar[] {
  return Array.from({ length: n }, (_, i) => {
    const v = Math.sin(i / 3.2) * 38 + Math.sin(i / 1.3) * 14 + (i > 26 ? 30 : 0)
    return {
      up: Math.max(2, Math.round(Math.max(v, 0) * 1.4)),
      down: Math.max(2, Math.round(Math.max(-v, 0) * 1.4)),
      home: v > 0,
      away: v < 0,
    }
  })
}

// Where to go once signed in: back where the user came from (the
// session-expired dialog sends ?redirect=), else home. Only a path of this
// site: "//other.site" is not one
export function nextPath(redirect: unknown): string {
  return typeof redirect === 'string' && redirect.startsWith('/') && !redirect.startsWith('//')
    ? redirect
    : '/'
}

// The tries left that a wrong code's 422 tells; null for any other error
export function attemptsLeft(err: unknown): number | null {
  if (!(err instanceof ApiError) || err.status !== 422) return null
  const left = (err.problem as WrongCode | null)?.attemptsLeft
  return typeof left === 'number' ? left : null
}

// The line under the cells after a wrong code
export function triesLeft(n: number): string {
  if (n <= 0) return 'Попытки закончились, отправь новый код.'
  if (n === 1) return 'Неверный код. Осталась 1 попытка.'
  return `Неверный код. Осталось ${n} ${n < 5 ? 'попытки' : 'попыток'}.`
}

// "0:42"
export function countdown(seconds: number): string {
  return `${Math.floor(seconds / 60)}:${String(seconds % 60).padStart(2, '0')}`
}
