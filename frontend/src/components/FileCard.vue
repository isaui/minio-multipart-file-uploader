<template>
  <div class="bg-white rounded-lg shadow-md overflow-hidden transition-transform hover:shadow-lg hover:-translate-y-1 border border-gray-100">
    <div class="p-4">
      <!-- File icon -->
      <div class="mb-3 flex justify-center">
        <div :class="getIconClass(file.filename)" class="text-4xl text-indigo-500"></div>
      </div>
      
      <!-- File details -->
      <div class="text-center mb-3">
        <h3 class="font-medium text-gray-900 truncate" :title="file.filename">{{ file.filename }}</h3>
        <p class="text-sm text-gray-500">{{ formatSize(file.size) }}</p>
        <p class="text-xs text-gray-400">{{ formatDate(file.uploaded_at) }}</p>
      </div>
      
      <!-- Actions -->
      <div class="flex justify-center space-x-2 mt-4">
        <!-- Download button -->
        <button 
          @click="handleDownload"
          class="p-2 text-indigo-600 hover:bg-indigo-50 rounded-full" 
          title="Download file"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
          </svg>
        </button>
        
        <!-- Share button -->
        <button 
          @click="handleShare"
          class="p-2 text-indigo-600 hover:bg-indigo-50 rounded-full" 
          title="Share file"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z" />
          </svg>
        </button>
        
        <!-- Delete button -->
        <button 
          @click="handleDelete"
          class="p-2 text-red-500 hover:bg-red-50 rounded-full" 
          title="Delete file"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits } from 'vue';
import type { FileItem } from '../services/fileService';
import { formatFileSize as formatSize } from '../services/fileService';

// Props
const props = defineProps<{
  file: FileItem;
}>();

// Emits
const emit = defineEmits<{
  (e: 'download', id: number): void;
  (e: 'share', id: number): void;
  (e: 'delete', id: number): void;
}>();

// Methods
function getIconClass(filename: string): string {
  const extension = filename.split('.').pop()?.toLowerCase() || '';
  let iconClass = 'far fa-file';
  
  // Check extension and assign appropriate icon class
  if (['jpg', 'jpeg', 'png', 'gif', 'webp', 'svg'].includes(extension)) {
    iconClass = 'far fa-file-image';
  } else if (['pdf'].includes(extension)) {
    iconClass = 'far fa-file-pdf';
  } else if (['doc', 'docx'].includes(extension)) {
    iconClass = 'far fa-file-word';
  } else if (['xls', 'xlsx'].includes(extension)) {
    iconClass = 'far fa-file-excel';
  } else if (['ppt', 'pptx'].includes(extension)) {
    iconClass = 'far fa-file-powerpoint';
  } else if (['zip', 'rar', '7z'].includes(extension)) {
    iconClass = 'far fa-file-archive';
  } else if (['mp3', 'wav', 'ogg'].includes(extension)) {
    iconClass = 'far fa-file-audio';
  } else if (['mp4', 'mov', 'avi'].includes(extension)) {
    iconClass = 'far fa-file-video';
  }

  return iconClass;
}

function formatDate(dateString: string): string {
  const date = new Date(dateString);
  return new Intl.DateTimeFormat('id-ID', { 
    year: 'numeric', 
    month: 'short', 
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date);
}

function handleDownload(): void {
  emit('download', props.file.id);
}

function handleShare(): void {
  emit('share', props.file.id);
}

function handleDelete(): void {
  emit('delete', props.file.id);
}
</script>
