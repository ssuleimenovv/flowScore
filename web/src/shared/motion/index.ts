import type { Plugin } from 'vue'
import { gsap } from 'gsap'
import { ScrollTrigger } from 'gsap/ScrollTrigger'
import 'lenis/dist/lenis.css'

import { DURATION, EASE, EASE_SWITCH, STAGGER } from './tokens'
import { MQ } from './media'
import { createSmoothScroll } from './smoothScroll'

export { DURATION, EASE, EASE_SWITCH, STAGGER, MQ }

export const motion: Plugin = {
  install() {
    gsap.registerPlugin(ScrollTrigger)
    gsap.defaults({ ease: EASE, duration: DURATION.enter })

    gsap.matchMedia().add(`${MQ.desktop} and ${MQ.motionOk}`, () => {
      return createSmoothScroll()
    })
  },
}

// gsap.matchMedia() - central figure here
// condition - true - function return & lenis on
// condition - false - gsap returns cleanup & lenis off
