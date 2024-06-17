<script setup lang="ts">
import Avatar from '@/components/Avatar.vue';
import 'moment/dist/locale/ja';
import moment from 'moment-timezone';
import { computed, ref } from 'vue';
import ConditionalLink from '@/components/ConditionalLink.vue';
import { Icon } from '@iconify/vue';
import type { Reaction } from '@/components/PostReactions.vue';
import PostReactions from '@/components/PostReactions.vue';

const props = defineProps<{
  id: string;
  name: string;
  date: Date;
  content: string;
  originalContent: string;
  detail?: boolean;
  reactions: Reaction[];
}>();
const emit = defineEmits<{
  (e: 'react'): void;
}>();

function getDateText() {
  return moment(props.date).fromNow();
}

const dateText = ref(getDateText());

const shareText = computed(() =>
  encodeURIComponent(
    `:@${props.name}: < ${props.content}\n\n[:fire: 発火村の投稿より :fire:](https://hakka-mura.trap.show/posts/${props.id})`,
  ),
);
</script>

<template>
  <ConditionalLink :condition="!detail" :to="`/posts/${id}/`" class="post-link">
    <div class="post">
      <router-link :to="`/users/${name}`" class="post-author-icon">
        <Avatar size="48px" :name="name" />
      </router-link>
      <div class="post-content">
        <div class="post-header">
          <router-link :to="`/users/${name}`" class="post-author">@{{ name }}</router-link>
          <span class="post-date">{{ dateText }}</span>
        </div>
        <div class="post-message-container">
          <div class="post-message">
            {{ content }}
          </div>
          <div v-if="!detail" class="original-message">{{ originalContent }}</div>
          <div v-if="detail" class="detail-original-message">{{ originalContent }}</div>
        </div>
        <div class="foot-action-container">
          <PostReactions :postId="id" :reactions="reactions" @react="emit('react')" />
          <div class="share-traQ" title="traQで共有">
            <a
              :href="`https://q.trap.jp/share-target?text=${shareText}`"
              target="_blank"
              rel="noreferrer noopener"
              @click="(e) => e.stopPropagation()"
            >
              <Icon icon="mdi:share-variant" />
            </a>
          </div>
        </div>
      </div>
    </div>
  </ConditionalLink>
</template>

<style lang="scss" scoped>
.post-link {
  text-decoration: none;
  color: inherit;
}

.post {
  display: flex;
  padding: 16px;
  color: var(--primary-text-color);
}

.post-author-icon {
  padding-right: 8px;
}

.post-content {
  flex: 1;
  min-width: 0;
}

.post-header {
  margin-bottom: 8px;

  .post-author {
    margin-right: 6px;
    font-weight: bold;
    text-decoration: none;
    color: inherit;

    &:hover {
      text-decoration: underline;
    }
  }

  .post-date {
    color: var(--dimmed-text-color);
  }
}

.post-message-container {
  position: relative;
}

.post-message {
  max-width: 100%;
  overflow-wrap: break-word;
  margin-bottom: 8px;
  position: relative;
}

.original-message {
  font-size: 11px;
  position: absolute;
  padding: 8px 16px;
  border-radius: 8px;
  background-color: #000a;
  color: white;
  visibility: hidden;
  opacity: 0%;
  transition: all 0.2s ease-out;
  z-index: 1;
  bottom: 12px;
  left: 0;
  transform: translateY(100%);
}

.post-message:hover + .original-message {
  visibility: visible;
  opacity: 100%;
  bottom: -4px;
}

.detail-original-message {
  color: var(--dimmed-text-color);
  margin-bottom: 8px;
}

.foot-action-container {
  display: flex;
  justify-content: space-between;
}

.share-traQ {
  padding-right: 16px;

  a {
    display: inline-grid;
    place-items: center;
    width: 32px;
    height: 32px;
    font-size: 1.2rem;
    border-radius: 50%;
    transition: all 0.2s;
    color: var(--dimmed-text-color);

    &:hover {
      background-color: #0001;
      color: var(--primary-text-color);
    }
  }
}
</style>
