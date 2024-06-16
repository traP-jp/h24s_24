<script setup lang="ts">
import MainLayout from '@/components/MainLayout.vue';
import { onBeforeRouteUpdate, useRoute } from 'vue-router';
import { getPost, type GetPostResponse } from '@/features/api';
import { ref } from 'vue';
import Post from '@/components/Post.vue';
import NewPostSection from '@/components/NewPostSection.vue';
import { convertReactions } from '@/features/reactions';

const route = useRoute();
if (route.params.id == undefined) {
  // TODO: error, id not found
}
// const id = route.params.id as string;
const postContent = ref<GetPostResponse>();
const loadPost = (id: string) => {
  getPost(id).then((e) => (postContent.value = e));
};
loadPost(useRoute().params.id as string);

onBeforeRouteUpdate((to) => {
  loadPost(to.params.id as string);
});
</script>

<template>
  <MainLayout>
    <div class="post-view-container">
      <div v-if="postContent != undefined">
        <div
          v-for="ancestor in postContent.ancestors"
          :key="ancestor.post.id"
          class="ancestor-post"
        >
          <Post
            :content="ancestor.post.converted_message"
            :originalContent="ancestor.post.original_message"
            :date="new Date(ancestor.post.created_at)"
            :name="ancestor.post.user_name"
            :reactions="convertReactions(ancestor.post.reactions, ancestor.post.my_reactions)"
            :id="ancestor.post.id"
            @react="() => loadPost(useRoute().params.id as string)"
          />
          <div class="ancestor-bar" />
        </div>
        <Post
          :content="postContent.converted_message"
          :originalContent="postContent.original_message"
          :date="new Date(postContent.created_at)"
          :name="postContent.user_name"
          :reactions="convertReactions(postContent.reactions, postContent.my_reactions)"
          :id="postContent.id"
          detail
          @react="() => loadPost(useRoute().params.id as string)"
        />
        <hr />
        <NewPostSection
          name=""
          :parent-id="postContent.id"
          @submit="() => loadPost(useRoute().params.id as string)"
        />
        <hr />
        <div v-for="child in postContent.children" :key="child.post.id">
          <Post
            :content="child.post.converted_message"
            :originalContent="child.post.original_message"
            :date="new Date(child.post.created_at)"
            :name="child.post.user_name"
            :reactions="convertReactions(child.post.reactions, child.post.my_reactions)"
            :id="child.post.id"
            @react="() => loadPost(useRoute().params.id as string)"
          />
          <hr />
        </div>
      </div>
    </div>
  </MainLayout>
</template>

<style scoped>
.post-view-container {
  padding-bottom: 50vh;
}

hr {
  border: none;
  border-top: 1px solid var(--dimmed-border-color);
}

.ancestor-post {
  position: relative;
}

.ancestor-bar {
  position: absolute;
  top: 72px;
  left: 36px;
  width: 4px;
  height: calc(100% - 68px);
  border-radius: 2px;
  background-color: var(--dimmed-border-color);
  z-index: -1;
}
</style>
