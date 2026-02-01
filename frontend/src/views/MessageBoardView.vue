<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { ChatBubbleLeftEllipsisIcon, UserIcon, PaperAirplaneIcon } from '@heroicons/vue/24/outline';


const messages = ref([]);
const newMessage = ref({ user: '', content: '' });
const isSubmitting = ref(false);

const fetchMessages = async () => {
  try {
    const res = await axios.get('/api/messages');
    messages.value = res.data;
  } catch (e) {
    console.error("無法取得留言", e);
  }
};

const postMessage = async () => {
  if (!newMessage.value.user || !newMessage.value.content) return;
  isSubmitting.value = true;
  try {
    await axios.post('/api/messages', newMessage.value);
    newMessage.value.content = ''; // 清空內容但保留使用者名稱
    await fetchMessages();
  } catch (e) {
    alert("發送失敗");
  } finally {
    isSubmitting.value = false;
  }
};

const formatTime = (timeStr) => {
  const date = new Date(timeStr);
  return date.toLocaleString('zh-TW', { hour12: false });
};

onMounted(fetchMessages);
</script>

<template>
  <div class="flex flex-col h-full gap-6 max-w-4xl mx-auto w-full">
    <div class="flex-1 overflow-y-auto custom-scrollbar space-y-4 pr-2">
      <div v-for="msg in messages" :key="msg.id"
        class="bg-white/[0.01] border-l-2 border-orange-500/30 border-y border-r border-white/5 p-4 transition-all hover:bg-white/[0.03]">
        <div class="flex justify-between items-start mb-2">
          <span class="font-mono text-orange-400 text-xs font-bold">{{ msg.user }}</span>
          <span class="font-mono text-gray-600 text-[10px]">{{ formatTime(msg.created_at) }}</span>
        </div>
        <p class="text-gray-300 text-sm leading-relaxed whitespace-pre-wrap">{{ msg.content }}</p>
      </div>
    </div>
    <div class="bg-white/[0.02] border border-white/10 p-6 rounded-lg">
      <h3 class="font-mono text-orange-500 text-sm mb-4 flex items-center gap-2">
        <ChatBubbleLeftEllipsisIcon class="w-4 h-4" /> NEW_POST
      </h3>
      <div class="space-y-4">
        <div class="relative">
          <UserIcon class="w-4 h-4 absolute left-3 top-3 text-gray-500" />
          <input v-model="newMessage.user" type="text" placeholder="USER_ID"
            class="w-full bg-black/50 border border-white/10 rounded pl-10 pr-4 py-2 text-sm font-mono focus:border-orange-500/50 outline-none" />
        </div>
        <textarea v-model="newMessage.content" rows="3" placeholder="MESSAGE_CONTENT..."
          class="w-full bg-black/50 border border-white/10 rounded p-4 text-sm font-mono focus:border-orange-500/50 outline-none resize-none"></textarea>
        <button @click="postMessage" :disabled="isSubmitting"
          class="w-full bg-orange-600 hover:bg-orange-500 text-black font-mono font-bold py-2 rounded flex items-center justify-center gap-2 transition-colors">
          <PaperAirplaneIcon class="w-4 h-4" /> SEND_DATA
        </button>
      </div>
    </div>

  </div>
</template>
