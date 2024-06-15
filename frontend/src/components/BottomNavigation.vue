<script setup lang="ts">
import { Icon } from '@iconify/vue';
</script>

<template>
  <nav class="bottom-navigation">
    <router-link to="/" class="bottom-navigation-link" exact-active-class="current">
      <Icon icon="mdi:home" class="bottom-navigation-icon" />
      <span>ホーム</span>
    </router-link>
    <router-link to="/trending" class="bottom-navigation-link" exact-active-class="current">
      <Icon icon="mdi:thunder" class="bottom-navigation-icon" />
      <span>トレンド</span>
    </router-link>
    <router-link to="/profile" class="bottom-navigation-link" exact-active-class="current">
      <Icon icon="mdi:account" class="bottom-navigation-icon" />
      <span>プロフィール</span>
    </router-link>
  </nav>
</template>

<style scoped lang="scss">
.bottom-navigation {
  display: flex;
  justify-content: space-around;
  align-items: center;
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background-color: #fff;
  border-top: 1px solid var(--dimmed-border-color);
}

.bottom-navigation-link {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 0.5rem;
  color: #333;
  text-decoration: none;

  span {
    font-size: 0.75rem;
  }

  &.current {
    color: var(--accent-color);

    span {
      font-weight: bold;
    }
  }

  .bottom-navigation-icon {
    font-size: 1.25rem;
  }
}
</style>
