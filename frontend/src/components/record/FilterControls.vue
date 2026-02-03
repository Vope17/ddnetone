<script setup>
import {
  MagnifyingGlassIcon,
  StarIcon as StarIconOutline,
  CalendarIcon,
  ArrowsUpDownIcon
} from '@heroicons/vue/24/outline';
import { StarIcon as StarIconSolid } from '@heroicons/vue/24/solid';

const props = defineProps({
  statusFilter: String,
  searchQuery: String,
  sortType: String,
  sortOrder: String
});

const emit = defineEmits([
  'update:statusFilter',
  'update:searchQuery',
  'update:sortType',
  'update:sortOrder'
]);

const toggleSort = (type) => {
  if (props.sortType === type) {
    emit('update:sortType', 'DEFAULT');
  } else {
    emit('update:sortType', type);
  }
};

const toggleOrder = () => {
  emit('update:sortOrder', props.sortOrder === 'ASC' ? 'DESC' : 'ASC');
};
</script>


<template>
  <div class="flex flex-col md:flex-row items-center justify-between gap-4">
  </div>

  <div class="bg-white/[0.02] border-y border-white/10 p-3 flex flex-col gap-4">

    <div class="flex flex-col md:flex-row items-center justify-between gap-4">

      <slot name="stats"></slot>

      <div class="flex bg-black/50 rounded p-1 border border-white/10 w-full md:w-auto overflow-x-auto">
        <button v-for="status in ['All', 'Completed', 'InProgress', 'Incomplete']" :key="status"
          @click="emit('update:statusFilter', status)" :class="['px-4 py-1.5 text-xs font-mono rounded whitespace-nowrap transition-colors',
            statusFilter === status
              ? (status === 'All' ? 'bg-gray-700 text-white' :
                status === 'Completed' ? 'bg-cyan-900/50 text-cyan-300 border-cyan-500/20' :
                  status === 'InProgress' ? 'bg-yellow-900/50 text-yellow-300 border-yellow-500/20' :
                    'bg-gray-800 text-gray-300 border-gray-600/50') + ' shadow-sm border'
              : 'text-gray-500 hover:text-gray-300']">
          {{ status === 'Incomplete' ? 'TODO' : (status === 'InProgress' ? 'WIP' : (status === 'Completed' ? 'DONE' :
          'ALL')) }}
        </button>
      </div>

    </div>

    <div
      class="flex flex-col md:flex-row items-center justify-between gap-4 border-t border-white/5 pt-4 md:border-t-0 md:pt-0">
      <div class="relative w-full group flex-1 md:max-w-md">
        <MagnifyingGlassIcon
          class="w-4 h-4 text-gray-500 absolute left-3 top-2.5 group-focus-within:text-cyan-400 transition-colors" />
        <input :value="searchQuery" @input="emit('update:searchQuery', $event.target.value)" type="text"
          placeholder="SEARCH MAP..."
          class="w-full bg-black/50 border border-white/10 rounded pl-9 pr-3 py-2 text-xs font-mono text-white focus:border-cyan-500/50 focus:outline-none transition-all placeholder:text-gray-700" />
      </div>

      <div class="flex items-center gap-2 w-full md:w-auto overflow-x-auto">
        <div class="flex bg-black/50 rounded border border-white/10 p-0.5">
          <button @click="toggleSort('STARS')"
            :class="['px-3 py-1.5 text-xs font-mono rounded flex items-center gap-1 transition-all', sortType === 'STARS' ? 'bg-yellow-500/20 text-yellow-300' : 'text-gray-500 hover:text-gray-300']">
            <component :is="sortType === 'STARS' ? StarIconSolid : StarIconOutline" class="w-3 h-3" />
            STARS
          </button>


          <button @click="toggleSort('POINTS')"
            :class="['px-3 py-1.5 text-xs font-mono rounded transition-all', sortType === 'POINTS' ? 'bg-cyan-500/20 text-cyan-300' : 'text-gray-500 hover:text-gray-300']">
            PTS
          </button>

          <div class="w-px bg-white/10 mx-1"></div>

          <button @click="toggleSort('DATE')"
            :class="['px-3 py-1.5 text-xs font-mono rounded flex items-center gap-1 transition-all', sortType === 'DATE' ? 'bg-purple-500/20 text-purple-300' : 'text-gray-500 hover:text-gray-300']">
            <CalendarIcon class="w-3 h-3" />
            DATE

          </button>


          <button @click="toggleOrder" class="px-3 py-1.5 text-gray-400 hover:text-white transition-colors">
            <ArrowsUpDownIcon class="w-3 h-3" />
          </button>

        </div>
      </div>
    </div>
  </div>

</template>
