<script setup lang="ts">
import Button from './Button.vue';
import Avatar from './Avatar.vue';
import Loader from '@/components/Loader.vue';

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
const loading = ref(false);

const post = async () => {
  loading.value = true;
  await createPost({
    message: inputContent.value,
    parent_id: props.parentId,
  });
  loading.value = false;
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
      <div v-if="loading" class="new-post-input-section-loader">
        <Loader />
      </div>
      <textarea
        type="text"
        :placeholder="`${parentId == undefined ? '投稿' : '返信'}する内容を入力（投稿時に自動で変換されます)`"
        v-model="inputContent"
        :disabled="loading"
      />
      <div class="post-footer">
        <span :class="{ 'post-charcount-warn': !canPost }" v-if="inputContent.length > 0"
          >{{ inputContent.length }}/280文字</span
        >
        <span class="post-button">
          <Button :disabled="!canPost || loading" :onclick="post">
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
    position: relative;

    .new-post-input-section-loader {
      position: absolute;
      inset: 0;
      display: flex;
      justify-content: center;
      align-items: center;
      background-color: #fff8;
    }

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
