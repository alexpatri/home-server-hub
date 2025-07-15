<template>
  <div class="h-screen w-screen bg-light flex flex-col items-center gap-5">
    <DeafaultHeader />
    <main class="px-3 sm:px-10 w-full h-full overflow-y-auto py-5">
      <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-6">
        <AppCard
          v-for="app in applications"
          :key="app.id"
          :application="app"
          defaultImage="/default-app.png"
        />
      </div>
      <BaseModal v-if="showConfigModal" title="Configurações" @close="showConfigModal = false"/>
    </main>
    <DefaultFooter @configBtnClick="showConfigModal = true"/>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import DefaultFooter from '@/components/DefaultFooter.vue'
import DeafaultHeader from '@/components/DeafaultHeader.vue'
import AppCard from '@/components/AppCard.vue'
import BaseModal from '@/components/BaseModal.vue'

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

const applications = ref<Application[]>([])
const showConfigModal = ref<boolean>(false)

const fetchApplications = async () => {
  try {
    const res = await axios.get(`${import.meta.env.VITE_API_BASE_URL}/applications`)
    applications.value = res.data.applications
  } catch (err) {
    console.error('Erro ao buscar aplicações:', err)
  }
}

onMounted(fetchApplications)
</script>
