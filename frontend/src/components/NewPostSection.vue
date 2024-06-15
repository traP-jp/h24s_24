<script setup lang="ts">

import Button from "./Button.vue"
import Avatar from "./Avatar.vue"
import {ref} from "vue"

const inputContext = ref('')
let ButtonVisible = ref(true)

function checkLength() {
  if (inputContext.value.length > 280) {
    ButtonVisible.value = false
  } else {
    ButtonVisible.value = true
  }
}

</script>

<template>
  <div class="newPostSection">
    <div class="author-icon">
      <Avatar name="" size="48px"></Avatar>
    </div>
    <div class="newPostInputSection">
      <input type="text" placeholder="投稿する内容を入力（投稿時に自動で変換されます)" v-model="inputContext"
             @input="checkLength">
      <p v-if="!ButtonVisible">文字数の上限は280文字です</p>
      <div class="post-button" align="right">
        <Button id="canPost" :enabled="ButtonVisible">投稿する</Button>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.newPostSection {
  width: 480px;
  padding: 16px;
  display: flex;

  .newPostInputSection {
    width: 100%;
    margin-left: 8px;

    input {
      display: block;
      width: 100%;
      border: none;
      padding-top: 16px;
      padding-bottom: 8px;

      &:focus-visible {
        outline: none;
      }
    }
    .post-button {
      margin: 8px;
    }
  }
}

</style>
