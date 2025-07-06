<template>
  <div
    @click="handleClick"
    class="cursor-pointer rounded-md shadow-lg hover:shadow-xl bg-white dark:bg-gray-800 transition transform hover:-translate-y-1 duration-200"
  >
    <img
      :src="imageSrc"
      alt="Application image"
      class="w-full h-32 object-cover rounded-t-md"
    />
    <div class="p-4">
      <h2 class="text-lg font-semibold text-gray-800 dark:text-white">
        {{ application.name }}
      </h2>
      <p v-if="application.tags" class="text-sm text-gray-500 mt-1">
        {{ application.tags.join(', ') }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  application: {
    id: string
    name: string
    tags: string[] | null
    image: { name: string; data: string } | null
    container: string
    ip: string
    port: number
    url: string
    status: string
  }
  defaultImage?: string
}

const props = defineProps<Props>()

const imageSrc = computed(() => {
  if (props.application.image?.data) {
    return `data:image/png;base64,${props.application.image.data}`
  }
  return props.defaultImage ?? '/default-app.png'
})

const handleClick = () => {
  const target = props.application.url?.trim()
    ? props.application.url
    : `http://${props.application.ip}:${props.application.port}`
  window.open(target, '_blank')
}
</script>
