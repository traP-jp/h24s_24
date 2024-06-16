<script setup lang="ts">
import { vTwemojiObj } from '@/features/vTwemoji';
import { deleteReaction, postReaction } from '@/features/api';
import { ref, effect } from 'vue';
import { reactionIcons } from '@/features/reactions';

export type Reaction = { id: number; count: number; clicked: boolean };
const props = defineProps<{
  postId: string;
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

async function toggleReaction(reaction: Reaction) {
  const r = copiedReactions.value.find((r) => r.id == reaction.id)!;
  if (reaction.clicked) {
    r.clicked = false;
    r.count--;
    await deleteReaction(props.postId, reaction.id);
    emits('react');
  } else {
    r.clicked = true;
    r.count++;
    newReaction.value = reaction.id;
    await postReaction(props.postId, reaction.id);
    emits('react');
  }
}

const vTwemoji = vTwemojiObj;
</script>

<template>
  <div class="post-reactions">
    <button
      v-for="reaction in copiedReactions"
      :key="reaction.id"
      class="post-reaction"
      :class="{ clicked: reaction.clicked, ripple: newReaction === reaction.id }"
      @click="
        (e) => {
          toggleReaction(reaction);
          e.stopPropagation();
          e.preventDefault();
        }
      "
    >
      <span class="post-reaction-icon" v-twemoji>{{ reactionIcons[reaction.id] }}</span>
      <span class="post-reaction-count">{{ reaction.count }}</span>
    </button>
  </div>
</template>

<style lang="scss" scoped>
.post-reactions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

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

  & > * {
    opacity: 40%;
  }

  &:hover {
    background-color: #0001;
  }

  &.clicked {
    background-color: var(--accent-color-10);

    & > * {
      opacity: 100%;
    }

    .post-reaction-count {
      color: var(--accent-color);
      font-weight: bold;
    }
  }
}

.post-reaction-icon {
  padding: 0 4px;
}

:global(.post-reaction-icon .twemoji) {
  height: 1em;
  width: 1em;
  margin: 0 0.05em 0 0.1em;
  vertical-align: -0.1em;
}

.post-reaction-count {
  color: var(--dimmed-text-color);
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
