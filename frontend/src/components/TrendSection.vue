<script setup lang="ts">
import Avatar from './Avatar.vue'
import { ref } from 'vue'

const buttonCls = ref('reaction-button')
function changeStatus(){
    
}
</script>

<template>
    <div class="trendSection" >
        <button :class="buttonCls" @click="changeStatus"> ❤️ </button>
        <button :class="buttonCls"> 🔥 </button>
        <button :class="buttonCls"> 💧 </button>
        <button :class="buttonCls"> 😭 </button>
        <button :class="buttonCls"> 🤔 </button>
    </div>
</template>

<style lang="scss">
.reaction-button{
    width: 50px;
    border:none;
}
.selected-button{
    width:50px;
    background-color:;
}

</style>