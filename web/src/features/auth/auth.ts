import { ApiError } from '@/shared/api/http'

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
    default:
      return 'Сервер ответил ошибкой, попробуй ещё раз'
  }
}
