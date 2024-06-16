<script setup lang="ts">
import MainLayout from '@/components/MainLayout.vue';
import { useRoute } from 'vue-router';
import {getMe, getPost, type GetPostResponse} from '@/features/api';
import { ref } from 'vue';
import Post from '@/components/Post.vue';
import { convertReactions } from '@/features/reactions';
import NewPostSection from '@/components/NewPostSection.vue';

const route = useRoute();
if (route.params.id == undefined) {
  // TODO: error, id not found
}
const id = route.params.id as string;
const postContent = ref<GetPostResponse>();
const loadPost = () => {
  getPost(id).then((e) => (postContent.value = e));
};
loadPost();

const username = ref('');
getMe().then((me) => {
  username.value = me.user_name;
});
</script>

<template>
  <MainLayout>
    <div class="post-view-container">
      <div v-if="postContent != undefined">
        <div v-for="ancestor in postContent.ancestors" :key="ancestor.post.id">
          <Post
            :content="ancestor.post.converted_message"
            :date="new Date(ancestor.post.created_at)"
            :name="ancestor.post.user_name"
            :reactions="convertReactions(ancestor.post.reactions, ancestor.post.my_reactions)"
            :id="ancestor.post.id"
            @react="loadPost"
          />
        </div>
        <hr />
        <Post
          :content="postContent.converted_message"
          :date="new Date(postContent.created_at)"
          :name="postContent.user_name"
          :reactions="convertReactions(postContent.reactions, postContent.my_reactions)"
          :id="postContent.id"
          @react="loadPost"
        />
        <hr />
        <NewPostSection :name="username" :parent-id="postContent.id" @submit="loadPost" />
        <!-- TODO: -->
        <div v-for="child in postContent.children" :key="child.post.id">
          <Post
            :content="child.post.converted_message"
            :date="new Date(child.post.created_at)"
            :name="child.post.user_name"
            :reactions="convertReactions(child.post.reactions, child.post.my_reactions)"
            :id="child.post.id"
            @react="loadPost"
          />
        </div>
      </div>
    </div>
  </MainLayout>
</template>

<style scoped>
.post-view-container {
  padding-bottom: 50vh;
}
</style>
