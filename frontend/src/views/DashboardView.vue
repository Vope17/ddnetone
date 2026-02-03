<script setup>
// 引入拆分後的組件
import CurrentPointsCard from '../components/dashboard/CurrentPointsCard.vue';
import MapsProgressCard from '../components/dashboard/MapsProgressCard.vue';
import PointsProgressCard from '../components/dashboard/PointsProgressCard.vue';
import GrowthChart from '../components/dashboard/GrowthChart.vue';
import PlayerLeaderboard from '../components/dashboard/PlayerLeaderboard.vue';

import { watch } from 'vue';

// 接收來自 App.vue 或 Router 的 props
const props = defineProps({
  summary: Object,
  players: Array,
  progressPercent: String,
  chartData: Object,
  growthData: Array
});
</script>


<template>
  <div class="flex flex-col h-full w-full gap-6">

    <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-3 md:gap-6 flex-shrink-0">

      <CurrentPointsCard :current-score="summary?.current_score" />

      <MapsProgressCard :completed-maps="summary?.completed_maps" :growth-data="growthData" />

      <PointsProgressCard :progress-percent="progressPercent" :target-score="summary?.target_score"
        :growth-data="growthData" />

    </div>

    <div class="flex-1 min-h-0 grid grid-cols-1 lg:grid-cols-4 gap-4 md:gap-6 overflow-y-auto lg:overflow-visible">

      <GrowthChart :chart-data="chartData" />

      <PlayerLeaderboard :players="players" />

    </div>
  </div>
</template>
