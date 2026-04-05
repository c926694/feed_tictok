<template>
  <section class="auth-page">
    <LoginForm />
  </section>
</template>

<script setup lang="ts">
import LoginForm from "@/components/auth/LoginForm.vue";
</script>

<style scoped>
.auth-page {
  min-height: 100svh;
  display: grid;
  place-items: center;
  padding: 20px;
}
</style>
