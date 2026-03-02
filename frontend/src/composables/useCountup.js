import { ref, watch } from 'vue';

export function useCountup(source, duration = 800) {
  const display = ref(typeof source.value === 'number' ? source.value : 0);

  watch(source, (newVal, oldVal) => {
    const start = oldVal ?? 0;
    const end = newVal ?? 0;
    if (start === end) return;

    const startTime = performance.now();

    const easeOutQuart = (t) => 1 - Math.pow(1 - t, 4);

    const tick = (now) => {
      const elapsed = now - startTime;
      const progress = Math.min(elapsed / duration, 1);
      display.value = Math.round(start + (end - start) * easeOutQuart(progress));
      if (progress < 1) requestAnimationFrame(tick);
    };

    requestAnimationFrame(tick);
  });

  return display;
}
