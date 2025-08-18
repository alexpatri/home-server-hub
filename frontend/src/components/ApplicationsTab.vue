<template>
  <div class="h-full w-full">
    <!-- Conteúdo principal da tab -->
    <div v-if="!showEditForm" class="h-full w-full">
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-2xl font-bold">Gerenciar Aplicações</h2>
      </div>

      <!-- Lista de aplicações atuais -->
      <div class="mb-8 w-full h-96">
        <h3 class="text-lg font-semibold mb-4">Aplicações Registradas</h3>
        <div
          v-if="!applications || applications.length === 0"
          class="text-gray-500 py-8 h-full flex justify-center items-center"
        >
          Nenhuma aplicação encontrada
        </div>
        <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div
            v-for="app in applications"
            :key="app.id"
            class="bg-white border rounded-lg p-4 hover:shadow-md transition-shadow"
          >
            <div class="flex items-center mb-3">
              <div class="w-12 h-12 bg-gray-200 rounded-lg flex items-center justify-center mr-3">
                <img
                  v-if="app.image?.data"
                  :src="`data:image/png;base64,${app.image.data}`"
                  :alt="app.name"
                  class="w-full h-full object-cover rounded-lg"
                />
                <font-awesome-icon v-else icon="cube" class="text-gray-400" />
              </div>
              <div>
                <h4 class="font-semibold">{{ app.name }}</h4>
                <p class="text-sm text-gray-600">{{ app.ip }}:{{ app.port }}</p>
              </div>
            </div>
            <div class="flex items-center justify-between">
              <span
                :class="[
                  'px-2 py-1 rounded-full text-xs',
                  app.status === 'active'
                    ? 'bg-green-100 text-green-800'
                    : 'bg-red-100 text-red-800',
                ]"
              >
                {{ app.status === 'active' ? 'Ativa' : 'Inativa' }}
              </span>
              <button @click="editApplication(app)" class="text-blue-600 hover:text-blue-800 p-1">
                <font-awesome-icon icon="edit" />
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Aplicações descobertas -->
      <div class="mb-8 w-full space-y-4">
        <div class="flex justify-between">
          <h3 class="text-lg font-semibold mb-4">
            Aplicações Disponíveis ({{ discoveredApps.length }})
          </h3>
          <button
            @click="discoverApplications"
            :disabled="isDiscovering"
            class="bg-primary text-white px-4 rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
          >
            <font-awesome-icon
              :icon="isDiscovering ? 'spinner' : 'search'"
              :class="{ 'animate-spin': isDiscovering }"
              class="mr-2"
            />
            {{ isDiscovering ? 'Buscando...' : 'Buscar Aplicações' }}
          </button>
        </div>
        <div v-if="discoveredApps.length > 0" class="flex gap-4 overflow-x-auto w-full pb-5">
          <DiscoveredAppCard
            v-for="(app, index) in discoveredApps"
            :key="index"
            :app="app"
            :is-saving="isSaving"
            @click="openEditForm(app)"
          />
        </div>
      </div>
    </div>

    <!-- Formulário de edição -->
    <div v-else class="h-full w-full">
      <div class="mb-6">
        <button
          @click="closeEditForm"
          class="flex items-center text-blue-600 hover:text-blue-800 mb-4"
        >
          <font-awesome-icon icon="arrow-left" class="mr-2" />
          Voltar
        </button>
        <h2 class="text-2xl font-bold">Configurar Aplicação</h2>
      </div>

      <div class="bg-white rounded-lg shadow-md p-6 max-w-2xl">
        <form @submit.prevent="saveApplication">
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
              <h3 class="text-lg font-semibold">{{ editingApp?.name }}</h3>
              <p class="text-gray-600">{{ editingApp?.ip }}</p>
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
              @click="closeEditForm"
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
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import DiscoveredAppCard from '@/components/DiscoveredAppCard.vue'

