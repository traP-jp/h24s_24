<script setup lang="ts">
import Avatar from '@/components/Avatar.vue';
import 'moment/dist/locale/ja';
import moment from 'moment-timezone';
import { ref } from 'vue';
import { reactionIcons } from '@/features/reactions';
import { deleteReaction, postReaction } from '@/features/api';

type Reaction = { id: number; count: number; clicked: boolean };

const props = defineProps<{
  id: string;
  name: string;
  date: Date;
  content: string;
  reactions: Reaction[];
}>();
const emits = defineEmits<{
  (e: 'react'): void;
}>();

function getDateText() {
  return moment(props.date).fromNow();
}

const dateText = ref(getDateText());

async function toggleReaction(reaction: Reaction) {
  if (reaction.clicked) {
    await deleteReaction(props.id, reaction.id);
    emits('react');
  } else {
    await postReaction(props.id, reaction.id);
    emits('react');
  }
}
</script>

<template>
  <div class="post">
    <div class="post-author-icon">
      <Avatar size="48px" :name="name" />
    </div>
    <div class="post-content">
      <div class="post-header">
        <span class="post-author">@{{ name }}</span>
        <span class="post-date">{{ dateText }}</span>
      </div>
      <div class="post-message">
        {{ content }}
      </div>
      <div class="post-reactions">
        <div
          v-for="reaction in reactions"
          :key="reaction.id"
          class="post-reaction"
          :class="{ clicked: reaction.clicked }"
          @click="() => toggleReaction(reaction)"
        >
          <span class="post-reaction-icon">{{ reactionIcons[reaction.id] }}</span>
          <span class="post-reaction-count">{{ reaction.count }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.post {
  display: flex;
  padding: 16px;
  color: var(--primary-text-color);

  &-author-icon {
    padding-right: 8px;
  }

  &-content {
    width: 480px;

    .post-header {
      margin-bottom: 8px;

      .post-author {
        margin-right: 6px;
        font-weight: bold;
      }

      .post-date {
        color: var(--dimmed-text-color); // TODO
      }
    }

    .post-message {
      max-width: 100%;
      overflow-wrap: break-word;
      margin-bottom: 8px;
    }

    .post-reactions {
      display: flex;

      .post-reaction {
        margin-right: 8px;
        opacity: 40%;

        .post-reaction-icon {
          padding: 4px;
        }

        .post-reaction-count {
          color: var(--dimmed-text-color);
        }

        &.clicked {
          opacity: 100%;

          .post-reaction-count {
            color: var(--accent-color);
            font-weight: bold;
          }
        }
      }
    }
  }
}
</style>
