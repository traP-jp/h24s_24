<script setup lang="ts">
import Avatar from "@/components/Avatar.vue";

const props = defineProps<{
  name: string,
  date: Date,
  content: string,
  stamp_counts: number[]
}>()

</script>

<template>
  <div class="post">
    <div class="post-author-icon">
      <Avatar size="50px" :name="name"/>
    </div>
    <div class="post-content">
      <div class="post-header">
        <span class="post-author">@{{ name }}</span>
        <span class="post-date">10分前</span>
      </div>
      <div class="post-message">
        {{ content }}
      </div>
      <div class="post-stamps">
        <div class="post-stamp">
          <span class="post-stamp-icon">❤️</span>
          <span class="post-stamp-count">{{ stamp_counts[0] }}</span>
        </div>
        <div class="post-stamp">
          <span class="post-stamp-icon">🔥</span>
          <span class="post-stamp-count">{{ stamp_counts[1] }}</span>
        </div>
        <div class="post-stamp">
          <span class="post-stamp-icon">💧</span>
          <span class="post-stamp-count">{{ stamp_counts[2] }}</span>
        </div>
        <div class="post-stamp">
          <span class="post-stamp-icon">😢</span>
          <span class="post-stamp-count">{{ stamp_counts[3] }}</span>
        </div>
        <div class="post-stamp clicked">
          <span class="post-stamp-icon">🤔</span>
          <span class="post-stamp-count">{{ stamp_counts[4] }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
$disabled: #ddd;
$highlight: orange;
.post {
  display: flex;
  padding: 10px;

  &-author-icon {
    padding-right: 10px;
  }

  &-content {
    width: 450px;

    .post-header {
      margin-bottom: 5px;

      .post-author {
        margin-right: 5px;
        font-weight: bold;
      }

      .post-date {
        color: gray; // TODO
      }
    }

    .post-message {
      max-width: 100%;
      overflow-wrap: break-word;
    }

    .post-stamps {
      display: flex;
      margin-top: 5px;

      .post-stamp {
        padding: 0 5px;

        .post-stamp-icon {
          padding-right: 5px;
        }

        .post-stamp-count {
          color: $disabled;
        }

        &.clicked {
          .post-stamp-count {
            color: $highlight;
            font-weight: bold;
          }
        }
      }
    }
  }
}
</style>
