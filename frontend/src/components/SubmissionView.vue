<script setup>
import { ref, watch, computed } from 'vue';
import axios from 'axios';
import {
  PaperAirplaneIcon,
  ChevronUpDownIcon,
  MagnifyingGlassIcon
} from '@heroicons/vue/24/outline';

// 難度列表
const difficulties = [

  'NOVICE', 'MODERATE', 'BRUTAL', 'INSANE',

  'DUMMY', 'SOLO', 'RACE', 'OLDSCHOOL',

  'DDMAX.EASY', 'DDMAX.NEXT', 'DDMAX.PRO', 'DDMAX.NUT',

  'EVENT', 'FUN'
];

const form = ref({
  difficulty: 'INSANE',
  map_name: '',
  runner: '',

  score: null,
  points: null,
  note: ''
});

// UI 狀態
const status = ref('idle');
const message = ref('');

const showDropdown = ref(false);
const mapOptions = ref([]);
const searchQuery = ref('');

// --- 1. 自動抓取地圖 API ---
const fetchMapOptions = async () => {

  try {
    const res = await axios.get(`/api/map-options?difficulty=${form.value.difficulty}`);
    mapOptions.value = res.data;

    searchQuery.value = '';
    form.value.map_name = '';
  } catch (e) {
    console.error("無法取得地圖列表", e);
  }
};

watch(() => form.value.difficulty, fetchMapOptions, { immediate: true });

// --- 2. 前端搜尋過濾 ---
const filteredOptions = computed(() => {
  if (!searchQuery.value) return mapOptions.value;
  const query = searchQuery.value.toLowerCase();
  return mapOptions.value.filter(m =>
    m.map_name.toLowerCase().includes(query)
  );
});

// --- 3. 選擇地圖 ---
const selectMap = (map) => {
  form.value.map_name = map.map_name;

  searchQuery.value = map.map_name;

  const autoScore = map.points || map.score;
  if (autoScore > 0) {

    form.value.score = autoScore;
    form.value.points = autoScore;
  }
  showDropdown.value = false;
};

// --- 4. 提交 ---
const submitForm = async () => {
  if (!form.value.map_name || !form.value.runner || !form.value.score) {
    alert("請填寫完整資訊");
    return;
  }
  status.value = 'submitting';

  try {
    await axios.post('/api/records', form.value);
    status.value = 'success';
    message.value = 'UPLOAD COMPLETE';

    // Reset
    form.value.map_name = '';
    searchQuery.value = '';
    form.value.score = null;

    form.value.note = '';
    fetchMapOptions();

    setTimeout(() => { status.value = 'idle'; }, 3000);
  } catch (e) {
    status.value = 'error';
    message.value = 'UPLOAD FAILED';
  }
};
</script>

<template>
  <div
    class="w-full flex-1 bg-[#0a0a0a] border border-cyan-500/30 p-4 md:p-8 relative overflow-y-auto custom-scrollbar group flex flex-col">

    <div class="absolute top-0 left-0 w-3 h-3 border-l-2 border-t-2 border-cyan-500"></div>
    <div class="absolute top-0 right-0 w-3 h-3 border-r-2 border-t-2 border-cyan-500"></div>
    <div class="absolute bottom-0 left-0 w-3 h-3 border-l-2 border-b-2 border-cyan-500"></div>
    <div class="absolute bottom-0 right-0 w-3 h-3 border-r-2 border-b-2 border-cyan-500"></div>

    <h2 class="text-2xl font-mono font-bold text-white mb-8 flex items-center gap-2 shrink-0">
      <span class="w-2 h-6 bg-cyan-500 block"></span> NEW_RECORD
    </h2>

    <form @submit.prevent="submitForm" class="w-full">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-x-8 gap-y-8">

        <div class="space-y-2">
          <label class="text-xs font-mono text-cyan-500/70">DIFFICULTY_LEVEL</label>
          <select v-model="form.difficulty"
            class="w-full bg-black/50 border border-white/20 text-white p-3 md:p-4 font-mono focus:border-cyan-500 focus:outline-none transition-colors appearance-none cursor-pointer">
            <option v-for="diff in difficulties" :key="diff" :value="diff">{{ diff }}</option>
          </select>
        </div>

        <div class="space-y-2">
          <label class="text-xs font-mono text-cyan-500/70">RUNNER</label>
          <input v-model="form.runner" type="text" placeholder="Player Name" required
            class="w-full bg-black/50 border border-white/20 text-white p-3 md:p-4 font-mono focus:border-cyan-500 focus:outline-none transition-colors" />
          <p class="text-[10px] text-gray-500 font-mono pt-1">
            For multiple runners, separate with comma (,) or (&).
          </p>
        </div>

        <div class="space-y-2 relative">
          <label class="text-xs font-mono text-cyan-500/70">MAP_NAME (AUTO-SEARCH)</label>
          <div class="relative">
            <input type="text" v-model="searchQuery" @focus="showDropdown = true" @input="form.map_name = searchQuery"
              placeholder="Type to search..."
              class="w-full bg-black/50 border border-white/20 text-white p-3 md:p-4 pl-12 font-mono focus:border-cyan-500 focus:outline-none transition-colors"
              autocomplete="off" />
            <MagnifyingGlassIcon class="w-6 h-6 text-gray-500 absolute left-4 top-4" />

            <button type="button" @click="showDropdown = !showDropdown"
              class="absolute right-3 top-4 text-gray-500 hover:text-white">
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

        <div class="space-y-2 relative">

          <label class="text-xs font-mono text-cyan-500/70">SCORE</label>
          <input v-model.number="form.score" type="number" placeholder="0" required readonly
            class="w-full bg-black/50 border border-white/20 text-cyan-400 p-3 md:p-4 font-mono text-xl font-bold focus:border-cyan-500 focus:outline-none transition-colors" />
        </div>

        <div class="space-y-2 md:col-span-2">
          <label class="text-xs font-mono text-cyan-500/70">NOTE</label>
          <textarea v-model="form.note" rows="3" placeholder="..."
            class="w-full bg-black/50 border border-white/20 text-gray-300 p-3 md:p-4 font-mono focus:border-cyan-500 focus:outline-none transition-colors"></textarea>
        </div>

        <div class="md:col-span-2 pt-4">
          <button type="submit" :disabled="status === 'submitting'" :class="['w-full py-5 font-mono font-bold tracking-widest flex items-center justify-center gap-2 transition-all duration-300 text-cyan-500/70',
            status === 'success' ? 'bg-green-600 text-cyan-300' :
              status === 'error' ? 'bg-red-600 text-white' :
                'bg-cyan-600 hover:bg-cyan-500 text-black']">
            <span v-if="status === 'idle'">SEND_DATA</span>
            <span v-if="status === 'submitting'" class="animate-pulse">PROCESSING...</span>
            <span v-if="status === 'success'">DONE</span>
            <PaperAirplaneIcon v-if="status === 'idle'" class="w-6 h-6" />
          </button>
        </div>

      </div>
    </form>
  </div>
</template>
