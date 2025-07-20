<template>
  <BaseModal title="Configurações" @close="$emit('close')">
    <div class="flex h-full min-h-[600px]">
      <!-- Menu lateral -->
      <div class="w-64 pr-4 border-r-2 border-gray-300">
        <ul class="space-y-2">
          <li>
            <button
              @click="activeTab = 'personalization'"
              :class="[
                'w-full text-left px-4 py-3 rounded-lg transition-colors cursor-pointer',
                activeTab === 'personalization'
                  ? 'bg-blue-100 text-blue-700 font-medium'
                  : 'text-gray-600 hover:bg-gray-100',
              ]"
            >
              <font-awesome-icon icon="palette" class="mr-1" />
              Personalização
            </button>
          </li>
          <li>
            <button
              @click="activeTab = 'applications'"
              :class="[
                'w-full text-left px-4 py-3 rounded-lg transition-colors cursor-pointer',
                activeTab === 'applications'
                  ? 'bg-blue-100 text-blue-700 font-medium'
                  : 'text-gray-600 hover:bg-gray-100',
              ]"
            >
              <font-awesome-icon icon="cube" class="mr-1" />
              Aplicações
            </button>
          </li>
        </ul>
      </div>

      <!-- Conteúdo -->
      <div class="flex-1 pl-4">
        <!-- Tab Personalização -->
        <div v-if="activeTab === 'personalization'" class="h-full">
          <h2 class="text-2xl font-bold mb-6">Personalização</h2>
          <div class="flex items-center justify-center h-64 text-gray-500">
            <div class="text-center">
              <font-awesome-icon icon="palette" class="text-6xl mb-4" />
              <p>Configurações de personalização em desenvolvimento</p>
            </div>
          </div>
        </div>

        <!-- Tab Aplicações -->
        <div v-if="activeTab === 'applications'" class="h-full">
          <div class="flex justify-between items-center mb-6">
            <h2 class="text-2xl font-bold">Gerenciar Aplicações</h2>
            <button
              @click="discoverApplications"
              :disabled="isDiscovering"
              class="bg-primary text-white px-4 py-2 rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
            >
              <font-awesome-icon
                :icon="isDiscovering ? 'spinner' : 'search'"
                :class="{ 'animate-spin': isDiscovering }"
                class="mr-2"
              />
              {{ isDiscovering ? 'Buscando...' : 'Buscar Novas Aplicações' }}
            </button>
          </div>

          <!-- Lista de aplicações atuais -->
          <div class="mb-8">
            <h3 class="text-lg font-semibold mb-4">Aplicações Instaladas</h3>
            <div
              v-if="!applications || applications.length === 0"
              class="text-gray-500 text-center py-8"
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
                  <div
                    class="w-12 h-12 bg-gray-200 rounded-lg flex items-center justify-center mr-3"
                  >
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
                  <button
                    @click="editApplication(app)"
                    class="text-blue-600 hover:text-blue-800 p-1"
                  >
                    <font-awesome-icon icon="edit" />
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- Aplicações descobertas -->
          <div v-if="discoveredApps.length > 0" class="mb-8">
            <h3 class="text-lg font-semibold mb-4">
              Aplicações Descobertas ({{ discoveredApps.length }})
            </h3>
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              <div
                v-for="(app, index) in discoveredApps"
                :key="index"
                class="bg-yellow-50 border border-yellow-200 rounded-lg p-4"
              >
                <div class="flex items-center mb-3">
                  <div
                    class="w-12 h-12 bg-yellow-200 rounded-lg flex items-center justify-center mr-3"
                  >
                    <font-awesome-icon icon="cube" class="text-yellow-600" />
                  </div>
                  <div class="flex-1">
                    <input
                      v-model="app.editableName"
                      class="font-semibold bg-transparent border-b border-yellow-300 focus:border-yellow-500 outline-none w-full"
                      :placeholder="app.name"
                    />
                    <p class="text-sm text-gray-600">{{ app.ip }}:{{ app.port }}</p>
                  </div>
                </div>

                <div class="mb-3">
                  <label class="block text-sm font-medium text-gray-700 mb-1">Tags:</label>
                  <input
                    v-model="app.editableTags"
                    class="w-full p-2 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    placeholder="tag1, tag2, tag3"
                  />
                </div>

                <div class="flex items-center justify-between">
                  <button
                    @click="saveApplication(app)"
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
            </div>
          </div>

          <!-- Estado vazio para aplicações descobertas -->
          <div v-else-if="hasSearched && !isDiscovering" class="text-center py-8 text-gray-500">
            <font-awesome-icon icon="search" class="text-4xl mb-4" />
            <p>Nenhuma nova aplicação encontrada</p>
          </div>
        </div>
      </div>
    </div>
  </BaseModal>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import BaseModal from '@/components/BaseModal.vue'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

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
  editableName: string
  editableTags: string
}

interface DiscoverResponse {
  discovered: DiscoveredApp[]
  total: number
}

// Emits
const emit = defineEmits<{
  (e: 'close'): void
  (e: 'applicationsUpdated'): void
}>()

// Estado reativo
const activeTab = ref<'personalization' | 'applications'>('applications')
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
      editableTags: app.tags.join(', '),
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
    const tags = app.editableTags
      .split(',')
      .map((tag) => tag.trim())
      .filter((tag) => tag.length > 0)

    const payload = {
      name: app.editableName || app.name,
      tags,
      ip: app.ip,
      port: app.port,
    }

    await axios.post(`${import.meta.env.VITE_API_BASE_URL}/applications`, payload, {
      params: {
        container_id: app.container,
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
