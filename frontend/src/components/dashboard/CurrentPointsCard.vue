<script setup>
import { SparklesIcon } from '@heroicons/vue/24/outline';


defineProps({
  currentScore: {
    type: Number,
    default: 0
  }
});
</script>

<template>
  <div class="bg-gradient-to-br from-gray-900 to-black border border-white/10 p-5 relative overflow-hidden group">
    <div class="absolute right-0 top-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
      <SparklesIcon class="w-16 h-16 text-cyan-500" />
    </div>
    <div class="text-xs font-mono text-cyan-500/70 mb-1 flex items-center gap-2">
      <span class="w-1 h-1 bg-cyan-500"></span> CURRENT_POINTS
    </div>
    <div class="text-3xl font-bold text-white font-mono tracking-tight">
      {{ currentScore.toLocaleString() }}
    </div>

    <div class="mt-2 w-full bg-gray-800 h-1 rounded-full overflow-hidden">
      <div class="h-full bg-cyan-500 w-[60%]"></div>
    </div>
  </div>
</template>
