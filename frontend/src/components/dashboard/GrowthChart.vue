<script setup>
import { Line } from 'vue-chartjs';
import { SignalIcon } from '@heroicons/vue/24/outline';
import {

  Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend, Filler
} from 'chart.js';

// 註冊 ChartJS 元件
ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend, Filler);

defineProps({
  chartData: Object
});

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false },
    tooltip: {

      backgroundColor: 'rgba(0, 0, 0, 0.95)',
      titleColor: '#22d3ee',
      titleFont: { family: "'JetBrains Mono', monospace", size: 13 },
      bodyColor: '#e2e8f0',
      bodyFont: { family: "'JetBrains Mono', monospace", size: 12 },
      borderColor: '#334155',
      borderWidth: 1,
      padding: 12,
      displayColors: false,
      callbacks: {
        title: (tooltipItems) => `> POINTS: ${tooltipItems[0].parsed.y}`,
        label: (context) => {

          const index = context.dataIndex;
          const dataset = context.dataset;
          const data = dataset.sourceData ? dataset.sourceData[index] : null;
          if (data) {
            return [
              `> RUNNER : ${data.runner || 'UNKNOWN'}`,
              `> MAP    : ${data.map_name || 'UNKNOWN'}`,
              `> MAP_PTS: ${data.map_points || 0}`,
            ];
          }
          return `> DATA: SYSTEM_SYNC`;
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
    x: {
      grid: { display: false },
      ticks: { color: '#475569', font: { family: "'JetBrains Mono', monospace", size: 10 } }
    }
  }
};
</script>

<template>
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
      <Line v-if="chartData" :data="chartData" :options="chartOptions" />
    </div>
  </div>
</template>
