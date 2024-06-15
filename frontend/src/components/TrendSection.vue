<script setup lang="ts">
import Avatar from './Avatar.vue'
import { ref } from 'vue'

</script>

<template>
    <div class="iconArea" >
        <input class="radioButton" id="radio0" type="radio" name="icons" value="0">
        <label for="radio0">❤️</label>
        <input class="radioButton" id="radio1" type="radio" name="icons" value="1">
        <label for="radio1">🔥</label>
        <input class="radioButton" id="radio2" type="radio" name="icons" value="2">
        <label for="radio2">💧</label>
        <input class="radioButton" id="radio3" type="radio" name="icons" value="3">
        <label for="radio3">🥲</label>
        <input class="radioButton" id="radio4" type="radio" name="icons" value="4">
        <label for="radio4">🤔</label>
    </div>
</template>

<style lang="scss">
.radioButton {
    display: none;
}
label {
  background-color: white;
  padding: 10px 20px;
}
.radioButton:checked + label {
  background-color: pink;
}

</style>