<script setup>
import { PaperAirplaneIcon } from '@heroicons/vue/24/outline';

defineProps({
  status: {
    type: String, // 'idle' | 'submitting' | 'success' | 'error' | 'closed'
    default: 'idle'
  }
});
</script>


<template>
  <button type="submit" :disabled="status === 'submitting' || status === 'closed'" :class="['w-full py-5 font-mono font-bold tracking-widest flex items-center justify-center gap-2 transition-all duration-300 text-cyan-500/70',
    status === 'success' ? 'bg-green-600 text-cyan-300' :
      status === 'error' ? 'bg-red-600 text-white' :
        status === 'closed' ? 'bg-gray-800 text-gray-500 cursor-not-allowed opacity-60' :
          'bg-cyan-600 hover:bg-cyan-500 text-black']">

    <span v-if="status === 'idle'">SUBMIT</span>
    <span v-if="status === 'submitting'" class="animate-pulse">PROCESSING...</span>
    <span v-if="status === 'success'">DONE</span>
    <span v-if="status === 'error'">ERROR</span>
    <span v-if="status === 'closed'">SUBMISSIONS_CLOSED</span>

    <PaperAirplaneIcon v-if="status === 'idle'" class="w-6 h-6" />
  </button>
</template>
