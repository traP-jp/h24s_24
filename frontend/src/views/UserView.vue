<script setup lang="ts">
import UserViewMain from '@/components/UserViewMain.vue';
</script>

<template>
  <UserViewMain :username="$route.params.username as string" />
</template>
