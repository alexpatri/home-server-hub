<template>
  <div class="h-full w-full">
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-2xl font-bold">Gerenciar Aplicações</h2>
    </div>

    <!-- Lista de aplicações atuais -->
    <div class="mb-8 w-full h-96">
      <h3 class="text-lg font-semibold mb-4">Aplicações Registradas</h3>
      <div v-if="!applications || applications.length === 0" class="text-gray-500 py-8 h-full flex justify-center items-center">
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
                app.status === 'active' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800',
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
          @click="console.log('test')"
        />
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
      editableName: app.name,
      editablePort: app.port,
      editableUrl: '',
      editableImage: null,
    }))

    hasSearched.value = true
  } catch (error) {
    console.error('Erro ao descobrir aplicações:', error)
  } finally {
    isDiscovering.value = false
  }
}

const saveApplication = async (app: DiscoveredApp): Promise<void> => {
  isSaving.value = true

  try {
    // Criar FormData para incluir a imagem se houver
    const formData = new FormData()
    formData.append('name', app.name)
    formData.append('ip', app.ip)
    formData.append('port', app.port.toString())

    if (app.url) {
      formData.append('url', app.url)
    }

    if (app.image) {
      formData.append('image', app.image)
    }

    await axios.post(`${import.meta.env.VITE_API_BASE_URL}/applications`, formData, {
      params: {
        container_id: app.container,
      },
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })

    // Remover da lista de descobertas
    const index = discoveredApps.value.findIndex(
      (discovered) => discovered.container === app.container,
    )
    if (index > -1) {
      discoveredApps.value.splice(index, 1)
    }

    // Atualizar lista de aplicações
    await fetchApplications()
    emit('applicationsUpdated')
  } catch (error) {
    console.error('Erro ao salvar aplicação:', error)
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
