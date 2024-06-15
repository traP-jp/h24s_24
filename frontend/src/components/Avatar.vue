<script setup lang="ts">
defineProps<{
  name: string,
  size: string
}>()
const endpoint = "https://q.trap.jp/api/v3"
</script>

<template>
<img class="avatar" :src="`${endpoint}/public/icon/${name}`" :style="{width: size, height: size}" />
</template>

<style scoped>
.avatar {
  border-radius: 50%;
}
</style>
