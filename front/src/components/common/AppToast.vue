<template>
  <transition name="toast-fade">
    <div v-if="visible" class="toast">{{ message }}</div>
  </transition>
</template>

<script setup lang="ts">
import { useToast } from "@/composables/useToast";

const { visible, message } = useToast();
</script>

<style scoped>
.toast {
  position: fixed;
  left: 50%;
  bottom: calc(96px + var(--safe-bottom));
  transform: translateX(-50%);
  z-index: 120;
  padding: 10px 14px;
  background: rgba(12, 15, 25, 0.92);
  border: 1px solid rgba(255, 255, 255, 0.18);
  border-radius: 12px;
  font-size: 13px;
  color: #fff;
}

.toast-fade-enter-active,
.toast-fade-leave-active {
  transition: all 0.2s ease;
}

.toast-fade-enter-from,
.toast-fade-leave-to {
  opacity: 0;
  transform: translate(-50%, 8px);
}
</style>
