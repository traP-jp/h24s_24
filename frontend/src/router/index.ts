import { createRouter, createWebHistory } from 'vue-router';
import MainView from '@/views/MainView.vue';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'main',
      component: MainView,
    },
    {
      path: '/trending',
      name: 'trending',
      component: () => import('@/views/TrendingView.vue'),
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('@/views/ProfileView.vue'),
    },
    {
      path: '/posts/:id',
      name: 'post',
      component: () => import('@/views/PostView.vue'),
    },
  ],
});

export default router;
