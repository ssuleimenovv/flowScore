import { describe, expect, it } from 'vitest'
import { ApiError } from '@/shared/api/http'
import { authError, formReady, passwordOk } from '../auth'

describe('passwordOk', () => {
  it('follows the rule the form tells', () => {
    expect(passwordOk('flowscore1')).toBe(true)
    expect(passwordOk('короткий1')).toBe(true) // 9 letters, counted as letters
    expect(passwordOk('short1')).toBe(false)
    expect(passwordOk('no digits here')).toBe(false)
  })
})

describe('formReady', () => {
  it('asks more of a sign-up than of a sign-in', () => {
    expect(formReady('login', '', 'a@b.kz', 'x')).toBe(true)
    expect(formReady('reg', '', 'a@b.kz', 'flowscore1')).toBe(false) // no name
    expect(formReady('reg', 'Аян', 'a@b.kz', 'short')).toBe(false)
    expect(formReady('reg', 'Аян', 'a@b.kz', 'flowscore1')).toBe(true)
  })
})

describe('authError', () => {
  it('tells each answer in words', () => {
    expect(authError(new ApiError(401, null))).toBe('Неверная почта или пароль')
    expect(authError(new ApiError(409, null))).toBe('Эта почта уже зарегистрирована')
    expect(authError(new ApiError(0, null))).toBe('Нет связи с сервером')
    expect(authError(new Error('boom'))).toBe('Что-то пошло не так')
  })
})
