<script setup>
import { WrenchScrewdriverIcon, UserPlusIcon } from '@heroicons/vue/24/outline';


defineProps({
  isWip: Boolean,
  hasDummy: Boolean
});

const emit = defineEmits(['update:isWip', 'update:hasDummy']);
</script>

<template>
  <div
    class="md:col-span-2 flex flex-col sm:flex-row justify-start gap-6 bg-white/[0.03] p-4 border border-white/10 rounded">


    <label class="flex items-center gap-3 cursor-pointer group select-none">
      <div class="relative flex items-center">
        <input type="checkbox" :checked="isWip" @change="emit('update:isWip', $event.target.checked)"
          class="peer sr-only" />
        <div
          class="w-5 h-5 border-2 border-gray-500 rounded peer-checked:bg-yellow-500 peer-checked:border-yellow-500 transition-all">
        </div>
        <WrenchScrewdriverIcon
          class="w-3.5 h-3.5 text-black absolute left-0.5 top-0.5 opacity-0 peer-checked:opacity-100 transition-opacity" />
      </div>
      <div>
        <span class="block text-sm font-mono font-bold text-gray-300 group-hover:text-white transition-colors">WORK IN
          PROGRESS</span>
        <span class="text-[10px] text-gray-500">Status will be WIP (No Points Awarded)</span>

      </div>
    </label>

    <div class="w-px bg-white/10 hidden sm:block"></div>

    <label class="flex items-center gap-3 cursor-pointer group select-none">
      <div class="relative flex items-center">
        <input type="checkbox" :checked="hasDummy" @change="emit('update:hasDummy', $event.target.checked)"
          class="peer sr-only" />
        <div
          class="w-5 h-5 border-2 border-gray-500 rounded peer-checked:bg-purple-500 peer-checked:border-purple-500 transition-all">
        </div>
        <UserPlusIcon
          class="w-3.5 h-3.5 text-white absolute left-0.5 top-0.5 opacity-0 peer-checked:opacity-100 transition-opacity" />
      </div>
      <div>
        <span
          class="block text-sm font-mono font-bold text-gray-300 group-hover:text-white transition-colors">DUMMY?</span>
        <span class="text-[10px] text-gray-500">Is dummy alive?</span>
      </div>
    </label>

  </div>
</template>
