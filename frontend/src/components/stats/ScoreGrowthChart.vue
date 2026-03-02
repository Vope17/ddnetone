<script setup>
import { computed, ref, onMounted, onUnmounted } from 'vue';

const props = defineProps({
  dailyActivity: Array // [{ date: 'YYYY-MM-DDT...', maps: N, score: N }]
});

const W = 800;
const H = 180;
const PAD = { top: 16, right: 20, bottom: 32, left: 52 };

// 累計分數資料
const points = computed(() => {
  const sorted = [...(props.dailyActivity || [])]
    .map(d => ({ date: d.date.slice(0, 10), score: Number(d.score) || 0 }))
    .sort((a, b) => a.date.localeCompare(b.date));

  let cum = 0;
  return sorted.map(d => ({ date: d.date, score: (cum += d.score) }));
});

const chartW = computed(() => W - PAD.left - PAD.right);
const chartH = computed(() => H - PAD.top - PAD.bottom);

const maxScore = computed(() => Math.max(1, ...points.value.map(p => p.score)));

const svgPoints = computed(() => {
  const n = points.value.length;
  if (n === 0) return [];
  return points.value.map((p, i) => ({
    x: PAD.left + (i / Math.max(n - 1, 1)) * chartW.value,
    y: PAD.top + chartH.value - (p.score / maxScore.value) * chartH.value,
    date: p.date,
    score: p.score,
  }));
});

// 平滑貝茲曲線 path
const linePath = computed(() => {
  const pts = svgPoints.value;
  if (pts.length < 2) return '';
  let d = `M ${pts[0].x} ${pts[0].y}`;
  for (let i = 1; i < pts.length; i++) {
    const cp1x = (pts[i - 1].x + pts[i].x) / 2;
    d += ` C ${cp1x} ${pts[i - 1].y}, ${cp1x} ${pts[i].y}, ${pts[i].x} ${pts[i].y}`;
  }
  return d;
});

// 填色 area path
const areaPath = computed(() => {
  const pts = svgPoints.value;
  if (pts.length < 2) return '';
  const bottom = PAD.top + chartH.value;
  let d = `M ${pts[0].x} ${bottom} L ${pts[0].x} ${pts[0].y}`;
  for (let i = 1; i < pts.length; i++) {
    const cp1x = (pts[i - 1].x + pts[i].x) / 2;
    d += ` C ${cp1x} ${pts[i - 1].y}, ${cp1x} ${pts[i].y}, ${pts[i].x} ${pts[i].y}`;
  }
  d += ` L ${pts[pts.length - 1].x} ${bottom} Z`;
  return d;
});

// Y 軸刻度
const yTicks = computed(() => {
  const max = maxScore.value;
  const steps = 4;
  return Array.from({ length: steps + 1 }, (_, i) => {
    const val = Math.round((max / steps) * i);
    const y = PAD.top + chartH.value - (val / max) * chartH.value;
    return { val, y };
  });
});

// X 軸月份標籤
const xLabels = computed(() => {
  const pts = svgPoints.value;
  const seen = new Set();
  return pts.filter(p => {
    const m = p.date.slice(0, 7);
    if (seen.has(m)) return false;
    seen.add(m);
    return true;
  }).map(p => ({
    x: p.x,
    label: new Date(p.date + 'T00:00:00').toLocaleString('en', { month: 'short', year: '2-digit' }),
  }));
});

function fmtScore(n) {
  if (n >= 1000) return (n / 1000).toFixed(1) + 'k';
  return String(n);
}

// Hover
const hovered = ref(null);
function onMouseMove(e) {
  const svg = e.currentTarget;
  const rect = svg.getBoundingClientRect();
  const mx = (e.clientX - rect.left) * (W / rect.width);
  const pts = svgPoints.value;
  if (!pts.length) return;
  let closest = pts[0];
  let minDist = Math.abs(pts[0].x - mx);
  for (const p of pts) {
    const d = Math.abs(p.x - mx);
    if (d < minDist) { minDist = d; closest = p; }
  }
  hovered.value = closest;
}
function onMouseLeave() { hovered.value = null; }
</script>

<template>
  <div class="bg-gradient-to-br from-gray-900 to-black border border-white/10 p-5">
    <div class="text-xs font-mono text-green-400/70 mb-4 flex items-center gap-2">
      <span class="w-1.5 h-1.5 rounded-full bg-green-400 animate-pulse"></span>
      SCORE_GROWTH
      <span class="text-gray-600 ml-1">— CUMULATIVE POINTS OVER TIME</span>
    </div>

    <div v-if="!points.length" class="flex items-center justify-center h-[180px] text-gray-600 font-mono text-xs">
      NO DATA
    </div>

    <svg v-else :viewBox="`0 0 ${W} ${H}`" preserveAspectRatio="xMidYMid meet"
      class="w-full" :style="{ height: H + 'px' }"
      @mousemove="onMouseMove" @mouseleave="onMouseLeave">

      <defs>
        <linearGradient id="scoreGrad" x1="0" y1="0" x2="0" y2="1">
          <stop offset="0%" stop-color="rgba(34,197,94,0.25)" />
          <stop offset="100%" stop-color="rgba(34,197,94,0)" />
        </linearGradient>
      </defs>

      <!-- Y 軸格線 + 刻度 -->
      <g v-for="t in yTicks" :key="t.val">
        <line :x1="PAD.left" :y1="t.y" :x2="W - PAD.right" :y2="t.y"
          stroke="rgba(255,255,255,0.05)" stroke-width="1" />
        <text :x="PAD.left - 6" :y="t.y + 4" text-anchor="end"
          font-size="9" font-family="'JetBrains Mono', monospace" fill="rgba(156,163,175,0.6)">
          {{ fmtScore(t.val) }}
        </text>
      </g>

      <!-- X 軸月份標籤 -->
      <text v-for="l in xLabels" :key="l.x"
        :x="l.x" :y="H - 6" text-anchor="middle"
        font-size="9" font-family="'JetBrains Mono', monospace" fill="rgba(156,163,175,0.5)">
        {{ l.label }}
      </text>

      <!-- Area fill -->
      <path :d="areaPath" fill="url(#scoreGrad)" />

      <!-- Line -->
      <path :d="linePath" fill="none" stroke="rgba(34,197,94,0.9)" stroke-width="2"
        stroke-linecap="round" stroke-linejoin="round" />

      <!-- Hover vertical line -->
      <template v-if="hovered">
        <line :x1="hovered.x" :y1="PAD.top" :x2="hovered.x" :y2="PAD.top + chartH"
          stroke="rgba(34,197,94,0.3)" stroke-width="1" stroke-dasharray="3 3" />
        <circle :cx="hovered.x" :cy="hovered.y" r="4"
          fill="#22c55e" stroke="rgba(0,0,0,0.8)" stroke-width="2" />

        <!-- Tooltip box -->
        <g :transform="`translate(${Math.min(hovered.x + 8, W - 110)}, ${Math.max(hovered.y - 40, PAD.top)})`">
          <rect x="0" y="0" width="100" height="38" rx="4"
            fill="rgba(3,7,18,0.95)" stroke="rgba(34,197,94,0.3)" stroke-width="1" />
          <text x="8" y="14" font-size="9" font-family="'JetBrains Mono', monospace" fill="rgba(156,163,175,0.8)">
            {{ hovered.date }}
          </text>
          <text x="8" y="29" font-size="11" font-weight="bold"
            font-family="'JetBrains Mono', monospace" fill="#4ade80">
            {{ hovered.score.toLocaleString() }} pts
          </text>
        </g>
      </template>
    </svg>
  </div>
</template>
