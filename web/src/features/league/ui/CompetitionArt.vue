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
    v-else-if="art === 'uefa.europa'"
    viewBox="0 0 1440 430"
    preserveAspectRatio="xMaxYMid slice"
    aria-hidden="true"
  >
    <defs>
      <linearGradient id="bart1" x1="0" y1="0" x2="1" y2="1">
        <stop offset="0" stop-color="#0a0503" />
        <stop offset="0.6" stop-color="#2a0f02" />
        <stop offset="1" stop-color="#5a2205" />
      </linearGradient>
      <linearGradient id="rart1" x1="0" y1="0" x2="0" y2="1">
        <stop offset="0" stop-color="#ffb066" stop-opacity="0.55" />
        <stop offset="1" stop-color="#ff6a00" stop-opacity="0" />
      </linearGradient>
      <filter id="fart1" x="-20%" y="-20%" width="140%" height="140%">
        <feGaussianBlur stdDeviation="8" />
      </filter>
      <linearGradient id="sart1" x1="0" y1="0" x2="1" y2="0">
        <stop offset="0" stop-color="#000" stop-opacity="0.45" />
        <stop offset="0.55" stop-color="#000" stop-opacity="0" />
      </linearGradient>
    </defs>
    <rect width="1440" height="430" fill="url(#bart1)" />
    <polygon points="892.8,-10 918.8,-10 822.4,430 662.4,430" fill="url(#rart1)" opacity="0.55" />
    <polygon points="1094.4,-10 1120.4,-10 1024.0,430 864.0,430" fill="url(#rart1)" opacity="0.4" />
    <polygon
      points="1296.0,-10 1322.0,-10 1225.6,430 1065.6,430"
      fill="url(#rart1)"
      opacity="0.5"
    />
    <g filter="url(#fart1)">
      <ellipse
        cx="1152.0"
        cy="266.6"
        rx="288.0"
        ry="68.8"
        fill="none"
        stroke="#ff8a1f"
        stroke-width="3"
        opacity="0.8"
        transform="rotate(-12 1152.0 266.6)"
      />
      <ellipse
        cx="1152.0"
        cy="266.6"
        rx="360.0"
        ry="90.30"
        fill="none"
        stroke="#ff8a1f"
        stroke-width="2"
        opacity="0.58"
        transform="rotate(-12 1152.0 266.6)"
      />
      <ellipse
        cx="1152.0"
        cy="266.6"
        rx="432.00"
        ry="111.8"
        fill="none"
        stroke="#ff8a1f"
        stroke-width="1"
        opacity="0.36"
        transform="rotate(-12 1152.0 266.6)"
      />
    </g>
    <ellipse
      cx="1152.0"
      cy="266.6"
      rx="288.0"
      ry="68.8"
      fill="none"
      stroke="#ff8a1f"
      stroke-width="3"
      opacity="0.8"
      transform="rotate(-12 1152.0 266.6)"
    />
    <ellipse
      cx="1152.0"
      cy="266.6"
      rx="360.0"
      ry="90.30"
      fill="none"
      stroke="#ff8a1f"
      stroke-width="2"
      opacity="0.58"
      transform="rotate(-12 1152.0 266.6)"
    />
    <ellipse
      cx="1152.0"
      cy="266.6"
      rx="432.00"
      ry="111.8"
      fill="none"
      stroke="#ff8a1f"
      stroke-width="1"
      opacity="0.36"
      transform="rotate(-12 1152.0 266.6)"
    />
    <circle cx="485" cy="105" r="1" fill="#fff" opacity="0.28" />
    <circle cx="946" cy="132" r="1.6" fill="#fff" opacity="0.50" />
    <circle cx="224" cy="182" r="1.2" fill="#fff" opacity="0.27" />
    <circle cx="843" cy="89" r="0.8" fill="#fff" opacity="0.38" />
    <circle cx="236" cy="47" r="1" fill="#fff" opacity="0.37" />
    <circle cx="380" cy="173" r="0.8" fill="#fff" opacity="0.67" />
    <circle cx="507" cy="119" r="1.2" fill="#fff" opacity="0.52" />
    <circle cx="280" cy="69" r="1.6" fill="#fff" opacity="0.62" />
    <circle cx="444" cy="146" r="1" fill="#fff" opacity="0.69" />
    <circle cx="874" cy="40" r="1.2" fill="#fff" opacity="0.79" />
    <circle cx="538" cy="10" r="1.2" fill="#fff" opacity="0.46" />
    <circle cx="1248" cy="7" r="1.6" fill="#fff" opacity="0.41" />
    <circle cx="399" cy="60" r="0.8" fill="#fff" opacity="0.75" />
    <circle cx="983" cy="46" r="1.2" fill="#fff" opacity="0.46" />
    <circle cx="1271" cy="76" r="1.6" fill="#fff" opacity="0.31" />
    <circle cx="276" cy="103" r="1" fill="#fff" opacity="0.52" />
    <circle cx="237" cy="79" r="1.2" fill="#fff" opacity="0.68" />
    <circle cx="188" cy="16" r="1.2" fill="#fff" opacity="0.40" />
    <circle cx="537" cy="193" r="1" fill="#fff" opacity="0.79" />
    <circle cx="1272" cy="75" r="1.2" fill="#fff" opacity="0.76" />
    <circle cx="1375" cy="138" r="1.2" fill="#fff" opacity="0.44" />
    <circle cx="633" cy="117" r="1.2" fill="#fff" opacity="0.49" />
    <circle cx="1062" cy="74" r="1.2" fill="#fff" opacity="0.44" />
    <circle cx="1272" cy="44" r="1.6" fill="#fff" opacity="0.29" />
    <circle cx="1304" cy="129" r="0.8" fill="#fff" opacity="0.38" />
    <circle cx="688" cy="123" r="1" fill="#fff" opacity="0.33" />
    <circle cx="1165" cy="38" r="1.6" fill="#fff" opacity="0.79" />
    <circle cx="599" cy="142" r="0.8" fill="#fff" opacity="0.35" />
    <circle cx="859" cy="78" r="0.8" fill="#fff" opacity="0.68" />
    <circle cx="1107" cy="113" r="1.6" fill="#fff" opacity="0.67" />
    <rect width="1440" height="430" fill="url(#sart1)" />
  </svg>
  <svg
    v-else-if="art === 'uefa.europa.conf'"
    viewBox="0 0 1440 430"
    preserveAspectRatio="xMaxYMid slice"
    aria-hidden="true"
  >
    <defs>
      <linearGradient id="bart2" x1="0" y1="0" x2="1" y2="0">
        <stop offset="0" stop-color="#02140a" />
        <stop offset="0.65" stop-color="#04351a" />
        <stop offset="1" stop-color="#087a36" />
      </linearGradient>
      <radialGradient id="gart2" cx="0.85" cy="0.4" r="0.5">
        <stop offset="0" stop-color="#3dff8a" stop-opacity="0.45" />
        <stop offset="1" stop-color="#3dff8a" stop-opacity="0" />
      </radialGradient>
      <linearGradient id="sart2" x1="0" y1="0" x2="1" y2="0">
        <stop offset="0" stop-color="#000" stop-opacity="0.45" />
        <stop offset="0.55" stop-color="#000" stop-opacity="0" />
      </linearGradient>
    </defs>
    <rect width="1440" height="430" fill="url(#bart2)" />
    <rect width="1440" height="430" fill="url(#gart2)" />
    <polygon
      points="764.6,25.5 720.5,51.0 676.4,25.5 676.4,-25.5 720.5,-51.0 764.6,-25.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.05"
    />
    <polygon
      points="854.7,25.5 810.6,51.0 766.4,25.5 766.4,-25.5 810.6,-51.0 854.7,-25.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.11"
    />
    <polygon
      points="944.8,25.5 900.6,51.0 856.5,25.5 856.5,-25.5 900.6,-51.0 944.8,-25.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.18"
    />
    <polygon
      points="1034.8,25.5 990.7,51.0 946.6,25.5 946.6,-25.5 990.7,-51.0 1034.8,-25.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.24"
    />
    <polygon
      points="1124.9,25.5 1080.8,51.0 1036.6,25.5 1036.6,-25.5 1080.8,-51.0 1124.9,-25.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.30"
    />
    <polygon
      points="1215.0,25.5 1170.8,51.0 1126.7,25.5 1126.7,-25.5 1170.8,-51.0 1215.0,-25.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.36"
    />
    <polygon
      points="1305.0,25.5 1260.9,51.0 1216.8,25.5 1216.8,-25.5 1260.9,-51.0 1305.0,-25.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.43"
    />
    <polygon
      points="1395.1,25.5 1351.0,51.0 1306.8,25.5 1306.8,-25.5 1351.0,-51.0 1395.1,-25.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.49"
    />
    <polygon
      points="1485.2,25.5 1441.0,51.0 1396.9,25.5 1396.9,-25.5 1441.0,-51.0 1485.2,-25.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.50"
    />
    <polygon
      points="719.6,103.5 675.5,129.0 631.3,103.5 631.3,52.5 675.5,27.0 719.6,52.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.02"
    />
    <polygon
      points="809.7,103.5 765.5,129.0 721.4,103.5 721.4,52.5 765.5,27.0 809.7,52.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.08"
    />
    <polygon
      points="899.7,103.5 855.6,129.0 811.5,103.5 811.5,52.5 855.6,27.0 899.7,52.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.14"
    />
    <polygon
      points="989.8,103.5 945.7,129.0 901.5,103.5 901.5,52.5 945.7,27.0 989.8,52.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.21"
    />
    <polygon
      points="1079.9,103.5 1035.7,129.0 991.6,103.5 991.6,52.5 1035.7,27.0 1079.9,52.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.27"
    />
    <polygon
      points="1169.9,103.5 1125.8,129.0 1081.7,103.5 1081.7,52.5 1125.8,27.0 1169.9,52.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.33"
    />
    <polygon
      points="1260.0,103.5 1215.9,129.0 1171.7,103.5 1171.7,52.5 1215.9,27.0 1260.0,52.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.39"
    />
    <polygon
      points="1350.1,103.5 1305.9,129.0 1261.8,103.5 1261.8,52.5 1305.9,27.0 1350.1,52.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.46"
    />
    <polygon
      points="1440.1,103.5 1396.0,129.0 1351.9,103.5 1351.9,52.5 1396.0,27.0 1440.1,52.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.50"
    />
    <polygon
      points="1530.2,103.5 1486.1,129.0 1441.9,103.5 1441.9,52.5 1486.1,27.0 1530.2,52.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.50"
    />
    <polygon
      points="764.6,181.5 720.5,207.0 676.4,181.5 676.4,130.5 720.5,105.0 764.6,130.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.05"
    />
    <polygon
      points="854.7,181.5 810.6,207.0 766.4,181.5 766.4,130.5 810.6,105.0 854.7,130.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.11"
    />
    <polygon
      points="944.8,181.5 900.6,207.0 856.5,181.5 856.5,130.5 900.6,105.0 944.8,130.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.18"
    />
    <polygon
      points="1034.8,181.5 990.7,207.0 946.6,181.5 946.6,130.5 990.7,105.0 1034.8,130.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.24"
    />
    <polygon
      points="1124.9,181.5 1080.8,207.0 1036.6,181.5 1036.6,130.5 1080.8,105.0 1124.9,130.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.30"
    />
    <polygon
      points="1215.0,181.5 1170.8,207.0 1126.7,181.5 1126.7,130.5 1170.8,105.0 1215.0,130.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.36"
    />
    <polygon
      points="1305.0,181.5 1260.9,207.0 1216.8,181.5 1216.8,130.5 1260.9,105.0 1305.0,130.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.43"
    />
    <polygon
      points="1395.1,181.5 1351.0,207.0 1306.8,181.5 1306.8,130.5 1351.0,105.0 1395.1,130.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.49"
    />
    <polygon
      points="1485.2,181.5 1441.0,207.0 1396.9,181.5 1396.9,130.5 1441.0,105.0 1485.2,130.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.50"
    />
    <polygon
      points="719.6,259.5 675.5,285.0 631.3,259.5 631.3,208.5 675.5,183.0 719.6,208.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.02"
    />
    <polygon
      points="809.7,259.5 765.5,285.0 721.4,259.5 721.4,208.5 765.5,183.0 809.7,208.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.08"
    />
    <polygon
      points="899.7,259.5 855.6,285.0 811.5,259.5 811.5,208.5 855.6,183.0 899.7,208.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.14"
    />
    <polygon
      points="989.8,259.5 945.7,285.0 901.5,259.5 901.5,208.5 945.7,183.0 989.8,208.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.21"
    />
    <polygon
      points="1079.9,259.5 1035.7,285.0 991.6,259.5 991.6,208.5 1035.7,183.0 1079.9,208.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.27"
    />
    <polygon
      points="1169.9,259.5 1125.8,285.0 1081.7,259.5 1081.7,208.5 1125.8,183.0 1169.9,208.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.33"
    />
    <polygon
      points="1260.0,259.5 1215.9,285.0 1171.7,259.5 1171.7,208.5 1215.9,183.0 1260.0,208.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.39"
    />
    <polygon
      points="1350.1,259.5 1305.9,285.0 1261.8,259.5 1261.8,208.5 1305.9,183.0 1350.1,208.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.46"
    />
    <polygon
      points="1440.1,259.5 1396.0,285.0 1351.9,259.5 1351.9,208.5 1396.0,183.0 1440.1,208.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.50"
    />
    <polygon
      points="1530.2,259.5 1486.1,285.0 1441.9,259.5 1441.9,208.5 1486.1,183.0 1530.2,208.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.50"
    />
    <polygon
      points="764.6,337.5 720.5,363.0 676.4,337.5 676.4,286.5 720.5,261.0 764.6,286.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.05"
    />
    <polygon
      points="854.7,337.5 810.6,363.0 766.4,337.5 766.4,286.5 810.6,261.0 854.7,286.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.11"
    />
    <polygon
      points="944.8,337.5 900.6,363.0 856.5,337.5 856.5,286.5 900.6,261.0 944.8,286.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.18"
    />
    <polygon
      points="1034.8,337.5 990.7,363.0 946.6,337.5 946.6,286.5 990.7,261.0 1034.8,286.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.24"
    />
    <polygon
      points="1124.9,337.5 1080.8,363.0 1036.6,337.5 1036.6,286.5 1080.8,261.0 1124.9,286.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.30"
    />
    <polygon
      points="1215.0,337.5 1170.8,363.0 1126.7,337.5 1126.7,286.5 1170.8,261.0 1215.0,286.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.36"
    />
    <polygon
      points="1305.0,337.5 1260.9,363.0 1216.8,337.5 1216.8,286.5 1260.9,261.0 1305.0,286.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.43"
    />
    <polygon
      points="1395.1,337.5 1351.0,363.0 1306.8,337.5 1306.8,286.5 1351.0,261.0 1395.1,286.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.49"
    />
    <polygon
      points="1485.2,337.5 1441.0,363.0 1396.9,337.5 1396.9,286.5 1441.0,261.0 1485.2,286.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.50"
    />
    <polygon
      points="719.6,415.5 675.5,441.0 631.3,415.5 631.3,364.5 675.5,339.0 719.6,364.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.02"
    />
    <polygon
      points="809.7,415.5 765.5,441.0 721.4,415.5 721.4,364.5 765.5,339.0 809.7,364.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.08"
    />
    <polygon
      points="899.7,415.5 855.6,441.0 811.5,415.5 811.5,364.5 855.6,339.0 899.7,364.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.14"
    />
    <polygon
      points="989.8,415.5 945.7,441.0 901.5,415.5 901.5,364.5 945.7,339.0 989.8,364.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.21"
    />
    <polygon
      points="1079.9,415.5 1035.7,441.0 991.6,415.5 991.6,364.5 1035.7,339.0 1079.9,364.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.27"
    />
    <polygon
      points="1169.9,415.5 1125.8,441.0 1081.7,415.5 1081.7,364.5 1125.8,339.0 1169.9,364.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.33"
    />
    <polygon
      points="1260.0,415.5 1215.9,441.0 1171.7,415.5 1171.7,364.5 1215.9,339.0 1260.0,364.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.39"
    />
    <polygon
      points="1350.1,415.5 1305.9,441.0 1261.8,415.5 1261.8,364.5 1305.9,339.0 1350.1,364.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.46"
    />
    <polygon
      points="1440.1,415.5 1396.0,441.0 1351.9,415.5 1351.9,364.5 1396.0,339.0 1440.1,364.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.50"
    />
    <polygon
      points="1530.2,415.5 1486.1,441.0 1441.9,415.5 1441.9,364.5 1486.1,339.0 1530.2,364.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.50"
    />
    <polygon
      points="764.6,493.5 720.5,519.0 676.4,493.5 676.4,442.5 720.5,417.0 764.6,442.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.05"
    />
    <polygon
      points="854.7,493.5 810.6,519.0 766.4,493.5 766.4,442.5 810.6,417.0 854.7,442.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.11"
    />
    <polygon
      points="944.8,493.5 900.6,519.0 856.5,493.5 856.5,442.5 900.6,417.0 944.8,442.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.18"
    />
    <polygon
      points="1034.8,493.5 990.7,519.0 946.6,493.5 946.6,442.5 990.7,417.0 1034.8,442.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.24"
    />
    <polygon
      points="1124.9,493.5 1080.8,519.0 1036.6,493.5 1036.6,442.5 1080.8,417.0 1124.9,442.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.30"
    />
    <polygon
      points="1215.0,493.5 1170.8,519.0 1126.7,493.5 1126.7,442.5 1170.8,417.0 1215.0,442.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.36"
    />
    <polygon
      points="1305.0,493.5 1260.9,519.0 1216.8,493.5 1216.8,442.5 1260.9,417.0 1305.0,442.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.43"
    />
    <polygon
      points="1395.1,493.5 1351.0,519.0 1306.8,493.5 1306.8,442.5 1351.0,417.0 1395.1,442.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.49"
    />
    <polygon
      points="1485.2,493.5 1441.0,519.0 1396.9,493.5 1396.9,442.5 1441.0,417.0 1485.2,442.5"
      fill="none"
      stroke="#5dffa0"
      stroke-width="1.2"
      opacity="0.50"
    />
    <rect width="1440" height="430" fill="url(#sart2)" />
  </svg>
  <svg
    v-else-if="art === 'esp.1'"
    viewBox="0 0 1440 430"
    preserveAspectRatio="xMaxYMid slice"
    aria-hidden="true"
  >
    <defs>
      <linearGradient id="bart3" x1="0" y1="1" x2="1" y2="0">
        <stop offset="0" stop-color="#8f0f0c" />
        <stop offset="0.45" stop-color="#e2380f" />
        <stop offset="1" stop-color="#ff9a1f" />
      </linearGradient>
      <linearGradient id="sart3" x1="0" y1="0" x2="1" y2="0">
        <stop offset="0" stop-color="#000" stop-opacity="0.45" />
        <stop offset="0.55" stop-color="#000" stop-opacity="0" />
      </linearGradient>
    </defs>
    <rect width="1440" height="430" fill="url(#bart3)" />
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
    <rect width="1440" height="430" fill="url(#sart3)" />
  </svg>
  <svg
    v-else-if="art === 'eng.1'"
    viewBox="0 0 1440 430"
    preserveAspectRatio="xMaxYMid slice"
    aria-hidden="true"
  >
    <defs>
      <radialGradient id="bart4" cx="0.78" cy="0.35" r="0.85">
        <stop offset="0" stop-color="#7a2fd0" />
        <stop offset="0.5" stop-color="#3d195b" />
        <stop offset="1" stop-color="#14051f" />
      </radialGradient>
      <pattern id="dart4" width="14" height="14" patternUnits="userSpaceOnUse">
        <circle cx="7" cy="7" r="2" fill="#f4b8ff" />
      </pattern>
      <linearGradient id="mart4" x1="0" y1="0" x2="1" y2="0">
        <stop offset="0.35" stop-color="#fff" stop-opacity="0" />
        <stop offset="1" stop-color="#fff" stop-opacity="1" />
      </linearGradient>
      <mask id="kart4"><rect width="1440" height="430" fill="url(#mart4)" /></mask>
      <linearGradient id="sart4" x1="0" y1="0" x2="1" y2="0">
        <stop offset="0" stop-color="#000" stop-opacity="0.45" />
        <stop offset="0.55" stop-color="#000" stop-opacity="0" />
      </linearGradient>
    </defs>
    <rect width="1440" height="430" fill="url(#bart4)" />
    <polygon
      points="648.0,430 1080.0,0 1440,0 1440,430"
      fill="url(#dart4)"
      opacity="0.22"
      mask="url(#kart4)"
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
    <rect width="1440" height="430" fill="url(#sart4)" />
  </svg>
  <svg
    v-else-if="art === 'ger.1'"
    viewBox="0 0 1440 430"
    preserveAspectRatio="xMaxYMid slice"
    aria-hidden="true"
  >
    <defs>
      <linearGradient id="bart5" x1="0" y1="0" x2="1" y2="1">
        <stop offset="0" stop-color="#4a0309" />
        <stop offset="0.6" stop-color="#b8050f" />
        <stop offset="1" stop-color="#ff2b37" />
      </linearGradient>
      <linearGradient id="sart5" x1="0" y1="0" x2="1" y2="0">
        <stop offset="0" stop-color="#000" stop-opacity="0.45" />
        <stop offset="0.55" stop-color="#000" stop-opacity="0" />
      </linearGradient>
    </defs>
    <rect width="1440" height="430" fill="url(#bart5)" />
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
    <rect width="1440" height="430" fill="url(#sart5)" />
  </svg>
  <svg
    v-else-if="art === 'ita.1'"
    viewBox="0 0 1440 430"
    preserveAspectRatio="xMaxYMid slice"
    aria-hidden="true"
  >
    <defs>
      <linearGradient id="bart6" x1="0" y1="0" x2="1" y2="0">
        <stop offset="0" stop-color="#061650" />
        <stop offset="0.6" stop-color="#0f37b8" />
        <stop offset="1" stop-color="#2f74ff" />
      </linearGradient>
      <linearGradient id="sart6" x1="0" y1="0" x2="1" y2="0">
        <stop offset="0" stop-color="#000" stop-opacity="0.45" />
        <stop offset="0.55" stop-color="#000" stop-opacity="0" />
      </linearGradient>
    </defs>
    <rect width="1440" height="430" fill="url(#bart6)" />
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
    <rect width="1440" height="430" fill="url(#sart6)" />
  </svg>
  <svg
    v-else-if="art === 'fra.1'"
    viewBox="0 0 1440 430"
    preserveAspectRatio="xMaxYMid slice"
    aria-hidden="true"
  >
    <defs>
      <linearGradient id="bart7" x1="0" y1="0" x2="1" y2="1">
        <stop offset="0" stop-color="#070900" />
        <stop offset="1" stop-color="#1c2600" />
      </linearGradient>
      <filter id="fart7" x="-10%" y="-10%" width="120%" height="120%">
        <feGaussianBlur stdDeviation="6" />
      </filter>
      <linearGradient id="sart7" x1="0" y1="0" x2="1" y2="0">
        <stop offset="0" stop-color="#000" stop-opacity="0.45" />
        <stop offset="0.55" stop-color="#000" stop-opacity="0" />
      </linearGradient>
    </defs>
    <rect width="1440" height="430" fill="url(#bart7)" />
    <g opacity="0.9">
      <line
        x1="720.0"
        y1="430"
        x2="1007.99"
        y2="0"
        stroke="#d4f542"
        stroke-width="10"
        filter="url(#fart7)"
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
        filter="url(#fart7)"
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
        filter="url(#fart7)"
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
        filter="url(#fart7)"
      />
      <line x1="1152.0" y1="430" x2="1440.0" y2="0" stroke="#efff9a" stroke-width="2" />
    </g>
    <rect width="1440" height="430" fill="url(#sart7)" />
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
