<template>
  <Teleport to="body">
    <transition name="fade">
      <div
        class="fixed inset-0 bg-black/40 flex items-center justify-center z-50"
        @click="$emit('close')"
      >
        <div
          class="bg-white mx-3 p-6 rounded-lg min-w-3/4 h-5/6 shadow-xl relative w-full max-w-md"
          @click.stop
        >
          <div class="absolute top-0 left-0 w-full min-h-16 py-2 px-4 bg-white shadow-lg rounded-t-lg flex items-center">
            <h1 class="font-bold text-2xl text-dark">{{ title }}</h1>

            <button
              @click="$emit('close')"
              class="absolute top-2 right-2 text-gray-500 hover:text-primary hover:cursor-pointer"
            >
              <font-awesome-icon :icon="['fas', 'xmark']" />
            </button>
          </div>

          <slot> </slot>
        </div>
      </div>
    </transition>
  </Teleport>
</template>

<script setup lang="ts">
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

interface Props {
  title?: string
}

const props = withDefaults(defineProps<Props>(), {
  title: '',
})

defineEmits(['close'])
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
