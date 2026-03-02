<script setup>
import { computed, ref } from 'vue';

const props = defineProps({
  dailyActivity: Array // [{ date: 'YYYY-MM-DD', maps: N, score: N }]
});

const tooltip = ref(null); // { date, maps, score, top, left }

const toTaipeiDate = (d) =>
  d.toLocaleDateString('en-CA', { timeZone: 'Asia/Taipei' }); // YYYY-MM-DD

const weeks = computed(() => {
  const activityMap = {};
  (props.dailyActivity || []).forEach(d => {
    const key = d.date.slice(0, 10);
    activityMap[key] = d;
  });

  const todayStr = toTaipeiDate(new Date());
  const today = new Date(todayStr + 'T00:00:00');
  const dayOfWeek = today.getDay();
  const start = new Date(today);
  start.setDate(today.getDate() - dayOfWeek - 3 * 7);

  const allWeeks = [];
  let week = [];
  const cursor = new Date(start);

  while (cursor <= today) {
    const iso = toTaipeiDate(cursor);
    const data = activityMap[iso];
    week.push({ date: iso, maps: data?.maps ?? 0, score: data?.score ?? 0 });
    if (week.length === 7) { allWeeks.push(week); week = []; }
    cursor.setDate(cursor.getDate() + 1);
  }
  if (week.length > 0) allWeeks.push(week);
  return allWeeks;
});

const maxMaps = computed(() => Math.max(1, ...(props.dailyActivity || []).map(d => d.maps)));

const cellStyle = (maps) => {
  if (maps === 0) return { background: 'rgba(255,255,255,0.04)', border: '1px solid rgba(255,255,255,0.04)' };
  const ratio = maps / maxMaps.value;
  const alpha = 0.25 + ratio * 0.75;
  return {
    background: `rgba(34,211,238,${alpha})`,
    border: `1px solid rgba(34,211,238,${Math.min(1, alpha + 0.2)})`,
    boxShadow: ratio > 0.6 ? `0 0 6px rgba(34,211,238,${ratio * 0.6})` : 'none',
  };
};

const DAY_LABELS = ['Sun', '', 'Tue', '', 'Thu', '', 'Sat'];

const monthLabels = computed(() => {
  const labels = [];
  let lastMonth = null;
  weeks.value.forEach((week, wi) => {
    const m = week[0]?.date?.slice(0, 7);
    if (m && m !== lastMonth) {
      labels.push({ week: wi, label: new Date(week[0].date + 'T00:00:00').toLocaleString('en', { month: 'short' }) });
      lastMonth = m;
    }
  });
  return labels;
});

function showTooltip(event, day, rowIndex) {
  const rect = event.target.getBoundingClientRect();
  const containerRect = event.target.closest('.heatmap-container').getBoundingClientRect();
  const showBelow = rowIndex < 3; // top rows: show tooltip below the cell
  tooltip.value = {
    date: day.date,
    maps: day.maps,
    score: day.score,
    top: showBelow
      ? rect.bottom - containerRect.top + 4
      : rect.top - containerRect.top - 64,
    left: rect.left - containerRect.left + 6,
  };
}
function hideTooltip() { tooltip.value = null; }
</script>

<template>
  <div class="bg-gradient-to-br from-gray-900 to-black border border-white/10 p-5">
    <div class="text-xs font-mono text-cyan-500/70 mb-4 flex items-center gap-2">
      <span class="w-1.5 h-1.5 rounded-full bg-cyan-500 animate-pulse"></span>
      ACTIVITY_HEATMAP
      <span class="text-gray-600 ml-1">— PAST MONTH</span>
    </div>

    <div class="overflow-x-auto heatmap-container relative">
      <!-- JS Tooltip -->
      <Transition name="fade">
        <div v-if="tooltip" class="absolute z-50 pointer-events-none"
          :style="{ top: tooltip.top + 'px', left: tooltip.left + 'px' }">
          <div class="bg-gray-950 border border-cyan-500/30 px-2.5 py-1.5 rounded-md text-[10px] font-mono whitespace-nowrap shadow-2xl shadow-cyan-500/10">
            <div class="text-gray-400 mb-0.5">{{ tooltip.date }}</div>
            <div class="text-cyan-300 font-bold">{{ tooltip.maps }} maps</div>
            <div v-if="tooltip.score" class="text-green-400">{{ tooltip.score.toLocaleString() }} pts</div>
          </div>
        </div>
      </Transition>

      <div class="inline-flex flex-col gap-1 min-w-max">
        <!-- 月份標籤 -->
        <div class="flex gap-[4px] ml-8">
          <template v-for="(w, wi) in weeks" :key="wi">
            <div class="w-[13px] text-[9px] font-mono text-gray-500 overflow-visible whitespace-nowrap">
              {{ monthLabels.find(m => m.week === wi)?.label ?? '' }}
            </div>
          </template>
        </div>

        <!-- 格子主體 -->
        <div class="flex gap-[4px]">
          <!-- 星期標籤 -->
          <div class="flex flex-col gap-[4px] mr-1 w-6">
            <div v-for="(label, li) in DAY_LABELS" :key="li"
              class="h-[13px] text-[9px] font-mono text-gray-600 leading-none flex items-center justify-end pr-0.5">
              {{ label }}
            </div>
          </div>

          <!-- 週格子 -->
          <div v-for="(week, wi) in weeks" :key="wi" class="flex flex-col gap-[4px]">
            <div v-for="(day, di) in week" :key="day.date"
              class="w-[13px] h-[13px] rounded-[2px] cursor-default transition-all duration-100 hover:scale-125 hover:z-10"
              :style="cellStyle(day.maps)"
              @mouseenter="showTooltip($event, day, di)"
              @mouseleave="hideTooltip">
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 圖例 -->
    <div class="flex items-center gap-2 mt-3 text-[9px] font-mono text-gray-600">
      <span>LESS</span>
      <div v-for="v in [0, 0.25, 0.5, 0.75, 1]" :key="v"
        class="w-[13px] h-[13px] rounded-[2px]"
        :style="cellStyle(v === 0 ? 0 : v * maxMaps)">
      </div>
      <span>MORE</span>
    </div>
  </div>
</template>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.1s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
