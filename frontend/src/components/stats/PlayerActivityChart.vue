<script setup>
import { computed, ref } from 'vue';

const props = defineProps({
  maps: Array
});

const DAYS = 30;
const tooltip = ref(null);

const toTaipeiDate = (d) =>
  d.toLocaleDateString('en-CA', { timeZone: 'Asia/Taipei' }); // YYYY-MM-DD

// 建立過去 30 天的日期列表（Taipei 時區）
const dates = computed(() => {
  const result = [];
  const todayStr = toTaipeiDate(new Date());
  const today = new Date(todayStr + 'T00:00:00');
  for (let i = DAYS - 1; i >= 0; i--) {
    const d = new Date(today);
    d.setDate(today.getDate() - i);
    result.push(toTaipeiDate(d));
  }
  return result;
});

const playerRows = computed(() => {
  const completed = (props.maps || []).filter(m => m.status === 2 && m.runner && m.finish_time);
  const dateSet = new Set(dates.value);

  const grid = {};
  completed.forEach(m => {
    const d = toTaipeiDate(new Date(m.finish_time)); // UTC → Taipei date
    if (!dateSet.has(d)) return;
    const runners = m.runner.split(/[,&]/).map(r => r.trim()).filter(Boolean);
    runners.forEach(r => {
      if (!grid[r]) grid[r] = {};
      grid[r][d] = (grid[r][d] || 0) + 1;
    });
  });

  return Object.entries(grid)
    .map(([name, dayMap]) => ({
      name,
      total: Object.values(dayMap).reduce((a, b) => a + b, 0),
      days: dates.value.map(d => ({ date: d, count: dayMap[d] || 0 })),
    }))
    .sort((a, b) => b.total - a.total);
});

const maxCount = computed(() =>
  Math.max(1, ...playerRows.value.flatMap(r => r.days.map(d => d.count)))
);

function cellStyle(count) {
  if (count === 0) return { background: 'rgba(255,255,255,0.04)', border: '1px solid rgba(255,255,255,0.03)' };
  const ratio = count / maxCount.value;
  const alpha = 0.2 + ratio * 0.8;
  return {
    background: `rgba(251,191,36,${alpha})`,
    border: `1px solid rgba(251,191,36,${Math.min(1, alpha + 0.15)})`,
    boxShadow: ratio > 0.5 ? `0 0 4px rgba(251,191,36,${ratio * 0.5})` : 'none',
  };
}

const dateLabels = computed(() =>
  dates.value.map((d, i) => ({
    d,
    label: i % 7 === 0 ? new Date(d + 'T00:00:00').toLocaleString('en', { month: 'short', day: 'numeric' }) : '',
  }))
);

function showTooltip(e, playerName, day, rowIndex) {
  if (day.count === 0) { tooltip.value = null; return; }
  const rect = e.target.getBoundingClientRect();
  const container = e.target.closest('.pa-container').getBoundingClientRect();
  const showBelow = rowIndex < 2;
  tooltip.value = {
    name: playerName,
    date: day.date,
    count: day.count,
    top: showBelow
      ? rect.bottom - container.top + 4
      : rect.top - container.top - 58,
    left: Math.min(rect.left - container.left + 6, container.width - 130),
  };
}
function hideTooltip() { tooltip.value = null; }
</script>

<template>
  <div class="bg-gradient-to-br from-gray-900 to-black border border-white/10 p-5 h-full flex flex-col">
    <div class="text-xs font-mono text-amber-400/70 mb-4 flex items-center gap-2 flex-shrink-0">
      <span class="w-1.5 h-1.5 rounded-full bg-amber-400 animate-pulse"></span>
      PLAYER_ACTIVITY
      <span class="text-gray-600 ml-1">— LAST 30 DAYS</span>
    </div>

    <div v-if="!playerRows.length" class="text-gray-600 font-mono text-xs text-center py-6">NO DATA</div>

    <div v-else class="pa-container relative overflow-x-auto overflow-y-auto max-h-96 custom-scrollbar">
      <!-- Tooltip -->
      <Transition name="fade">
        <div v-if="tooltip" class="absolute z-50 pointer-events-none"
          :style="{ top: tooltip.top + 'px', left: tooltip.left + 'px' }">
          <div class="bg-gray-950 border border-amber-500/30 px-2.5 py-1.5 rounded-md text-[10px] font-mono whitespace-nowrap shadow-xl">
            <div class="text-amber-300 font-bold">{{ tooltip.name }}</div>
            <div class="text-gray-400">{{ tooltip.date }}</div>
            <div class="text-amber-200">{{ tooltip.count }} maps</div>
          </div>
        </div>
      </Transition>

      <div class="min-w-max">
        <!-- 日期標籤 -->
        <div class="flex mb-1 ml-28">
          <div v-for="({ d, label }) in dateLabels" :key="d"
            class="w-[14px] text-[8px] font-mono text-gray-600 overflow-visible whitespace-nowrap">
            {{ label }}
          </div>
        </div>

        <!-- 每個玩家一行 -->
        <div v-for="(row, ri) in playerRows" :key="row.name" class="flex items-center gap-1 mb-1">
          <!-- 玩家名 -->
          <div class="w-24 flex-shrink-0 text-right pr-2 text-[10px] font-mono text-gray-400 truncate hover:text-amber-300 transition-colors">
            {{ row.name }}
          </div>
          <!-- 格子 -->
          <div v-for="day in row.days" :key="day.date"
            class="w-[14px] h-[14px] rounded-[2px] flex-shrink-0 cursor-default transition-all duration-100 hover:scale-125"
            :style="cellStyle(day.count)"
            @mouseenter="showTooltip($event, row.name, day, ri)"
            @mouseleave="hideTooltip">
          </div>
          <!-- 總數 -->
          <div class="ml-1.5 text-[10px] font-mono text-gray-600">{{ row.total }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.1s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
