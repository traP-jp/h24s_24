<script setup lang="ts">
import Avatar from './Avatar.vue'
import { ref } from 'vue'
</script>

<template>
    <div class="icon-area" >
        <input class="radio-button" id="radio0" type="radio" name="icons" value="iine">
        <label for="radio0">❤️</label>
        <input class="radio-button" id="radio1" type="radio" name="icons" value="enjyou">
        <label for="radio1">🔥</label>
        <input class="radio-button" id="radio2" type="radio" name="icons" value="water">
        <label for="radio2">💧</label>
        <input class="radio-button" id="radio3" type="radio" name="icons" value="sad">
        <label for="radio3">🥲</label>
        <input class="radio-button" id="radio4" type="radio" name="icons" value="thinking">
        <label for="radio4">🤔</label>
    </div>
</template>

<style lang="scss">
.icon-area{
    display: flex;
}
    .radio-button {
        display: none;
    }
    label {
    flex:1;
    display: flex;
    align-items: center;
    justify-content: center;
    text-align: center;
    background-color: white;
    padding: 10px 0;
    position: relative;
    }
    .radio-button:checked + label::after{
    content: '';
    display: block;
    position: absolute;
    bottom: 0;
    left: 0;
    width: 100%;
    height: 6px;
    background-color: orange;
    border-radius: 6px;
    }
</style>