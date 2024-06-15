<script setup lang="ts">
import Avatar from './Avatar.vue'
import { ref } from 'vue'

const buttonCls = ref('reaction-button')

</script>

<template>
    <div class="trendSection" >
        <button :class="buttonCls"> ❤️ </button>
        <button :class="buttonCls"> 🔥 </button>
        <button :class="buttonCls"> 💧 </button>
        <button :class="buttonCls"> 😭 </button>
        <button :class="buttonCls"> 🤔 </button>
    </div>
</template>

<style>
.reaction-button{
    width: 50px;
    border:none;
}
.selected-button{
    width:50px;
}

</style>