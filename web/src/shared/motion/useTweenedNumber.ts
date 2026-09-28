import { onScopeDispose, ref, watch, type Ref } from 'vue'
import { gsap } from 'gsap'
import { DURATION, EASE } from './tokens'
import { MQ } from './media'

// A number that counts to its new value instead of jumping, like the Flow
// figures on the mockup. Starts from 0, so the first value counts up too.
export function useTweenedNumber(source: Ref<number>) {
  const shown = ref(0)
  const state = { value: 0 }

  watch(
    source,
    (to) => {
      if (matchMedia(MQ.reduce).matches) {
        gsap.killTweensOf(state)
        state.value = to
        shown.value = Math.round(to)
        return
      }
      gsap.to(state, {
        value: to,
        duration: DURATION.count,
        ease: EASE,
        overwrite: true,
        onUpdate: () => (shown.value = Math.round(state.value)),
      })
    },
    { immediate: true },
  )

  onScopeDispose(() => gsap.killTweensOf(state))

  return shown
}
