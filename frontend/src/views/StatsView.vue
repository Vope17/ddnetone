<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import EtaCard from '../components/dashboard/EtaCard.vue';
import DailyActivityChart from '../components/stats/DailyActivityChart.vue';
import PlayerDifficultyChart from '../components/stats/PlayerDifficultyChart.vue';
import ScoreGrowthChart from '../components/stats/ScoreGrowthChart.vue';
import PlayerActivityChart from '../components/stats/PlayerActivityChart.vue';

const props = defineProps({
  growthData: Array,
  summary: Object,
  maps: Array
});

const dailyActivity = ref([]);

onMounted(async () => {
  try {
    const res = await axios.get('/api/daily-activity');
    dailyActivity.value = res.data;
  } catch (e) {
    console.warn('daily-activity fetch failed');
  }
});
</script>

<template>
  <div class="flex flex-col h-full w-full gap-4 overflow-y-auto custom-scrollbar">

    <div class="text-xs font-mono text-gray-600 tracking-widest flex-shrink-0">
      STATS_VIEW // ANALYTICS &amp; INSIGHTS
    </div>

    <!-- Row 1: ETA + 分數折線圖 -->
    <div class="grid grid-cols-1 lg:grid-cols-4 gap-4 flex-shrink-0">
      <div class="lg:col-span-1">
        <EtaCard :growth-data="growthData" :summary="summary" />
      </div>
      <div class="lg:col-span-3">
        <ScoreGrowthChart :daily-activity="dailyActivity" />
      </div>
    </div>

    <!-- Row 2: Daily Activity + 玩家活躍度 -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 flex-shrink-0">
      <DailyActivityChart :daily-activity="dailyActivity" />
      <PlayerActivityChart :maps="maps" />
    </div>

    <!-- Row 3: 玩家難度貢獻圖 -->
    <div class="flex-shrink-0">
      <PlayerDifficultyChart :maps="maps" />
    </div>

  </div>
</template>
