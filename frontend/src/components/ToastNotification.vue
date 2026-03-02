<script setup>
import { ref } from 'vue';

const toasts = ref([]);
let nextId = 0;

const addToast = (toast) => {
  const id = nextId++;
  toasts.value.push({ id, ...toast });
  setTimeout(() => {
    toasts.value = toasts.value.filter(t => t.id !== id);
  }, 4000);
};

defineExpose({ addToast });
</script>

<template>
  <div class="fixed bottom-6 right-6 z-[100] flex flex-col gap-2 pointer-events-none">
    <transition-group name="toast">
      <div v-for="t in toasts" :key="t.id"
        class="pointer-events-auto flex items-start gap-3 bg-gray-950 border px-4 py-3 shadow-2xl min-w-[260px] max-w-[340px]"
        :class="t.type === 'success'
          ? 'border-cyan-500/50 shadow-cyan-500/10'
          : 'border-red-500/50 shadow-red-500/10'">

        <!-- 左側指示條 -->
        <div class="w-[2px] self-stretch flex-shrink-0 rounded-full"
          :class="t.type === 'success' ? 'bg-cyan-500' : 'bg-red-500'"></div>

        <div class="flex flex-col gap-0.5 min-w-0">
          <div class="text-[9px] font-mono tracking-widest"
            :class="t.type === 'success' ? 'text-cyan-500/70' : 'text-red-500/70'">
            {{ t.type === 'success' ? 'MAP_CLEARED' : 'SYSTEM_ALERT' }}
          </div>
          <div class="text-white font-mono text-xs font-bold truncate">{{ t.title }}</div>
          <div v-if="t.subtitle" class="text-gray-400 font-mono text-[10px] truncate">{{ t.subtitle }}</div>
        </div>

        <!-- 發光點 -->
        <div class="ml-auto flex-shrink-0 mt-0.5">
          <div class="w-1.5 h-1.5 rounded-full animate-pulse"
            :class="t.type === 'success' ? 'bg-cyan-500' : 'bg-red-500'"></div>
        </div>
      </div>
    </transition-group>
  </div>
</template>

<style scoped>
.toast-enter-active {
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.toast-leave-active {
  transition: all 0.25s ease-in;
}
.toast-enter-from {
  opacity: 0;
  transform: translateX(60px);
}
.toast-leave-to {
  opacity: 0;
  transform: translateX(60px);
}
</style>
