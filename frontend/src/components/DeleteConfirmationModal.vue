<template>
  <div v-if="show" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white rounded-lg max-w-md w-full p-6 shadow-xl">
      <h3 class="text-lg font-medium text-gray-900 mb-4">Confirm Delete</h3>
      <p class="text-sm text-gray-500 mb-6">
        Are you sure you want to delete this file? This action cannot be undone.
      </p>
      <div class="flex justify-end space-x-4">
        <button 
          @click="$emit('cancel')" 
          class="px-4 py-2 bg-gray-100 text-gray-700 rounded-md hover:bg-gray-200"
        >
          Cancel
        </button>
        <button 
          @click="$emit('confirm')" 
          class="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700"
        >
          Delete
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  show: boolean;
}>();

defineEmits<{
  (e: 'confirm'): void;
  (e: 'cancel'): void;
}>();
</script>
