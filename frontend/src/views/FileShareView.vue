<template>
  <div class="container mx-auto px-4 py-8 max-w-3xl">
    <!-- Loading state -->
    <div v-if="loading" class="flex justify-center items-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-indigo-500"></div>
    </div>
    
    <!-- Error state -->
    <div v-else-if="error" class="bg-white p-8 rounded-lg shadow-md text-center">
      <svg xmlns="http://www.w3.org/2000/svg" class="mx-auto h-16 w-16 text-red-500 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
      </svg>
      <h2 class="text-xl font-medium text-gray-900 mb-2">Error</h2>
      <p class="text-gray-500 mb-6">{{ error }}</p>
      <router-link to="/" class="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-md">
        Go Home
      </router-link>
    </div>
    
    <!-- File content -->
    <div v-else class="bg-white rounded-lg shadow-md overflow-hidden">
      <!-- File header -->
      <div class="p-6 border-b">
        <div class="flex items-center">
          <div class="bg-gray-100 p-3 rounded-lg mr-4">
            <svg class="h-8 w-8 text-gray-500" :class="fileIconClass" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </div>
          <div>
            <h1 class="text-xl font-medium text-gray-900">{{ fileData.filename }}</h1>
            <p class="text-gray-500 text-sm">{{ formatFileSize(fileData.size) }} • Shared {{ formatDate(fileData.shared_at) }}</p>
          </div>
        </div>
      </div>
      
      <!-- File preview (if supported) -->
      <div v-if="isPreviewable" class="p-6 border-b bg-gray-50">
        <div class="rounded-lg overflow-hidden shadow-sm">
          <!-- Image preview -->
          <img 
            v-if="isImage" 
            :src="`${apiUrl}/shares/${fileId}/content`" 
            class="max-w-full h-auto mx-auto"
            alt="File preview"
          />
          
          <!-- PDF preview -->
          <iframe 
            v-else-if="isPdf"
            :src="`${apiUrl}/shares/${fileId}/content`" 
            class="w-full h-96"
            title="PDF preview"
          ></iframe>
          
          <!-- Video preview -->
          <video 
            v-else-if="isVideo"
            controls
            class="w-full"
          >
            <source :src="`${apiUrl}/shares/${fileId}/content`" :type="fileData.content_type">
            Your browser does not support the video tag.
          </video>
          
          <!-- Audio preview -->
          <audio 
            v-else-if="isAudio"
            controls
            class="w-full"
          >
            <source :src="`${apiUrl}/shares/${fileId}/content`" :type="fileData.content_type">
            Your browser does not support the audio tag.
          </audio>
        </div>
      </div>
      
      <!-- Download section -->
      <div class="p-6">
        <div class="text-center">
          <p class="text-gray-500 mb-4">
            This file was shared with you. Click the button below to download it.
          </p>
          <a 
            :href="`${apiUrl}/shares/${fileId}/download`" 
            class="bg-indigo-600 hover:bg-indigo-700 text-white px-5 py-2.5 rounded-md inline-flex items-center"
            download
          >
            <svg class="w-5 h-5 mr-2" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
            </svg>
            Download File
          </a>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { formatFileSize } from '../services/fileService';

// Get route information
const route = useRoute();
const fileId = computed(() => route.params.fileId as string);

// API URL
const apiUrl = import.meta.env.VITE_API_URL || '/api';

// State
const loading = ref(true);
const error = ref<string | null>(null);
const fileData = ref<{
  id: number;
  filename: string;
  size: number;
  content_type: string;
  shared_at: string;
}>({
  id: 0,
  filename: '',
  size: 0,
  content_type: '',
  shared_at: ''
});

// Computed properties for file type checking
const isImage = computed(() => 
  fileData.value.content_type.startsWith('image/')
);

const isPdf = computed(() => 
  fileData.value.content_type === 'application/pdf'
);

const isVideo = computed(() => 
  fileData.value.content_type.startsWith('video/')
);

const isAudio = computed(() => 
  fileData.value.content_type.startsWith('audio/')
);

const isPreviewable = computed(() => 
  isImage.value || isPdf.value || isVideo.value || isAudio.value
);

const fileIconClass = computed(() => {
  if (isImage.value) return 'text-blue-500';
  if (isPdf.value) return 'text-red-500';
  if (isVideo.value) return 'text-green-500';
  if (isAudio.value) return 'text-purple-500';
  return 'text-gray-500';
});

// Format date
function formatDate(dateString: string): string {
  const date = new Date(dateString);
  return new Intl.DateTimeFormat('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  }).format(date);
}

// Fetch shared file data
async function fetchSharedFile() {
  try {
    loading.value = true;
    error.value = null;
    
    const response = await fetch(`${apiUrl}/shares/${fileId.value}`);
    
    if (!response.ok) {
      if (response.status === 404) {
        throw new Error('The file you are looking for does not exist or has been removed.');
      }
      throw new Error('Failed to load the shared file');
    }
    
    const data = await response.json();
    fileData.value = data;
    
    // Record view
    try {
      await fetch(`${apiUrl}/shares/${fileId.value}/view`, {
        method: 'POST'
      });
    } catch (err) {
      // Silently fail if view recording fails
      console.error('Failed to record view:', err);
    }
    
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'An unexpected error occurred';
    console.error('Error loading shared file:', err);
  } finally {
    loading.value = false;
  }
}

// Fetch data on component mount
onMounted(fetchSharedFile);
</script>
