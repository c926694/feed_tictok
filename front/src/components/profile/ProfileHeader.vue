<template>
  <header class="profile-header">
    <img v-if="user.avatar" :src="user.avatar" alt="avatar" />
    <div v-else class="avatar-fallback">{{ user.nickname.slice(0, 1) }}</div>
    <h2>{{ user.nickname }}</h2>
    <p>@{{ user.username }}</p>
    <small>{{ user.bio }}</small>

    <div class="stats">
      <div class="stat">
        <strong>{{ user.followCount }}</strong>
        <span>关注</span>
      </div>
      <div class="stat">
        <strong>{{ user.followerCount }}</strong>
        <span>粉丝</span>
      </div>
      <div class="stat">
        <strong>{{ videoCount }}</strong>
        <span>视频</span>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import type { User } from "@/types/domain";

defineProps<{
  user: User;
  videoCount: number;
}>();
</script>

<style scoped>
.profile-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 24px 14px 18px;
}

img,
.avatar-fallback {
  width: 82px;
  height: 82px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid rgba(255, 255, 255, 0.2);
}

.avatar-fallback {
  display: grid;
  place-items: center;
  background: #1a2238;
  font-size: 30px;
}

h2 {
  margin: 0;
}

p {
  margin: 0;
  color: var(--text-muted);
}

small {
  color: rgba(255, 255, 255, 0.8);
}

.stats {
  margin-top: 8px;
  width: min(420px, 90vw);
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8px;
}

.stat {
  padding: 10px 8px;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.08);
  text-align: center;
}

.stat strong {
  display: block;
  font-size: 18px;
  line-height: 1.1;
}

.stat span {
  color: var(--text-muted);
  font-size: 12px;
}
</style>
