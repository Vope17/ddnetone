<script setup>
import {
  StarIcon as StarIconOutline,

  UserGroupIcon,
  CalendarIcon
} from '@heroicons/vue/24/outline';
import { StarIcon as StarIconSolid } from '@heroicons/vue/24/solid';

const props = defineProps({
  map: Object

});


const getStatusColor = (status) => {
  if (status === 2) return 'text-cyan-400 border-cyan-500/50 bg-cyan-500/10';
  if (status === 1) return 'text-yellow-400 border-yellow-500/50 bg-yellow-500/10';
  return 'text-gray-500 border-white/10 bg-black/40';
};

const formatDate = (dateStr) => {
  if (!dateStr) return '';
  const date = new Date(dateStr);
  if (isNaN(date.getTime())) return '';
  return date.toLocaleDateString('zh-TW', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  });
};
</script>

<template>
  <div
    :class="['group relative border p-4 transition-all duration-300 hover:translate-y-[-2px] flex flex-col', getStatusColor(map.status)]">
    <div class="flex justify-between items-start mb-2 pl-2">
      <div class="flex items-center gap-0.5" title="Map Difficulty Stars">
        <template v-for="n in 5" :key="n">
          <StarIconSolid v-if="n <= (map.stars || 0)" class="w-3 h-3 text-yellow-500" />
          <StarIconOutline v-else class="w-3 h-3 text-gray-700" />
        </template>
      </div>
      <div class="text-right">
        <span :class="['font-mono text-xl font-bold', map.status === 2 ? 'text-white' : 'text-gray-600']">
          {{ map.points }}
        </span>
        <span class="text-[9px] text-gray-600 ml-1">PTS</span>
      </div>
    </div>

    <div v-if="map.has_dummy" class="flex items-center gap-1 mt-1 text-purple-400" title="Dummy Alive">
      <UserGroupIcon class="w-4 h-4" />
      <span class="text-[9px] font-mono border border-purple-500/30 px-1 rounded bg-purple-500/10">DUMMY</span>
    </div>

    <h4 :class="['font-bold text-md truncate mb-1', map.status === 2 ? 'text-white' : 'text-gray-400']">
      {{ map.map_name }}
    </h4>

    <div class="flex flex-col gap-1 mb-3">
      <div class="flex items-center gap-2">
        <span class="text-[9px] font-mono border border-white/10 px-1 rounded text-gray-500">{{ map.difficulty }}</span>
        <span v-if="map.runner" class="text-[9px] font-mono text-gray-400">BY {{ map.runner }}</span>
      </div>
      <div v-if="map.finish_time" class="flex items-center gap-1.5 text-[9px] font-mono text-cyan-500/70 mt-1">
        <CalendarIcon class="w-3 h-3" />
        <span>{{ formatDate(map.finish_time) }}</span>
      </div>
    </div>

    <div v-if="map.note" class="mt-auto pt-2 border-t border-dashed border-white/10">
      <p class="text-[10px] text-gray-500 italic truncate">{{ map.note }}</p>
    </div>
  </div>
</template>
