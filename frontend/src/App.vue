<script setup>
import { ref } from 'vue';
import { useGameData } from './composables/useGameData';
import NavBar from './components/NavBar.vue';
import DashboardView from './components/DashboardView.vue';
import RecordsView from './components/RecordsView.vue';
import SubmissionView from './components/SubmissionView.vue';
// import MessageBoardView from './components/MessageBoardView.vue';

const activeView = ref('dashboard');
const { summary, players, maps, progressPercent, chartData, fetchData } = useGameData();
</script>

<template>
  <div
    class="h-[100dvh] w-full bg-[#050505] text-gray-300 font-sans overflow-hidden flex flex-col selection:bg-cyan-500/30 selection:text-cyan-200">
    <div class="fixed inset-0 z-0 pointer-events-none select-none">
      <div
        class="absolute inset-0 bg-[linear-gradient(rgba(18,16,16,0)_50%,rgba(0,0,0,0.25)_50%),linear-gradient(90deg,rgba(255,0,0,0.06),rgba(0,255,0,0.02),rgba(0,0,255,0.06))] z-[2] bg-[length:100%_4px,3px_100%] pointer-events-none">
      </div>
      <div
        class="absolute top-[-20%] left-[-10%] w-[60%] h-[60%] bg-cyan-900/10 rounded-full blur-[30px] transform-gpu will-change-transform">
      </div>
      <div
        class="absolute bottom-[-20%] right-[-10%] w-[60%] h-[60%] bg-purple-900/10 rounded-full blur-[30px] transform-gpu will-change-transform">
      </div>
      <div
        class="absolute inset-0 bg-[linear-gradient(rgba(255,255,255,0.03)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.03)_1px,transparent_1px)] bg-[size:50px_50px] opacity-30 transform-gpu will-change-transform">
      </div>
    </div>

    <NavBar :activeView="activeView" @update:activeView="val => activeView = val" />

    <main class="relative z-10 flex-1 flex flex-col overflow-hidden custom-scrollbar px-3 pb-6 pt-6">

      <transition name="fade-slide" mode="out-in">
        <DashboardView v-if="activeView === 'dashboard'" :summary="summary" :players="players"
          :progressPercent="progressPercent" :chartData="chartData" />

        <RecordsView v-else-if="activeView === 'maps'" :maps="maps" @record-deleted="fetchData" />

        <SubmissionView v-else-if="activeView === 'submission'" />

        <!-- <MessageBoardView v-else-if="activeView === 'board'" /> -->
      </transition>
    </main>
  </div>
</template>

<style>
/* 樣式保持不變 */
body,
html {
  margin: 0;
  padding: 0;
  overflow: hidden;
  height: 100%;
  width: 100%;
}

@import url('https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;700&display=swap');

body {
  font-family: 'Inter', sans-serif;
}

.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

.custom-scrollbar::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.3);
}

.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #1f2937;
  border-radius: 3px;
}

.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #4b5563;
}

.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.3s ease;
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
