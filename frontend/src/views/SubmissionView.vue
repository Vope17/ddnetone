<script setup>
import { ref } from 'vue';
import axios from 'axios';

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

const form = ref({
  difficulty: 'ALL',
  map_name: '',
  runner: '',
  score: null,
  isWip: false,
  hasDummy: false,
  note: ''
});

// UI status
const status = ref('idle');

const selectedMapDifficulty = ref('');
const handleMapSelect = (map) => {
  const autoScore = map.points || map.score;
  if (autoScore > 0) {
    form.value.score = autoScore;
  }
  selectedMapDifficulty.value = map.difficulty || '';
};

// Submission is open
const SUBMISSION_OPEN = true;

const submitForm = async () => {
  if (!SUBMISSION_OPEN) {
    status.value = 'closed';
    return;
  }
  if (!form.value.map_name || !form.value.runner || (!form.value.isWip && !form.value.score)) {
    alert("Please fill in all required fields");
    return;
  }
  if (form.value.difficulty === 'ALL' && !selectedMapDifficulty.value) {
    alert("Please select a map from the dropdown");
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

    form.value.map_name = '';
    form.value.score = null;
    form.value.note = '';
    form.value.isWip = false;
    form.value.hasDummy = false;
    selectedMapDifficulty.value = '';

    if (mapInputRef.value) {
      mapInputRef.value.refresh();
    }

    setTimeout(() => { status.value = 'idle'; }, 3000);

  } catch (e) {
    status.value = 'error';
    console.error(e);
    setTimeout(() => { status.value = 'idle'; }, 3000);
  }
};
</script>

<template>

  <div class="w-full flex-1 bg-[#0a0a0a] relative flex flex-col min-h-0 group border border-cyan-500/30">

    <div class="absolute top-0 left-0 w-3 h-3 border-l-2 border-t-2 z-20 border-cyan-500"></div>
    <div class="absolute top-0 right-0 w-3 h-3 border-r-2 border-t-2 z-20 border-cyan-500"></div>
    <div class="absolute bottom-0 left-0 w-3 h-3 border-l-2 border-b-2 z-20 border-cyan-500"></div>
    <div class="absolute bottom-0 right-0 w-3 h-3 border-r-2 border-b-2 z-20 border-cyan-500"></div>

    <!-- SUBMIT form -->
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
  </div>
</template>
