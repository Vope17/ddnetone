<script setup>
import { Line } from 'vue-chartjs';
import { computed } from 'vue';
import {
  SparklesIcon, MapIcon, SignalIcon, TrophyIcon
} from '@heroicons/vue/24/outline';
import {
  Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend, Filler
} from 'chart.js';

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend, Filler);

const props = defineProps({
  summary: Object,
  players: Array,
  progressPercent: String,
  chartData: Object
});
// 1. 設定地圖的總目標 (因為 summary 可能只有分數的 target，這裡手動設一個)
const TARGET_MAPS = 2403;

// 2. 設定地圖的里程碑 (手動輸入數值與文字)
const manualMapLabels = [
  { count: 0, text: '2024/09/30 19:30 開始傳奇的一刻' },
  { count: 100, text: '2024/10/01' },
  { count: 300, text: '2024/10/06' },
  { count: 400, text: '2026/01/24' },
];

// --- 計算邏輯區 ---

// 計算地圖進度百分比
const mapsPercent = computed(() => {
  const current = props.summary?.completed_maps || 0;
  return Math.min((current / TARGET_MAPS) * 100, 100).toFixed(1);
});

// 計算地圖里程碑的位置
const mapMilestones = computed(() => {
  return manualMapLabels.map(item => ({
    ...item,
    left: (item.count / TARGET_MAPS) * 100 + '%'
  }));
});

const manualLabels = [
  { score: 0, text: '2024/09/30 19:30 開始傳奇的一刻' },
  { score: 1000, text: '2024/10/01 通不到6小時破千 666' },
  { score: 2000, text: '2024/10/02 01:32' },
  { score: 2500, text: '2024/10/02 23:25' },
  { score: 3975, text: '2026/01/23 22:20 假裝破4000分 fake news' },
  { score: 4000, text: '2026/01/24 00:20' },
  { score: 5000, text: '2026/01/31 21:00' },
];

const milestones = computed(() => {
  // 防呆：如果還沒有目標分數，預設 10000 避免除以 0
  const target = props.summary?.target_score || 10000;

  return manualLabels.map(item => ({
    ...item,
    // 計算百分比位置
    left: (item.score / target) * 100 + '%'
  }));
});
// Chart Options 可以放在這裡或獨立出去

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false },
    tooltip: {
      backgroundColor: 'rgba(0, 0, 0, 0.95)', //稍微加深背景
      titleColor: '#22d3ee', // Cyan 標題 (地圖名)
      titleFont: { family: "'JetBrains Mono', monospace", size: 13 },
      bodyColor: '#e2e8f0', // 淺灰內容
      bodyFont: { family: "'JetBrains Mono', monospace", size: 12 },
      borderColor: '#334155',
      borderWidth: 1,
      padding: 12,
      displayColors: false, // 隱藏前面的小色塊

      callbacks: {
        // 1. 標題：改為顯示總分 (原本的 context.parsed.y 在這裡要從 tooltipItems 拿)
        title: (tooltipItems) => {
          const score = tooltipItems[0].parsed.y;
          return `> POINTS: ${score}`;
        },

        // 2. 內容：顯示玩家、地圖、以及地圖滿分

        label: (context) => {
          const index = context.dataIndex;
          const dataset = context.dataset;

          const data = dataset.sourceData ? dataset.sourceData[index] : null;

          if (data) {
            // 回傳陣列會自動換行
            return [
              `> RUNNER : ${data.runner || 'UNKNOWN'}`,
              `> MAP    : ${data.map_name || 'UNKNOWN'}`,
              `> MAP_PTS: ${data.map_points || 0}`,

            ];
          } else {

            // 備用方案，若無資料則顯示基本的點數資訊
            return `> DATA: SYSTEM_SYNC`;
          }
        }
      }
    }
  },
  scales: {
    y: {
      grid: { color: 'rgba(255, 255, 255, 0.05)', tickLength: 0 },
      ticks: { color: '#475569', font: { family: "'JetBrains Mono', monospace", size: 10 } },
      border: { display: false }
    },

    x: { grid: { display: false }, ticks: { color: '#475569', font: { family: "'JetBrains Mono', monospace", size: 10 } } }
  }
};
</script>

