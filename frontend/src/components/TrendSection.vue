<script lang="ts" setup>
import { reactionIcons } from '@/features/reactions';
import { vTwemojiObj } from '@/features/vTwemoji';
import { ref } from 'vue';

const emits = defineEmits<{
  (e: 'change', id: number): void;
}>();
const value = ref<number>(0);

const vTwemoji = vTwemojiObj;
</script>

<template>
  <div class="icon-area">
    <div v-for="(reaction, id) in reactionIcons" class="icon-part" :key="id">
      <input
        class="radio-button"
        :id="`radio${id}`"
        type="radio"
        name="icons"
        :value="id"
        v-model="value"
        @change="emits('change', value)"
      />
      <label :for="`radio${id}`" v-twemoji>{{ reaction }}️</label>
    </div>
  </div>
</template>

<style lang="scss" scoped>
:global(.icon-area .twemoji) {
  height: 1.2em;
  width: 1.2em;
  margin: 0 0.05em 0 0.1em;
  vertical-align: -0.1em;
}

.icon-area {
  display: flex;
}

.radio-button {
  display: none;
}

.icon-part {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
  background-color: white;
  position: relative;
  border-bottom: 1px solid lightgray;
}

label {
  cursor: pointer;
  display: block;
  padding: 10px;
  position: relative;
  width: 100%;
}

.radio-button:checked + label::after {
  content: '';
  display: block;
  position: absolute;
  bottom: 0;
  left: 50%;
  width: 100%;
  transform: translateX(-50%);
  height: 6px;
  background-color: orange;
  border-radius: 6px;
}
</style>
