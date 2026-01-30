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
  isWip: false,
  hasDummy: false,
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
  // 如果是 WIP，可以允許 Score 為 0，否則必須有 Score
  if (!form.value.map_name || !form.value.runner || (!form.value.isWip && !form.value.score)) {
    alert("請填寫完整資訊");
    return;
  }
  status.value = 'submitting';

  try {
    // 準備 payload
    const payload = {
      difficulty: form.value.difficulty,
      map_name: form.value.map_name,
      runner: form.value.runner,
      score: form.value.score,
      note: form.value.note,
      // ★ 轉換狀態：勾選 WIP -> status 1, 否則 -> status 2
      status: form.value.isWip ? 1 : 2,
      // 確保欄位名稱對應後端 json tag
      has_dummy: form.value.hasDummy
    };

    await axios.post('/api/records', payload);

    status.value = 'success';
    message.value = 'UPLOAD COMPLETE';

    // Reset
    form.value.map_name = '';
    searchQuery.value = '';
    form.value.score = null;
    form.value.note = '';
    form.value.isWip = false;   // Reset checkbox
    form.value.hasDummy = false; // Reset checkbox

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

        <div class="space-y-2 relative">

          <label class="text-xs font-mono text-cyan-500/70">SCORE</label>
          <input v-model.number="form.score" type="number" placeholder="0" required readonly
            class="w-full bg-black/50 border border-white/20 text-cyan-400 p-3 md:p-4 font-mono text-xl font-bold focus:border-cyan-500 focus:outline-none transition-colors" />
        </div>
        <div
          class="md:col-span-2 flex flex-col sm:flex-row justify-start gap-6 bg-white/[0.03] p-4 border border-white/10 rounded">

          <label class="flex items-center gap-3 cursor-pointer group select-none">
            <div class="relative flex items-center">
              <input type="checkbox" v-model="form.isWip" class="peer sr-only" />
              <div
                class="w-5 h-5 border-2 border-gray-500 rounded peer-checked:bg-yellow-500 peer-checked:border-yellow-500 transition-all">
              </div>
              <WrenchScrewdriverIcon
                class="w-3.5 h-3.5 text-black absolute left-0.5 top-0.5 opacity-0 peer-checked:opacity-100 transition-opacity" />
            </div>
            <div>
              <span
                class="block text-sm font-mono font-bold text-gray-300 group-hover:text-white transition-colors">WORK IN
                PROGRESS</span>
              <span class="text-[10px] text-gray-500">Status will be WIP (No Points Awarded)</span>
            </div>
          </label>

          <div class="w-px bg-white/10 hidden sm:block"></div>

          <label class="flex items-center gap-3 cursor-pointer group select-none">
            <div class="relative flex items-center">
              <input type="checkbox" v-model="form.hasDummy" class="peer sr-only" />
              <div
                class="w-5 h-5 border-2 border-gray-500 rounded peer-checked:bg-purple-500 peer-checked:border-purple-500 transition-all">
              </div>
              <UserPlusIcon
                class="w-3.5 h-3.5 text-white absolute left-0.5 top-0.5 opacity-0 peer-checked:opacity-100 transition-opacity" />
            </div>
            <div>
              <span class="block text-sm font-mono font-bold text-gray-300 group-hover:text-white transition-colors">
                DUMMY?</span>
              <span class="text-[10px] text-gray-500">Is dummy alive?</span>
            </div>
          </label>

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
            <span v-if="status === 'idle'">SUBMIT</span>
            <span v-if="status === 'submitting'" class="animate-pulse">PROCESSING...</span>
            <span v-if="status === 'success'">DONE</span>
            <PaperAirplaneIcon v-if="status === 'idle'" class="w-6 h-6" />
          </button>
        </div>

      </div>
    </form>
  </div>
</template>
