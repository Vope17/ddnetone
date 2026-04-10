<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';


// 引入子組件
import MapSearchInput from '../components/submission/MapSearchInput.vue';
import PlayerSearchInput from '../components/submission/PlayerSearchInput.vue';
import SubmissionFlags from '../components/submission/SubmissionFlags.vue';
import StatusButton from '../components/submission/StatusButton.vue';

const mapInputRef = ref(null);

const difficulties = [
  'ALL',
  'NOVICE', 'MODERATE', 'BRUTAL', 'INSANE',
  'DUMMY', 'SOLO', 'RACE', 'OLDSCHOOL',
  'DDMAX.EASY', 'DDMAX.NEXT', 'DDMAX.PRO', 'DDMAX.NUT',
  'EVENT', 'FUN'
];

// 模式切換：submit / load
const activeMode = ref('submit');

// 表單資料
const form = ref({
  difficulty: 'ALL',
  map_name: '',
  runner: '',
  score: null,
  isWip: false,
  hasDummy: false,
  note: ''
});

// UI 狀態
const status = ref('idle');

// 當子組件選擇了地圖，回填分數並記住實際難度
const selectedMapDifficulty = ref('');
const handleMapSelect = (map) => {
  const autoScore = map.points || map.score;
  if (autoScore > 0) {
    form.value.score = autoScore;
  }
  selectedMapDifficulty.value = map.difficulty || '';
};

// 提交邏輯

const submitForm = async () => {

  // 驗證
  if (!form.value.map_name || !form.value.runner || (!form.value.isWip && !form.value.score)) {
    alert("請填寫完整資訊");
    return;
  }
  if (form.value.difficulty === 'ALL' && !selectedMapDifficulty.value) {
    alert("請從下拉選單選擇地圖");
    return;
  }

  status.value = 'submitting';

  try {
    const payload = {
      difficulty: form.value.difficulty === 'ALL' ? selectedMapDifficulty.value : form.value.difficulty,
      map_name: form.value.map_name,
      runner: form.value.runner,
      score: form.value.score,
      note: form.value.note,
      status: form.value.isWip ? 1 : 2,
      has_dummy: form.value.hasDummy
    };

    await axios.post('/api/records', payload);

    status.value = 'success';

    // Reset Form
    form.value.map_name = '';
    form.value.score = null;
    form.value.note = '';
    form.value.isWip = false;
    form.value.hasDummy = false;
    selectedMapDifficulty.value = '';

    if (mapInputRef.value) {
      mapInputRef.value.refresh();
    }

    // 3秒後恢復按鈕
    setTimeout(() => { status.value = 'idle'; }, 3000);

  } catch (e) {
    status.value = 'error';
    console.error(e);
    // 錯誤時也建議過幾秒恢復 idle，讓使用者重試

    setTimeout(() => { status.value = 'idle'; }, 3000);
  }
};

// ── LOAD MODE ──
const loadOptions = ref([]);
const loadStatus = ref('idle');
const selectedLoadId = ref(null);

const fetchLoadOptions = async () => {
  try {
    const res = await axios.get('/api/load-options');
    loadOptions.value = res.data;
  } catch (e) {
    console.error('Failed to fetch load options', e);
  }
};

const switchMode = (mode) => {
  activeMode.value = mode;
  loadStatus.value = 'idle';
  selectedLoadId.value = null;
  if (mode === 'load') fetchLoadOptions();
};

const selectedLoadMap = () => loadOptions.value.find(m => m.id === selectedLoadId.value) || null;

const submitLoad = async () => {
  const map = selectedLoadMap();
  if (!map) { alert('請選擇要加載的地圖'); return; }
  loadStatus.value = 'submitting';
  try {
    await axios.post('/api/records/load', { map_name: map.map_name, difficulty: map.difficulty });
    loadStatus.value = 'success';
    selectedLoadId.value = null;
    await fetchLoadOptions();
    setTimeout(() => { loadStatus.value = 'idle'; }, 3000);
  } catch (e) {
    const msg = e.response?.data?.error || '加載失敗';
    alert(msg);
    loadStatus.value = 'error';
    setTimeout(() => { loadStatus.value = 'idle'; }, 3000);
  }
};
</script>

