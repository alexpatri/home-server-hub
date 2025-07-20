<template>
  <Teleport to="body">
    <transition name="fade">
      <div
        class="fixed inset-0 bg-black/40 flex items-center justify-center z-50"
        @click="$emit('close')"
      >
        <div
          class="bg-white mx-3 rounded-lg min-w-3/4 h-5/6 shadow-xl w-full max-w-md flex flex-col"
          @click.stop
        >
          <div
            class="relative w-full min-h-16 py-2 px-4 bg-white shadow-lg rounded-t-lg flex items-center"
          >
            <h1 class="font-bold text-2xl text-dark">{{ title }}</h1>

            <button
              @click="$emit('close')"
              class="absolute top-2 right-2 text-gray-500 hover:text-primary hover:cursor-pointer"
            >
              <font-awesome-icon :icon="['fas', 'xmark']" />
            </button>
          </div>

          <div class="w-full h-full p-4">
            <slot> </slot>
          </div>
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
