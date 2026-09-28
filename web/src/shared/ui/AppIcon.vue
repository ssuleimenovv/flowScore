<script setup lang="ts">
// Stroke icons from the mockup (24×24 grid). In the native shell they
// become SF Symbols (see the Tokens board).
export type IconName =
  | 'home'
  | 'trophy'
  | 'sliders'
  | 'search'
  | 'person'
  | 'bell'
  | 'sun'
  | 'moon'
  | 'back'
  | 'arrow'
  | 'check'
  | 'close'
  | 'retry'
  | 'ball'
  | 'star'
  | 'warning'
  | 'offline'
  | 'lock'
  | 'pulse'

interface Shape {
  paths: string[]
  circles?: Array<[cx: number, cy: number, r: number]>
}

const icons: Record<IconName, Shape> = {
  home: { paths: ['M3 11l9-7 9 7v9a1 1 0 0 1-1 1h-5v-6h-6v6H4a1 1 0 0 1-1-1z'] },
  trophy: {
    paths: [
      'M8 21h8M12 17v4M7 4h10v5a5 5 0 0 1-10 0z',
      'M17 5h3v2a3 3 0 0 1-3 3M7 5H4v2a3 3 0 0 0 3 3',
    ],
  },
  sliders: {
    paths: ['M4 6h10M18 6h2M4 12h4M12 12h8M4 18h12'],
    circles: [
      [16, 6, 2],
      [10, 12, 2],
      [18, 18, 2],
    ],
  },
  search: { paths: ['M20 20l-4-4'], circles: [[11, 11, 7]] },
  person: { paths: ['M4 21c1.5-4 4.5-6 8-6s6.5 2 8 6'], circles: [[12, 8, 4]] },
  bell: { paths: ['M6 9a6 6 0 1 1 12 0c0 6 3 8 3 8H3s3-2 3-8', 'M10.3 21a2 2 0 0 0 3.4 0'] },
  sun: {
    paths: [
      'M12 2v2M12 20v2M4.9 4.9l1.4 1.4M17.7 17.7l1.4 1.4M2 12h2M20 12h2M4.9 19.1l1.4-1.4M17.7 6.3l1.4-1.4',
    ],
    circles: [[12, 12, 4]],
  },
  moon: { paths: ['M20 14.5A8 8 0 0 1 9.5 4a8 8 0 1 0 10.5 10.5z'] },
  back: { paths: ['M15 5l-7 7 7 7'] },
  arrow: { paths: ['M5 12h14M13 616 6-6 6'] },
  check: { paths: ['M5 12.5l4.5 4.5L19 7.5'] },
  close: { paths: ['M6 6l12 12M18 6L6 18'] },
  retry: { paths: ['M4 12a8 8 0 1 0 2.3-5.7M4 4v4h4'] },
  ball: { paths: ['M12 7l4 3-1.5 4.5h-5L8 10z'], circles: [[12, 12, 9]] },
  star: { paths: ['M12 3l2.8 5.8 6.2.8-4.5 4.4 1.1 6.2L12 17.3l-5.6 2.9 1.1-6.2L3 9.6l6.2-.8z'] },
  warning: { paths: ['M12 3l9.5 17h-19z', 'M12 10v4M12 17h.01'] },
  offline: {
    paths: [
      'M2 8.5a15 15 0 0 1 20 0M5.5 12a10 10 0 0 1 13 0M9 15.5a5 5 0 0 1 6 0M12 19h.01',
      'M3 3l18 18',
    ],
  },
  // The mockup draws the lock body as <rect x=4 y=11 w=16 h=10 rx=2>; same shape as a path
  lock: {
    paths: [
      'M6 11h12a2 2 0 0 1 2 2v6a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2v-6a2 2 0 0 1 2-2z',
      'M8 11V7a4 4 0 0 1 8 0v4',
    ],
  },
  pulse: { paths: ['M3 12h4l3-8 4 16 3-8h4'] },
}

const { name, size = 24 } = defineProps<{ name: IconName; size?: number }>()
</script>

<template>
  <svg
    :width="size"
    :height="size"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="1.8"
    stroke-linecap="round"
    stroke-linejoin="round"
    aria-hidden="true"
  >
    <path v-for="d in icons[name].paths" :key="d" :d="d" />
    <circle v-for="[cx, cy, r] in icons[name].circles ?? []" :key="`${cx}-${cy}`" :cx :cy :r />
  </svg>
</template>
