<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import Loader from '@/components/Loader.vue'
</script>

<template>
  <RouterView />
  <Loader />
</template>

<style scoped></style>
