<script setup>
import { ref, computed } from 'vue';
import axios from 'axios';

const adminKey = ref('');
const isAuthenticated = ref(false);
const authError = ref('');
const records = ref([]);
const loading = ref(false);
const searchQuery = ref('');
const editingId = ref(null);
const editForm = ref({ note: '', runner: '' });
const actionStatus = ref({});

const activeTab = ref('records'); // 'records' | 'maps'

const mapDifficulties = [
  'NOVICE', 'MODERATE', 'BRUTAL', 'INSANE',
  'DUMMY', 'SOLO', 'RACE', 'OLDSCHOOL',
  'DDMAX.EASY', 'DDMAX.NEXT', 'DDMAX.PRO', 'DDMAX.NUT',
  'EVENT', 'FUN'
];

const mapForm = ref({ map_name: '', difficulty: 'NOVICE', points: 0, stars: 0 });
const mapStatus = ref('idle'); // 'idle' | 'loading' | 'done' | 'error'
const mapError = ref('');

const login = async () => {
  try {
    const res = await axios.get('/api/admin/records', {
      headers: { 'X-Admin-Key': adminKey.value }
    });
    records.value = res.data;
    isAuthenticated.value = true;
    authError.value = '';
  } catch {
    authError.value = 'Invalid admin key';
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

const switchTab = (tab) => {
  activeTab.value = tab;
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
  if (!confirm(`Revert [${record.map_name}] to incomplete?\nThis will deduct ${record.runner}'s points.`)) return;
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

const createMap = async () => {
  mapStatus.value = 'loading';
  mapError.value = '';
  try {
    await axios.post('/api/admin/maps', mapForm.value, {
      headers: { 'X-Admin-Key': adminKey.value }
    });
    mapStatus.value = 'done';
    mapForm.value = { map_name: '', difficulty: 'NOVICE', points: 0, stars: 0 };
  } catch (err) {
    mapStatus.value = 'error';
    mapError.value = err.response?.data?.error || 'Failed to create map';
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
          @click="switchTab('maps')"
          class="font-mono text-xs px-4 py-2 transition-colors"
          :class="activeTab === 'maps' ? 'text-red-400 border-b-2 border-red-500/60' : 'text-gray-500 hover:text-gray-300'">
          MAPS
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
                  <td class="py-1.5 px-2 max-w-[120px] truncate text-white">{{ r.map_name }}</td>
                  <td class="py-1.5 px-2 text-gray-400 hidden sm:table-cell">{{ r.difficulty }}</td>
                  <td class="py-1.5 px-2 text-cyan-400">{{ r.runner }}</td>
                  <td class="py-1.5 px-2 text-green-400 hidden md:table-cell">{{ r.score }}</td>
                  <td class="py-1.5 px-2 text-gray-400 max-w-[100px] truncate">{{ r.note || '-' }}</td>
                  <td class="py-1.5 px-2 text-gray-600 hidden lg:table-cell">
                    {{ r.finish_time ? new Date(r.finish_time).toLocaleDateString('sv-SE') : '-' }}
                  </td>
                  <td class="py-1.5 px-2 text-right whitespace-nowrap">
                    <button v-if="editingId !== r.id" @click="startEdit(r)"
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

      <!-- Maps Tab -->
      <template v-if="activeTab === 'maps'">
        <div class="text-xs font-mono text-red-500/70 tracking-widest flex-shrink-0">ADD MAP</div>
        <div class="flex-1 overflow-y-auto custom-scrollbar">
          <div class="bg-gray-900/40 border border-white/10 p-4 flex flex-col gap-3 max-w-md">
            <div class="flex flex-col gap-1">
              <label class="text-gray-500 font-mono text-xs">MAP NAME</label>
              <input
                v-model="mapForm.map_name"
                placeholder="map name..."
                class="bg-black border border-white/20 text-white font-mono text-xs px-3 py-2 outline-none focus:border-red-500/50 w-full"
              />
            </div>
            <div class="flex flex-col gap-1">
              <label class="text-gray-500 font-mono text-xs">DIFFICULTY</label>
              <select
                v-model="mapForm.difficulty"
                class="bg-black border border-white/20 text-white font-mono text-xs px-3 py-2 outline-none focus:border-red-500/50 w-full">
                <option v-for="d in mapDifficulties" :key="d" :value="d">{{ d }}</option>
              </select>
            </div>
            <div class="flex gap-3">
              <div class="flex flex-col gap-1 flex-1">
                <label class="text-gray-500 font-mono text-xs">POINTS</label>
                <input
                  v-model.number="mapForm.points"
                  type="number" min="0"
                  class="bg-black border border-white/20 text-white font-mono text-xs px-3 py-2 outline-none focus:border-red-500/50 w-full"
                />
              </div>
              <div class="flex flex-col gap-1 flex-1">
                <label class="text-gray-500 font-mono text-xs">STARS (0-5)</label>
                <input
                  v-model.number="mapForm.stars"
                  type="number" min="0" max="5"
                  class="bg-black border border-white/20 text-white font-mono text-xs px-3 py-2 outline-none focus:border-red-500/50 w-full"
                />
              </div>
            </div>
            <button
              @click="createMap"
              :disabled="mapStatus === 'loading' || !mapForm.map_name"
              class="bg-red-900/40 border border-red-500/40 text-red-300 font-mono text-xs px-4 py-2 hover:bg-red-900/60 transition-colors disabled:opacity-40 disabled:cursor-not-allowed">
              {{ mapStatus === 'loading' ? 'ADDING...' : 'ADD MAP' }}
            </button>
            <p v-if="mapStatus === 'done'" class="text-green-400 font-mono text-xs">✓ Map created successfully</p>
            <p v-if="mapStatus === 'error'" class="text-red-400 font-mono text-xs">✗ {{ mapError }}</p>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>
