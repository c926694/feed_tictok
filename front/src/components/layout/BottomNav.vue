<template>
  <nav class="bottom-nav">
    <RouterLink
      v-for="item in navItems"
      :key="item.to"
      :to="item.to"
      class="nav-item"
      :class="{ active: route.path === item.to }"
    >
      <span class="icon">{{ item.icon }}</span>
      <span>{{ item.label }}</span>
    </RouterLink>
  </nav>
</template>

<script setup lang="ts">
import { useRoute } from "vue-router";

const route = useRoute();

const navItems = [
  { label: "首页", to: "/feed", icon: "⌂" },
  { label: "发布", to: "/upload", icon: "＋" },
  { label: "我的", to: "/profile", icon: "我" }
];
</script>

<style scoped>
.bottom-nav {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 80;
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  background: linear-gradient(180deg, transparent, rgba(5, 7, 13, 0.95) 35%);
  padding: 10px 16px calc(12px + var(--safe-bottom));
}

.nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  color: var(--text-muted);
  font-size: 12px;
  transition: color 0.2s ease;
}

.icon {
  font-size: 18px;
  line-height: 1;
}

.active {
  color: var(--text-primary);
}
</style>
