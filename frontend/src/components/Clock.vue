<!-- components/Clock.vue -->
<template>
  <div class="flex flex-col items-center text-black dark:text-white">
    <h1 class="text-[64px] font-bold leading-none">{{ time }}</h1>
    <div class="flex items-center space-x-2 mt-2 text-lg">
      <font-awesome-icon :icon="['fas', 'cloud']" />
      <span>23° Brasília, DF</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'


const time = ref('00:00')

function updateTime() {
  const now = new Date()
  const hours = String(now.getHours()).padStart(2, '0')
  const minutes = String(now.getMinutes()).padStart(2, '0')
  time.value = `${hours}:${minutes}`
}

onMounted(() => {
  updateTime()
  const interval = setInterval(updateTime, 1000)
  onUnmounted(() => clearInterval(interval))
})
</script>
