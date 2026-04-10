<script setup>
import { toRef, computed } from 'vue';
import { SparklesIcon } from '@heroicons/vue/24/outline';
import { useCountup } from '../../composables/useCountup';

const props = defineProps({
  currentScore: {
    type: Number,
    default: 0
  },
  loadedMaps: {
    type: Number,
    default: 0
  }
});

const isLoadingPhase = computed(() => props.loadedMaps > 0);
const animatedScore = useCountup(toRef(props, 'currentScore'));
</script>

<template>
  <div class="bg-gradient-to-br from-gray-900 to-black border p-5 relative overflow-hidden group"
    :class="isLoadingPhase ? 'border-amber-500/30' : 'border-white/10'">
    <div class="absolute right-0 top-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
      <SparklesIcon class="w-16 h-16" :class="isLoadingPhase ? 'text-amber-500' : 'text-cyan-500'" />
    </div>
    <div class="text-xs font-mono mb-1 flex items-center gap-2"
      :class="isLoadingPhase ? 'text-amber-500/70' : 'text-cyan-500/70'">
      <span class="w-1 h-1" :class="isLoadingPhase ? 'bg-amber-500' : 'bg-cyan-500'"></span>
      {{ isLoadingPhase ? 'LOADING_PHASE' : 'CURRENT_POINTS' }}
    </div>
    <div class="text-3xl font-bold font-mono tracking-tight"
      :class="isLoadingPhase ? 'text-amber-300' : 'text-white'">
      {{ animatedScore.toLocaleString() }}
    </div>
    <div v-if="isLoadingPhase" class="text-[10px] font-mono text-amber-500/50 mt-1">
      {{ loadedMaps }} MAP{{ loadedMaps > 1 ? 'S' : '' }} LOADED
    </div>

    <div class="mt-2 w-full bg-gray-800 h-1 rounded-full overflow-hidden">
      <div class="h-full w-[60%]" :class="isLoadingPhase ? 'bg-amber-500' : 'bg-cyan-500'"></div>
    </div>
  </div>
</template>
