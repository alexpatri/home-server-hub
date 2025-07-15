<template>
  <div
    v-if="isLoading"
    class="fixed inset-0 bg-black/40 backdrop-blur-sm flex items-center justify-center z-50"
  >
    <div
      class="animate-spin w-14 h-14 border-4 border-light border-t-transparent rounded-full"
    ></div>
  </div>
</template>

<script setup lang="ts">
import { useLoading } from '@/stores/loadingStore'

const { isLoading } = useLoading()
</script>
