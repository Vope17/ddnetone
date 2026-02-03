<script setup>
import { computed } from 'vue';

const props = defineProps({
  progressPercent: [String, Number],
  targetScore: {
    type: Number,
    default: 10000
  },
  growthData: Array // 保留接口給未來使用
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

// 2. 自動計算里程碑 (間隔 1000)
const autoMilestones = computed(() => {
  // 如果沒有歷史資料，就回傳空陣列
  if (!props.growthData || props.growthData.length === 0) return [];

  const results = [];
  const targetMax = props.targetScore || 10000;

  // 使用 for 迴圈，從 1000 開始，每次加 1000，直到目標分數
  for (let s = 5000; s <= targetMax; s += 1000) {

    // 【注意】這裡假設 growthData 裡的分數欄位叫 "score"
    // 如果你的 API 回傳欄位是 "points" 或 "total_score"，請將 d.score 改成對應名稱
    const match = props.growthData.find(d => d.points >= s);

    if (match) {
      // 日期格式化邏輯 (YYYY/MM/DD)
      const timeStr = match.timestamp
        ? new Date(match.timestamp).toLocaleString('sv-SE').replace(/-/g, '/') : 'UNKNOWN';

      results.push({
        score: match.points, // 顯示實際分數
        text: timeStr// 顯示格式化日期
      });
    }
  }

  return results;
});

// 3. 合併並計算位置
const milestones = computed(() => {
  const target = props.targetScore || 10000;

  // 合併手動與自動產生的里程碑
  const allMilestones = [...manualLabels, ...autoMilestones.value];

  return allMilestones.map(item => ({
    ...item,
    // 計算 CSS left 百分比
    left: (item.score / target) * 100 + '%'
  }));
});
</script>

<template>
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
              <span class="text-[10px] text-white font-bold font-mono leading-none">{{ m.score }}</span>
              <span
                class="text-[8px] text-cyan-400 font-mono mt-0.5 whitespace-nowrap border-t border-white/10 pt-0.5 w-full text-center">{{
                  m.text }}</span>
            </div>
          </div>
        </div>
      </div>


      <div class="flex justify-between mt-8 text-[10px] text-gray-600 font-mono">
        <span>0</span>
        <span>TARGET: {{ targetScore.toLocaleString() }}</span>
      </div>
    </div>
  </div>
</template>
