<script setup>
import { ref, computed } from 'vue';
import {
  CpuChipIcon,
  FlagIcon,
  CalculatorIcon,
  ChevronDownIcon,
  MagnifyingGlassIcon,
  StarIcon as StarIconOutline,
  ArrowsUpDownIcon,
  CalendarIcon,
  UserGroupIcon
} from '@heroicons/vue/24/outline';
import { StarIcon as StarIconSolid } from '@heroicons/vue/24/solid';
import { CATEGORIES } from '../constants/categories';

const props = defineProps({
  maps: Array
});

// --- State ---
const currentTab = ref('All');
const statusFilter = ref('Completed');
const searchQuery = ref('');
const sortType = ref('DEFAULT'); // DEFAULT (Score), STARS, POINTS
const sortOrder = ref('DESC');   // DESC, ASC
const isCategoryMenuOpen = ref(false);

// --- 1. 基礎分類篩選 ---
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

// --- 2. 統計數據 ---
const categoryStats = computed(() => {
  const list = mapsInCurrentCategory.value;
  const currentScore = list.reduce((acc, m) => acc + (m.status === 2 ? (m.points || 0) : 0), 0);
  const totalScore = list.reduce((acc, m) => acc + (m.points || 0), 0);
  const completedCount = list.filter(m => m.status === 2).length;
  const totalCount = list.length;

  return {
    currentScore,
    totalScore,
    totalCount,
    completedCount,

    completionRate: totalCount > 0 ? Math.round((completedCount / totalCount) * 100) : 0
  };
});

// --- 3. 進階篩選與排序邏輯 (核心修復) ---
const finalFilteredMaps = computed(() => {
  // A. 取得當前分類的地圖

  let list = mapsInCurrentCategory.value;

  // B. 狀態篩選
  if (statusFilter.value === 'Completed') list = list.filter(m => m.status === 2);
  else if (statusFilter.value === 'InProgress') list = list.filter(m => m.status === 1);
  else if (statusFilter.value === 'Incomplete') list = list.filter(m => m.status === 0 || !m.status);

  // C. 搜尋篩選
  if (searchQuery.value.trim()) {

    const query = searchQuery.value.toLowerCase();
    list = list.filter(m =>
      m.map_name.toLowerCase().includes(query) ||
      (m.runner && m.runner.toLowerCase().includes(query)) ||
      (m.note && m.note.toLowerCase().includes(query))
    );
  }

  // D. 排序邏輯 ★★★ 修正重點 ★★★
  // 使用 [...list] 建立複本，避免 .sort() 原地修改導致 Vue 響應式失效或資料錯亂
  return [...list].sort((a, b) => {
    let valA = 0, valB = 0;

    if (sortType.value === 'STARS') {
      valA = a.stars || 0;
      valB = b.stars || 0;
    } else if (sortType.value === 'POINTS') {
      valA = a.points || 0;
      valB = b.points || 0;
    } else if (sortType.value === 'DATE') {
      // ★ 新增：時間排序邏輯
      // 將時間字串轉為 Timestamp (毫秒) 進行比較
      // 如果沒有時間，視為 0 (排在最後)
      valA = a.finish_time ? new Date(a.finish_time).getTime() : 0;
      valB = b.finish_time ? new Date(b.finish_time).getTime() : 0;
    } else {
      // DEFAULT: 預設依據「獲得分數」排序
      // 如果是 TODO 列表，獲得分數通常都是 0，這時次要排序可以用 Points 或 MapName
      valA = a.score || 0;
      valB = b.score || 0;

      // 如果主要分數一樣 (例如都是0)，加入次要排序讓列表不混亂 (這裡用 Points)
      if (valA === valB) {
        valA = a.points || 0;
        valB = b.points || 0;
      }
    }

    if (sortOrder.value === 'DESC') return valB - valA;

    else return valA - valB;
  });
});

// --- Actions ---
const selectCategory = (tab) => {
  currentTab.value = tab;
  isCategoryMenuOpen.value = false;
};

// 切換排序 (支援取消)
const toggleSort = (type) => {
  if (sortType.value === type) {
    sortType.value = 'DEFAULT'; // 取消選取，回到預設
  } else {
    sortType.value = type;

  }
};

const getStatusColor = (status) => {
  if (status === 2) return 'text-cyan-400 border-cyan-500/50 bg-cyan-500/10';
  if (status === 1) return 'text-yellow-400 border-yellow-500/50 bg-yellow-500/10';
  return 'text-gray-500 border-white/10 bg-black/40';
};

// ★ 新增：日期格式化函式
const formatDate = (dateStr) => {
  if (!dateStr) return '';
  const date = new Date(dateStr);
  // 檢查是否為有效日期，避免顯示 "Invalid Date"
  if (isNaN(date.getTime())) return '';

  return date.toLocaleDateString('zh-TW', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  });
};

</script>

