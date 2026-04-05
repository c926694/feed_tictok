<template>
  <section class="feed-page">
    <aside class="left-rail">
      <div class="brand">
        <span class="logo-dot"></span>
        <strong>Feed Video</strong>
      </div>

      <RouterLink class="rail-link active" to="/feed">推荐</RouterLink>
      <RouterLink class="rail-link" to="/upload">发布</RouterLink>
      <RouterLink class="rail-link" to="/profile">我的</RouterLink>
    </aside>

    <section class="main-stage">
      <header class="top-bar">
        <div class="feed-tabs">
          <button :class="{ active: tab === 'recommend' }" @click="changeTab('recommend')">推荐</button>
          <button :class="{ active: tab === 'follow' }" @click="changeTab('follow')">关注</button>
          <button :class="{ active: tab === 'hot' }" @click="changeTab('hot')">热门</button>
          <button v-if="tab === 'hot' && hotPlayMode" class="playback-back" @click="backToHotBoard">返回热榜</button>
        </div>

        <div class="quick-nav">
          <RouterLink to="/upload">发布</RouterLink>
          <RouterLink to="/profile">我的</RouterLink>
        </div>
      </header>

      <HotRankBoard v-if="tab === 'hot' && !hotPlayMode" @play="openHotPlayback" />
      <VideoFeed
        v-else-if="tab === 'hot'"
        :tab="tab"
        :initial-hot-videos="hotSeedVideos"
        :initial-hot-video-id="hotSeedVideoId"
      />
      <VideoFeed v-else :tab="tab" />
    </section>
  </section>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";
import HotRankBoard from "@/components/feed/HotRankBoard.vue";
import VideoFeed from "@/components/feed/VideoFeed.vue";
import type { Video } from "@/types/domain";

const tab = ref<"recommend" | "follow" | "hot">("recommend");
const hotPlayMode = ref(false);
const hotSeedVideos = ref<Video[]>([]);
const hotSeedVideoId = ref(0);

function changeTab(nextTab: "recommend" | "follow" | "hot") {
  tab.value = nextTab;
}

function openHotPlayback(payload: { videoId: number; videos: Video[] }) {
  hotSeedVideoId.value = payload.videoId;
  hotSeedVideos.value = payload.videos.slice();
  hotPlayMode.value = true;
}

function backToHotBoard() {
  hotPlayMode.value = false;
}

watch(
  () => tab.value,
  (value) => {
    if (value !== "hot") {
      hotPlayMode.value = false;
      hotSeedVideoId.value = 0;
      hotSeedVideos.value = [];
    }
  }
);
</script>

<style scoped>
.feed-page {
  min-height: 100svh;
  background: linear-gradient(150deg, #0e1324 0%, #131b31 35%, #170f22 100%);
  display: grid;
  grid-template-columns: 124px 1fr;
  gap: 16px;
  padding: 12px;
}

.left-rail {
  border-radius: 16px;
  background: rgba(11, 14, 24, 0.88);
  border: 1px solid rgba(255, 255, 255, 0.09);
  padding: 16px 10px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.brand {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #fff;
  margin-bottom: 8px;
}

.logo-dot {
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background: linear-gradient(135deg, #17d7d1, #ff4d6d);
}

.rail-link {
  display: block;
  padding: 10px 12px;
  border-radius: 10px;
  color: #b8c0db;
}

.rail-link.active,
.rail-link:hover {
  color: #fff;
  background: rgba(255, 255, 255, 0.09);
}

.main-stage {
  min-width: 0;
}

.top-bar {
  height: 56px;
  border-radius: 14px;
  background: rgba(11, 14, 24, 0.88);
  border: 1px solid rgba(255, 255, 255, 0.09);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 12px;
  margin-bottom: 12px;
}

.feed-tabs {
  display: inline-flex;
  gap: 8px;
}

.feed-tabs button {
  border: none;
  border-radius: 999px;
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.06);
  color: #bdc6e2;
  cursor: pointer;
}

.feed-tabs button.active {
  color: #fff;
  background: linear-gradient(90deg, #2f6df5, #5f48e6);
}

.playback-back {
  border: none;
  border-radius: 999px;
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.09);
  color: #dfe7ff;
  cursor: pointer;
}

.quick-nav {
  display: inline-flex;
  gap: 10px;
}

.quick-nav a {
  color: #d9e2ff;
  font-size: 14px;
}

@media (max-width: 900px) {
  .feed-page {
    grid-template-columns: 1fr;
    padding: 0;
    gap: 0;
  }

  .left-rail {
    display: none;
  }

  .top-bar {
    border-radius: 0;
    margin-bottom: 0;
    border-left: none;
    border-right: none;
  }
}
</style>
