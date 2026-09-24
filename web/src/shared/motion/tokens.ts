import { gsap } from 'gsap'
import { CustomEase } from 'gsap/CustomEase'

gsap.registerPlugin(CustomEase)

// Keep in sync with --fs-ease / --fs-duration-* in styles/tokens.css

// CustomEase - gsap has no clue css cubic-bezier
export const EASE = CustomEase.create('fs', '0.2, 0.7, 0.2, 1')
export const EASE_SWITCH = CustomEase.create('fs-switch', '0.3, 0.7, 0.3, 1')

export const DURATION = {
  press: 0.14,
  fade: 0.15,
  color: 0.2,
  toggle: 0.22,
  nudge: 0.25,
  lift: 0.3,
  grow: 0.5,
  enter: 0.7,
  bar: 0.8,
  barH: 0.9,
  count: 1.2,
} as const

export const STAGGER = 0.06
