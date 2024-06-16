<script lang="ts" setup>
import {reactionIcons} from "@/features/reactions";
import {ref} from "vue";

const emits = defineEmits<{
  (e: 'change', id: number): void;
}>();
const value = ref<number>(0);
</script>

<template>
  <div class="icon-area">
    <div v-for="(reaction, id) in reactionIcons" class="icon-part" :key="id">
      <input class="radio-button" :id="`radio${id}`" type="radio" name="icons" :value="id" v-model="value" @change="emits('change', value)">
      <label :for="`radio${id}`">{{ reaction }}️</label>
    </div>
  </div>
</template>

<style lang="scss" scoped>
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
  padding: 10px 0;
  position: relative;
  border-bottom: 1px solid lightgray;
}

.radio-button:checked + label::after {
  content: '';
  display: block;
  position: absolute;
  bottom: 0;
  left: 50%;
  width: 80%;
  transform: translateX(-50%);
  height: 6px;
  background-color: orange;
  border-radius: 6px;
}
</style>