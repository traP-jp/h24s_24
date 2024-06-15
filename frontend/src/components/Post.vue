<script setup lang="ts">
import Avatar from "@/components/Avatar.vue";
import 'moment/dist/locale/ja';
import moment from 'moment-timezone';
import {ref} from "vue";

const props = defineProps<{
  name: string,
  date: Date,
  content: string,
  stamp_counts: number[],
  stamp_clicked: boolean[]
}>()

function getDateText() {
  return moment(props.date).fromNow();
}

const date_text = ref(getDateText());

setInterval(() => {
  date_text.value = getDateText();
}, 1000);

</script>

<template>
  <div class="post">
    <div class="post-author-icon">
      <Avatar size="50px" :name="name"/>
    </div>
    <div class="post-content">
      <div class="post-header">
        <span class="post-author">@{{ name }}</span>
        <span class="post-date">{{date_text}}</span>
      </div>
      <div class="post-message">
        {{ content }}
      </div>
      <div class="post-stamps">
        <div class="post-stamp" :class="stamp_clicked[0] ? ['clicked'] : undefined">
          <span class="post-stamp-icon">❤️</span>
          <span class="post-stamp-count">{{ stamp_counts[0] }}</span>
        </div>
        <div class="post-stamp" :class="stamp_clicked[1] ? ['clicked'] : undefined">
          <span class="post-stamp-icon">🔥</span>
          <span class="post-stamp-count">{{ stamp_counts[1] }}</span>
        </div>
        <div class="post-stamp" :class="stamp_clicked[2] ? ['clicked'] : undefined">
          <span class="post-stamp-icon">💧</span>
          <span class="post-stamp-count">{{ stamp_counts[2] }}</span>
        </div>
        <div class="post-stamp" :class="stamp_clicked[3] ? ['clicked'] : undefined">
          <span class="post-stamp-icon">😢</span>
          <span class="post-stamp-count">{{ stamp_counts[3] }}</span>
        </div>
        <div class="post-stamp" :class="stamp_clicked[4] ? ['clicked'] : undefined">
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
