<script setup lang="ts">
// The artwork a competition's page stands on, from the CompetitionBackdropsExtended
// board: abstract and our own, no league's logo or imagery. Generated from the
// board's drawing code (scratchpad boards/gen_art.py), so the page matches it.
// The left side is held dark: the white text above it keeps its contrast.
import type { ArtKey } from '../accents'

const { art } = defineProps<{ art: ArtKey }>()
</script>

<template>
  <svg
    v-if="art === 'uefa.champions'"
    viewBox="0 0 1440 430"
    preserveAspectRatio="xMaxYMid slice"
    aria-hidden="true"
  >
    <defs>
      <linearGradient id="bart0" x1="0" y1="0" x2="1" y2="0.4">
        <stop offset="0" stop-color="#040824" />
        <stop offset="0.55" stop-color="#0a1a6e" />
        <stop offset="1" stop-color="#1638d6" />
      </linearGradient>
      <radialGradient id="gart0" cx="0.82" cy="0.55" r="0.55">
        <stop offset="0" stop-color="#3f7bff" stop-opacity="0.9" />
        <stop offset="1" stop-color="#3f7bff" stop-opacity="0" />
      </radialGradient>
      <linearGradient id="iart0" x1="0" y1="0" x2="1" y2="0">
        <stop offset="0" stop-color="#5ff2ff" />
        <stop offset="0.35" stop-color="#7a6bff" />
        <stop offset="0.65" stop-color="#ff5fd2" />
        <stop offset="1" stop-color="#7dffb0" />
      </linearGradient>
      <linearGradient id="sart0" x1="0" y1="0" x2="1" y2="0">
        <stop offset="0" stop-color="#000" stop-opacity="0.45" />
        <stop offset="0.55" stop-color="#000" stop-opacity="0" />
      </linearGradient>
    </defs>
    <rect width="1440" height="430" fill="url(#bart0)" />
    <rect width="1440" height="430" fill="url(#gart0)" />
    <circle cx="1403" cy="136" r="0.8" fill="#fff" opacity="0.43" />
    <circle cx="359" cy="94" r="1.6" fill="#fff" opacity="0.35" />
    <circle cx="297" cy="137" r="0.8" fill="#fff" opacity="0.57" />
    <circle cx="1234" cy="78" r="1.2" fill="#fff" opacity="0.51" />
    <circle cx="1038" cy="276" r="0.8" fill="#fff" opacity="0.80" />
    <circle cx="651" cy="190" r="1.2" fill="#fff" opacity="0.58" />
    <circle cx="487" cy="245" r="1.2" fill="#fff" opacity="0.57" />
    <circle cx="180" cy="8" r="1.6" fill="#fff" opacity="0.69" />
    <circle cx="71" cy="64" r="1.2" fill="#fff" opacity="0.78" />
    <circle cx="1192" cy="164" r="1.6" fill="#fff" opacity="0.54" />
    <circle cx="1248" cy="14" r="1.6" fill="#fff" opacity="0.72" />
    <circle cx="1246" cy="42" r="0.8" fill="#fff" opacity="0.74" />
    <circle cx="1416" cy="186" r="1" fill="#fff" opacity="0.65" />
    <circle cx="66" cy="17" r="1.2" fill="#fff" opacity="0.54" />
    <circle cx="101" cy="37" r="1" fill="#fff" opacity="0.47" />
    <circle cx="1317" cy="193" r="0.8" fill="#fff" opacity="0.45" />
    <circle cx="586" cy="290" r="1.6" fill="#fff" opacity="0.53" />
    <circle cx="93" cy="264" r="0.8" fill="#fff" opacity="0.67" />
    <circle cx="648" cy="183" r="0.8" fill="#fff" opacity="0.67" />
    <circle cx="215" cy="273" r="1.2" fill="#fff" opacity="0.75" />
    <circle cx="734" cy="105" r="0.8" fill="#fff" opacity="0.52" />
    <circle cx="262" cy="271" r="1" fill="#fff" opacity="0.43" />
    <circle cx="1243" cy="203" r="1.6" fill="#fff" opacity="0.80" />
    <circle cx="1388" cy="128" r="1" fill="#fff" opacity="0.66" />
    <circle cx="625" cy="283" r="1" fill="#fff" opacity="0.28" />
    <circle cx="624" cy="72" r="1.2" fill="#fff" opacity="0.73" />
    <circle cx="440" cy="33" r="1" fill="#fff" opacity="0.67" />
    <circle cx="376" cy="269" r="0.8" fill="#fff" opacity="0.54" />
    <circle cx="1024" cy="177" r="1" fill="#fff" opacity="0.29" />
    <circle cx="275" cy="211" r="1.6" fill="#fff" opacity="0.77" />
    <circle cx="1433" cy="278" r="1.2" fill="#fff" opacity="0.59" />
    <circle cx="1276" cy="115" r="1" fill="#fff" opacity="0.28" />
    <circle cx="1263" cy="148" r="1.2" fill="#fff" opacity="0.61" />
    <circle cx="329" cy="42" r="1.6" fill="#fff" opacity="0.57" />
    <circle cx="437" cy="199" r="1" fill="#fff" opacity="0.34" />
    <circle cx="56" cy="184" r="1" fill="#fff" opacity="0.42" />
    <circle cx="644" cy="167" r="0.8" fill="#fff" opacity="0.48" />
    <circle cx="1224" cy="161" r="1" fill="#fff" opacity="0.72" />
    <circle cx="377" cy="91" r="1.6" fill="#fff" opacity="0.78" />
    <circle cx="958" cy="292" r="1.6" fill="#fff" opacity="0.46" />
    <circle cx="1198" cy="288" r="0.8" fill="#fff" opacity="0.76" />
    <circle cx="1111" cy="105" r="1.6" fill="#fff" opacity="0.69" />
    <circle cx="1410" cy="28" r="1" fill="#fff" opacity="0.58" />
    <circle cx="1011" cy="152" r="1.2" fill="#fff" opacity="0.49" />
    <circle cx="320" cy="77" r="1.6" fill="#fff" opacity="0.79" />
    <circle cx="525" cy="135" r="0.8" fill="#fff" opacity="0.46" />
    <circle cx="1117" cy="61" r="1" fill="#fff" opacity="0.77" />
    <circle cx="12" cy="27" r="1.6" fill="#fff" opacity="0.65" />
    <circle cx="75" cy="134" r="1.2" fill="#fff" opacity="0.68" />
    <circle cx="235" cy="34" r="1" fill="#fff" opacity="0.46" />
    <circle cx="637" cy="92" r="0.8" fill="#fff" opacity="0.65" />
    <circle cx="1102" cy="108" r="1.6" fill="#fff" opacity="0.78" />
    <circle cx="1313" cy="15" r="1.6" fill="#fff" opacity="0.45" />
    <circle cx="230" cy="5" r="1" fill="#fff" opacity="0.78" />
    <circle cx="694" cy="144" r="0.8" fill="#fff" opacity="0.33" />
    <circle cx="863" cy="212" r="1.6" fill="#fff" opacity="0.49" />
    <circle cx="610" cy="179" r="0.8" fill="#fff" opacity="0.72" />
    <circle cx="492" cy="29" r="1.6" fill="#fff" opacity="0.48" />
    <circle cx="92" cy="70" r="1" fill="#fff" opacity="0.79" />
    <circle cx="817" cy="230" r="1.6" fill="#fff" opacity="0.70" />
    <circle cx="813" cy="208" r="1.2" fill="#fff" opacity="0.63" />
    <circle cx="1205" cy="127" r="1.2" fill="#fff" opacity="0.43" />
    <circle cx="321" cy="150" r="1" fill="#fff" opacity="0.35" />
    <circle cx="1259" cy="97" r="1.6" fill="#fff" opacity="0.39" />
    <circle cx="942" cy="99" r="0.8" fill="#fff" opacity="0.41" />
    <circle cx="64" cy="103" r="0.8" fill="#fff" opacity="0.47" />
    <circle cx="3" cy="258" r="1.2" fill="#fff" opacity="0.41" />
    <circle cx="540" cy="159" r="1.6" fill="#fff" opacity="0.29" />
    <circle cx="1174" cy="111" r="1" fill="#fff" opacity="0.74" />
    <circle cx="95" cy="157" r="1" fill="#fff" opacity="0.51" />
    <path
      d="M547.2 451.5 Q 1036.8 21.5 1555.2 266.6"
      fill="none"
      stroke="url(#iart0)"
      stroke-width="5"
      stroke-linecap="round"
      opacity="0.95"
    />
    <path
      d="M587.2 451.5 Q 1036.8 51.6 1555.2 245.09"
      fill="none"
      stroke="url(#iart0)"
      stroke-width="4"
      stroke-linecap="round"
      opacity="0.77"
    />
    <path
      d="M627.2 451.5 Q 1036.8 81.7 1555.2 223.6"
      fill="none"
      stroke="url(#iart0)"
      stroke-width="3"
      stroke-linecap="round"
      opacity="0.59"
    />
    <path
      d="M667.2 451.5 Q 1036.8 111.8 1555.2 202.1"
      fill="none"
      stroke="url(#iart0)"
      stroke-width="2"
      stroke-linecap="round"
      opacity="0.40"
    />
    <ellipse cx="1180.8" cy="438.6" rx="489.6" ry="94.6" fill="#2a55ff" opacity="0.35" />
    <ellipse
      cx="1180.8"
      cy="438.6"
      rx="489.6"
      ry="94.6"
      fill="none"
      stroke="#9fd0ff"
      stroke-width="1.5"
      opacity="0.6"
    />
    <rect width="1440" height="430" fill="url(#sart0)" />
  </svg>
  <svg
    v-else-if="art === 'esp.1'"
    viewBox="0 0 1440 430"
    preserveAspectRatio="xMaxYMid slice"
    aria-hidden="true"
  >
    <defs>
      <linearGradient id="bart1" x1="0" y1="1" x2="1" y2="0">
        <stop offset="0" stop-color="#8f0f0c" />
        <stop offset="0.45" stop-color="#e2380f" />
        <stop offset="1" stop-color="#ff9a1f" />
      </linearGradient>
      <linearGradient id="sart1" x1="0" y1="0" x2="1" y2="0">
        <stop offset="0" stop-color="#000" stop-opacity="0.45" />
        <stop offset="0.55" stop-color="#000" stop-opacity="0" />
      </linearGradient>
    </defs>
    <rect width="1440" height="430" fill="url(#bart1)" />
    <line x1="576.0" y1="430" x2="736.0" y2="0" stroke="#fff" stroke-width="1" opacity="0.12" />
    <line x1="602.0" y1="430" x2="762.0" y2="0" stroke="#fff" stroke-width="1" opacity="0.12" />
    <line x1="628.0" y1="430" x2="788.0" y2="0" stroke="#fff" stroke-width="1" opacity="0.12" />
    <line x1="654.0" y1="430" x2="814.0" y2="0" stroke="#fff" stroke-width="1" opacity="0.12" />
    <line x1="680.0" y1="430" x2="840.0" y2="0" stroke="#fff" stroke-width="1" opacity="0.12" />
    <line x1="706.0" y1="430" x2="866.0" y2="0" stroke="#fff" stroke-width="1" opacity="0.12" />
    <line x1="732.0" y1="430" x2="892.0" y2="0" stroke="#fff" stroke-width="1" opacity="0.12" />
    <line x1="758.0" y1="430" x2="918.0" y2="0" stroke="#fff" stroke-width="1" opacity="0.12" />
    <line x1="784.0" y1="430" x2="944.0" y2="0" stroke="#fff" stroke-width="1" opacity="0.12" />
    <line x1="810.0" y1="430" x2="970.0" y2="0" stroke="#fff" stroke-width="1" opacity="0.12" />
    <line x1="836.0" y1="430" x2="996.0" y2="0" stroke="#fff" stroke-width="1" opacity="0.12" />
    <line x1="862.0" y1="430" x2="1022.0" y2="0" stroke="#fff" stroke-width="1" opacity="0.12" />
    <line x1="888.0" y1="430" x2="1048.0" y2="0" stroke="#fff" stroke-width="1" opacity="0.12" />
    <line x1="914.0" y1="430" x2="1074.0" y2="0" stroke="#fff" stroke-width="1" opacity="0.12" />
    <line x1="940.0" y1="430" x2="1100.0" y2="0" stroke="#fff" stroke-width="1" opacity="0.12" />
    <line x1="966.0" y1="430" x2="1126.0" y2="0" stroke="#fff" stroke-width="1" opacity="0.12" />
    <line x1="992.0" y1="430" x2="1152.0" y2="0" stroke="#fff" stroke-width="1" opacity="0.12" />
    <line x1="1018.0" y1="430" x2="1178.0" y2="0" stroke="#fff" stroke-width="1" opacity="0.12" />
    <line x1="1044.0" y1="430" x2="1204.0" y2="0" stroke="#fff" stroke-width="1" opacity="0.12" />
    <line x1="1070.0" y1="430" x2="1230.0" y2="0" stroke="#fff" stroke-width="1" opacity="0.12" />
    <polygon points="691.19,430 911.19,0 991.19,0 771.19,430" fill="#fff" opacity="0.1" />
    <polygon points="892.8,430 1112.8,0 1192.8,0 972.8,430" fill="#fff" opacity="0.16" />
    <polygon points="1065.6,430 1285.6,0 1365.6,0 1145.6,430" fill="#fff" opacity="0.08" />
    <polygon points="1238.4,430 1458.4,0 1538.4,0 1318.4,430" fill="#fff" opacity="0.13" />
    <circle cx="1296.0" cy="86.0" r="215.0" fill="#ffd27a" opacity="0.18" />
    <rect width="1440" height="430" fill="url(#sart1)" />
  </svg>
  <svg
    v-else-if="art === 'eng.1'"
    viewBox="0 0 1440 430"
    preserveAspectRatio="xMaxYMid slice"
    aria-hidden="true"
  >
    <defs>
      <radialGradient id="bart2" cx="0.78" cy="0.35" r="0.85">
        <stop offset="0" stop-color="#7a2fd0" />
        <stop offset="0.5" stop-color="#3d195b" />
        <stop offset="1" stop-color="#14051f" />
      </radialGradient>
      <pattern id="dart2" width="14" height="14" patternUnits="userSpaceOnUse">
        <circle cx="7" cy="7" r="2" fill="#f4b8ff" />
      </pattern>
      <linearGradient id="mart2" x1="0" y1="0" x2="1" y2="0">
        <stop offset="0.35" stop-color="#fff" stop-opacity="0" />
        <stop offset="1" stop-color="#fff" stop-opacity="1" />
      </linearGradient>
      <mask id="kart2"><rect width="1440" height="430" fill="url(#mart2)" /></mask>
      <linearGradient id="sart2" x1="0" y1="0" x2="1" y2="0">
        <stop offset="0" stop-color="#000" stop-opacity="0.45" />
        <stop offset="0.55" stop-color="#000" stop-opacity="0" />
      </linearGradient>
    </defs>
    <rect width="1440" height="430" fill="url(#bart2)" />
    <polygon
      points="648.0,430 1080.0,0 1440,0 1440,430"
      fill="url(#dart2)"
      opacity="0.22"
      mask="url(#kart2)"
    />
    <path
      d="M720.0 430 C 1007.99 258.0, 1152.0 129.0, 1512.0 43.0"
      fill="none"
      stroke="#f4b8ff"
      stroke-width="2"
      opacity="0.5"
    />
    <path
      d="M720.0 430 C 1007.99 258.0, 1152.0 129.0, 1512.0 94.6"
      fill="none"
      stroke="#f4b8ff"
      stroke-width="2"
      opacity="0.38"
    />
    <path
      d="M720.0 430 C 1007.99 258.0, 1152.0 129.0, 1512.0 146.2"
      fill="none"
      stroke="#f4b8ff"
      stroke-width="2"
      opacity="0.26"
    />
    <path
      d="M720.0 430 C 1007.99 258.0, 1152.0 129.0, 1512.0 197.79"
      fill="none"
      stroke="#f4b8ff"
      stroke-width="2"
      opacity="0.14"
    />
    <rect width="1440" height="430" fill="url(#sart2)" />
  </svg>
  <svg
    v-else-if="art === 'ger.1'"
    viewBox="0 0 1440 430"
    preserveAspectRatio="xMaxYMid slice"
    aria-hidden="true"
  >
    <defs>
      <linearGradient id="bart3" x1="0" y1="0" x2="1" y2="1">
        <stop offset="0" stop-color="#4a0309" />
        <stop offset="0.6" stop-color="#b8050f" />
        <stop offset="1" stop-color="#ff2b37" />
      </linearGradient>
      <linearGradient id="sart3" x1="0" y1="0" x2="1" y2="0">
        <stop offset="0" stop-color="#000" stop-opacity="0.45" />
        <stop offset="0.55" stop-color="#000" stop-opacity="0" />
      </linearGradient>
    </defs>
    <rect width="1440" height="430" fill="url(#bart3)" />
    <polygon points="1007.99,0 1440,0 1440,301.0" fill="#000" opacity="0.18" />
    <polyline
      points="648.0,86.0 864.0,21.5 1080.0,150.5 1296.0,51.6 1512.0,172.0"
      fill="none"
      stroke="#fff"
      stroke-width="6"
      stroke-linejoin="miter"
      opacity="0.22"
    />
    <polyline
      points="648.0,163.4 864.0,98.89 1080.0,227.9 1296.0,129.0 1512.0,249.40"
      fill="none"
      stroke="#fff"
      stroke-width="5"
      stroke-linejoin="miter"
      opacity="0.18"
    />
    <polyline
      points="648.0,240.8 864.0,176.29 1080.0,305.3 1296.0,206.4 1512.0,326.8"
      fill="none"
      stroke="#fff"
      stroke-width="4"
      stroke-linejoin="miter"
      opacity="0.14"
    />
    <polyline
      points="648.0,318.2 864.0,253.70 1080.0,382.7 1296.0,283.8 1512.0,404.20"
      fill="none"
      stroke="#fff"
      stroke-width="3"
      stroke-linejoin="miter"
      opacity="0.1"
    />
    <rect width="1440" height="430" fill="url(#sart3)" />
  </svg>
  <svg
    v-else-if="art === 'ita.1'"
    viewBox="0 0 1440 430"
    preserveAspectRatio="xMaxYMid slice"
    aria-hidden="true"
  >
    <defs>
      <linearGradient id="bart4" x1="0" y1="0" x2="1" y2="0">
        <stop offset="0" stop-color="#061650" />
        <stop offset="0.6" stop-color="#0f37b8" />
        <stop offset="1" stop-color="#2f74ff" />
      </linearGradient>
      <linearGradient id="sart4" x1="0" y1="0" x2="1" y2="0">
        <stop offset="0" stop-color="#000" stop-opacity="0.45" />
        <stop offset="0.55" stop-color="#000" stop-opacity="0" />
      </linearGradient>
    </defs>
    <rect width="1440" height="430" fill="url(#bart4)" />
    <path
      d="M432.0 236.50 C 720.0 86.0, 1007.99 387.0, 1512.0 150.5"
      fill="none"
      stroke="#cfe0ff"
      stroke-width="2"
      opacity="0.45"
    />
    <path
      d="M432.0 279.5 C 720.0 129.00, 1007.99 408.50, 1512.0 193.49"
      fill="none"
      stroke="#cfe0ff"
      stroke-width="3"
      opacity="0.38"
    />
    <path
      d="M432.0 322.5 C 720.0 172.0, 1007.99 430.0, 1512.0 236.50"
      fill="none"
      stroke="#cfe0ff"
      stroke-width="2"
      opacity="0.31"
    />
    <path
      d="M432.0 365.50 C 720.0 215.0, 1007.99 451.5, 1512.0 279.5"
      fill="none"
      stroke="#cfe0ff"
      stroke-width="3"
      opacity="0.24"
    />
    <path
      d="M432.0 408.50 C 720.0 258.00, 1007.99 473.00, 1512.0 322.5"
      fill="none"
      stroke="#cfe0ff"
      stroke-width="2"
      opacity="0.16"
    />
    <path
      d="M432.0 451.5 C 720.0 301.0, 1007.99 494.49, 1512.0 365.5"
      fill="none"
      stroke="#cfe0ff"
      stroke-width="3"
      opacity="0.09"
    />
    <rect width="1440" height="430" fill="url(#sart4)" />
  </svg>
  <svg
    v-else-if="art === 'fra.1'"
    viewBox="0 0 1440 430"
    preserveAspectRatio="xMaxYMid slice"
    aria-hidden="true"
  >
    <defs>
      <linearGradient id="bart5" x1="0" y1="0" x2="1" y2="1">
        <stop offset="0" stop-color="#070900" />
        <stop offset="1" stop-color="#1c2600" />
      </linearGradient>
      <filter id="fart5" x="-10%" y="-10%" width="120%" height="120%">
        <feGaussianBlur stdDeviation="6" />
      </filter>
      <linearGradient id="sart5" x1="0" y1="0" x2="1" y2="0">
        <stop offset="0" stop-color="#000" stop-opacity="0.45" />
        <stop offset="0.55" stop-color="#000" stop-opacity="0" />
      </linearGradient>
    </defs>
    <rect width="1440" height="430" fill="url(#bart5)" />
    <g opacity="0.9">
      <line
        x1="720.0"
        y1="430"
        x2="1007.99"
        y2="0"
        stroke="#d4f542"
        stroke-width="10"
        filter="url(#fart5)"
      />
      <line x1="720.0" y1="430" x2="1007.99" y2="0" stroke="#efff9a" stroke-width="2" />
    </g>
    <g opacity="0.75">
      <line
        x1="864.0"
        y1="430"
        x2="1152.0"
        y2="0"
        stroke="#d4f542"
        stroke-width="10"
        filter="url(#fart5)"
      />
      <line x1="864.0" y1="430" x2="1152.0" y2="0" stroke="#efff9a" stroke-width="2" />
    </g>
    <g opacity="0.60">
      <line
        x1="1007.99"
        y1="430"
        x2="1295.99"
        y2="0"
        stroke="#d4f542"
        stroke-width="10"
        filter="url(#fart5)"
      />
      <line x1="1007.99" y1="430" x2="1295.99" y2="0" stroke="#efff9a" stroke-width="2" />
    </g>
    <g opacity="0.45">
      <line
        x1="1152.0"
        y1="430"
        x2="1440.0"
        y2="0"
        stroke="#d4f542"
        stroke-width="10"
        filter="url(#fart5)"
      />
      <line x1="1152.0" y1="430" x2="1440.0" y2="0" stroke="#efff9a" stroke-width="2" />
    </g>
    <rect width="1440" height="430" fill="url(#sart5)" />
  </svg>
</template>

<style scoped>
svg {
  position: absolute;
  inset: 0;
  display: block;
  width: 100%;
  height: 100%;
}
</style>
