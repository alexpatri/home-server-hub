<template>
  <div class="h-full w-full">
    <div class="mb-6">
      <button @click="handleBack" class="flex items-center text-blue-600 hover:text-blue-800 mb-4">
        <font-awesome-icon icon="arrow-left" class="mr-2" />
        Voltar
      </button>
      <h2 class="text-2xl font-bold">Configurar Aplicação</h2>
    </div>

    <div class="pt-6">
      <form @submit.prevent="handleSave">
        <!-- Preview da imagem e informações básicas -->
        <div class="flex items-center mb-6">
          <div
            class="w-16 h-16 bg-gray-200 rounded-lg flex items-center justify-center mr-4 relative overflow-hidden"
          >
            <img
              v-if="imagePreview"
              :src="imagePreview"
              alt="Preview"
              class="w-full h-full object-cover rounded-lg"
            />
            <font-awesome-icon v-else icon="cube" class="text-gray-400 text-2xl" />
          </div>
          <div>
            <h3 class="text-lg font-semibold">{{ app.name }}</h3>
            <p class="text-gray-600">{{ app.ip }}</p>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <!-- Nome da aplicação -->
          <div class="md:col-span-2">
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Nome da aplicação <span class="text-red-500">*</span>
            </label>
            <input
              v-model="formData.name"
              type="text"
              required
              class="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              placeholder="Digite o nome da aplicação"
            />
          </div>

          <!-- Porta -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Porta <span class="text-red-500">*</span>
            </label>
            <input
              v-model.number="formData.port"
              type="number"
              required
              min="1"
              max="65535"
              class="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              placeholder="8080"
            />
          </div>

          <!-- URL pública -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              URL pública (opcional)
            </label>
            <input
              v-model="formData.url"
              type="url"
              class="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              placeholder="https://exemplo.com"
            />
          </div>

          <!-- Tags -->
          <div class="md:col-span-2">
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Tags (separadas por vírgula)
            </label>
            <input
              v-model="formData.tags"
              type="text"
              class="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              placeholder="web, dashboard, monitoramento"
            />
          </div>

          <!-- Upload de imagem -->
          <div class="md:col-span-2">
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Imagem da aplicação (opcional)
            </label>
            <div
              class="mt-1 flex justify-center px-6 pt-5 pb-6 border-2 border-gray-300 border-dashed rounded-lg hover:border-gray-400 transition-colors"
            >
              <div class="space-y-1 text-center">
                <font-awesome-icon
                  icon="cloud-upload-alt"
                  class="mx-auto h-12 w-12 text-gray-400"
                />
                <div class="flex text-sm text-gray-600">
                  <label
                    class="relative cursor-pointer bg-white rounded-md font-medium text-blue-600 hover:text-blue-500 focus-within:outline-none focus-within:ring-2 focus-within:ring-offset-2 focus-within:ring-blue-500"
                  >
                    <span>Selecionar arquivo</span>
                    <input
                      ref="fileInput"
                      type="file"
                      accept="image/*"
                      @change="handleImageUpload"
                      class="sr-only"
                    />
                  </label>
                  <p class="pl-1">ou arraste e solte</p>
                </div>
                <p class="text-xs text-gray-500">PNG, JPG, GIF até 5MB</p>
                <button
                  v-if="formData.image || imagePreview"
                  @click="clearImage"
                  type="button"
                  class="text-red-600 hover:text-red-800 text-sm mt-2"
                >
                  <font-awesome-icon icon="trash" class="mr-1" />
                  Remover imagem
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Botões de ação -->
        <div class="flex justify-end space-x-4 mt-8">
          <button
            @click="handleBack"
            type="button"
            class="px-6 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors"
          >
            Cancelar
          </button>
          <button
            type="submit"
            :disabled="isSaving"
            class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
          >
            <font-awesome-icon
              :icon="isSaving ? 'spinner' : 'save'"
              :class="{ 'animate-spin': isSaving }"
              class="mr-2"
            />
            {{ isSaving ? 'Salvando...' : 'Salvar Aplicação' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

// Interfaces
interface DiscoveredApp {
  container: string
  exists: boolean
  ip: string
  name: string
  port: number
  tags: string[]
  url: string
  image: File | null
}

interface FormData {
  name: string
  port: number
  url: string
  tags: string
  image: File | null
}

// Props
const props = defineProps<{
  app: DiscoveredApp
  isSaving?: boolean
}>()

// Emits
const emit = defineEmits<{
  (e: 'back'): void
  (e: 'save', formData: FormData): void
}>()

// Estado reativo
const imagePreview = ref<string | null>(null)
const fileInput = ref<HTMLInputElement | null>(null)

const formData = ref<FormData>({
  name: '',
  port: 8080,
  url: '',
  tags: '',
  image: null,
})

// Watchers para sincronizar com as props
watch(
  () => props.app,
  (newApp) => {
    if (newApp) {
      formData.value = {
        name: newApp.name,
        port: newApp.port,
        url: newApp.url || '',
        tags: newApp.tags.join(', '),
        image: newApp.image,
      }

      // Limpar preview anterior
      imagePreview.value = null
    }
  },
  { immediate: true },
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

    formData.value.image = file

    // Criar preview
    const reader = new FileReader()
    reader.onload = (e) => {
      imagePreview.value = e.target?.result as string
    }
    reader.readAsDataURL(file)
  }
}

const clearImage = (): void => {
  formData.value.image = null
  imagePreview.value = null
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

const handleBack = (): void => {
  emit('back')
}

const handleSave = (): void => {
  emit('save', { ...formData.value })
}

// Inicializar dados quando o componente for montado
onMounted(() => {
  if (props.app) {
    formData.value = {
      name: props.app.name,
      port: props.app.port,
      url: props.app.url || '',
      tags: props.app.tags.join(', '),
      image: props.app.image,
    }
  }
})
</script>
