<script setup>
import { ref } from 'vue';
import {
  Bars3Icon, XMarkIcon, ChartBarIcon, MapIcon, SparklesIcon
} from '@heroicons/vue/24/outline';

const props = defineProps({
  activeView: String
});

const emit = defineEmits(['update:activeView']);
const isMobileMenuOpen = ref(false);


const switchView = (view) => {
  emit('update:activeView', view);
  isMobileMenuOpen.value = false;
};
</script>


<template>
  <nav class="w-full z-50 top-0 h-16 border-b border-white/10 bg-[#050505]/80 backdrop-blur-md">
    <div class="max-w-[1920px] mx-auto px-6 h-full flex items-center justify-between">

      <div class="flex items-center gap-4">
        <div class="w-8 h-8 flex items-center justify-center bg-cyan-500/10 border border-cyan-500/50 rounded-sm">
          <SparklesIcon class="w-5 h-5 text-cyan-400" />
        </div>
        <div>
          <div class="font-black text-lg tracking-[0.2em] text-white">DDNETONE<span class="text-cyan-500"></span></div>
        </div>
      </div>

      <div class="hidden md:flex items-center gap-1 bg-white/5 p-1 rounded-lg border border-white/5">
        <button @click="switchView('dashboard')"
          :class="['px-6 py-1.5 rounded-md text-sm font-mono transition-all duration-300 flex items-center gap-2',
            activeView === 'dashboard' ? 'bg-cyan-500/20 text-cyan-300 shadow-[0_0_15px_rgba(6,182,212,0.2)]' : 'hover:text-white text-gray-500']">
          <ChartBarIcon class="w-4 h-4" /> DASHBOARD
        </button>
        <button @click="switchView('maps')"
          :class="['px-6 py-1.5 rounded-md text-sm font-mono transition-all duration-300 flex items-center gap-2',
            activeView === 'maps' ? 'bg-purple-500/20 text-purple-300 shadow-[0_0_15px_rgba(168,85,247,0.2)]' : 'hover:text-white text-gray-500']">
          <MapIcon class="w-4 h-4" /> RECORDS
        </button>
        <button @click="switchView('submission')"
          :class="['px-6 py-1.5 rounded-md text-sm font-mono transition-all duration-300 flex items-center gap-2',
            activeView === 'submission' ? 'bg-green-500/20 text-green-300 shadow-[0_0_15px_rgba(34,197,94,0.2)]' : 'hover:text-white text-gray-500']">
          <span class="text-lg font-bold leading-none">+</span> SUBMISSION
        </button>
      </div>

      <div class="flex items-center gap-4">
        <div class="hidden sm:flex items-center gap-2 text-xs font-mono text-gray-500">
          <span class="w-2 h-2 rounded-full bg-green-500 animate-pulse"></span> SYSTEM ONLINE
        </div>
        <button @click="isMobileMenuOpen = !isMobileMenuOpen" class="md:hidden text-gray-400 hover:text-white">
          <Bars3Icon v-if="!isMobileMenuOpen" class="w-6 h-6" />
          <XMarkIcon v-else class="w-6 h-6" />
        </button>
      </div>
    </div>
  </nav>
</template>
