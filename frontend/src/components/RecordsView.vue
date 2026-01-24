<script setup>
import { ref, computed } from 'vue';
import {
  CpuChipIcon,
  FlagIcon,
  PlayCircleIcon,
  StopCircleIcon,

  CalculatorIcon,
  ChevronDownIcon // ★ 新增：下拉選單箭頭圖示
} from '@heroicons/vue/24/outline';
import { CATEGORIES } from '../constants/categories';

const props = defineProps({
  maps: Array

});


const currentTab = ref('Insane');
const statusFilter = ref('All');
const isCategoryMenuOpen = ref(false); // ★ 新增：控制手機版分類選單開關

// 1. 先篩選出當前分類的地圖
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

// categoryStats 計算邏輯維持不變
const categoryStats = computed(() => {
  const list = mapsInCurrentCategory.value;
  const currentScore = list.reduce((acc, m) => acc + (m.status === 2 ? (m.points || 0) : 0), 0);
  const totalScore = list.reduce((acc, m) => acc + (m.points || 0), 0);
  const completedCount = list.filter(m => m.status === 2).length;
  const inProgressCount = list.filter(m => m.status === 1).length;
  const totalCount = list.length;

  return {
    currentScore,
    totalScore,
    totalCount,
    completedCount,
    inProgressCount,
    completionRate: totalCount > 0 ? Math.round((completedCount / totalCount) * 100) : 0
  };
});

// finalFilteredMaps 邏輯維持不變
const finalFilteredMaps = computed(() => {
  let list = mapsInCurrentCategory.value;
  if (statusFilter.value === 'Completed') return list.filter(m => m.status === 2);
  else if (statusFilter.value === 'InProgress') return list.filter(m => m.status === 1);
  else if (statusFilter.value === 'Incomplete') return list.filter(m => m.status === 0 || !m.status);
  return list;
});

// Helper: 選擇分類後自動關閉選單
const selectCategory = (tab) => {
  currentTab.value = tab;
  isCategoryMenuOpen.value = false;
};

const getStatusColor = (status) => {
  if (status === 2) return 'text-cyan-400 border-cyan-500/50 bg-cyan-500/10';
  if (status === 1) return 'text-yellow-400 border-yellow-500/50 bg-yellow-500/10';
  return 'text-gray-500 border-white/10 bg-black/40';
};
</script>

