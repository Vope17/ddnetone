<script setup>
import { computed } from 'vue';
import { ClockIcon } from '@heroicons/vue/24/outline';

const props = defineProps({
  growthData: Array,
  summary: Object
});

const TARGET_MAPS = 2403;

// 線性迴歸：給定 (x[], y[]) 回傳斜率 (單位/ms)
function linearSlope(xs, ys) {
  const n = xs.length;
  if (n < 2) return null;
  const mx = xs.reduce((a, b) => a + b, 0) / n;
  const my = ys.reduce((a, b) => a + b, 0) / n;
  const num = xs.reduce((s, x, i) => s + (x - mx) * (ys[i] - my), 0);
  const den = xs.reduce((s, x) => s + (x - mx) ** 2, 0);
  return den === 0 ? null : num / den;
}

const stats = computed(() => {
  const data = props.growthData;
  if (!data || data.length < 2) return null;

  const xs = data.map(d => new Date(d.timestamp).getTime());
  const mapYs = data.map(d => d.maps);
  const ptYs = data.map(d => d.points);

  const mapSlope = linearSlope(xs, mapYs); // maps per ms
  const ptSlope = linearSlope(xs, ptYs);   // points per ms

  const nowTs = Date.now();
  const currentMaps = props.summary?.completed_maps ?? data[data.length - 1]?.maps ?? 0;
  const currentPts = props.summary?.current_score ?? data[data.length - 1]?.points ?? 0;
  const targetPts = props.summary?.target_score ?? 10000;

  const msPerDay = 86400000;

  const mapsPerDay = mapSlope ? mapSlope * msPerDay : null;
  const ptsPerDay = ptSlope ? ptSlope * msPerDay : null;

  const mapsRemaining = TARGET_MAPS - currentMaps;
  const ptsRemaining = targetPts - currentPts;

  const mapsEta = (mapSlope && mapsRemaining > 0)
    ? new Date(nowTs + mapsRemaining / mapSlope)
    : null;
  const ptsEta = (ptSlope && ptsRemaining > 0)
    ? new Date(nowTs + ptsRemaining / ptSlope)
    : null;

  const fmt = (d) => d ? d.toLocaleDateString('sv-SE').replace(/-/g, '/') : '—';

  return {
    mapsPerDay: mapsPerDay ? mapsPerDay.toFixed(1) : '—',
    ptsPerDay: ptsPerDay ? Math.round(ptsPerDay).toLocaleString() : '—',
    mapsEta: fmt(mapsEta),
    ptsEta: fmt(ptsEta),
    mapsComplete: currentMaps >= TARGET_MAPS,
    ptsComplete: currentPts >= targetPts,
  };
});
</script>

<template>
  <div v-if="stats"
    class="bg-gradient-to-br from-gray-900 to-black border border-white/10 p-5 relative overflow-hidden">
    <div class="absolute right-0 top-0 p-4 opacity-10">
      <ClockIcon class="w-14 h-14 text-amber-400" />
    </div>

    <div class="text-xs font-mono text-amber-400/70 mb-3 flex items-center gap-2">
      <span class="w-1 h-1 bg-amber-400"></span> ETA_ANALYSIS
    </div>

    <div class="grid grid-cols-2 gap-x-4 gap-y-3">
      <!-- Maps -->
      <div>
        <div class="text-[9px] font-mono text-gray-500 tracking-widest mb-0.5">MAPS / DAY</div>
        <div class="text-lg font-bold font-mono text-white">{{ stats.mapsPerDay }}</div>
      </div>
      <div>
        <div class="text-[9px] font-mono text-gray-500 tracking-widest mb-0.5">MAPS ETA</div>
        <div class="text-sm font-bold font-mono"
          :class="stats.mapsComplete ? 'text-green-400' : 'text-amber-300'">
          {{ stats.mapsComplete ? 'COMPLETE' : stats.mapsEta }}
        </div>
      </div>

      <!-- Points -->
      <div>
        <div class="text-[9px] font-mono text-gray-500 tracking-widest mb-0.5">PTS / DAY</div>
        <div class="text-lg font-bold font-mono text-white">{{ stats.ptsPerDay }}</div>
      </div>
      <div>
        <div class="text-[9px] font-mono text-gray-500 tracking-widest mb-0.5">SCORE ETA</div>
        <div class="text-sm font-bold font-mono"
          :class="stats.ptsComplete ? 'text-green-400' : 'text-amber-300'">
          {{ stats.ptsComplete ? 'COMPLETE' : stats.ptsEta }}
        </div>
      </div>
    </div>

    <div class="mt-3 text-[9px] font-mono text-gray-600">
      BASED ON {{ growthData?.length ?? 0 }} DATA POINTS (7D)
    </div>
  </div>
</template>
