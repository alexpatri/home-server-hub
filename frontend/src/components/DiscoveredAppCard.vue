<template>
  <div class="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
    <div class="flex items-center mb-3">
      <div
        class="w-12 h-12 bg-yellow-200 rounded-lg flex items-center justify-center mr-3 relative overflow-hidden"
      >
        <img
          v-if="imagePreview"
          :src="imagePreview"
          alt="Preview"
          class="w-full h-full object-cover rounded-lg"
        />
        <font-awesome-icon v-else icon="cube" class="text-yellow-600" />
      </div>
      <div class="flex-1">
        <input
          v-model="editableName"
          class="font-semibold bg-transparent border-b border-yellow-300 focus:border-yellow-500 outline-none w-full"
          :placeholder="app.name"
        />
        <p class="text-sm text-gray-600">{{ app.ip }}:{{ editablePort }}</p>
      </div>
    </div>

    <div class="space-y-3 mb-3">
      <!-- Campo Port -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">Porta:</label>
        <input
          v-model.number="editablePort"
          type="number"
          class="w-full p-2 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          :placeholder="app.port.toString()"
          min="1"
          max="65535"
        />
      </div>

      <!-- Campo URL (opcional) -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">URL pública (opcional):</label>
        <input
          v-model="editableUrl"
          type="url"
          class="w-full p-2 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          placeholder="https://exemplo.com"
        />
      </div>

      <!-- Campo Imagem (opcional) -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">Imagem (opcional):</label>
        <input
          ref="fileInput"
          type="file"
          accept="image/*"
          @change="handleImageUpload"
          class="w-full p-2 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500 focus:border-blue-500 file:mr-4 file:py-1 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
        />
        <p class="text-xs text-gray-500 mt-1">PNG, JPG, GIF até 5MB</p>
      </div>
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

      <button
        v-if="editableImage || imagePreview"
        @click="clearImage"
        class="text-red-600 hover:text-red-800 text-sm"
      >
        <font-awesome-icon icon="trash" class="mr-1" />
        Remover imagem
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

// Emits
const emit = defineEmits<{
  (e: 'save', app: DiscoveredApp): void
}>()

// Estado reativo local
const editableName = ref<string>(props.app.name)
const editablePort = ref<number>(props.app.port)
const editableUrl = ref<string>(props.app.url)
const editableImage = ref<File | null>(props.app.image)
const imagePreview = ref<string | null>(null)
const fileInput = ref<HTMLInputElement | null>(null)

// Watchers para sincronizar mudanças das props
watch(
  () => props.app.name,
  (newValue) => {
    editableName.value = newValue
  },
)

watch(
  () => props.app.port,
  (newValue) => {
    editablePort.value = newValue
  },
)

watch(
  () => props.app.url,
  (newValue) => {
    editableUrl.value = newValue
  },
)

watch(
  () => props.app.image,
  (newValue) => {
    editableImage.value = newValue
  },
)

// Métodos
const handleImageUpload = (event: Event): void => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]

  if (file) {
    // Validar tipo de arquivo
    if (!file.type.startsWith('image/')) {
      alert('Por favor, selecione apenas arquivos de imagem')
      return
    }

    // Validar tamanho (5MB máximo)
    if (file.size > 5 * 1024 * 1024) {
      alert('A imagem deve ter no máximo 5MB')
      return
    }

    editableImage.value = file

    // Criar preview
    const reader = new FileReader()
    reader.onload = (e) => {
      imagePreview.value = e.target?.result as string
    }
    reader.readAsDataURL(file)
  }
}

const clearImage = (): void => {
  editableImage.value = null
  imagePreview.value = null
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

const handleSave = (): void => {
  const updatedApp: DiscoveredApp = {
    name: editableName.value,
    container: props.app.container,
    ip: props.app.ip,
    tags: props.app.tags,
    port: editablePort.value,
    url: editableUrl.value,
    image: editableImage.value,
  }

  emit('save', updatedApp)
}
</script>
