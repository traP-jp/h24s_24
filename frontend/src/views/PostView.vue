<script setup lang="ts">
import MainLayout from '@/components/MainLayout.vue';
import { useRoute } from 'vue-router';
import { getPost, type GetPostResponse, type Reaction } from '@/features/api';
import { ref } from 'vue';
import Post from '@/components/Post.vue';
import { reactionIcons } from '@/features/reactions';
import NewPostSection from '@/components/NewPostSection.vue';

const route = useRoute();
if (route.params.id == undefined) {
  // TODO: error, id not found
}
const id = route.params.id as string;
const postContent = ref<GetPostResponse>();
getPost(id).then((e) => (postContent.value = e));

const convertReactions = (src: Reaction[], my: number[]) => {
  const dist: { id: number; count: number; clicked: boolean }[] = [];
  for (let i = 0; i < reactionIcons.length; i++) {
    const found = src.find((r) => r.id == i);
    if (found) {
      dist.push({
        id: i,
        count: found.count,
        clicked: my.find((m) => m == i) != undefined,
      });
    } else {
      dist.push({
        id: i,
        count: 0,
        clicked: false,
      });
    }
  }
  return dist;
};
</script>

<template>
  <MainLayout>
    <div v-if="postContent != undefined">
      <div v-for="ancestor in postContent.ancestors" :key="ancestor.post.id">
        <Post
          :content="ancestor.post.converted_message"
          :date="new Date(ancestor.post.created_at)"
          :name="ancestor.post.user_name"
          :reactions="convertReactions(ancestor.post.reactions, ancestor.post.my_reactions)"
          :id="ancestor.post.id"
        />
      </div>
      <hr />
      <Post
        :content="postContent.converted_message"
        :date="new Date(postContent.created_at)"
        :name="postContent.user_name"
        :reactions="convertReactions(postContent.reactions, postContent.my_reactions)"
        :id="postContent.id"
      />
      <hr />
      <NewPostSection name="" :parent-id="postContent.id" />
      <!-- TODO: -->
      <div v-for="child in postContent.children" :key="child.post.id">
        <Post
          :content="child.post.converted_message"
          :date="new Date(child.post.created_at)"
          :name="child.post.user_name"
          :reactions="convertReactions(child.post.reactions, child.post.my_reactions)"
          :id="child.post.id"
        />
      </div>
    </div>
  </MainLayout>
</template>
