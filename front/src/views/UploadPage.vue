<template>
  <section class="page">
    <header>
      <h1>发布</h1>
      <p>发布后将显示在视频流中。</p>
    </header>
    <UploadForm />
    <BottomNav />
  </section>
</template>

<script setup lang="ts">
import BottomNav from "@/components/layout/BottomNav.vue";
import UploadForm from "@/components/upload/UploadForm.vue";
</script>

<style scoped>
.page {
  min-height: 100svh;
}

header {
  padding: 22px 16px 0;
}

h1 {
  margin: 0;
}

p {
  margin: 8px 0 0;
  color: var(--text-muted);
}
</style>
