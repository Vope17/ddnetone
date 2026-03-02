<script setup>
import { computed, ref } from 'vue';

const props = defineProps({
  dailyActivity: Array // [{ date: 'YYYY-MM-DDT...', maps: N, score: N }]
});

const DAYS = 14;
const W = 600;
const H = 160;
const PAD = { top: 12, right: 12, bottom: 28, left: 36 };

const toTaipeiDate = (d) => d.toLocaleDateString('en-CA', { timeZone: 'Asia/Taipei' });

// 建立過去 30 天的每日資料
const daily = computed(() => {
  const activityMap = {};
  (props.dailyActivity || []).forEach(d => {
    activityMap[d.date.slice(0, 10)] = { maps: d.maps, score: Number(d.score) || 0 };
  });

  const todayStr = toTaipeiDate(new Date());
  const today = new Date(todayStr + 'T00:00:00');
  const result = [];
  for (let i = DAYS - 1; i >= 0; i--) {
    const d = new Date(today);
    d.setDate(today.getDate() - i);
    const key = toTaipeiDate(d);
    result.push({
      date: key,
      maps: activityMap[key]?.maps ?? 0,
      score: activityMap[key]?.score ?? 0,
      isToday: i === 0,
    });
  }
  return result;
});

const maxMaps = computed(() => Math.max(1, ...daily.value.map(d => d.maps)));

const chartW = computed(() => W - PAD.left - PAD.right);
const chartH = computed(() => H - PAD.top - PAD.bottom);
const barW = computed(() => Math.max(2, chartW.value / DAYS - 2));
const gap = computed(() => (chartW.value - barW.value * DAYS) / (DAYS - 1));

function barX(i) { return PAD.left + i * (barW.value + gap.value); }
function barH(maps) { return (maps / maxMaps.value) * chartH.value; }
function barY(maps) { return PAD.top + chartH.value - barH(maps); }

// Y 軸刻度
const yTicks = computed(() => {
  const max = maxMaps.value;
  const steps = 3;
  return Array.from({ length: steps + 1 }, (_, i) => {
    const val = Math.round((max / steps) * i);
    const y = PAD.top + chartH.value - (val / max) * chartH.value;
    return { val, y };
  });
});

// X 軸每7天標一次
const xLabels = computed(() =>
  daily.value
    .map((d, i) => ({ ...d, i }))
    .filter((_, i) => i % 7 === 0 || i === DAYS - 1)
    .map(d => ({
      x: barX(d.i) + barW.value / 2,
      label: new Date(d.date + 'T00:00:00').toLocaleString('en', { month: 'short', day: 'numeric' }),
    }))
);

const hovered = ref(null);
function onMouseMove(e) {
  const svg = e.currentTarget;
  const rect = svg.getBoundingClientRect();
  const mx = (e.clientX - rect.left) * (W / rect.width);
  let closest = null, minDist = Infinity;
  daily.value.forEach((d, i) => {
    const cx = barX(i) + barW.value / 2;
    const dist = Math.abs(cx - mx);
    if (dist < minDist) { minDist = dist; closest = { ...d, i }; }
  });
  hovered.value = closest && minDist < barW.value * 3 ? closest : null;
}
function onMouseLeave() { hovered.value = null; }
</script>

<template>
  <div class="bg-gradient-to-br from-gray-900 to-black border border-white/10 p-5 h-full flex flex-col">
    <div class="text-xs font-mono text-cyan-500/70 mb-3 flex items-center gap-2 flex-shrink-0">
      <span class="w-1.5 h-1.5 rounded-full bg-cyan-500 animate-pulse"></span>
      DAILY_ACTIVITY
      <span class="text-gray-600 ml-1">— MAPS COMPLETED PER DAY</span>
    </div>

    <div v-if="!daily.some(d => d.maps > 0)" class="flex-1 flex items-center justify-center text-gray-600 font-mono text-xs">
      NO DATA
    </div>

    <svg v-else :viewBox="`0 0 ${W} ${H}`" preserveAspectRatio="xMidYMid meet"
      class="w-full flex-1"
      @mousemove="onMouseMove" @mouseleave="onMouseLeave">

      <!-- Y 軸格線 + 刻度 -->
      <g v-for="t in yTicks" :key="t.val">
        <line :x1="PAD.left" :y1="t.y" :x2="W - PAD.right" :y2="t.y"
          stroke="rgba(255,255,255,0.05)" stroke-width="1" />
        <text v-if="t.val > 0" :x="PAD.left - 4" :y="t.y + 3" text-anchor="end"
          font-size="12" font-family="'JetBrains Mono', monospace" fill="rgba(156,163,175,0.5)">
          {{ t.val }}
        </text>
      </g>

      <!-- X 軸標籤 -->
      <text v-for="l in xLabels" :key="l.x"
        :x="l.x" :y="H - 4" text-anchor="middle"
        font-size="12" font-family="'JetBrains Mono', monospace" fill="rgba(156,163,175,0.45)">
        {{ l.label }}
      </text>

      <!-- Bars -->
      <g v-for="(d, i) in daily" :key="d.date">
        <rect v-if="d.maps > 0"
          :x="barX(i)" :y="barY(d.maps)"
          :width="barW" :height="barH(d.maps)"
          :fill="d.isToday ? 'rgba(34,211,238,1)' : hovered?.i === i ? 'rgba(34,211,238,0.9)' : 'rgba(34,211,238,0.55)'"
          :filter="d.isToday || hovered?.i === i ? 'url(#glow)' : ''"
          rx="1" />
        <!-- 空格子底線 -->
        <line v-if="d.maps === 0"
          :x1="barX(i)" :y1="PAD.top + chartH" :x2="barX(i) + barW" :y2="PAD.top + chartH"
          stroke="rgba(255,255,255,0.06)" stroke-width="1" />
      </g>

      <!-- Hover 高亮線 + tooltip -->
      <template v-if="hovered && hovered.maps > 0">
        <line :x1="barX(hovered.i) + barW / 2" :y1="PAD.top"
          :x2="barX(hovered.i) + barW / 2" :y2="PAD.top + chartH"
          stroke="rgba(34,211,238,0.2)" stroke-width="1" stroke-dasharray="3 2" />
        <g :transform="`translate(${Math.min(barX(hovered.i) + barW + 4, W - 105)}, ${Math.max(barY(hovered.maps) - 10, PAD.top)})`">
          <rect x="0" y="0" width="110" height="44" rx="4"
            fill="rgba(3,7,18,0.95)" stroke="rgba(34,211,238,0.3)" stroke-width="1" />
          <text x="8" y="16" font-size="12" font-family="'JetBrains Mono', monospace" fill="rgba(156,163,175,0.8)">
            {{ hovered.date }}
          </text>
          <text x="8" y="34" font-size="14" font-weight="bold"
            font-family="'JetBrains Mono', monospace" fill="rgba(34,211,238,1)">
            {{ hovered.maps }} maps
          </text>
        </g>
      </template>

      <defs>
        <filter id="glow" x="-50%" y="-50%" width="200%" height="200%">
          <feGaussianBlur stdDeviation="2" result="blur" />
          <feMerge><feMergeNode in="blur"/><feMergeNode in="SourceGraphic"/></feMerge>
        </filter>
      </defs>
    </svg>
  </div>
</template>
