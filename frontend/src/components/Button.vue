<script setup lang="ts">
defineProps<{
  enabled: boolean
}>()
</script>

<template>
    <button type="button" :class="enabled ? 'enabled' : undefined"><slot></slot></button>
</template>

<style lang="scss" scoped>
button {
  color: white;
  padding: 8px 16px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: bold;
  text-align: center;
  cursor: pointer;
  &.enabled {
    background-color: var(--accent-color);
  }
}
</style>
