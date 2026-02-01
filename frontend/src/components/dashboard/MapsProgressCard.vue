<script setup>
import { computed } from 'vue';

const props = defineProps({
  completedMaps: {
    type: Number,
    default: 0
  },
  growthData: Array // 如果你需要用歷史資料，記得傳入這個
});

const TARGET_MAPS = 2403;

// 設定地圖里程碑
const manualMapLabels = [
  { count: 0, text: '2024/09/30 19:30 開始傳奇的一刻' },
  { count: 100, text: '2024/10/01' },
  { count: 300, text: '2024/10/06' },
  { count: 400, text: '2026/01/24' },
];

// 計算百分比
const mapsPercent = computed(() => {
  return Math.min((props.completedMaps / TARGET_MAPS) * 100, 100).toFixed(1);
});

// 計算位置 (這裡你可以隨時換成我們之前討論過的更強大的混合邏輯)
const mapMilestones = computed(() => {
  return manualMapLabels.map(item => ({
    ...item,
    left: (item.count / TARGET_MAPS) * 100 + '%'
  }));
});
</script>

<template>
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
          <div
            class="w-2 h-2 rounded-full bg-white/50 mt-0.5 transition-all duration-200 group-hover:bg-white group-hover:scale-125 group-hover:shadow-[0_0_10px_rgba(255,255,255,1)]">
          </div>

          <div
            class="absolute top-7 flex flex-col items-center opacity-0 group-hover:opacity-100 transform translate-y-1 group-hover:translate-y-0 transition-all duration-200 ease-out">
            <div
              class="w-0 h-0 border-l-[4px] border-l-transparent border-r-[4px] border-r-transparent border-b-[4px] border-b-gray-900">
            </div>
            <div
              class="bg-gray-900 border border-white/20 px-2 py-1.5 rounded shadow-xl flex flex-col items-center min-w-[60px]">
              <span class="text-[10px] text-white font-bold font-mono leading-none">{{ m.count }}</span>
              <span
                class="text-[8px] text-violet-300 font-mono mt-0.5 whitespace-nowrap border-t border-white/10 pt-0.5 w-full text-center">{{
                m.text }}</span>

            </div>
          </div>
        </div>
      </div>

      <div class="flex justify-between mt-8 text-[10px] text-gray-600 font-mono">
        <span>0</span>

        <span class="text-gray-500">PROGRESS: {{ completedMaps }} / {{ TARGET_MAPS }}</span>
      </div>
    </div>
  </div>
</template>
