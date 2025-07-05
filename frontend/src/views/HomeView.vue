<template>
  <div class="h-screen w-screen bg-light flex flex-col items-center gap-5">
    <header class="w-full bg-light py-3 px-3 sm:px-10 flex items-center justify-between shadow-md">
      <Clock />
      <SearchInput class="hidden md:inline" ref="searchRef" v-model="searchValue" placeholder="Buscar na Internet" />
    </header>
    <main class="px-3 sm:px-10 w-full h-full"></main>
    <DefaultFooter />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

import Clock from '@/components/Clock.vue'
import SearchInput from '@/components/SearchInput.vue'
import DefaultFooter from '@/components/DefaultFooter.vue'

const searchRef = ref()
const searchValue = ref<string>('')

const focusSearch = (): void => {
  searchRef.value?.focus()
}

const blurSearch = (): void => {
  searchRef.value?.blur()
}

// Keyboard shortcuts
const handleKeydown = (event: KeyboardEvent): void => {
  if ((event.ctrlKey || event.metaKey) && event.key === 'k') {
    event.preventDefault()
    focusSearch()
  }

  if (event.key === 'Escape') {
    blurSearch()
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})
</script>
