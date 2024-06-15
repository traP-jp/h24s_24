<script setup lang="ts">
defineProps<{
  disabled: boolean;
}>();
</script>

<template>
  <button type="button" :disabled="disabled">
    <slot />
  </button>
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
  background-color: var(--accent-color);

  &:disabled {
    background-color: var(--dimmed-border-color);
  }
}
</style>
