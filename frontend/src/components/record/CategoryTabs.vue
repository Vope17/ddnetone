<script setup>
import { ref } from 'vue';
import { ChevronDownIcon } from '@heroicons/vue/24/outline';
import { CATEGORIES } from '../../constants/categories'; // 確保路徑正確

const props = defineProps({
  modelValue: String // currentTab
});

const emit = defineEmits(['update:modelValue']);

const isCategoryMenuOpen = ref(false);

const selectCategory = (tab) => {
  emit('update:modelValue', tab);
  isCategoryMenuOpen.value = false;
};
</script>

<template>
  <div class="shrink-0">
    <div class="md:hidden relative z-30">
      <button @click="isCategoryMenuOpen = !isCategoryMenuOpen"
        class="w-full flex items-center justify-between px-4 py-3 bg-cyan-900/20 border border-cyan-500/50 text-cyan-400 font-mono font-bold uppercase tracking-widest rounded">
        <span>{{ modelValue }}</span>
        <ChevronDownIcon :class="['w-5 h-5 transition', isCategoryMenuOpen ? 'rotate-180' : '']" />
      </button>
      <transition enter-active-class="transition duration-100 ease-out" enter-from-class="transform scale-95 opacity-0"
        enter-to-class="transform scale-100 opacity-100" leave-active-class="transition duration-75 ease-in"
        leave-from-class="transform scale-100 opacity-100" leave-to-class="transform scale-95 opacity-0">
        <div v-if="isCategoryMenuOpen"
          class="absolute top-full left-0 w-full mt-2 bg-[#0a0a0a] border border-gray-800 rounded shadow-xl max-h-60 overflow-y-auto custom-scrollbar p-1 grid grid-cols-2 gap-1">
          <button v-for="tab in CATEGORIES" :key="tab" @click="selectCategory(tab)" :class="['px-3 py-2 text-xs font-mono text-left uppercase hover:bg-white/5 rounded',
            modelValue === tab ? 'text-cyan-400 bg-cyan-900/20 border border-cyan-500/30' : 'text-gray-400']">
            {{ tab }}
          </button>
        </div>
      </transition>
    </div>

    <div class="hidden md:flex flex-wrap gap-2">
      <button v-for="tab in CATEGORIES" :key="tab" @click="emit('update:modelValue', tab)"
        :class="['px-4 py-1.5 font-mono text-xs font-bold border transition-all uppercase tracking-widest clip-path-slant',
          modelValue === tab ? 'bg-cyan-500 text-purple-400 border-cyan-400' : 'bg-black/50 text-gray-500 border-gray-800 hover:border-gray-600']">
        {{ tab }}
      </button>
    </div>
  </div>
</template>

<style scoped>
.clip-path-slant {
  clip-path: polygon(10px 0, 100% 0, 100% calc(100% - 10px), calc(100% - 10px) 100%, 0 100%, 0 10px);
}
</style>
