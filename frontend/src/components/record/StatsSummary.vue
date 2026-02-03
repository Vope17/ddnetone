<script setup>
import { CalculatorIcon, FlagIcon } from '@heroicons/vue/24/outline';

defineProps({
  stats: {
    type: Object,
    required: true,
    default: () => ({ currentScore: 0, totalScore: 0, completedCount: 0, totalCount: 0 })
  }

});

</script>

<template>
  <div class="flex items-center gap-6 w-full md:w-auto justify-between md:justify-start">
    <div class="flex items-center gap-3">
      <div class="p-1.5 bg-cyan-500/10 rounded border border-cyan-500/30 hidden sm:block">
        <CalculatorIcon class="w-4 h-4 text-cyan-400" />
      </div>
      <div>
        <div class="text-[9px] text-gray-500 font-mono uppercase">Score</div>
        <div class="text-sm font-bold text-white font-mono">
          {{ stats.currentScore }} <span class="text-gray-600">/ {{ stats.totalScore }}</span>
        </div>
      </div>
    </div>

    <div class="w-px h-6 bg-white/10 hidden md:block"></div>

    <div class="flex items-center gap-3">

      <div class="p-1.5 bg-purple-500/10 rounded border border-purple-500/30 hidden sm:block">
        <FlagIcon class="w-4 h-4 text-purple-400" />
      </div>
      <div>
        <div class="text-[9px] text-gray-500 font-mono uppercase">Done</div>
        <div class="text-sm font-bold text-white font-mono">
          {{ stats.completedCount }} <span class="text-gray-600">/ {{ stats.totalCount }}</span>

        </div>
      </div>
    </div>
  </div>
</template>