<template>
  <div class="flex flex-col h-full w-full gap-6">
    <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-3 md:gap-6 flex-shrink-0">
      <div class="bg-gradient-to-br from-gray-900 to-black border border-white/10 p-5 relative overflow-hidden group">
        <div class="absolute right-0 top-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
          <SparklesIcon class="w-16 h-16 text-cyan-500" />
        </div>
        <div class="text-xs font-mono text-cyan-500/70 mb-1 flex items-center gap-2">
          <span class="w-1 h-1 bg-cyan-500"></span> CURRENT_POINTS
        </div>
        <div class="text-3xl font-bold text-white font-mono tracking-tight">{{ summary.current_score.toLocaleString() }}
        </div>
        <div class="mt-2 w-full bg-gray-800 h-1 rounded-full overflow-hidden">
          <div class="h-full bg-cyan-500 w-[60%]"></div>
        </div>
      </div>


      <div
        class="bg-gradient-to-br from-gray-900 to-black border border-white/10 p-5 relative flex flex-col justify-center">

        <div class="flex justify-between items-end mb-2">
          <div class="text-xs font-mono text-violet-400/70">MAPS_CLEARED</div>
          <div class="text-xl font-bold text-white font-mono">{{ mapsPercent }}%</div>
        </div>

        <div class="relative w-full pt-2">

          <div class="relative w-full h-3 bg-gray-800/50 rounded-sm overflow-hidden border border-white/5">
            <div class="absolute inset-0 bg-[url('https://grainy-gradients.vercel.app/noise.svg')] opacity-20"></div>

            <div
              class="h-full bg-gradient-to-r from-blue-600 via-violet-500 to-fuchsia-400 transition-all duration-1000 relative"
              :style="{ width: mapsPercent + '%' }">
              <div class="absolute right-0 top-0 bottom-0 w-[1px] bg-white shadow-[0_0_8px_#fff]"></div>
            </div>
          </div>

          <div class="absolute inset-0 pointer-events-none z-10">
            <div v-for="m in mapMilestones" :key="m.count"
              class="absolute top-0 h-full flex flex-col items-center group pointer-events-auto hover:z-30"
              :style="{ left: m.left, transform: 'translateX(-50%)' }">

              <div class="absolute -inset-y-2 w-6 bg-transparent cursor-help"></div>

              <div class="w-2 h-2 rounded-full bg-white/50 mt-0.5 transition-all duration-200 
                    group-hover:bg-white group-hover:scale-125 group-hover:shadow-[0_0_10px_rgba(255,255,255,1)]">
              </div>


              <div
                class="absolute top-7 flex flex-col items-center opacity-0 group-hover:opacity-100 transform translate-y-1 group-hover:translate-y-0 transition-all duration-200 ease-out">
                <div
                  class="w-0 h-0 border-l-[4px] border-l-transparent border-r-[4px] border-r-transparent border-b-[4px] border-b-gray-900">
                </div>


                <div
                  class="bg-gray-900 border border-white/20 px-2 py-1.5 rounded shadow-xl flex flex-col items-center min-w-[60px]">
                  <span class="text-[10px] text-white font-bold font-mono leading-none">

                    {{ m.count }}
                  </span>
                  <span
                    class="text-[8px] text-violet-300 font-mono mt-0.5 whitespace-nowrap border-t border-white/10 pt-0.5 w-full text-center">
                    {{ m.text }}
                  </span>
                </div>
              </div>

            </div>
          </div>

          <div class="flex justify-between mt-8 text-[10px] text-gray-600 font-mono">
            <span>0</span>
            <span class="text-gray-500">
              PROGRESS:{{ summary.completed_maps }} / {{ TARGET_MAPS }}
            </span>
          </div>

        </div>
      </div>

      <div
        class="bg-gradient-to-br from-gray-900 to-black border border-white/10 p-5 relative flex flex-col justify-center">
        <div class="flex justify-between items-end mb-2">
          <div class="text-xs font-mono text-green-500/70">POINTS_PROGRESS</div>
          <div class="text-xl font-bold text-white font-mono">{{ progressPercent }}%</div>
        </div>


        <div class="relative w-full pt-2">

          <div class="relative w-full h-3 bg-gray-800/50 rounded-sm overflow-hidden border border-white/5">
            <div class="absolute inset-0 bg-[url('https://grainy-gradients.vercel.app/noise.svg')] opacity-20"></div>
            <div
              class="h-full bg-gradient-to-r from-cyan-600 via-green-500 to-emerald-400 transition-all duration-1000 relative"
              :style="{ width: progressPercent + '%' }">
              <div class="absolute right-0 top-0 bottom-0 w-[1px] bg-white shadow-[0_0_8px_#fff]"></div>
            </div>
          </div>

          <div class="absolute inset-0 pointer-events-none z-10">
            <div v-for="m in milestones" :key="m.score"
              class="absolute top-0 h-full flex flex-col items-center group pointer-events-auto hover:z-30"
              :style="{ left: m.left, transform: 'translateX(-50%)' }">

              <div class="absolute -inset-y-2 w-6 bg-transparent cursor-help"></div>

              <div
                class="w-px h-[24px] bg-white/40 mt-0 transition-colors duration-200 group-hover:bg-white group-hover:shadow-[0_0_5px_rgba(255,255,255,0.8)]">
              </div>

              <div
                class="absolute top-7 flex flex-col items-center opacity-0 group-hover:opacity-100 transform translate-y-1 group-hover:translate-y-0 transition-all duration-200 ease-out">

                <div
                  class="w-0 h-0 border-l-[4px] border-l-transparent border-r-[4px] border-r-transparent border-b-[4px] border-b-gray-900">
                </div>

                <div
                  class="bg-gray-900 border border-white/20 px-2 py-1.5 rounded shadow-xl flex flex-col items-center min-w-[60px]">
                  <span class="text-[10px] text-white font-bold font-mono leading-none">
                    {{ m.score }}
                  </span>
                  <span
                    class="text-[8px] text-cyan-400 font-mono mt-0.5 whitespace-nowrap border-t border-white/10 pt-0.5 w-full text-center">
                    {{ m.text }}
                  </span>
                </div>
              </div>

            </div>
          </div>

          <div class="flex justify-between mt-8 text-[10px] text-gray-600 font-mono">
            <span>0</span>
            <span>TARGET: {{ summary.target_score?.toLocaleString() }}</span>
          </div>

        </div>
      </div>
    </div>

    <div class="flex-1 min-h-0 grid grid-cols-1 lg:grid-cols-4 gap-4 md:gap-6 overflow-y-auto lg:overflow-visible">
      <div class="lg:col-span-3 bg-[#0a0a0a] border border-white/10 relative flex flex-col">
        <div class="flex items-center justify-between p-4 border-b border-white/5 bg-white/[0.02]">
          <h3 class="font-mono text-sm text-cyan-400 flex items-center gap-2">
            <SignalIcon class="w-4 h-4" /> POINTS_GROWTH
          </h3>
          <div class="flex gap-2">
            <span class="w-2 h-2 bg-cyan-500 rounded-full animate-pulse"></span>
            <span class="text-[10px] font-mono text-gray-500">LIVE FEED</span>
          </div>
        </div>
        <div class="flex-1 w-full min-h-0 p-4 relative">
          <Line :data="chartData" :options="chartOptions" />
        </div>
      </div>

      <div class="lg:col-span-1 bg-[#0a0a0a] border border-yellow-500/20 relative flex flex-col h-full overflow-hidden">
        <div class="p-4 border-b border-white/5 bg-yellow-500/5 flex items-center justify-between shrink-0">
          <h3 class="font-mono text-sm text-yellow-500 flex items-center gap-2">
            <TrophyIcon class="w-4 h-4" /> TOP_PLAYERS
          </h3>
          <span class="text-[10px] border border-yellow-900 text-yellow-700 px-1 rounded">S-TIER</span>
        </div>
        <div class="flex-1 overflow-y-auto p-2 space-y-2 custom-scrollbar">
          <div v-for="(p, index) in players" :key="p.id"
            class="group relative p-3 border border-white/5 hover:border-yellow-500/30 bg-white/[0.02] hover:bg-yellow-500/[0.05] transition-all duration-200">
            <div class="flex items-center gap-3 relative z-10">

              <div class="flex flex-col items-center justify-center w-8">
                <div
                  :class="['font-mono text-lg font-bold italic leading-none',
                    index === 0 ? 'text-yellow-400' : index === 1 ? 'text-gray-300' : index === 2 ? 'text-orange-400' : 'text-gray-700']">
                  {{ index + 1 }}
                </div>

              </div>
              <div class="flex-1 min-w-0">
                <div class="font-bold text-sm text-gray-200 truncate group-hover:text-yellow-100 transition">{{ p.name
                  }}</div>
                <div class="flex items-center gap-2 text-[10px] text-gray-500 font-mono">
                  <span class="uppercase tracking-wider text-gray-600">{{ p.role }}</span>
                </div>
              </div>
              <div class="text-right">
                <div class="font-mono font-bold text-cyan-400 text-sm">{{ p.score_contrib }}</div>

                <div class="text-[9px] text-gray-600">{{ p.map_count }} MAPS</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
