<script setup>
import { ref } from 'vue';
import axios from 'axios';


// 引入子組件
import MapSearchInput from '../components/submission/MapSearchInput.vue';
import PlayerSearchInput from '../components/submission/PlayerSearchInput.vue';
import SubmissionFlags from '../components/submission/SubmissionFlags.vue';
import StatusButton from '../components/submission/StatusButton.vue';

// 難度列表
const difficulties = [
  'NOVICE', 'MODERATE', 'BRUTAL', 'INSANE',
  'DUMMY', 'SOLO', 'RACE', 'OLDSCHOOL',
  'DDMAX.EASY', 'DDMAX.NEXT', 'DDMAX.PRO', 'DDMAX.NUT',
  'EVENT', 'FUN'
];

// 表單資料
const form = ref({
  difficulty: 'INSANE',
  map_name: '',
  runner: '',
  score: null,
  isWip: false,
  hasDummy: false,
  note: ''
});

// UI 狀態
const status = ref('idle');

// 當子組件選擇了地圖，回填分數
const handleMapSelect = (map) => {
  const autoScore = map.points || map.score;
  if (autoScore > 0) {
    form.value.score = autoScore;
  }
};

// 提交邏輯

const submitForm = async () => {

  // 驗證
  if (!form.value.map_name || !form.value.runner || (!form.value.isWip && !form.value.score)) {
    alert("請填寫完整資訊");
    return;
  }

  status.value = 'submitting';

  try {
    const payload = {
      difficulty: form.value.difficulty,
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

    // 3秒後恢復按鈕
    setTimeout(() => { status.value = 'idle'; }, 3000);

  } catch (e) {
    status.value = 'error';
    console.error(e);
    // 錯誤時也建議過幾秒恢復 idle，讓使用者重試

    setTimeout(() => { status.value = 'idle'; }, 3000);
  }
};
</script>

<template>

  <div class="w-full flex-1 bg-[#0a0a0a] border border-cyan-500/30 relative flex flex-col min-h-0 group">

    <div class="absolute top-0 left-0 w-3 h-3 border-l-2 border-t-2 border-cyan-500 z-20"></div>
    <div class="absolute top-0 right-0 w-3 h-3 border-r-2 border-t-2 border-cyan-500 z-20"></div>
    <div class="absolute bottom-0 left-0 w-3 h-3 border-l-2 border-b-2 border-cyan-500 z-20"></div>
    <div class="absolute bottom-0 right-0 w-3 h-3 border-r-2 border-b-2 border-cyan-500 z-20"></div>


    <div class="flex-1 overflow-y-auto custom-scrollbar p-4 md:p-8 flex flex-col">

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

          <MapSearchInput v-model="form.map_name" :difficulty="form.difficulty" @select="handleMapSelect" />

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
  </div>
</template>
