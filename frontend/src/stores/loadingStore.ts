import { ref } from 'vue'

const isLoading = ref(false)

export function useLoading() {
  return {
    isLoading,
    show: () => (isLoading.value = true),
    hide: () => (isLoading.value = false),
  }
}
