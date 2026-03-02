<script setup>
import { ref, computed } from 'vue';
import axios from 'axios';

const adminKey = ref('');
const isAuthenticated = ref(false);
const authError = ref('');
const records = ref([]);
const players = ref([]);
const loading = ref(false);
const searchQuery = ref('');
const editingId = ref(null);
const editForm = ref({ note: '', runner: '' });
const actionStatus = ref({}); // { [id]: 'loading' | 'done' | 'error' }

const activeTab = ref('records'); // 'records' | 'players'
const editingPlayerId = ref(null);
const editPlayerForm = ref({ name: '', role: '', score_contrib: 0, map_count: 0 });
const playerActionStatus = ref({});
const playerSearch = ref('');

const login = async () => {
  try {
    const res = await axios.get('/api/admin/records', {
      headers: { 'X-Admin-Key': adminKey.value }
    });
    records.value = res.data;
    isAuthenticated.value = true;
    authError.value = '';
  } catch {
    authError.value = '密鑰錯誤';
  }
};

const fetchRecords = async () => {
  loading.value = true;
  try {
    const res = await axios.get('/api/admin/records', {
      headers: { 'X-Admin-Key': adminKey.value }
    });
    records.value = res.data;
  } finally {
    loading.value = false;
  }
};

const fetchPlayers = async () => {
  loading.value = true;
  try {
    const res = await axios.get('/api/admin/players', {
      headers: { 'X-Admin-Key': adminKey.value }
    });
    players.value = res.data;
  } finally {
    loading.value = false;
  }
};

const switchTab = async (tab) => {
  activeTab.value = tab;
  if (tab === 'players' && players.value.length === 0) {
    await fetchPlayers();
  }
};

const startEdit = (record) => {
  editingId.value = record.id;
  editForm.value = { note: record.note || '', runner: record.runner || '' };
};

const cancelEdit = () => {
  editingId.value = null;
};

const saveEdit = async (id) => {
  actionStatus.value[id] = 'loading';
  try {
    const res = await axios.put(`/api/admin/records/${id}`, editForm.value, {
      headers: { 'X-Admin-Key': adminKey.value }
    });
    const idx = records.value.findIndex(r => r.id === id);
    if (idx !== -1) records.value[idx] = res.data;
    editingId.value = null;
    actionStatus.value[id] = 'done';
  } catch {
    actionStatus.value[id] = 'error';
  }
};

const undoRecord = async (record) => {
  if (!confirm(`確定要將 [${record.map_name}] 改回未完成？\n此操作將扣除 ${record.runner} 的積分。`)) return;
  actionStatus.value[record.id] = 'loading';
  try {
    await axios.put(`/api/admin/records/${record.id}/undo`, {}, {
      headers: { 'X-Admin-Key': adminKey.value }
    });
    records.value = records.value.filter(r => r.id !== record.id);
    actionStatus.value[record.id] = 'done';
  } catch {
    actionStatus.value[record.id] = 'error';
  }
};

const startEditPlayer = (player) => {
  editingPlayerId.value = player.id;
  editPlayerForm.value = {
    name: player.name || '',
    role: player.role || '',
    score_contrib: player.score_contrib,
    map_count: player.map_count,
  };
};

const cancelEditPlayer = () => {
  editingPlayerId.value = null;
};

const saveEditPlayer = async (id) => {
  playerActionStatus.value[id] = 'loading';
  try {
    const res = await axios.put(`/api/admin/players/${id}`, editPlayerForm.value, {
      headers: { 'X-Admin-Key': adminKey.value }
    });
    const idx = players.value.findIndex(p => p.id === id);
    if (idx !== -1) players.value[idx] = res.data;
    editingPlayerId.value = null;
    playerActionStatus.value[id] = 'done';
  } catch {
    playerActionStatus.value[id] = 'error';
  }
};

const deletePlayer = async (player) => {
  if (!confirm(`確定要刪除玩家 [${player.name}]？\n此操作無法復原。`)) return;
  playerActionStatus.value[player.id] = 'loading';
  try {
    await axios.delete(`/api/admin/players/${player.id}`, {
      headers: { 'X-Admin-Key': adminKey.value }
    });
    players.value = players.value.filter(p => p.id !== player.id);
    playerActionStatus.value[player.id] = 'done';
  } catch {
    playerActionStatus.value[player.id] = 'error';
  }
};

