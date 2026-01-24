<script setup>
import { ref, computed } from 'vue';
import { 
  CpuChipIcon, 
  FlagIcon, 

  PlayCircleIcon, 

  StopCircleIcon,
  CalculatorIcon 
} from '@heroicons/vue/24/outline';
import { CATEGORIES } from '../constants/categories';

const props = defineProps({
  maps: Array
});

const currentTab = ref('Insane');
const statusFilter = ref('All'); // 新增: 'All', 'Completed', 'InProgress', 'Incomplete'

// 1. 先篩選出當前分類的地圖 (Base List)
const mapsInCurrentCategory = computed(() => {
  if (currentTab.value === 'All') return props.maps;
  
  const target = currentTab.value.toLowerCase();
  return props.maps.filter(m => {
    if (!m.difficulty) return false;
    const dbDiff = m.difficulty.toLowerCase();
    
    if (target === 'ddmax.easy' && dbDiff.includes('eazy')) return true;
    if (target === 'ddmax.nut' && (dbDiff.includes('ntr') || dbDiff.includes('nut'))) return true;
    return dbDiff === target;
  });
});

// 2. 計算該分類的統計數據 (基於 mapsInCurrentCategory)
const categoryStats = computed(() => {

  const list = mapsInCurrentCategory.value;
  const totalScore = list.reduce((acc, m) => acc + (m.points || 0), 0); // 假設後端有傳 points (地圖總分)，若無則需後端補充
  // 如果後端沒傳 points，暫時用 score 替代或需要後端補欄位，這裡假設 score 是獲得分數
  // 注意：通常地圖會有一個 "Max Points"，您可能需要檢查 maps API 是否有回傳地圖的滿分
  // 這裡假設 m.score 是獲得分數。如果 m.points 不存在，這行統計會不準。
  // **修正建議**：請確保後端 /api/maps 有回傳該地圖的滿分欄位 (例如 max_points)。
  // 這裡暫時演示：已獲得分數 / 總獲得分數 (這邏輯怪怪的，通常是 獲得/總分)
  
  const currentScore = list.reduce((acc, m) => acc + (m.score || 0), 0);
  // 暫時用: 總分 = 所有地圖如果都滿分 (這裡假設每張圖分數記錄在 score 若已完成，或需要額外欄位)
  // 為了演示，我們計算:
  const completedCount = list.filter(m => m.status === 2).length;
  const inProgressCount = list.filter(m => m.status === 1).length;
  const totalCount = list.length;
  
  return {
    currentScore,
    // 這裡顯示的是累積分數
    totalCount,
    completedCount,
    inProgressCount,
    completionRate: totalCount > 0 ? Math.round((completedCount / totalCount) * 100) : 0
  };
});

// 3. 最後根據狀態過濾 (顯示在列表用)
const finalFilteredMaps = computed(() => {
  let list = mapsInCurrentCategory.value;
  
  if (statusFilter.value === 'Completed') {
    return list.filter(m => m.status === 2);

  } else if (statusFilter.value === 'InProgress') {
    return list.filter(m => m.status === 1);
  } else if (statusFilter.value === 'Incomplete') {
    return list.filter(m => m.status === 0 || !m.status);
  }
  
  return list;
});

// 狀態顏色 Helper
const getStatusColor = (status) => {
  if (status === 2) return 'text-cyan-400 border-cyan-500/50 bg-cyan-500/10';
  if (status === 1) return 'text-yellow-400 border-yellow-500/50 bg-yellow-500/10';
  return 'text-gray-500 border-white/10 bg-black/40';
};
</script>

