<script setup lang="ts">
import BottomNavigation from '@/components/BottomNavigation.vue';
import SideNavigation from '@/components/SideNavigation.vue';
</script>

<template>
  <div class="main-layout">
    <SideNavigation class="side-navigation" />
    <main>
      <slot />
    </main>
    <BottomNavigation class="bottom-navigation" />
  </div>
</template>

<style scoped lang="scss">
.main-layout {
  height: 100vh;
  position: relative;
}

main {
  height: 100vh;
  margin: 0 auto;
}

.side-navigation {
  visibility: hidden;
}

@media screen and (min-width: 768px) {
  main {
    max-width: 600px;
    border-inline: 1px solid var(--dimmed-border-color);
  }

  .side-navigation {
    visibility: visible;
  }

  .bottom-navigation {
    visibility: hidden;
  }
}
</style>
