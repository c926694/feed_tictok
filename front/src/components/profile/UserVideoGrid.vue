<template>
  <section class="grid-wrap">
    <h3>我的视频</h3>
    <div v-if="videos.length" class="grid">
      <RouterLink
        v-for="video in videos"
        :key="video.id"
        class="item"
        :to="{ path: '/profile/videos', query: { videoId: video.id } }"
      >
        <img v-if="video.coverUrl" :src="video.coverUrl" :alt="video.title" />
        <div v-else class="fallback">暂无封面</div>
        <footer>
          <strong>{{ video.title }}</strong>
          <small>点赞 {{ video.likeCount }} · 评论 {{ video.commentCount }}</small>
        </footer>
      </RouterLink>
    </div>
    <p v-else class="empty">还没有可展示的视频，先去发布一个吧。</p>
  </section>
</template>

<script setup lang="ts">
import type { Video } from "@/types/domain";

defineProps<{
  videos: Video[];
}>();
</script>

<style scoped>
.grid-wrap {
  padding: 8px 14px 120px;
  max-width: 760px;
  margin: 0 auto;
}

h3 {
  margin: 0 0 10px;
}

.grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8px;
}

.item {
  background: rgba(255, 255, 255, 0.03);
  border-radius: 10px;
  overflow: hidden;
  text-decoration: none;
  color: inherit;
}

img,
.fallback {
  width: 100%;
  aspect-ratio: 3 / 4;
  object-fit: cover;
}

.fallback {
  display: grid;
  place-items: center;
  color: var(--text-muted);
  font-size: 12px;
}

footer {
  padding: 6px;
  display: flex;
  flex-direction: column;
  gap: 3px;
}

strong {
  font-size: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

small {
  color: var(--text-muted);
  font-size: 11px;
}

.empty {
  color: var(--text-muted);
}

@media (max-width: 520px) {
  .grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}
</style>