// Interfaces
interface Application {
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

interface DiscoverResponse {
  discovered: DiscoveredApp[]
  total: number
}

interface FormData {
  name: string
  port: number
  url: string
  tags: string
  image: File | null
}

// Emits
const emit = defineEmits<{
  (e: 'applicationsUpdated'): void
}>()

// Estado reativo
const applications = ref<Application[]>([])
const discoveredApps = ref<DiscoveredApp[]>([])
const isDiscovering = ref<boolean>(false)
const isSaving = ref<boolean>(false)
const hasSearched = ref<boolean>(false)

// Estado do formulário de edição
const showEditForm = ref<boolean>(false)
const editingApp = ref<DiscoveredApp | null>(null)
const imagePreview = ref<string | null>(null)
const fileInput = ref<HTMLInputElement | null>(null)

const formData = ref<FormData>({
  name: '',
  port: 8080,
  url: '',
  tags: '',
  image: null,
})

// Métodos
const fetchApplications = async (): Promise<void> => {
  try {
    const response = await axios.get<{ applications: Application[] }>(
      `${import.meta.env.VITE_API_BASE_URL}/applications`,
    )
    applications.value = response.data.applications
  } catch (error) {
    console.error('Erro ao buscar aplicações:', error)
  }
}

const discoverApplications = async (): Promise<void> => {
  isDiscovering.value = true
  hasSearched.value = false

  try {
    const response = await axios.get<DiscoverResponse>(
      `${import.meta.env.VITE_API_BASE_URL}/applications/discover`,
    )

    discoveredApps.value = response.data.discovered.map((app) => ({
      ...app,
      url: app.url || '',
      image: null,
    }))

    hasSearched.value = true
  } catch (error) {
    console.error('Erro ao descobrir aplicações:', error)
  } finally {
    isDiscovering.value = false
  }
}

const openEditForm = (app: DiscoveredApp): void => {
  editingApp.value = app
  formData.value = {
    name: app.name,
    port: app.port,
    url: app.url || '',
    tags: app.tags.join(', '),
    image: app.image,
  }

  // Limpar preview de imagem anterior
  imagePreview.value = null

  showEditForm.value = true
}

const closeEditForm = (): void => {
  showEditForm.value = false
  editingApp.value = null
  imagePreview.value = null

  // Reset form data
  formData.value = {
    name: '',
    port: 8080,
    url: '',
    tags: '',
    image: null,
  }

  // Clear file input
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

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

const saveApplication = async (): Promise<void> => {
  if (!editingApp.value) return

  isSaving.value = true

  try {
    const tags = formData.value.tags
      .split(',')
      .map((tag) => tag.trim())
      .filter((tag) => tag.length > 0)

    // Criar FormData para incluir a imagem se houver
    const formDataToSend = new FormData()
    formDataToSend.append('name', formData.value.name)
    formDataToSend.append('ip', editingApp.value.ip)
    formDataToSend.append('port', formData.value.port.toString())
    formDataToSend.append('tags', JSON.stringify(tags))

    if (formData.value.url) {
      formDataToSend.append('url', formData.value.url)
    }

    if (formData.value.image) {
      formDataToSend.append('image', formData.value.image)
    }

    await axios.post(`${import.meta.env.VITE_API_BASE_URL}/applications`, formDataToSend, {
      params: {
        container_id: editingApp.value.container,
      },
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })

    // Remover da lista de descobertas
    const index = discoveredApps.value.findIndex(
      (discovered) => discovered.container === editingApp.value?.container,
    )
    if (index > -1) {
      discoveredApps.value.splice(index, 1)
    }

    // Atualizar lista de aplicações
    await fetchApplications()
    emit('applicationsUpdated')

    // Fechar formulário
    closeEditForm()
  } catch (error) {
    console.error('Erro ao salvar aplicação:', error)
    alert('Erro ao salvar aplicação. Tente novamente.')
  } finally {
    isSaving.value = false
  }
}

const editApplication = (app: Application): void => {
  console.log('Editando aplicação:', app)
  // Implementar modal de edição se necessário
}

// Carregar aplicações ao montar o componente
onMounted(fetchApplications)
</script>
