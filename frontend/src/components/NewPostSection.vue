<script setup lang="ts">
import Button from "./Button.vue"
import Avatar from "./Avatar.vue"
import {ref, computed} from "vue"

const props = defineProps<{
  name: string,
}>()

const inputContent = ref('')
const canPost = computed(() => {
  return (inputContent.value.length != 0) && (inputContent.value.length <= 280)
})

</script>

<template>
  <div class="new-post-section">
    <div class="author-icon">
      <Avatar :name="name" size="48px"></Avatar>
    </div>
    <div class="new-post-input-section">
      <input type="text" placeholder="投稿する内容を入力（投稿時に自動で変換されます)" v-model="inputContent">
      <div class="post-footer">
        <span v-if="!canPost">文字数の上限は280文字です</span>
        <span class="post-button">
          <Button :disabled="!canPost">投稿する</Button>
        </span>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.new-post-section {
  width: 100%;
  padding: 16px;
  display: flex;

  .new-post-input-section {
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

    .post-footer {
      margin: 8px;
      text-align: right;

      .post-button {
        margin-left: 20px;
      }
    }
  }
}

</style>
