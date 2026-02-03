<script setup>
import { ref, computed, onMounted } from 'vue';
import axios from 'axios';
import { UserIcon } from '@heroicons/vue/24/outline';

const props = defineProps({
  modelValue: String
});

const emit = defineEmits(['update:modelValue']);

const showDropdown = ref(false);
const playerOptions = ref([]);
const searchQuery = ref(props.modelValue || '');

onMounted(async () => {
  try {
    const res = await axios.get('/api/player-options');
    playerOptions.value = res.data;
  } catch (e) {
    console.error("無法取得玩家列表", e);
  }
});

// --- 邏輯 1: 解析目前輸入框中「已經選定」的名字 ---
const existingPlayers = computed(() => {
  if (!searchQuery.value) return [];
  // 1. 依照分隔符號切分
  const parts = searchQuery.value.split(/[,&]/);
  // 2. 移除最後一個元素 (因為那是正在輸入中的，不算已選定)
  parts.pop();
  // 3. 清理空白並轉小寫 (方便比對)
  return parts.map(p => p.trim().toLowerCase()).filter(p => p);
});

// --- 邏輯 2: 取得最後一個正在輸入的部分 ---
const currentInputPart = computed(() => {
  const parts = searchQuery.value.split(/[,&]/);
  return parts[parts.length - 1].trim(); // 只要最後一段
});

// --- 邏輯 3: 過濾清單 (搜尋 + 防重複) ---
const filteredOptions = computed(() => {
  const query = currentInputPart.value.toLowerCase();

  // 如果輸入框是空的且沒打開選單，就不算
  // 但如果游標在輸入逗號後 (query 為空)，我們希望顯示所有「未選過」的玩家

  return playerOptions.value.filter(name => {
    const lowerName = name.toLowerCase();

    // A. 必須符合當前輸入的關鍵字
    const matchesQuery = lowerName.includes(query);


    // B. 必須「不在」已經選定的名單中 (防重複核心)
    const isNotSelected = !existingPlayers.value.includes(lowerName);

    return matchesQuery && isNotSelected;
  });
});

const handleInput = (e) => {
  const val = e.target.value;
  searchQuery.value = val;
  emit('update:modelValue', val);
  showDropdown.value = true;
};

// --- 邏輯 4: 選擇並自動補逗號 ---
const selectPlayer = (name) => {
  // 1. 使用正則保留分隔符號
  const parts = searchQuery.value.split(/([,&])/);

  // 2. 移除最後一個部分 (即使用者正在打的殘缺名字)
  // 注意：如果輸入的是 "Player1, P"，這裡的 parts 最後一項其實是 " P" (包含前面的空格)
  // 所以 pop() 之後，空格就不見了
  parts.pop();

  // 3. 重新拼接前面的部分
  let prefix = parts.join('');

  // ★★★ 核心修復開始 ★★★
  // 檢查 prefix 的最後一個字元是不是分隔符號 (逗號或&)
  // 如果是，代表我們剛才把跟隨在它後面的空格誤刪了，或者使用者根本沒打空格
  // 我們統一補上一個空格，保持格式整潔
  if (prefix && /[,&]$/.test(prefix)) {
    prefix += ' ';
  }
  // ★★★ 核心修復結束 ★★★

  // 4. 新的值 = 前綴 + 完整名字 + 自動逗號與空格
  const newValue = prefix + name;

  searchQuery.value = newValue;
  emit('update:modelValue', newValue);

  showDropdown.value = false;

  // 讓焦點回到 input，方便繼續打字
  // 如果你希望選完後選單不要關閉，可以把上面的 showDropdown.value = false 拿掉
};
</script>

<template>
  <div class="space-y-2 relative">
    <label class="text-xs font-mono text-cyan-500/70">PLAYERS (AUTO-SEARCH)</label>

    <div class="relative">
      <input type="text" :value="searchQuery" @focus="showDropdown = true" @input="handleInput"
        placeholder="Player1, Player2..."
        class="w-full bg-black/50 border border-white/20 text-white p-4 font-mono focus:border-cyan-500 focus:outline-none transition-colors pl-12"
        autocomplete="off" required />
      <UserIcon class="w-6 h-6 text-gray-500 absolute left-4 top-1/2 -translate-y-1/2" />

    </div>

    <div v-if="showDropdown && filteredOptions.length > 0"
      class="absolute z-50 w-full mt-1 max-h-48 overflow-y-auto bg-[#111] border border-cyan-500/50 shadow-xl custom-scrollbar">
      <div v-for="name in filteredOptions" :key="name" @click="selectPlayer(name)"
        class="p-3 hover:bg-cyan-900/30 cursor-pointer border-b border-white/5 flex items-center gap-2 group">
        <UserIcon class="w-3 h-3 text-gray-600 group-hover:text-cyan-400" />
        <span class="text-gray-200 font-mono group-hover:text-cyan-400">{{ name }}</span>

        <span class="ml-auto text-[9px] text-gray-600 font-mono group-hover:text-cyan-500/50">ADD</span>
      </div>
    </div>

    <div v-if="showDropdown" @click="showDropdown = false" class="fixed inset-0 z-40 bg-transparent"></div>

    <p class="text-[10px] text-gray-500 font-mono pt-1">
      Use <span class="text-cyan-400">,</span> or <span class="text-cyan-400">&</span> to separate multiple runners.
    </p>
  </div>
</template>
