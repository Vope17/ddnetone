// src/composables/useGameData.js
import { ref, onMounted, onUnmounted, computed, shallowRef } from 'vue';
import axios from 'axios';

export function useGameData() {
  const summary = ref({ current_score: 0, target_score: 32450, completed_maps: 0 });
  const players = ref([]);
  const maps = shallowRef([]);
  const growthData = ref([]);

  // 模擬數據
  const mockData = () => {
    // --- 模擬數據 (如果API沒通) ---
    summary.value = { current_score: 12450, target_score: 32450, completed_maps: 42 };
    players.value = [
      { id: 1, name: 'KAI_ZEN', role: 'DUELIST', score_contrib: 5200, map_count: 12 },
      { id: 2, name: 'NEON_X', role: 'INITIATOR', score_contrib: 4100, map_count: 10 },
      { id: 3, name: 'VOID_WALKER', role: 'CONTROLLER', score_contrib: 3800, map_count: 9 },
      { id: 4, name: 'CYPHER_99', role: 'SENTINEL', score_contrib: 2200, map_count: 6 },
      { id: 5, name: 'JETT_LAG', role: 'DUELIST', score_contrib: 1500, map_count: 4 },
      { id: 6, name: 'SOVA_DRONE', role: 'INITIATOR', score_contrib: 1200, map_count: 3 },
    ];
    maps.value = [
      { id: 1, map_name: 'ASCENT_X', difficulty: 'Insane', score: 980, runner: 'KAI_ZEN', note: 'Flawless run' },
      { id: 2, map_name: 'HAVEN_ZERO', difficulty: 'Insane', score: 850, runner: 'NEON_X', note: '' },
      { id: 3, map_name: 'BIND_PROTOCOL', difficulty: 'Brutal', score: 1200, runner: 'VOID', note: 'Glitch used' },
      { id: 4, map_name: 'SPLIT_REDUX', difficulty: 'Insane', score: 920, runner: 'KAI_ZEN', note: '' },
      { id: 5, map_name: 'LOTUS_PRIME', difficulty: 'Novice', score: 400, runner: 'NEWBIE', note: '' },
    ];
    growthData.value = Array.from({ length: 12 }, (_, i) => ({ hours: i * 2, points: 1000 + Math.random() * 5000 + (i * 1000) }));
  };

  const fetchData = async () => {
    try {
      // 為了避免畫面閃爍，這裡通常不需要清空數據，直接覆蓋即可
      const [sumRes, playRes, mapRes, growthRes] = await Promise.all([
        axios.get(`/api/summary`),
        axios.get(`/api/leaderboard`),
        axios.get(`/api/maps`),
        axios.get(`/api/growth`)
      ]);

      summary.value = sumRes.data;
      players.value = playRes.data;
      maps.value = mapRes.data;
      growthData.value = growthRes.data;

      // 可以在這裡 console.log("Auto Refreshed") 確認有沒有在跑
    } catch (e) {
      console.warn("API 連線失敗，保持現有數據或切換至展示模式");
      // 注意：自動更新失敗時，通常不建議切換 mockData，保留舊數據體驗較好
      if (players.value.length === 0) mockData();
    }
  };
  // 定義計時器變數
  let pollingTimer = null;

  onMounted(() => {
    // 1. 畫面載入時先抓一次
    fetchData();

    // 2. 設定輪詢：每 5000 毫秒 (5秒) 自動抓取一次
    // 您可以根據需求調整時間，例如 3000 (3秒) 或 10000 (10秒)
    pollingTimer = setInterval(fetchData, 5000);
  });

  // 3. 當組件卸載時，務必清除計時器 (雖然 App.vue 通常不會卸載，但這是好習慣)
  onUnmounted(() => {
    if (pollingTimer) clearInterval(pollingTimer);
  });

  const progressPercent = computed(() => {
    if (!summary.value.target_score) return 0;
    return Math.min(100, (summary.value.current_score / summary.value.target_score) * 100).toFixed(1);
  });

  // Chart Data 配置

  const chartData = computed(() => {
    return {
      labels: growthData.value.map(d => {
        if (!d.timestamp) return ''; // 防止舊資料沒時間
        const date = new Date(d.timestamp);
        // 回傳格式：1/29
        return `${date.getMonth() + 1}/${date.getDate()}`;
      }),
      datasets: [
        {
          label: 'SCORE',
          backgroundColor: (ctx) => {
            const canvas = ctx.chart.ctx;
            const gradient = canvas.createLinearGradient(0, 0, 0, 400);
            gradient.addColorStop(0, 'rgba(6, 182, 212, 0.4)');
            gradient.addColorStop(1, 'rgba(6, 182, 212, 0.0)');
            return gradient;
          },
          borderColor: '#22d3ee',
          borderWidth: 2,
          pointBackgroundColor: '#000',
          pointBorderColor: '#22d3ee',
          pointHoverBackgroundColor: '#fff',
          pointHoverRadius: 6,
          pointRadius: 3,

          fill: true,
          tension: 0.3,
          data: growthData.value.map(d => d.points),

        }
      ]
    };
  });

  return {
    summary,
    players,
    maps,
    growthData,
    progressPercent,
    chartData,
    fetchData
  };
}
