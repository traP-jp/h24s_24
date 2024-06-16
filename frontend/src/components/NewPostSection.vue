<script setup lang="ts">
import Button from './Button.vue';
import Avatar from './Avatar.vue';
import { ref, computed } from 'vue';
import { createPost } from '@/features/api';

const props = defineProps<{
  name: string;
  parentId?: string;
}>();
const emit = defineEmits<{
  (e: 'submit'): void;
}>();

const inputContent = ref('');
const canPost = computed(() => {
  return inputContent.value.length != 0 && inputContent.value.length <= 280;
});
const post = () => {
  createPost({
    message: inputContent.value,
    parent_id: props.parentId,
  });
  emit('submit');
  inputContent.value = '';
};
</script>

<template>
  <div class="new-post-section">
    <div class="author-icon">
      <Avatar :name="name" size="48px"></Avatar>
    </div>
    <div class="new-post-input-section">
      <textarea
        type="text"
        :placeholder="`${parentId == undefined ? '投稿' : '返信'}する内容を入力（投稿時に自動で変換されます)`"
        v-model="inputContent"
      />
      <div class="post-footer">
        <span :class="{ 'post-charcount-warn': !canPost }">{{ inputContent.length }}/280文字</span>
        <span class="post-button">
          <Button :disabled="!canPost" :onclick="post">
            {{ parentId == undefined ? '投稿' : '返信' }}する</Button
          >
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
    width: calc(100% - 56px);
    margin-left: 8px;

    textarea {
      display: block;
      width: 100%;
      max-width: 100%;
      border: none;
      padding-top: 16px;
      padding-bottom: 8px;
      resize: none;
      field-sizing: content;

      &:focus-visible {
        outline: none;
      }
    }

    .post-footer {
      margin: 8px;
      text-align: right;

      span.post-charcount-warn {
        color: red;
      }

      .post-button {
        margin-left: 20px;
      }
    }
  }
}
</style>
