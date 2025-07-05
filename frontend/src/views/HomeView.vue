<template>
  <div class="h-screen w-screen bg-light flex flex-col items-center gap-5">
    <DeafaultHeader />
    <main class="px-3 sm:px-10 w-full h-full"></main>
    <DefaultFooter />
  </div>
</template>

<script setup lang="ts">
import DefaultFooter from '@/components/DefaultFooter.vue'
import DeafaultHeader from '@/components/DeafaultHeader.vue'
</script>