<template>
  <div class="flex flex-col h-full overflow-hidden gap-4">

    <div class="shrink-0 space-y-4">

      <div class="md:hidden relative z-30">
        <button @click="isCategoryMenuOpen = !isCategoryMenuOpen"
          class="w-full flex items-center justify-between px-4 py-3 bg-cyan-900/20 border border-cyan-500/50 text-cyan-400 font-mono font-bold uppercase tracking-widest rounded">
          <span>{{ currentTab }}</span>
          <ChevronDownIcon
            :class="['w-5 h-5 transition-transform duration-300', isCategoryMenuOpen ? 'rotate-180' : '']" />
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
        <button v-for="tab in CATEGORIES" :key="tab" @click="currentTab = tab" :class="['px-4 py-1.5 font-mono text-xs font-bold border transition-all duration-200 uppercase tracking-widest clip-path-slant',
          currentTab === tab
            ? 'bg-cyan-500 text-black border-cyan-400'

            : 'bg-black/50 text-gray-500 border-gray-800 hover:border-gray-600 hover:text-gray-300']">
          {{ tab }}
        </button>
      </div>

      <div
        class="flex flex-col md:flex-row gap-4 items-stretch md:items-center bg-white/[0.02] border-y border-white/10 p-3">
        <div class="flex items-center gap-6 flex-1 justify-center md:justify-start">
          <div class="flex items-center gap-3">
            <div class="p-2 bg-cyan-500/10 rounded border border-cyan-500/30 hidden sm:block">
              <CalculatorIcon class="w-5 h-5 text-cyan-400" />

            </div>
            <div>
              <div class="text-[10px] text-gray-500 font-mono uppercase tracking-widest">Score</div>
              <div class="flex items-baseline gap-1">
                <span class="text-lg md:text-xl font-bold text-white font-mono">{{ categoryStats.currentScore }}</span>
                <span class="text-xs text-gray-600 font-mono">/ {{ categoryStats.totalScore }}</span>
              </div>
            </div>
          </div>

          <div class="w-px h-8 bg-white/10"></div>


          <div class="flex items-center gap-3">
            <div class="p-2 bg-purple-500/10 rounded border border-purple-500/30 hidden sm:block">
              <FlagIcon class="w-5 h-5 text-purple-400" />
            </div>
            <div>

              <div class="text-[10px] text-gray-500 font-mono uppercase tracking-widest">Progress</div>
              <div class="flex items-baseline gap-1">
                <span class="text-lg md:text-xl font-bold text-white font-mono">{{ categoryStats.completedCount
                  }}</span>
                <span class="text-xs text-gray-600 font-mono">/ {{ categoryStats.totalCount }}</span>
                <span class="ml-1 text-xs text-purple-400 font-mono hidden sm:inline">({{ categoryStats.completionRate
                  }}%)</span>
              </div>

            </div>
          </div>
        </div>

        <div class="flex bg-black/50 rounded-lg p-1 border border-white/10 overflow-x-auto custom-scrollbar">
          <button @click="statusFilter = 'All'"
            :class="['px-3 py-1 text-[10px] font-mono rounded transition-all whitespace-nowrap', statusFilter === 'All' ? 'bg-gray-700 text-white shadow' : 'text-gray-500 hover:text-gray-300']">
            ALL
          </button>
          <button @click="statusFilter = 'Completed'"
            :class="['px-3 py-1 text-[10px] font-mono rounded transition-all flex items-center gap-1 whitespace-nowrap', statusFilter === 'Completed' ? 'bg-cyan-900/50 text-cyan-300 shadow border border-cyan-500/30' : 'text-gray-500 hover:text-cyan-400']">
            <div class="w-1.5 h-1.5 rounded-full bg-cyan-400"></div> DONE
          </button>
          <button @click="statusFilter = 'InProgress'"
            :class="['px-3 py-1 text-[10px] font-mono rounded transition-all flex items-center gap-1 whitespace-nowrap', statusFilter === 'InProgress' ? 'bg-yellow-900/50 text-yellow-300 shadow border border-yellow-500/30' : 'text-gray-500 hover:text-yellow-400']">
            <div class="w-1.5 h-1.5 rounded-full bg-yellow-400 animate-pulse"></div> WIP
          </button>
          <button @click="statusFilter = 'Incomplete'"
            :class="['px-3 py-1 text-[10px] font-mono rounded transition-all flex items-center gap-1 whitespace-nowrap', statusFilter === 'Incomplete' ? 'bg-gray-800 text-gray-300 shadow border border-gray-600' : 'text-gray-500 hover:text-gray-300']">
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
        <div v-for="map in finalFilteredMaps" :key="map.id" :class="['group relative border p-5 transition-all duration-300 hover:translate-y-[-2px]',
          getStatusColor(map.status)]">

          <div class="absolute top-3 right-3">
            <span v-if="map.status === 2"
              class="flex items-center gap-1 text-[9px] border border-cyan-500/30 bg-cyan-500/10 px-1.5 py-0.5 rounded text-cyan-300 font-mono">
              <FlagIcon class="w-3 h-3" /> CLEARED
            </span>
            <span v-else-if="map.status === 1"
              class="flex items-center gap-1 text-[9px] border border-yellow-500/30 bg-yellow-500/10 px-1.5 py-0.5 rounded text-yellow-300 font-mono">
              <PlayCircleIcon class="w-3 h-3 animate-pulse" /> RUNNING
            </span>
            <span v-else
              class="flex items-center gap-1 text-[9px] border border-gray-700 bg-gray-800 px-1.5 py-0.5 rounded text-gray-500 font-mono">
              <StopCircleIcon class="w-3 h-3" /> PENDING
            </span>
          </div>

          <div class="flex justify-between items-start mb-6 mt-2">
            <div class="p-2 bg-white/5 rounded border border-white/5 group-hover:border-white/20 transition-colors">
              <CpuChipIcon class="w-6 h-6 opacity-70" />
            </div>

            <div class="text-right mt-8">

              <span
                :class="['font-mono text-2xl font-bold tracking-tight', map.status === 2 ? 'text-white drop-shadow-[0_0_5px_rgba(6,182,212,0.5)]' : 'text-gray-600']">{{
                  map.points }}</span>
              <span
                :class="['text-xs ml-1', map.status === 2 ? 'opacity-70 text-cyan-200' : 'text-gray-700']">PTS</span>
            </div>
          </div>

          <div class="space-y-1">
            <h4 :class="['font-bold text-lg truncate', map.status === 2 ? 'text-white' : 'text-gray-400']"
              :title="map.map_name">{{ map.map_name }}</h4>
            <div class="flex items-center gap-2">
              <span class="text-[10px] font-mono px-1.5 py-0.5 border border-white/10 rounded text-gray-500">{{
                map.difficulty }}</span>
              <span v-if="map.runner && map.runner !== '-'" class="text-[10px] text-gray-400 font-mono">BY {{ map.runner
                }}</span>
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
.clip-path-slant {
  clip-path: polygon(10px 0, 100% 0, 100% calc(100% - 10px), calc(100% - 10px) 100%, 0 100%, 0 10px);
}
</style>
