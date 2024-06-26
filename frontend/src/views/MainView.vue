<script setup lang="ts">
import MainLayout from '@/components/MainLayout.vue';
import Post from '@/components/Post.vue';
import IntersectionObserver from '@/components/IntersectionObserver.vue';
import Loader from '@/components/Loader.vue';
import NewPostSection from '@/components/NewPostSection.vue';

import { getPosts, type Post as PostType } from '@/features/api';
import { ref } from 'vue';
import { convertReactions } from '@/features/reactions';

const posts = ref<PostType[]>([]);
const isEnd = ref(false);
const loading = ref(false);

const fetchNew = async () => {
  try {
    const retrieved = await getPosts({ after: posts.value[0]?.id });
    posts.value.unshift(...retrieved);
    if (retrieved.length === 30) {
      await fetchNew();
    }
  } catch (err) {
    console.error(err);
  }
};

const fetchMore = async () => {
  if (isEnd.value) return;
  if (loading.value) return;

  try {
    loading.value = true;
    const retrieved = await getPosts({ before: posts.value.slice(-1)[0]?.id });
    if (retrieved.length === 0) {
      isEnd.value = true;
      return;
    }

    posts.value.push(...retrieved);
  } catch (err) {
    console.error(err);
  } finally {
    loading.value = false;
  }
};

fetchMore();
</script>

<template>
  <MainLayout>
    <div class="container">
      <NewPostSection @submit="fetchNew" />
      <hr />
      <div class="posts">
        <div v-for="post in posts" :key="post.id">
          <Post
            :id="post.id"
            :content="post.converted_message"
            :date="new Date(post.created_at)"
            :name="post.user_name"
            :originalContent="post.original_message"
            :reactions="convertReactions(post.reactions, post.my_reactions)"
          />
          <hr />
        </div>
      </div>
      <IntersectionObserver @intersect="fetchMore" />
      <div class="loader-container" v-if="!isEnd">
        <Loader />
      </div>
    </div>
  </MainLayout>
</template>

<style lang="scss" scoped>
.container {
  padding-bottom: 50vh;
}

hr {
  border: none;
  border-top: 1px solid var(--dimmed-border-color);
}

.posts {
  display: flex;
  flex-direction: column;
}

.loader-container {
  display: flex;
  justify-content: center;
  padding: 16px;
}
</style>
