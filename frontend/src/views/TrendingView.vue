<script setup lang="ts">
import MainLayout from '@/components/MainLayout.vue';
import TrendSection from '@/components/TrendSection.vue';
import { ref } from 'vue';
import Post from '@/components/Post.vue';
import { getTrend, type GetTrendResponse } from '@/features/api';
import { convertReactions } from '@/features/reactions';

const target = ref<number>(0);
const timeline = ref<GetTrendResponse>();
const loadPost = () => {
  getTrend(target.value).then((e) => (timeline.value = e));
};
loadPost();
</script>

<template>
  <MainLayout>
    <TrendSection @change="(id) => {
        target = id;
        loadPost();
      }
      " />
    <div v-for="post in timeline" :key="post.id">
      <Post class="trending-post" :content="post.converted_message" :date="new Date(post.created_at)"
        :name="post.user_name" :reactions="convertReactions(post.reactions, post.my_reactions)" :id="post.id"
        @react="loadPost" />
    </div>
  </MainLayout>
</template>

<style lang="scss" scoped>
.trending-post {
  border-bottom: 1px solid var(--dimmed-border-color);
}
</style>
