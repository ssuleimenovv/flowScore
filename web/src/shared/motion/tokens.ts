import { gsap } from 'gsap'
import { CustomEase } from 'gsap/CustomEase'

gsap.registerPlugin(CustomEase)

// Keep in sync with --fs-ease / --fs-duration-* in styles/tokens.css
export const EASE = CustomEase.create('fs', '0.2, 0.7, 0.2, 1')

// CustomEase - gsap has no clue css cubic-bezier
export const DURATION = {
  fast: 0.15,
  base: 0.3,
  slow: 0.7,
} as const
