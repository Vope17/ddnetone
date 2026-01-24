<script setup>
import { Line } from 'vue-chartjs';
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

// Chart Options 可以放在這裡或獨立出去

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false },
    tooltip: {
      backgroundColor: 'rgba(0, 0, 0, 0.9)',
      titleColor: '#22d3ee',
      titleFont: { family: "'JetBrains Mono', monospace" },
      bodyColor: '#fff',
      bodyFont: { family: "'JetBrains Mono', monospace" },
      borderColor: '#334155',
      borderWidth: 1,
      padding: 10,
      displayColors: false,
      callbacks: { label: (c) => `> OUTPUT: ${c.parsed.y}` }
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
  <div class="flex flex-col h-full gap-6">
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 flex-shrink-0">
      <div class="bg-gradient-to-br from-gray-900 to-black border border-white/10 p-5 relative overflow-hidden group">
        <div class="absolute right-0 top-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
          <SparklesIcon class="w-16 h-16 text-cyan-500" />
        </div>
        <div class="text-xs font-mono text-cyan-500/70 mb-1 flex items-center gap-2">
          <span class="w-1 h-1 bg-cyan-500"></span> CURRENT_SCORE
        </div>
        <div class="text-3xl font-bold text-white font-mono tracking-tight">{{ summary.current_score.toLocaleString() }}</div>
        <div class="mt-2 w-full bg-gray-800 h-1 rounded-full overflow-hidden">
          <div class="h-full bg-cyan-500 w-[60%]"></div>
        </div>
      </div>


      <div class="bg-gradient-to-br from-gray-900 to-black border border-white/10 p-5 relative overflow-hidden group">
        <div class="absolute right-0 top-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
          <MapIcon class="w-16 h-16 text-purple-500" />
        </div>

        <div class="text-xs font-mono text-purple-500/70 mb-1 flex items-center gap-2">
          <span class="w-1 h-1 bg-purple-500"></span> MAPS_CLEARED
        </div>

        <div class="text-3xl font-bold text-white font-mono tracking-tight">{{ summary.completed_maps }}</div>
        <div class="text-[10px] text-gray-500 font-mono mt-1">LATEST: ASCENT_X</div>

      </div>

      <div class="bg-gradient-to-br from-gray-900 to-black border border-white/10 p-5 relative flex flex-col justify-center">

        <div class="flex justify-between items-end mb-2">
          <div class="text-xs font-mono text-green-500/70">POINTS_PROGRESS</div>
          <div class="text-xl font-bold text-white font-mono">{{ progressPercent }}%</div>
        </div>

        <div class="relative w-full h-3 bg-gray-800/50 rounded-sm overflow-hidden border border-white/5">

          <div class="absolute inset-0 bg-[url('https://grainy-gradients.vercel.app/noise.svg')] opacity-20"></div>
          <div class="h-full bg-gradient-to-r from-cyan-600 via-green-500 to-emerald-400 transition-all duration-1000 relative" :style="{ width: progressPercent + '%' }">
             <div class="absolute right-0 top-0 bottom-0 w-[1px] bg-white shadow-[0_0_8px_#fff]"></div>
          </div>
        </div>
        <div class="flex justify-between mt-2 text-[10px] text-gray-600 font-mono">
          <span>0</span>
          <span>TARGET: {{ summary.target_score.toLocaleString() }}</span>
        </div>
      </div>
    </div>

    <div class="flex-1 min-h-0 grid grid-cols-1 lg:grid-cols-4 gap-6">
      <div class="lg:col-span-3 bg-[#0a0a0a] border border-white/10 relative flex flex-col">
         <div class="flex items-center justify-between p-4 border-b border-white/5 bg-white/[0.02]">
          <h3 class="font-mono text-sm text-cyan-400 flex items-center gap-2">
            <SignalIcon class="w-4 h-4" /> PERFORMANCE_GROWTH
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
                 <div :class="['font-mono text-lg font-bold italic leading-none', 
                  index === 0 ? 'text-yellow-400' : index === 1 ? 'text-gray-300' : index === 2 ? 'text-orange-400' : 'text-gray-700']">
                  {{ index + 1 }}
                 </div>

              </div>
              <div class="flex-1 min-w-0">
                <div class="font-bold text-sm text-gray-200 truncate group-hover:text-yellow-100 transition">{{ p.name }}</div>
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
