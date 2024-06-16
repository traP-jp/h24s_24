<script setup lang="ts">
import Avatar from '@/components/Avatar.vue';
import 'moment/dist/locale/ja';
import moment from 'moment-timezone';
import { effect, ref } from 'vue';
import { reactionIcons } from '@/features/reactions';
import { deleteReaction, postReaction } from '@/features/api';
import twemoji from 'twemoji'

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

const copiedReactions = ref(props.reactions);
effect(() => {
  copiedReactions.value = props.reactions;
});
const newReaction = ref<number>(-1);

function getDateText() {
  return moment(props.date).fromNow();
}

const dateText = ref(getDateText());

async function toggleReaction(reaction: Reaction) {
  const r = copiedReactions.value.find((r) => r.id == reaction.id)!;
  if (reaction.clicked) {
    r.clicked = false;
    r.count--;
    await deleteReaction(props.id, reaction.id);
    emits('react');
  } else {
    r.clicked = true;
    r.count++;
    newReaction.value = reaction.id;
    await postReaction(props.id, reaction.id);
    emits('react');
  }
}

const vTwemoji = {
  mounted: (el: HTMLElement) => {
    el.innerHTML = twemoji.parse(el.innerHTML, {
      className: 'twemoji',
    });
  },
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
        <button v-for="reaction in copiedReactions" :key="reaction.id" class="post-reaction"
          :class="{ clicked: reaction.clicked, ripple: newReaction === reaction.id }"
          @click="() => toggleReaction(reaction)">
          <span class="post-reaction-icon" v-twemoji>{{ reactionIcons[reaction.id] }}</span>
          <span class="post-reaction-count">{{ reaction.count }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
:global(.twemoji) {
  height: 1em;
  width: 1em;
  margin: 0 .05em 0 .1em;
  vertical-align: -0.1em;
}

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
      gap: 8px;

      .post-reaction {
        background-color: inherit;
        border: none;
        border-radius: 8px;
        display: flex;
        padding: 4px 12px 4px 4px;
        font-size: 1rem;
        display: flex;
        align-items: center;
        position: relative;
        transition: background-color 0.2s;

        &.ripple::before {
          content: '';
          position: absolute;
          width: 100%;
          aspect-ratio: 1/1;
          background-color: var(--accent-color);
          border-radius: 50%;
          z-index: -1;
          animation: ripple 0.5s ease-out forwards;
        }

        &>* {
          opacity: 40%;
        }

        &:hover {
          background-color: #0001;
        }

        .post-reaction-icon {
          padding: 0 4px;
        }

        .post-reaction-count {
          color: var(--dimmed-text-color);
        }

        &.clicked {
          background-color: var(--accent-color-10);

          &>* {
            opacity: 100%;
          }

          .post-reaction-count {
            color: var(--accent-color);
            font-weight: bold;
          }
        }
      }
    }
  }
}

@keyframes ripple {
  from {
    opacity: 0.2;
    transform: scale(0.5);
  }

  to {
    opacity: 0;
    transform: scale(1.5);
  }
}
</style>
