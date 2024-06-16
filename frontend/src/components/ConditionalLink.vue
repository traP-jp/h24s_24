<script setup lang="ts">
defineProps<{
  to: string;
  condition: boolean;
}>()
</script>

<template>
  <router-link v-if="condition" :to="to">
    <slot />
  </router-link>
  <template v-else>
    <slot />
  </template>
</template>
