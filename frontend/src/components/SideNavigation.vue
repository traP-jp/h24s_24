<script setup lang="ts">
import { Icon } from '@iconify/vue';
</script>

<template>
  <nav class="side-navigation">
    <router-link to="/" class="side-navigation-link" exact-active-class="current">
      <Icon icon="mdi:home" class="side-navigation-icon" />
      <span>ホーム</span>
    </router-link>
    <router-link to="/trending" class="side-navigation-link" exact-active-class="current">
      <Icon icon="mdi:thunder" class="side-navigation-icon" />
      <span>トレンド</span>
    </router-link>
    <router-link to="/profile" class="side-navigation-link" exact-active-class="current">
      <Icon icon="mdi:account" class="side-navigation-icon" />
      <span>プロフィール</span>
    </router-link>
  </nav>
</template>

<style scoped lang="scss">
.side-navigation {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  position: fixed;
  top: 0;
  left: calc(50% - 300px);
  transform: translateX(-100%);
  bottom: 0;
  background-color: #fff;
  padding-top: 16px;
}

.side-navigation-link {
  display: flex;
  gap: 8px;
  align-items: center;
  justify-content: center;
  padding: 0.5rem 1rem;
  color: var(--primary-text-color);
  text-decoration: none;
  border-radius: 8px;
  margin-right: 16px;
  transition: background-color 0.2s;

  &:hover {
    background-color: #0001;
  }

  .side-navigation-icon {
    font-size: 1.5rem;
  }

  &.current {
    color: var(--accent-color);

    span {
      font-weight: bold;
    }
  }
}

@media screen and (max-width: 1024px) {
  .side-navigation {
    .side-navigation-link {
      padding: 0.5rem;
    }

    span {
      display: none;
    }
  }
}
</style>
