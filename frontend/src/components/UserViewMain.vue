<script setup lang="ts">
import Loader from '@/components/Loader.vue';
import MainLayout from '@/components/MainLayout.vue';
import Avatar from '@/components/Avatar.vue';
import Post from '@/components/Post.vue';

import { getReactions } from '@/features/post';
import { getUser } from '@/features/api';
import { useFetcher } from '@/features/useFetcher';

const props = defineProps<{
  username: string;
}>();

const { data, loading } = useFetcher(() => getUser(props.username));
</script>

<template>
  <MainLayout>
    <div v-if="!loading" class="user-view">
      <div class="user-view-header">
        <div class="user-view-header-icon">
          <Avatar size="96px" :name="username" />
        </div>
        <div class="user-view-header-content">
          <div class="user-view-header-name">@{{ username }}</div>
          <div class="user-view-header-stat-container">
            <div class="user-view-header-stat">
              <span class="user-view-header-stat-count">{{ data?.post_count }}</span>
              <span class="user-view-header-stat-label">投稿</span>
            </div>
            <div class="user-view-header-stat">
              <span class="user-view-header-stat-count">{{ data?.reaction_count }}</span>
              <span class="user-view-header-stat-label">反応した</span>
            </div>
            <div class="user-view-header-stat">
              <span class="user-view-header-stat-count">{{ data?.get_reaction_count }}</span>
              <span class="user-view-header-stat-label">反応された</span>
            </div>
          </div>
        </div>
      </div>
      <div class="user-view-posts">
        <Post v-for="post in data?.posts" :key="post.id" :name="post.user_name" :date="new Date(post.created_at)"
          :content="post.converted_message" :reactions="getReactions(post)" />
      </div>
    </div>
    <div v-if="loading" class="user-view-loader">
      <Loader />
    </div>
  </MainLayout>
</template>

<style lang="scss" scoped>
.user-view {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.user-view-header {
  display: flex;
  gap: 16px;
  align-items: center;
  border-bottom: 1px solid var(--dimmed-border-color);
  padding: 16px;
}

.user-view-header-icon {
  flex-shrink: 0;
}

.user-view-header-content {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.user-view-header-name {
  font-size: 20px;
}

.user-view-header-stat-container {
  display: flex;
  gap: 16px;
}

.user-view-header-stat {
  display: flex;
  gap: 4px;
}

.user-view-header-stat-count {
  font-size: 14px;
  font-weight: bold;
}

.user-view-header-stat-label {
  font-size: 14px;
  color: var(--dimmed-text-color);
}

.user-view-posts {
  display: grid;
  gap: 16px;
}

.user-view-loader {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}
</style>
