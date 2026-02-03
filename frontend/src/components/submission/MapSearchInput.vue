<script setup>
import { ref, computed, watch } from 'vue';
import axios from 'axios';
import { MagnifyingGlassIcon, ChevronUpDownIcon } from '@heroicons/vue/24/outline';

const props = defineProps({
  modelValue: String, // map_name
  difficulty: String
});

const emit = defineEmits(['update:modelValue', 'select']);

const showDropdown = ref(false);
const mapOptions = ref([]);
const searchQuery = ref('');

// 當難度改變時，重新抓取地圖列表
watch(() => props.difficulty, async (newDiff) => {
  if (!newDiff) return;
  try {
    const res = await axios.get(`/api/map-options?difficulty=${newDiff}`);
    mapOptions.value = res.data;
    // 清空搜尋狀態
    searchQuery.value = '';
    emit('update:modelValue', '');
  } catch (e) {
    console.error("無法取得地圖列表", e);
  }
}, { immediate: true });

// 前端過濾

const filteredOptions = computed(() => {
  if (!searchQuery.value) return mapOptions.value;
  const query = searchQuery.value.toLowerCase();
  return mapOptions.value.filter(m =>
    m.map_name.toLowerCase().includes(query)
  );
});

// 處理輸入
const handleInput = (e) => {
  const val = e.target.value;
  searchQuery.value = val;
  emit('update:modelValue', val);
};

// 選擇地圖
const selectMap = (map) => {
  searchQuery.value = map.map_name;
  emit('update:modelValue', map.map_name);
  emit('select', map); // 將完整地圖物件傳回給父層 (為了拿分數)
  showDropdown.value = false;
};

</script>

<template>

  <div class="space-y-2 relative">
    <label class="text-xs font-mono text-cyan-500/70">MAP_NAME (AUTO-SEARCH)</label>
    <div class="relative">
      <input type="text" :value="searchQuery" @focus="showDropdown = true" @input="handleInput"
        placeholder="Type to search..."
        class="w-full bg-black/50 border border-white/20 text-white p-4 font-mono focus:border-cyan-500 focus:outline-none transition-colors pl-12"
        autocomplete="off" />

      <MagnifyingGlassIcon class="w-6 h-6 text-gray-500 absolute left-4 top-1/2 -translate-y-1/2" />

      <button type="button" @click="showDropdown = !showDropdown"
        class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-white">
        <ChevronUpDownIcon class="w-6 h-6" />
      </button>
    </div>

    <div v-if="showDropdown"
      class="absolute z-50 w-full mt-1 max-h-60 overflow-y-auto bg-[#111] border border-cyan-500/50 shadow-xl custom-scrollbar">
      <div v-if="filteredOptions.length === 0" class="p-3 text-gray-500 font-mono text-sm text-center">
        NO INCOMPLETE MAPS FOUND
      </div>
      <div v-for="map in filteredOptions" :key="map.id" @click="selectMap(map)"
        class="p-3 hover:bg-cyan-900/30 cursor-pointer border-b border-white/5 flex justify-between items-center group">
        <span class="text-gray-200 font-mono group-hover:text-cyan-400">{{ map.map_name }}</span>
        <span v-if="map.points || map.score"
          class="text-xs bg-gray-800 px-2 py-1 rounded text-gray-400 font-mono group-hover:text-cyan-300">

          {{ map.points || map.score }} PTS
        </span>
      </div>
    </div>

    <div v-if="showDropdown" @click="showDropdown = false" class="fixed inset-0 z-40 bg-transparent"></div>
  </div>
</template>
