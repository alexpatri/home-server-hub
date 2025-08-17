<template>
  <div class="bg-yellow-50 border border-yellow-200 rounded-lg p-4 cursor-pointer hover:shadow-md transition-all duration-200 md:w-64">
    <div class="flex items-center mb-3">
      <div
        class="w-12 h-12 bg-yellow-200 rounded-lg flex items-center justify-center mr-3 relative overflow-hidden"
      >
        <font-awesome-icon icon="cube" class="text-yellow-600" />
      </div>
      <div class="flex-1">
        <p
          class="font-semibold bg-transparent border-b border-yellow-300 focus:border-yellow-500 outline-none w-full"
        >
          {{ app.name }}
        </p>
        <p class="text-sm text-gray-600">{{ app.ip }}:{{ app.port }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

// Interfaces
interface DiscoveredApp {
  container: string
  ip: string
  name: string
  port: number
  tags: string[]
  url: string
  image: File | null
}

// Props
const props = defineProps<{
  app: DiscoveredApp
  isSaving?: boolean
}>()
</script>
