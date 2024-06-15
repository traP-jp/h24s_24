<script setup lang="ts">

import Button from "./Button.vue"
import Avatar from "./Avatar.vue"

</script>

<template>
    <div class="newPostSection">
        <div class="newPostInputSection">        
            <Avatar name="" size="30px"></Avatar>
            <input  type="text" placeholder="投稿する内容を入力（投稿時に自動で変換されます">
        </div>
        <div align="right">
            <Button>投稿する</Button>
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
Button {
    display: flex;
    align-items: flex-start;
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
    justify-content: flex-end;
}
</style>
