<script setup lang="ts">
import MainLayout from '@/components/MainLayout.vue';
import Post from '@/components/Post.vue';
import IntersectionObserver from '@/components/IntersectionObserver.vue';
import Loader from '@/components/Loader.vue';
import { getReactions } from '@/features/post';
import { getPosts, type Post as PostType } from '@/features/api';
import { ref } from 'vue';

const posts = ref<PostType[]>([]);
const oldestId = ref<string | undefined>(undefined);
const isEnd = ref(false);
const loading = ref(false);

const fetchNew = async () => {
  if (isEnd.value) return;
  if (loading.value) return;

  try {
    loading.value = true;
    const retrieved = await getPosts({ after: oldestId.value });
    if (retrieved.length === 0) {
      isEnd.value = true;
      return;
    }

    posts.value.push(...retrieved);
    console.log(posts.value);

    oldestId.value = retrieved.slice(-1)[0].id;
  } catch (err) {
    console.error(err);
  } finally {
    loading.value = false;
  }
};

fetchNew();
</script>

<template>
  <MainLayout>
    <div class="container">
      <div class="posts">
        <div v-for="post in posts" :key="post.id">
          <router-link :to="`/posts/${post.id}`" class="post-link">
            <Post
              :id="post.id"
              :content="post.converted_message"
              :date="new Date(post.created_at)"
              :name="post.user_name"
              :reactions="getReactions(post)"
            />
          </router-link>
        </div>
      </div>
      <IntersectionObserver @intersect="fetchNew" />
      <div class="loader-container" v-if="!isEnd">
        <Loader />
      </div>
    </div>
  </MainLayout>
</template>

<style scoped>
.container {
  padding-bottom: 50vh;
}

.posts {
  display: flex;
  flex-direction: column;
}

.post-link {
  text-decoration: none;
  color: inherit;

  &::after {
    content: '';
    display: block;
    width: 100%;
    height: 1px;
    background-color: var(--dimmed-border-color);
  }
}

.loader-container {
  display: flex;
  justify-content: center;
  padding: 16px;
}
</style>