<template>
  <div class="flex flex-col h-full overflow-hidden gap-4">
    <div class="shrink-0 space-y-3">
      <div class="md:hidden relative z-30">
        <button @click="isCategoryMenuOpen = !isCategoryMenuOpen"
          class="w-full flex items-center justify-between px-4 py-3 bg-cyan-900/20 border border-cyan-500/50 text-cyan-400 font-mono font-bold uppercase tracking-widest rounded">
          <span>{{ currentTab }}</span>
          <ChevronDownIcon :class="['w-5 h-5 transition', isCategoryMenuOpen ? 'rotate-180' : '']" />
        </button>
        <transition enter-active-class="transition duration-100 ease-out"
          enter-from-class="transform scale-95 opacity-0" enter-to-class="transform scale-100 opacity-100"
          leave-active-class="transition duration-75 ease-in" leave-from-class="transform scale-100 opacity-100"
          leave-to-class="transform scale-95 opacity-0">
          <div v-if="isCategoryMenuOpen"
            class="absolute top-full left-0 w-full mt-2 bg-[#0a0a0a] border border-gray-800 rounded shadow-xl max-h-60 overflow-y-auto custom-scrollbar p-1 grid grid-cols-2 gap-1">
            <button v-for="tab in CATEGORIES" :key="tab" @click="selectCategory(tab)" :class="['px-3 py-2 text-xs font-mono text-left uppercase hover:bg-white/5 rounded',
              currentTab === tab ? 'text-cyan-400 bg-cyan-900/20 border border-cyan-500/30' : 'text-gray-400']">
              {{ tab }}
            </button>
          </div>
        </transition>
      </div>

      <div class="hidden md:flex flex-wrap gap-2">
        <button v-for="tab in CATEGORIES" :key="tab" @click="currentTab = tab"
          :class="['px-4 py-1.5 font-mono text-xs font-bold border transition-all uppercase tracking-widest clip-path-slant',
            currentTab === tab ? 'bg-cyan-500 text-purple-400 border-cyan-400' : 'bg-black/50 text-gray-500 border-gray-800 hover:border-gray-600']">
          {{ tab }}
        </button>
      </div>

      <div class="bg-white/[0.02] border-y border-white/10 p-3 flex flex-col gap-4">
        <div class="flex flex-col md:flex-row items-center justify-between gap-4">
          <div class="flex items-center gap-6 w-full md:w-auto justify-between md:justify-start">
            <div class="flex items-center gap-3">
              <div class="p-1.5 bg-cyan-500/10 rounded border border-cyan-500/30 hidden sm:block">
                <CalculatorIcon class="w-4 h-4 text-cyan-400" />
              </div>
              <div>
                <div class="text-[9px] text-gray-500 font-mono uppercase">Score</div>
                <div class="text-sm font-bold text-white font-mono">{{ categoryStats.currentScore }} <span
                    class="text-gray-600">/ {{ categoryStats.totalScore }}</span></div>
              </div>
            </div>
            <div class="w-px h-6 bg-white/10 hidden md:block"></div>
            <div class="flex items-center gap-3">
              <div class="p-1.5 bg-purple-500/10 rounded border border-purple-500/30 hidden sm:block">
                <FlagIcon class="w-4 h-4 text-purple-400" />
              </div>
              <div>
                <div class="text-[9px] text-gray-500 font-mono uppercase">Done</div>
                <div class="text-sm font-bold text-white font-mono">{{ categoryStats.completedCount }} <span
                    class="text-gray-600">/ {{ categoryStats.totalCount }}</span></div>
              </div>
            </div>
          </div>

          <div class="flex bg-black/50 rounded p-1 border border-white/10 w-full md:w-auto overflow-x-auto">
            <button @click="statusFilter = 'All'"
              :class="['px-4 py-1.5 text-xs font-mono rounded whitespace-nowrap transition-colors', statusFilter === 'All' ? 'bg-gray-700 text-white shadow-sm' : 'text-gray-500 hover:text-gray-300']">ALL</button>
            <button @click="statusFilter = 'Completed'"
              :class="['px-4 py-1.5 text-xs font-mono rounded whitespace-nowrap transition-colors', statusFilter === 'Completed' ? 'bg-cyan-900/50 text-cyan-300 shadow-sm border border-cyan-500/20' : 'text-gray-500 hover:text-cyan-400']">DONE</button>
            <button @click="statusFilter = 'InProgress'"
              :class="['px-4 py-1.5 text-xs font-mono rounded whitespace-nowrap transition-colors', statusFilter === 'InProgress' ? 'bg-yellow-900/50 text-yellow-300 shadow-sm border border-yellow-500/20' : 'text-gray-500 hover:text-yellow-400']">WIP</button>
            <button @click="statusFilter = 'Incomplete'"
              :class="['px-4 py-1.5 text-xs font-mono rounded whitespace-nowrap transition-colors', statusFilter === 'Incomplete' ? 'bg-gray-800 text-gray-300 shadow-sm border border-gray-600/50' : 'text-gray-500 hover:text-gray-400']">TODO</button>
          </div>
        </div>

        <div
          class="flex flex-col md:flex-row items-center justify-between gap-4 border-t border-white/5 pt-4 md:border-t-0 md:pt-0">
          <div class="relative w-full group flex-1 md:max-w-md">
            <MagnifyingGlassIcon
              class="w-4 h-4 text-gray-500 absolute left-3 top-2.5 group-focus-within:text-cyan-400 transition-colors" />
            <input v-model="searchQuery" type="text" placeholder="SEARCH MAP..."
              class="w-full bg-black/50 border border-white/10 rounded pl-9 pr-3 py-2 text-xs font-mono text-white focus:border-cyan-500/50 focus:outline-none transition-all placeholder:text-gray-700" />
          </div>

          <div class="flex items-center gap-2 w-full md:w-auto overflow-x-auto">
            <div class="flex bg-black/50 rounded border border-white/10 p-0.5">
              <button @click="toggleSort('STARS')"
                :class="['px-3 py-1.5 text-xs font-mono rounded flex items-center gap-1 transition-all', sortType === 'STARS' ? 'bg-yellow-500/20 text-yellow-300' : 'text-gray-500 hover:text-gray-300']">
                <StarIconOutline v-if="sortType !== 'STARS'" class="w-3 h-3" />
                <StarIconSolid v-else class="w-3 h-3" />
                STARS
              </button>
              <button @click="toggleSort('POINTS')"
                :class="['px-3 py-1.5 text-xs font-mono rounded transition-all', sortType === 'POINTS' ? 'bg-cyan-500/20 text-cyan-300' : 'text-gray-500 hover:text-gray-300']">PTS</button>
              <div class="w-px bg-white/10 mx-1"></div>
              <button @click="toggleSort('DATE')"
                :class="['px-3 py-1.5 text-xs font-mono rounded flex items-center gap-1 transition-all', sortType === 'DATE' ? 'bg-purple-500/20 text-purple-300' : 'text-gray-500 hover:text-gray-300']">
                <CalendarIcon class="w-3 h-3" />
                DATE
              </button>
              <button @click="sortOrder = sortOrder === 'ASC' ? 'DESC' : 'ASC'"
                class="px-3 py-1.5 text-gray-400 hover:text-white transition-colors">
                <ArrowsUpDownIcon class="w-3 h-3" />
              </button>
            </div>
          </div>
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
          :class="['group relative border p-4 transition-all duration-300 hover:translate-y-[-2px] flex flex-col', getStatusColor(map.status)]">

          <div class="flex justify-between items-start mb-2 pl-2">
            <div class="flex items-center gap-0.5" title="Map Difficulty Stars">
              <template v-for="n in 5" :key="n">
                <StarIconSolid v-if="n <= (map.stars || 0)" class="w-3 h-3 text-yellow-500" />
                <StarIconOutline v-else class="w-3 h-3 text-gray-700" />
              </template>
            </div>
            <div class="text-right">
              <span :class="['font-mono text-xl font-bold', map.status === 2 ? 'text-white' : 'text-gray-600']">{{
                map.points
              }}</span>
              <span class="text-[9px] text-gray-600 ml-1">PTS</span>
            </div>
          </div>
          <div v-if="map.has_dummy" class="flex items-center gap-1 mt-1 text-purple-400" title="Dummy Alive">
            <UserGroupIcon class="w-4 h-4" />
            <span class="text-[9px] font-mono border border-purple-500/30 px-1 rounded bg-purple-500/10">DUMMY</span>
          </div>

          <h4 :class="['font-bold text-md truncate mb-1', map.status === 2 ? 'text-white' : 'text-gray-400']">{{
            map.map_name }}
          </h4>

          <div class="flex flex-col gap-1 mb-3">
            <div class="flex items-center gap-2">
              <span class="text-[9px] font-mono border border-white/10 px-1 rounded text-gray-500">{{ map.difficulty
              }}</span>
              <span v-if="map.runner" class="text-[9px] font-mono text-gray-400">BY {{ map.runner }}</span>
            </div>
            <div v-if="map.finish_time" class="flex items-center gap-1.5 text-[9px] font-mono text-cyan-500/70 mt-1">
              <CalendarIcon class="w-3 h-3" />
              <span>{{ formatDate(map.finish_time) }}</span>
            </div>
          </div>

          <div v-if="map.note" class="mt-auto pt-2 border-t border-dashed border-white/10">
            <p class="text-[10px] text-gray-500 italic truncate">{{ map.note }}</p>
          </div>

        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.clip-path-slant {
  clip-path: polygon(10px 0, 100% 0, 100% calc(100% - 10px), calc(100% - 10px) 100%, 0 100%, 0 10px);
}
</style>
