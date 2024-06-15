<script setup lang="ts">

import Button from "./Button.vue"
import Avatar from "./Avatar.vue"
import { ref } from "vue"

const inputContext = ref('')
let ButtonVisible = ref(true)

function checkLength() {
    if (inputContext.value.length > 10){
        ButtonVisible.value = false
    } else {
        ButtonVisible.value = true
    }
}

</script>

<template>
    <div class="newPostSection">
        <div class="newPostInputSection">        
            <Avatar name="" size="30px"></Avatar>
            <input  type="text" placeholder="投稿する内容を入力（投稿時に自動で変換されます)" v-model="inputContext" @input="checkLength">
            <p v-if="!ButtonVisible">文字数の上限は280文字です</p>
        </div>
        <div align="right">
            <Button v-if="ButtonVisible" id="canPost">投稿する</Button>
            <Button v-if="!ButtonVisible" id="cannotPost">投稿する</Button>
        </div>
    </div>
</template>

<style>
.newPostSection{
    width: 500px;
    height: 300px;
    background-color: white;
    margin: 10px
}
.newPostInputSection{
    width:250px;
    height:50px;
    background-color:white;
}
#canPost {
    background-color: #F4A344;
    color: white; 
    padding: 8px 16px;
    border: none;
    border-radius: 10px;
    font-size: 14px;
    font-weight: bold;
    text-align: center;
    cursor: pointer;
    margin: 10px;
}
#cannotPost {
    background-color: gainsboro;
    color: white; 
    padding: 8px 16px;
    border: none;
    border-radius: 10px;
    font-size: 14px;
    font-weight: bold;
    text-align: center;
    cursor: pointer;
    margin: 10px;
}
</style>