<template>

  <div class="w-full flex-1 bg-[#0a0a0a] relative flex flex-col min-h-0 group"
    :class="activeMode === 'load' ? 'border border-amber-500/30' : 'border border-cyan-500/30'">

    <div class="absolute top-0 left-0 w-3 h-3 border-l-2 border-t-2 z-20"
      :class="activeMode === 'load' ? 'border-amber-500' : 'border-cyan-500'"></div>
    <div class="absolute top-0 right-0 w-3 h-3 border-r-2 border-t-2 z-20"
      :class="activeMode === 'load' ? 'border-amber-500' : 'border-cyan-500'"></div>
    <div class="absolute bottom-0 left-0 w-3 h-3 border-l-2 border-b-2 z-20"
      :class="activeMode === 'load' ? 'border-amber-500' : 'border-cyan-500'"></div>
    <div class="absolute bottom-0 right-0 w-3 h-3 border-r-2 border-b-2 z-20"
      :class="activeMode === 'load' ? 'border-amber-500' : 'border-cyan-500'"></div>

    <!-- Mode tabs -->
    <div class="flex border-b border-white/10 flex-shrink-0">
      <button @click="switchMode('submit')"
        class="font-mono text-xs px-5 py-2.5 transition-colors"
        :class="activeMode === 'submit' ? 'text-cyan-400 border-b-2 border-cyan-500' : 'text-gray-500 hover:text-gray-300'">
        SUBMIT
      </button>
      <button @click="switchMode('load')"
        class="font-mono text-xs px-5 py-2.5 transition-colors"
        :class="activeMode === 'load' ? 'text-amber-400 border-b-2 border-amber-500' : 'text-gray-500 hover:text-gray-300'">
        LOAD
      </button>
    </div>

    <!-- SUBMIT mode -->
    <div v-if="activeMode === 'submit'" class="flex-1 overflow-y-auto custom-scrollbar p-4 md:p-8 flex flex-col">

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

          <PlayerSearchInput v-model="form.runner" />

          <MapSearchInput ref="mapInputRef" v-model="form.map_name" :difficulty="form.difficulty"
            @select="handleMapSelect" />

          <div class="space-y-2 relative">
            <label class="text-xs font-mono text-cyan-500/70">SCORE</label>
            <input v-model.number="form.score" type="number" placeholder="0" required readonly
              class="w-full bg-black/50 border border-white/20 text-cyan-400 p-3 md:p-4 font-mono text-xl font-bold focus:border-cyan-500 focus:outline-none transition-colors" />
          </div>

          <SubmissionFlags v-model:isWip="form.isWip" v-model:hasDummy="form.hasDummy" />

          <div class="space-y-2 md:col-span-2">
            <label class="text-xs font-mono text-cyan-500/70">NOTE</label>

            <textarea v-model="form.note" rows="3" placeholder="..."
              class="w-full bg-black/50 border border-white/20 text-gray-300 p-3 md:p-4 font-mono focus:border-cyan-500 focus:outline-none transition-colors"></textarea>
          </div>

          <div class="md:col-span-2 pt-4 pb-4">
            <StatusButton :status="status" />
          </div>

        </div>
      </form>
    </div>

    <!-- LOAD mode -->
    <div v-else class="flex-1 overflow-y-auto custom-scrollbar p-4 md:p-8 flex flex-col">
      <h2 class="text-2xl font-mono font-bold text-white mb-2 flex items-center gap-2 shrink-0">
        <span class="w-2 h-6 bg-amber-500 block"></span> LOAD_MAP
      </h2>
      <p class="text-xs font-mono text-amber-500/60 mb-8">將已完成的地圖標記為已加載，分數將從 10000 開始倒數。</p>

      <div class="flex flex-col gap-6 max-w-md">
        <div class="space-y-2">
          <label class="text-xs font-mono text-amber-500/70">SELECT_COMPLETED_MAP</label>
          <select v-model="selectedLoadId"
            class="w-full bg-black/50 border border-amber-500/30 text-white p-3 md:p-4 font-mono focus:border-amber-500 focus:outline-none transition-colors appearance-none cursor-pointer">
            <option :value="null" disabled>— 選擇已完成的地圖 —</option>
            <option v-for="m in loadOptions" :key="m.id" :value="m.id">
              [{{ m.difficulty }}] {{ m.map_name }} ({{ m.points }} pts) — {{ m.runner }}
            </option>
          </select>
          <p v-if="loadOptions.length === 0" class="text-xs font-mono text-gray-600">尚無可加載的地圖</p>
        </div>

        <div v-if="selectedLoadMap()" class="bg-amber-900/20 border border-amber-500/30 p-4 font-mono text-sm space-y-1">
          <div class="text-amber-300">{{ selectedLoadMap().map_name }}</div>
          <div class="text-gray-400 text-xs">{{ selectedLoadMap().difficulty }} · {{ selectedLoadMap().points }} pts · BY {{ selectedLoadMap().runner }}</div>
        </div>

        <button @click="submitLoad"
          :disabled="!selectedLoadId || loadStatus === 'submitting'"
          class="bg-amber-900/40 border border-amber-500/40 text-amber-300 font-mono text-sm px-6 py-3 hover:bg-amber-900/60 transition-colors disabled:opacity-40 disabled:cursor-not-allowed">
          {{ loadStatus === 'submitting' ? 'LOADING...' : 'CONFIRM_LOAD' }}
        </button>
        <p v-if="loadStatus === 'success'" class="text-green-400 font-mono text-xs">✓ 地圖已加載，分數已更新</p>
        <p v-if="loadStatus === 'error'" class="text-red-400 font-mono text-xs">✗ 加載失敗</p>
      </div>
    </div>
  </div>
</template>

