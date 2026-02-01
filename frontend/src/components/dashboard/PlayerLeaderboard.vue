<script setup>
import { TrophyIcon } from '@heroicons/vue/24/outline';

defineProps({
  players: Array
});
</script>

<template>
  <div class="lg:col-span-1 bg-[#0a0a0a] border border-yellow-500/20 relative flex flex-col h-full overflow-hidden">
    <div class="p-4 border-b border-white/5 bg-yellow-500/5 flex items-center justify-between shrink-0">
      <h3 class="font-mono text-sm text-yellow-500 flex items-center gap-2">
        <TrophyIcon class="w-4 h-4" /> TOP_PLAYERS
      </h3>
      <span class="text-[10px] border border-yellow-900 text-yellow-700 px-1 rounded">S-TIER</span>
    </div>
    <div class="flex-1 overflow-y-auto p-2 space-y-2 custom-scrollbar">
      <div v-for="(p, index) in players" :key="p.id"
        class="group relative p-3 border border-white/5 hover:border-yellow-500/30 bg-white/[0.02] hover:bg-yellow-500/[0.05] transition-all duration-200">
        <div class="flex items-center gap-3 relative z-10">
          <div class="flex flex-col items-center justify-center w-8">
            <div
              :class="['font-mono text-lg font-bold italic leading-none',
                index === 0 ? 'text-yellow-400' : index === 1 ? 'text-gray-300' : index === 2 ? 'text-orange-400' : 'text-gray-700']">
              {{ index + 1 }}
            </div>

          </div>
          <div class="flex-1 min-w-0">
            <div class="font-bold text-sm text-gray-200 truncate group-hover:text-yellow-100 transition">{{ p.name }}
            </div>
            <div class="flex items-center gap-2 text-[10px] text-gray-500 font-mono">
              <span class="uppercase tracking-wider text-gray-600">{{ p.role }}</span>
            </div>
          </div>
          <div class="text-right">
            <div class="font-mono font-bold text-cyan-400 text-sm">{{ p.score_contrib }}</div>
            <div class="text-[9px] text-gray-600">{{ p.map_count }} MAPS</div>
          </div>
        </div>
      </div>
    </div>
  </div>

</template>
