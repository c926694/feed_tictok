<template>
  <transition name="drawer">
    <section v-if="open" class="mask" @click.self="$emit('close')">
      <div class="sheet">
        <header>
          <h4>评论 {{ comments.length }}</h4>
          <button @click="$emit('close')">关闭</button>
        </header>

        <ul class="list">
          <li v-for="item in comments" :key="item.id">
            <div class="author-row">
              <img v-if="item.author.avatar" :src="item.author.avatar" alt="avatar" />
              <span>{{ item.author.username || item.author.nickname }}</span>
            </div>
            <p>{{ item.content }}</p>
            <small>
              <button class="like-btn" :class="{ active: item.liked }" @click="onLike(item.id)">
                <svg viewBox="0 0 24 24" class="icon-svg" aria-hidden="true">
                  <path d="M12 20.2 4.9 13.7a4.9 4.9 0 0 1 6.9-7L12 7.9l.2-.2a4.9 4.9 0 0 1 6.9 7L12 20.2Z" />
                </svg>
                <span>{{ item.likeCount }}</span>
              </button>
            </small>
          </li>
        </ul>

        <form class="composer" @submit.prevent="submitComment">
          <input v-model.trim="draft" placeholder="写下你的评论..." />
          <button :disabled="sending || !draft" type="submit">{{ sending ? "发送中" : "发送" }}</button>
        </form>
      </div>
    </section>
  </transition>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";
import { createComment, fetchCommentList, switchCommentLike } from "@/api";
import type { Comment } from "@/types/domain";
import { useToast } from "@/composables/useToast";

const props = defineProps<{
  open: boolean;
  videoId: number;
}>();

defineEmits<{
  (e: "close"): void;
}>();

const { showToast } = useToast();

const comments = ref<Comment[]>([]);
const draft = ref("");
const sending = ref(false);

watch(
  () => props.open,
  async (isOpen) => {
    if (!isOpen || !props.videoId) return;
    comments.value = await fetchCommentList(props.videoId);
  },
  { immediate: true }
);

async function submitComment() {
  if (!draft.value || !props.videoId) return;
  sending.value = true;
  try {
    await createComment({
      video_id: props.videoId,
      content: draft.value
    });
    draft.value = "";
    comments.value = await fetchCommentList(props.videoId);
    showToast("评论成功");
  } finally {
    sending.value = false;
  }
}

async function onLike(commentId: number) {
  const targetLiked = await switchCommentLike(commentId);
  comments.value = comments.value.map((item) =>
    item.id === commentId
      ? {
          ...item,
          liked: targetLiked,
          likeCount: Math.max(0, item.likeCount + ((targetLiked ? 1 : 0) - (item.liked ? 1 : 0)))
        }
      : item
  );
}
</script>

<style scoped>
.mask {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.42);
  z-index: 90;
}

.sheet {
  position: absolute;
  top: 0;
  right: 0;
  width: min(420px, 34vw);
  height: 100svh;
  background: rgba(15, 19, 34, 0.96);
  border-left: 1px solid var(--line);
  box-shadow: -24px 0 60px rgba(0, 0, 0, 0.32);
  display: flex;
  flex-direction: column;
}

header {
  padding: 14px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid var(--line);
}

h4 {
  margin: 0;
}

header button {
  border: none;
  background: transparent;
  color: var(--text-muted);
}

.list {
  flex: 1;
  margin: 0;
  padding: 0 14px;
  overflow: auto;
  list-style: none;
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.list::-webkit-scrollbar {
  display: none;
}

li {
  padding: 12px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.author-row {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text-muted);
  font-size: 12px;
}

.author-row img {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  object-fit: cover;
}

p {
  margin: 8px 0;
}

.like-btn {
  border: none;
  color: var(--text-muted);
  background: transparent;
  padding: 0;
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.like-btn.active {
  color: #ff5d7a;
}

.icon-svg {
  width: 16px;
  height: 16px;
  fill: currentColor;
}

.composer {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 8px;
  padding: 12px 14px calc(12px + var(--safe-bottom));
  border-top: 1px solid var(--line);
}

input {
  border: 1px solid var(--line);
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.04);
  color: #fff;
  padding: 10px 12px;
}

.composer button {
  border: none;
  border-radius: 999px;
  background: var(--accent);
  color: #fff;
  padding: 0 16px;
}

.drawer-enter-active,
.drawer-leave-active {
  transition: opacity 0.2s ease;
}

.drawer-enter-from,
.drawer-leave-to {
  opacity: 0;
}

.drawer-enter-active .sheet,
.drawer-leave-active .sheet {
  transition: transform 0.24s ease;
}

.drawer-enter-from .sheet,
.drawer-leave-to .sheet {
  transform: translateX(100%);
}

@media (max-width: 900px) {
  .sheet {
    width: min(100vw, 380px);
  }
}

@media (max-width: 640px) {
  .sheet {
    width: 100vw;
  }
}
</style>
