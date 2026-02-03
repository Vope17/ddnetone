<script setup>
import { ref, computed } from 'vue';
import { CpuChipIcon } from '@heroicons/vue/24/outline';

// 引入子組件
import CategoryTabs from '../components/record/CategoryTabs.vue';
import StatsSummary from '../components/record/StatsSummary.vue';
import FilterControls from '../components/record/FilterControls.vue';
import MapCard from '../components/record/MapCard.vue';

const props = defineProps({
  maps: Array
});

// --- State ---
const currentTab = ref('All');
const statusFilter = ref('Completed');
const searchQuery = ref('');
const sortType = ref('DEFAULT');

const sortOrder = ref('DESC');


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
  return { currentScore, totalScore, totalCount, completedCount };
});

// --- 3. 進階篩選與排序邏輯 ---
const finalFilteredMaps = computed(() => {
  let list = mapsInCurrentCategory.value;

  // 狀態篩選
  if (statusFilter.value === 'Completed') list = list.filter(m => m.status === 2);
  else if (statusFilter.value === 'InProgress') list = list.filter(m => m.status === 1);
  else if (statusFilter.value === 'Incomplete') list = list.filter(m => m.status === 0 || !m.status);

  // 搜尋
  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase();
    list = list.filter(m =>
      m.map_name.toLowerCase().includes(query) ||
      (m.runner && m.runner.toLowerCase().includes(query)) ||
      (m.note && m.note.toLowerCase().includes(query))
    );
  }

  // 排序
  return [...list].sort((a, b) => {
    let valA = 0, valB = 0;
    if (sortType.value === 'STARS') {
      valA = a.stars || 0; valB = b.stars || 0;
    } else if (sortType.value === 'POINTS') {
      valA = a.points || 0; valB = b.points || 0;
    } else if (sortType.value === 'DATE') {
      valA = a.finish_time ? new Date(a.finish_time).getTime() : 0;
      valB = b.finish_time ? new Date(b.finish_time).getTime() : 0;
    } else {
      valA = a.score || 0; valB = b.score || 0;
      if (valA === valB) { valA = a.points || 0; valB = b.points || 0; }
    }
    return sortOrder.value === 'DESC' ? valB - valA : valA - valB;
  });
});
</script>

<template>
  <div class="flex flex-col h-full overflow-hidden gap-4">
    <div class="shrink-0 space-y-3">

      <CategoryTabs v-model="currentTab" />

      <FilterControls v-model:statusFilter="statusFilter" v-model:searchQuery="searchQuery" v-model:sortType="sortType"
        v-model:sortOrder="sortOrder">
        <template #stats>
          <StatsSummary :stats="categoryStats" />
        </template>
      </FilterControls>
    </div>

    <div class="flex-1 overflow-y-auto custom-scrollbar pb-20 px-1">
      <div v-if="finalFilteredMaps.length === 0" class="h-64 flex flex-col items-center justify-center text-gray-600">
        <CpuChipIcon class="w-12 h-12 mb-4 opacity-50" />
        <div class="font-mono text-sm">NO RECORDS FOUND</div>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-4">
        <MapCard v-for="map in finalFilteredMaps" :key="map.id" :map="map" />
      </div>
    </div>
  </div>
</template>
