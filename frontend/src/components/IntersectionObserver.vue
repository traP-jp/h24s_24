<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';

const props = defineProps<{
  onIntersect: (isIntersecting: boolean) => void;
}>()

const target = ref(null);
const isIntersecting = ref(false);

const callback: IntersectionObserverCallback = (entries) => {
  entries.forEach(entry => {
    isIntersecting.value = entry.isIntersecting;
    console.log(isIntersecting.value);

    if (entry.isIntersecting) {
      props.onIntersect(isIntersecting.value);
    }
  });
};

onMounted(() => {
  const observer = new IntersectionObserver(callback, {
    root: null, // viewportをルートとする
    rootMargin: '0px',
    threshold: 0.1 // 10% が画面内に入ったらコールバックを呼ぶ
  });

  if (target.value) {
    observer.observe(target.value);
  }

  onUnmounted(() => {
    if (target.value) {
      observer.unobserve(target.value);
    }
  });
});
</script>

<template>
  <div ref="target" class="observed-element">
    <slot></slot>
  </div>
</template>

<style scoped>
.observed-element {
  height: 1px;
}
</style>
