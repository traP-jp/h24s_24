<script setup lang="ts">
defineProps<{
  name: string,
  size: string
}>()
const endpoint = "https://q.trap.jp/api/v3"
</script>

<template>
<img class="avatar" :src="`${endpoint}/public/icon/${name}`" />
</template>

<style scoped>
.avatar {
  width: v-bind(size);
  height: v-bind(size);
  border-radius: 50%;
}
</style>
