<script setup>
import { computed } from 'vue';

const props = defineProps({
  maps: Array
});

const DIFF_CONFIG = {
  'NOVICE':    { color: '#22c55e', glow: 'rgba(34,197,94,0.4)' },
  'MODERATE':  { color: '#3b82f6', glow: 'rgba(59,130,246,0.4)' },
  'BRUTAL':    { color: '#a855f7', glow: 'rgba(168,85,247,0.4)' },
  'INSANE':    { color: '#ef4444', glow: 'rgba(239,68,68,0.4)' },
  'DUMMY':     { color: '#f97316', glow: 'rgba(249,115,22,0.4)' },
  'SOLO':      { color: '#14b8a6', glow: 'rgba(20,184,166,0.4)' },
  'RACE':      { color: '#eab308', glow: 'rgba(234,179,8,0.4)' },
  'OLDSCHOOL': { color: '#94a3b8', glow: 'rgba(148,163,184,0.4)' },
};

const playerData = computed(() => {
  const completed = (props.maps || []).filter(m => m.status === 2 && m.runner);

  const playerMap = {};
  completed.forEach(m => {
    const runners = m.runner.split(/[,&]/).map(r => r.trim()).filter(Boolean);
    const diff = (m.difficulty || 'OTHER').toUpperCase();
    const score = Number(m.score) || 0;
    runners.forEach(r => {
      if (!playerMap[r]) playerMap[r] = { total: 0, score: 0, diffs: {} };
      playerMap[r].total += 1;
      playerMap[r].score += score;
      playerMap[r].diffs[diff] = (playerMap[r].diffs[diff] || 0) + 1;
    });
  });

  return Object.entries(playerMap)
    .sort((a, b) => b[1].score - a[1].score)
    .map(([name, data], i) => ({
      name, rank: i + 1,
      total: data.total,
      score: data.score,
      diffs: Object.entries(data.diffs)
        .sort((a, b) => b[1] - a[1])
        .map(([diff, count]) => ({
          diff, count,
          pct: (count / data.total) * 100,
          color: DIFF_CONFIG[diff]?.color ?? '#64748b',
          glow: DIFF_CONFIG[diff]?.glow ?? 'rgba(100,116,139,0.4)',
        }))
    }));
});

const RANK_COLORS = ['#f59e0b', '#94a3b8', '#cd7c4a'];

function fmtScore(n) {
  return n >= 1000 ? (n / 1000).toFixed(1) + 'k' : String(n);
}
</script>

<template>
  <div class="bg-gradient-to-br from-gray-900 to-black border border-white/10 p-5">
    <div class="text-xs font-mono text-violet-400/70 mb-5 flex items-center gap-2">
      <span class="w-1.5 h-1.5 rounded-full bg-violet-400 animate-pulse"></span>
      PLAYER_CONTRIBUTION
      <span class="text-gray-600 ml-1">— RANKED BY SCORE</span>
    </div>

    <div v-if="!playerData.length" class="text-gray-600 font-mono text-xs text-center py-8">NO DATA</div>

    <div class="grid grid-cols-2 md:grid-cols-3 xl:grid-cols-4 gap-3">
      <div v-for="player in playerData" :key="player.name"
        class="bg-black/40 border border-white/5 rounded-lg p-4 hover:border-white/15 transition-all duration-300 group flex flex-col">

        <!-- Header: rank + name + score -->
        <div class="flex items-center gap-2 mb-3">
          <span class="text-sm font-bold font-mono w-5 flex-shrink-0 leading-none"
            :style="{ color: player.rank <= 3 ? RANK_COLORS[player.rank - 1] : '#6b7280' }">
            #{{ player.rank }}
          </span>
          <span class="text-white text-xs font-mono truncate flex-1 group-hover:text-violet-300 transition-colors">
            {{ player.name }}
          </span>
          <span class="text-violet-300 text-[10px] font-mono flex-shrink-0">{{ fmtScore(player.score) }}</span>
        </div>

        <!-- 難度比例橫條 -->
        <div class="flex h-2 rounded-full overflow-hidden gap-px mb-3">
          <div v-for="d in player.diffs" :key="d.diff"
            class="h-full transition-all duration-500 flex-shrink-0"
            :style="{ width: d.pct + '%', background: d.color, boxShadow: `0 0 4px ${d.glow}` }">
          </div>
        </div>

        <!-- 難度明細 -->
        <div class="space-y-1 flex-1">
          <div v-for="d in player.diffs" :key="d.diff"
            class="flex items-center gap-1.5 text-[10px] font-mono">
            <span class="w-1.5 h-1.5 rounded-full flex-shrink-0" :style="{ background: d.color }"></span>
            <span class="text-gray-500 flex-1 truncate">{{ d.diff }}</span>
            <span class="font-bold" :style="{ color: d.color }">{{ d.count }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
