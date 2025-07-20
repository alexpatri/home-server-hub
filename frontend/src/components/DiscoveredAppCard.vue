<template>
  <div class="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
    <div class="flex items-center mb-3">
      <div class="w-12 h-12 bg-yellow-200 rounded-lg flex items-center justify-center mr-3">
        <font-awesome-icon icon="cube" class="text-yellow-600" />
      </div>
      <div class="flex-1">
        <input
          v-model="editableName"
          class="font-semibold bg-transparent border-b border-yellow-300 focus:border-yellow-500 outline-none w-full"
          :placeholder="app.name"
        />
        <p class="text-sm text-gray-600">{{ app.ip }}:{{ app.port }}</p>
      </div>
    </div>

    <div class="mb-3">
      <label class="block text-sm font-medium text-gray-700 mb-1">Tags:</label>
      <input
        v-model="editableTags"
        class="w-full p-2 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
        placeholder="tag1, tag2, tag3"
      />
    </div>

    <div class="flex items-center justify-between">
      <button
        @click="handleSave"
        :disabled="isSaving"
        class="bg-green-600 text-white px-3 py-1 rounded hover:bg-green-700 disabled:bg-gray-400 text-sm"
      >
        <font-awesome-icon
          :icon="isSaving ? 'spinner' : 'save'"
          :class="{ 'animate-spin': isSaving }"
          class="mr-1"
        />
        Salvar
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

// Interfaces
interface DiscoveredApp {
  container: string
  exists: boolean
  ip: string
  name: string
  port: number
  tags: string[]
  editableName: string
  editableTags: string
}

// Props
const props = defineProps<{
  app: DiscoveredApp
  isSaving?: boolean
}>()

// Emits
const emit = defineEmits<{
  (e: 'save', app: DiscoveredApp): void
}>()

// Estado reativo local
const editableName = ref<string>(props.app.editableName)
const editableTags = ref<string>(props.app.editableTags)

// Watchers para sincronizar mudanças das props
watch(
  () => props.app.editableName,
  (newValue) => {
    editableName.value = newValue
  },
)

watch(
  () => props.app.editableTags,
  (newValue) => {
    editableTags.value = newValue
  },
)

// Métodos
const handleSave = (): void => {
  const updatedApp: DiscoveredApp = {
    ...props.app,
    editableName: editableName.value,
    editableTags: editableTags.value,
  }

  emit('save', updatedApp)
}
</script>