const filtered = computed(() => {
  if (!searchQuery.value) return records.value;
  const q = searchQuery.value.toLowerCase();
  return records.value.filter(r =>
    r.map_name?.toLowerCase().includes(q) ||
    r.runner?.toLowerCase().includes(q) ||
    r.note?.toLowerCase().includes(q)
  );
});

const filteredPlayers = computed(() => {
  if (!playerSearch.value) return players.value;
  const q = playerSearch.value.toLowerCase();
  return players.value.filter(p =>
    p.name?.toLowerCase().includes(q) ||
    p.role?.toLowerCase().includes(q)
  );
});
</script>

<template>
  <div class="flex flex-col h-full w-full">
    <!-- 登入畫面 -->
    <div v-if="!isAuthenticated" class="flex flex-col items-center justify-center h-full gap-4">
      <div class="text-xs font-mono text-red-500/70 tracking-widest">ADMIN_ACCESS</div>
      <div class="bg-gray-900 border border-white/10 p-6 w-full max-w-sm flex flex-col gap-3">
        <input
          v-model="adminKey"
          type="password"
          placeholder="ADMIN KEY"
          class="bg-black border border-white/20 text-white font-mono text-sm px-3 py-2 outline-none focus:border-red-500/50 w-full"
          @keyup.enter="login"
        />
        <p v-if="authError" class="text-red-400 text-xs font-mono">{{ authError }}</p>
        <button @click="login"
          class="bg-red-900/40 border border-red-500/40 text-red-300 font-mono text-xs px-4 py-2 hover:bg-red-900/60 transition-colors">
          AUTHENTICATE
        </button>
      </div>
    </div>

    <!-- 管理面板 -->
    <div v-else class="flex flex-col h-full gap-3 overflow-hidden">
      <!-- Tab 切換 -->
      <div class="flex items-center gap-0 flex-shrink-0 border-b border-white/10">
        <button
          @click="switchTab('records')"
          class="font-mono text-xs px-4 py-2 transition-colors"
          :class="activeTab === 'records' ? 'text-red-400 border-b-2 border-red-500/60' : 'text-gray-500 hover:text-gray-300'">
          RECORDS
        </button>
        <button
          @click="switchTab('players')"
          class="font-mono text-xs px-4 py-2 transition-colors"
          :class="activeTab === 'players' ? 'text-red-400 border-b-2 border-red-500/60' : 'text-gray-500 hover:text-gray-300'">
          PLAYERS
        </button>
      </div>

      <!-- Records Tab -->
      <template v-if="activeTab === 'records'">
        <div class="flex items-center justify-between flex-shrink-0">
          <div class="text-xs font-mono text-red-500/70 tracking-widest">RECORDS // {{ records.length }}</div>
          <div class="flex items-center gap-2">
            <input
              v-model="searchQuery"
              placeholder="SEARCH..."
              class="bg-black border border-white/10 text-white font-mono text-xs px-3 py-1.5 outline-none focus:border-red-500/30 w-40"
            />
            <button @click="fetchRecords"
              class="bg-white/5 border border-white/10 text-gray-400 font-mono text-xs px-3 py-1.5 hover:text-white transition-colors">
              REFRESH
            </button>
          </div>
        </div>

        <div class="flex-1 overflow-y-auto custom-scrollbar">
          <table class="w-full text-xs font-mono border-collapse">
            <thead class="sticky top-0 bg-[#050505] z-10">
              <tr class="text-gray-600 border-b border-white/5">
                <th class="text-left py-2 px-2 w-8">#</th>
                <th class="text-left py-2 px-2">MAP</th>
                <th class="text-left py-2 px-2 hidden sm:table-cell">DIFF</th>
                <th class="text-left py-2 px-2">RUNNER</th>
                <th class="text-left py-2 px-2 hidden md:table-cell">SCORE</th>
                <th class="text-left py-2 px-2">NOTE</th>
                <th class="text-left py-2 px-2 hidden lg:table-cell">FINISH</th>
                <th class="text-right py-2 px-2">OPS</th>
              </tr>
            </thead>
            <tbody>
              <template v-for="r in filtered" :key="r.id">
                <tr v-if="editingId !== r.id"
                  class="border-b border-white/5 hover:bg-white/2 transition-colors"
                  :class="{ 'opacity-50': actionStatus[r.id] === 'loading' }">
                  <td class="py-1.5 px-2 text-gray-600">{{ r.id }}</td>
                  <td class="py-1.5 px-2 text-white max-w-[120px] truncate">{{ r.map_name }}</td>
                  <td class="py-1.5 px-2 text-gray-400 hidden sm:table-cell">{{ r.difficulty }}</td>
                  <td class="py-1.5 px-2 text-cyan-400">{{ r.runner }}</td>
                  <td class="py-1.5 px-2 text-green-400 hidden md:table-cell">{{ r.score }}</td>
                  <td class="py-1.5 px-2 text-gray-400 max-w-[100px] truncate">{{ r.note || '-' }}</td>
                  <td class="py-1.5 px-2 text-gray-600 hidden lg:table-cell">
                    {{ r.finish_time ? new Date(r.finish_time).toLocaleDateString('sv-SE') : '-' }}
                  </td>
                  <td class="py-1.5 px-2 text-right whitespace-nowrap">
                    <button @click="startEdit(r)"
                      class="text-violet-400 hover:text-violet-200 px-2 py-0.5 border border-violet-400/20 hover:border-violet-400/60 transition-colors mr-1">
                      EDIT
                    </button>
                    <button @click="undoRecord(r)"
                      class="text-red-400 hover:text-red-200 px-2 py-0.5 border border-red-400/20 hover:border-red-400/60 transition-colors">
                      UNDO
                    </button>
                  </td>
                </tr>
                <tr v-else class="border-b border-violet-500/20 bg-violet-900/10">
                  <td class="py-1.5 px-2 text-gray-600">{{ r.id }}</td>
                  <td class="py-1.5 px-2 text-white">{{ r.map_name }}</td>
                  <td class="py-1.5 px-2 text-gray-400 hidden sm:table-cell">{{ r.difficulty }}</td>
                  <td class="py-1.5 px-2">
                    <input v-model="editForm.runner"
                      class="bg-black border border-violet-500/30 text-cyan-300 font-mono text-xs px-2 py-0.5 w-full outline-none focus:border-violet-400" />
                  </td>
                  <td class="py-1.5 px-2 text-green-400 hidden md:table-cell">{{ r.score }}</td>
                  <td class="py-1.5 px-2">
                    <input v-model="editForm.note"
                      class="bg-black border border-violet-500/30 text-gray-300 font-mono text-xs px-2 py-0.5 w-full outline-none focus:border-violet-400" />
                  </td>
                  <td class="py-1.5 px-2 text-gray-600 hidden lg:table-cell">
                    {{ r.finish_time ? new Date(r.finish_time).toLocaleDateString('sv-SE') : '-' }}
                  </td>
                  <td class="py-1.5 px-2 text-right whitespace-nowrap">
                    <button @click="saveEdit(r.id)"
                      class="text-green-400 hover:text-green-200 px-2 py-0.5 border border-green-400/20 hover:border-green-400/60 transition-colors mr-1">
                      SAVE
                    </button>
                    <button @click="cancelEdit"
                      class="text-gray-400 hover:text-gray-200 px-2 py-0.5 border border-gray-400/20 hover:border-gray-400/60 transition-colors">
                      CANCEL
                    </button>
                  </td>
                </tr>
              </template>
            </tbody>
          </table>
        </div>
      </template>

      <!-- Players Tab -->
      <template v-if="activeTab === 'players'">
        <div class="flex items-center justify-between flex-shrink-0">
          <div class="text-xs font-mono text-red-500/70 tracking-widest">PLAYERS // {{ players.length }}</div>
          <div class="flex items-center gap-2">
            <input
              v-model="playerSearch"
              placeholder="SEARCH..."
              class="bg-black border border-white/10 text-white font-mono text-xs px-3 py-1.5 outline-none focus:border-red-500/30 w-40"
            />
            <button @click="fetchPlayers"
              class="bg-white/5 border border-white/10 text-gray-400 font-mono text-xs px-3 py-1.5 hover:text-white transition-colors">
              REFRESH
            </button>
          </div>
        </div>

        <div class="flex-1 overflow-y-auto custom-scrollbar">
          <table class="w-full text-xs font-mono border-collapse">
            <thead class="sticky top-0 bg-[#050505] z-10">
              <tr class="text-gray-600 border-b border-white/5">
                <th class="text-left py-2 px-2 w-8">#</th>
                <th class="text-left py-2 px-2">NAME</th>
                <th class="text-left py-2 px-2">ROLE</th>
                <th class="text-left py-2 px-2">SCORE</th>
                <th class="text-left py-2 px-2">MAPS</th>
                <th class="text-right py-2 px-2">OPS</th>
              </tr>
            </thead>
            <tbody>
              <template v-for="p in filteredPlayers" :key="p.id">
                <!-- 正常行 -->
                <tr v-if="editingPlayerId !== p.id"
                  class="border-b border-white/5 hover:bg-white/2 transition-colors"
                  :class="{ 'opacity-50': playerActionStatus[p.id] === 'loading' }">
                  <td class="py-1.5 px-2 text-gray-600">{{ p.id }}</td>
                  <td class="py-1.5 px-2 text-cyan-400">{{ p.name }}</td>
                  <td class="py-1.5 px-2 text-yellow-400/80">{{ p.role || '-' }}</td>
                  <td class="py-1.5 px-2 text-green-400">{{ p.score_contrib }}</td>
                  <td class="py-1.5 px-2 text-gray-400">{{ p.map_count }}</td>
                  <td class="py-1.5 px-2 text-right">
                    <button @click="startEditPlayer(p)"
                      class="text-violet-400 hover:text-violet-200 px-2 py-0.5 border border-violet-400/20 hover:border-violet-400/60 transition-colors mr-1">
                      EDIT
                    </button>
                    <button @click="deletePlayer(p)"
                      class="text-red-400 hover:text-red-200 px-2 py-0.5 border border-red-400/20 hover:border-red-400/60 transition-colors">
                      DEL
                    </button>
                  </td>
                </tr>
                <!-- 編輯行 -->
                <tr v-else class="border-b border-violet-500/20 bg-violet-900/10">
                  <td class="py-1.5 px-2 text-gray-600">{{ p.id }}</td>
                  <td class="py-1.5 px-2">
                    <input v-model="editPlayerForm.name"
                      class="bg-black border border-violet-500/30 text-cyan-300 font-mono text-xs px-2 py-0.5 w-full outline-none focus:border-violet-400" />
                  </td>
                  <td class="py-1.5 px-2">
                    <input v-model="editPlayerForm.role"
                      class="bg-black border border-violet-500/30 text-yellow-300 font-mono text-xs px-2 py-0.5 w-full outline-none focus:border-violet-400" />
                  </td>
                  <td class="py-1.5 px-2">
                    <input v-model.number="editPlayerForm.score_contrib" type="number" step="0.01"
                      class="bg-black border border-violet-500/30 text-green-300 font-mono text-xs px-2 py-0.5 w-24 outline-none focus:border-violet-400" />
                  </td>
                  <td class="py-1.5 px-2">
                    <input v-model.number="editPlayerForm.map_count" type="number"
                      class="bg-black border border-violet-500/30 text-gray-300 font-mono text-xs px-2 py-0.5 w-16 outline-none focus:border-violet-400" />
                  </td>
                  <td class="py-1.5 px-2 text-right whitespace-nowrap">
                    <button @click="saveEditPlayer(p.id)"
                      class="text-green-400 hover:text-green-200 px-2 py-0.5 border border-green-400/20 hover:border-green-400/60 transition-colors mr-1">
                      SAVE
                    </button>
                    <button @click="cancelEditPlayer"
                      class="text-gray-400 hover:text-gray-200 px-2 py-0.5 border border-gray-400/20 hover:border-gray-400/60 transition-colors">
                      CANCEL
                    </button>
                  </td>
                </tr>
              </template>
            </tbody>
          </table>
        </div>
      </template>
    </div>
  </div>
</template>
