<template>
  <div class="w-full max-w-2xl">
    <div class="relative">
      <div
        class="flex items-center w-full px-4 py-3 bg-white border border-gray-300 rounded-lg shadow-sm hover:shadow-md transition-shadow duration-200 focus-within:shadow-md focus-within:border-primary"
      >
        <div class="flex items-center justify-center w-5 h-5 text-dark mr-3">
          <font-awesome-icon :icon="['fas', 'magnifying-glass']" />
        </div>

        <input
          ref="searchInput"
          v-model="searchQuery"
          type="text"
          :placeholder="placeholder"
          class="flex-1 text-dark placeholder-gray-500 bg-transparent border-none outline-none text-base"
          @focus="handleFocus"
          @blur="handleBlur"
          @keydown.enter="handleSearch"
          @input="handleInput"
        />

        <button
          v-if="searchQuery"
          @click="clearSearch"
          class="flex items-center justify-center w-6 h-6 text-gray-400 hover:text-dark hover:cursor-pointer transition-colors duration-200 ml-2"
        >
          <font-awesome-icon :icon="['fas', 'x']" />
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

interface Props {
  modelValue?: string
  placeholder?: string
  suggestions?: string[]
  showSuggestions?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  placeholder: 'Buscar na Internet',
  suggestions: () => [],
  showSuggestions: false,
})

interface Emits {
  (e: 'update:modelValue', value: string): void
  (e: 'search', query: string): void
  (e: 'input', value: string): void
  (e: 'focus'): void
  (e: 'blur'): void
}

const emit = defineEmits<Emits>()

const searchInput = ref<HTMLInputElement>()
const searchQuery = ref<string>(props.modelValue)
const isFocused = ref<boolean>(false)

watch(
  () => props.modelValue,
  (newVal: string) => {
    searchQuery.value = newVal
  },
)

watch(searchQuery, (newVal: string) => {
  emit('update:modelValue', newVal)
})

const handleFocus = (): void => {
  isFocused.value = true
  emit('focus')
}

const handleBlur = (): void => {
  isFocused.value = false
  emit('blur')
}

const handleInput = (): void => {
  emit('input', searchQuery.value)
}

const handleSearch = (): void => {
  if (searchQuery.value.trim()) {
    emit('search', searchQuery.value.trim())
    performSearch(searchQuery.value.trim())
  }
}

const clearSearch = async (): Promise<void> => {
  searchQuery.value = ''
  await nextTick()
  searchInput.value?.focus()
}

const performSearch = (query: string): void => {
  window.open(`https://www.google.com/search?q=${encodeURIComponent(query)}`, '_blank')
}

defineExpose({
  focus: () => searchInput.value?.focus(),
  blur: () => searchInput.value?.blur(),
  clear: clearSearch,
})
</script>
