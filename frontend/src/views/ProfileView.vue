<script setup lang="ts">
import MainLayout from '@/components/MainLayout.vue';
import UserViewMain from '@/components/UserViewMain.vue';
import Loader from '@/components/Loader.vue';

import { getMe } from '@/features/api';
import { ref } from 'vue';

const username = ref<string | null>(null);
getMe().then((me) => {
  username.value = me.user_name;
});
</script>

<template>
  <MainLayout>
    <div class="profile-loader" v-if="username === null">
      <Loader />
    </div>
    <UserViewMain :username="username" v-if="username !== null" />
  </MainLayout>
</template>

<style scoped>
.profile-loader {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}
</style>
