<template>
  <BaseModal title="Configurações" @close="$emit('close')">
    <div class="flex h-full w-full min-h-[600px]">
      <!-- Menu lateral -->
      <div class="w-64 pr-4 border-r-2 border-gray-300">
        <ul class="space-y-2">
          <li v-for="tab in tabs" :key="tab.id">
            <button
              @click="activeTab = tab.id"
              :class="[
                'w-full text-left px-4 py-3 rounded-lg transition-colors cursor-pointer',
                activeTab === tab.id
                  ? 'bg-blue-100 text-blue-700 font-medium'
                  : 'text-gray-600 hover:bg-gray-100',
              ]"
            >
              <font-awesome-icon :icon="tab.icon" class="mr-1" />
              {{ tab.name }}
            </button>
          </li>
        </ul>
      </div>

      <!-- Conteúdo -->
      <div class="flex-1 pl-4 w-1/2">
        <component :is="currentTabComponent" @applicationsUpdated="handleApplicationsUpdated" />
      </div>
    </div>
  </BaseModal>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import BaseModal from '@/components/BaseModal.vue'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import CustomizationTab from '@/components/CustomizationTab.vue'
import ApplicationsTab from '@/components/ApplicationsTab.vue'

// Interfaces
interface Tab {
  id: string
  name: string
  icon: string
  component: any
}

// Emits
const emit = defineEmits<{
  (e: 'close'): void
  (e: 'applicationsUpdated'): void
}>()

// Configuração das tabs
const tabs: Tab[] = [
  {
    id: 'personalization',
    name: 'Personalização',
    icon: 'palette',
    component: CustomizationTab,
  },
  {
    id: 'applications',
    name: 'Aplicações',
    icon: 'cube',
    component: ApplicationsTab,
  },
]

// Estado reativo
const activeTab = ref<string>('applications')

// Computed
const currentTabComponent = computed(() => {
  const currentTab = tabs.find((tab) => tab.id === activeTab.value)
  return currentTab?.component
})

// Métodos
const handleApplicationsUpdated = (): void => {
  emit('applicationsUpdated')
}
</script>
