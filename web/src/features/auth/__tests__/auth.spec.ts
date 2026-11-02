import { describe, expect, it } from 'vitest'
import { ApiError } from '@/shared/api/http'
import {
  attemptsLeft,
  authError,
  countdown,
  flowBars,
  formReady,
  nextPath,
  passwordOk,
  triesLeft,
} from '../auth'

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

describe('flowBars', () => {
  it('draws the wave of the Auth board', () => {
    const bars = flowBars()
    expect(bars).toHaveLength(40)
    // Every bar shows, and a bar is never both teams'
    expect(bars.every((b) => b.up >= 2 && b.down >= 2 && !(b.home && b.away))).toBe(true)
    // Home pulls ahead after bar 26, as the board draws it; the wave dips
    // back to away only at the very end
    expect(bars.slice(27, 37).every((b) => b.home)).toBe(true)
    expect(bars[37]!.away).toBe(true)
  })
})

describe('the code screen', () => {
  it('reads the tries left from a wrong code', () => {
    const wrong = new ApiError(422, {
      type: 'about:blank',
      title: 'Wrong code',
      status: 422,
      attemptsLeft: 3,
    } as never)
    expect(attemptsLeft(wrong)).toBe(3)
    expect(attemptsLeft(new ApiError(410, null))).toBeNull()
    expect(triesLeft(3)).toBe('Неверный код. Осталось 3 попытки.')
    expect(triesLeft(1)).toBe('Неверный код. Осталась 1 попытка.')
    expect(triesLeft(0)).toBe('Попытки закончились, отправь новый код.')
  })

  it('counts down and goes back only within the site', () => {
    expect(countdown(42)).toBe('0:42')
    expect(countdown(60)).toBe('1:00')
    expect(nextPath('/matches/1')).toBe('/matches/1')
    expect(nextPath('//evil.example')).toBe('/')
    expect(nextPath(undefined)).toBe('/')
  })
})
