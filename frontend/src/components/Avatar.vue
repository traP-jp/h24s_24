<script setup lang="ts">
defineProps<{
  name: string
}>()
const endpoint = "https://q.trap.jp/api/v3"
</script>

<template>
<img class="avatar" :src="`${endpoint}/public/icon/${name}`" />
</template>

<style scoped>
.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  margin: 10px;
}
</style>