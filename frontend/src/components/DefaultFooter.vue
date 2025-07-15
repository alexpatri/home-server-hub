<template>
  <footer class="bg-dark px-3 sm:px-10 h-16 w-full flex items-center justify-between">
    <div class="w-10/12">
      <!--TODO: observability com prometheus-->
    </div>
    <button
      class="text-light text-2xl hover:cursor-pointer hover:text-primary transition-all"
      @click="$emit('configBtnClick')"
    >
      <font-awesome-icon :icon="['fas', 'gear']" />
    </button>
  </footer>
</template>

<script setup lang="ts">
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

defineEmits(['configBtnClick'])
</script>