<template>
  <div class="flex flex-col h-full overflow-hidden gap-4">
    
    <div class="shrink-0 space-y-4">
      <div class="flex flex-wrap gap-2">
         <button v-for="tab in CATEGORIES" :key="tab" 
          @click="currentTab = tab"
          :class="['px-4 py-1.5 font-mono text-xs font-bold border transition-all duration-200 uppercase tracking-widest clip-path-slant', 
            currentTab === tab 
            ? 'bg-cyan-500 text-black border-cyan-400' 
            : 'bg-black/50 text-gray-500 border-gray-800 hover:border-gray-600 hover:text-gray-300']">
          {{ tab }}
        </button>
      </div>

      <div class="flex flex-col md:flex-row gap-4 items-stretch md:items-center bg-white/[0.02] border-y border-white/10 p-3">

        
        <div class="flex items-center gap-6 flex-1">
          <div class="flex items-center gap-3">
            <div class="p-2 bg-cyan-500/10 rounded border border-cyan-500/30">
              <CalculatorIcon class="w-5 h-5 text-cyan-400" />

            </div>
            <div>
              <div class="text-[10px] text-gray-500 font-mono uppercase tracking-widest">Score Collected</div>
              <div class="flex items-baseline gap-1">
                <span class="text-xl font-bold text-white font-mono">{{ categoryStats.currentScore }}</span>
                <span class="text-xs text-gray-600 font-mono">PTS</span>
              </div>
            </div>
          </div>

          <div class="w-px h-8 bg-white/10"></div>

          <div class="flex items-center gap-3">
             <div class="p-2 bg-purple-500/10 rounded border border-purple-500/30">
              <FlagIcon class="w-5 h-5 text-purple-400" />
            </div>
            <div>
              <div class="text-[10px] text-gray-500 font-mono uppercase tracking-widest">Completion</div>
              <div class="flex items-baseline gap-1">
                <span class="text-xl font-bold text-white font-mono">{{ categoryStats.completedCount }}</span>
                <span class="text-xs text-gray-600 font-mono">/ {{ categoryStats.totalCount }}</span>
                <span class="ml-1 text-xs text-purple-400 font-mono">({{ categoryStats.completionRate }}%)</span>
              </div>
            </div>
          </div>
        </div>

        <div class="flex bg-black/50 rounded-lg p-1 border border-white/10">

          <button @click="statusFilter = 'All'"
            :class="['px-3 py-1 text-[10px] font-mono rounded transition-all', statusFilter === 'All' ? 'bg-gray-700 text-white shadow' : 'text-gray-500 hover:text-gray-300']">
            ALL
          </button>
          <button @click="statusFilter = 'Completed'"
            :class="['px-3 py-1 text-[10px] font-mono rounded transition-all flex items-center gap-1', statusFilter === 'Completed' ? 'bg-cyan-900/50 text-cyan-300 shadow border border-cyan-500/30' : 'text-gray-500 hover:text-cyan-400']">
            <div class="w-1.5 h-1.5 rounded-full bg-cyan-400"></div> DONE
          </button>
          <button @click="statusFilter = 'InProgress'"
            :class="['px-3 py-1 text-[10px] font-mono rounded transition-all flex items-center gap-1', statusFilter === 'InProgress' ? 'bg-yellow-900/50 text-yellow-300 shadow border border-yellow-500/30' : 'text-gray-500 hover:text-yellow-400']">
            <div class="w-1.5 h-1.5 rounded-full bg-yellow-400 animate-pulse"></div> WIP
          </button>
          <button @click="statusFilter = 'Incomplete'"
            :class="['px-3 py-1 text-[10px] font-mono rounded transition-all flex items-center gap-1', statusFilter === 'Incomplete' ? 'bg-gray-800 text-gray-300 shadow border border-gray-600' : 'text-gray-500 hover:text-gray-300']">
            <div class="w-1.5 h-1.5 rounded-full bg-gray-600"></div> TODO
          </button>
        </div>

      </div>

    </div>

    <div class="flex-1 overflow-y-auto custom-scrollbar pb-20 px-1">
      <div v-if="finalFilteredMaps.length === 0" class="h-64 flex flex-col items-center justify-center text-gray-600">
         <CpuChipIcon class="w-12 h-12 mb-4 opacity-50" />
         <div class="font-mono text-sm">NO RECORDS FOUND</div>
      </div>


      <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-4">
        <div v-for="map in finalFilteredMaps" :key="map.id" 
          :class="['group relative border p-5 transition-all duration-300 hover:translate-y-[-2px]', 
            getStatusColor(map.status)]">

          
          <div class="absolute top-3 right-3">
             <span v-if="map.status === 2" class="flex items-center gap-1 text-[9px] border border-cyan-500/30 bg-cyan-500/10 px-1.5 py-0.5 rounded text-cyan-300 font-mono">
               <FlagIcon class="w-3 h-3" /> CLEARED
             </span>
             <span v-else-if="map.status === 1" class="flex items-center gap-1 text-[9px] border border-yellow-500/30 bg-yellow-500/10 px-1.5 py-0.5 rounded text-yellow-300 font-mono">
               <PlayCircleIcon class="w-3 h-3 animate-pulse" /> RUNNING
             </span>
             <span v-else class="flex items-center gap-1 text-[9px] border border-gray-700 bg-gray-800 px-1.5 py-0.5 rounded text-gray-500 font-mono">
               <StopCircleIcon class="w-3 h-3" /> PENDING
             </span>
          </div>

          <div class="flex justify-between items-start mb-6 mt-2">
            <div class="p-2 bg-white/5 rounded border border-white/5 group-hover:border-white/20 transition-colors">
              <CpuChipIcon class="w-6 h-6 opacity-70" />
            </div>
            <div v-if="map.status === 2" class="text-right mt-8">
                <span class="font-mono text-2xl font-bold tracking-tight drop-shadow-[0_0_5px_rgba(6,182,212,0.5)]">{{ map.score }}</span>
                <span class="text-xs ml-1 opacity-70">PTS</span>
            </div>
             <div v-else class="text-right mt-8 opacity-30">

                <span class="font-mono text-xl font-bold">---</span>
            </div>
          </div>
          
          <div class="space-y-1">
            <h4 :class="['font-bold text-lg truncate', map.status === 2 ? 'text-white' : 'text-gray-400']" :title="map.map_name">
              {{ map.map_name }}
            </h4>
            <div class="flex items-center gap-2">
              <span class="text-[10px] font-mono px-1.5 py-0.5 border border-white/10 rounded text-gray-500">
                {{ map.difficulty }}
              </span>
              <span v-if="map.runner && map.runner !== '-'" class="text-[10px] text-gray-400 font-mono">
                BY {{ map.runner }}
              </span>
            </div>
          </div>

          <div v-if="map.note" class="mt-4 pt-3 border-t border-dashed border-white/10">
            <p class="text-xs opacity-60 italic line-clamp-2">"{{ map.note }}"</p>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<style scoped>
/* 斜角按鈕效果 */
.clip-path-slant {
  clip-path: polygon(10px 0, 100% 0, 100% calc(100% - 10px), calc(100% - 10px) 100%, 0 100%, 0 10px);
}
</style>
