<template>
  <div class="container mx-auto px-4 py-6">
    <h1 class="text-2xl font-bold text-gray-800 mb-6">Shared Files</h1>
    
    <!-- Loading state -->
    <div v-if="loading" class="flex justify-center items-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-indigo-500"></div>
    </div>
    
    <!-- Error state -->
    <div v-else-if="error" class="bg-red-50 border-l-4 border-red-500 p-4 mb-6">
      <div class="flex">
        <div class="flex-shrink-0">
          <svg class="h-5 w-5 text-red-500" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
          </svg>
        </div>
        <div class="ml-3">
          <p class="text-sm text-red-700">{{ error }}</p>
        </div>
      </div>
    </div>
    
    <!-- No shared files -->
    <div v-else-if="sharedFiles.length === 0" class="bg-gray-50 p-8 text-center rounded-lg border border-gray-200">
      <svg xmlns="http://www.w3.org/2000/svg" class="mx-auto h-12 w-12 text-gray-400 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
      </svg>
      <h3 class="text-lg font-medium text-gray-900 mb-2">No shared files</h3>
      <p class="text-gray-500 mb-4">You haven't shared any files yet.</p>
      <router-link to="/" class="text-indigo-600 hover:text-indigo-800 font-medium">
        Go to My Files
      </router-link>
    </div>
    
    <!-- Shared files list -->
    <div v-else class="bg-white shadow-sm rounded-lg overflow-hidden">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              File
            </th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Shared On
            </th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Link
            </th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Views
            </th>
            <th scope="col" class="relative px-6 py-3">
              <span class="sr-only">Actions</span>
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="file in sharedFiles" :key="file.id">
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="flex items-center">
                <div class="flex-shrink-0 h-10 w-10 flex items-center justify-center rounded bg-gray-100">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                </div>
                <div class="ml-4">
                  <div class="text-sm font-medium text-gray-900">
                    {{ file.filename }}
                  </div>
                  <div class="text-sm text-gray-500">
                    {{ formatFileSize(file.size) }}
                  </div>
                </div>
              </div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm text-gray-900">{{ formatDate(file.shared_at) }}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="flex items-center">
                <input 
                  type="text" 
                  :value="getShareLink(file.share_id)" 
                  readonly
                  class="text-sm text-gray-500 border border-gray-300 rounded-l-md px-3 py-1 focus:outline-none"
                />
                <button 
                  @click="copyShareLink(file.share_id)"
                  class="bg-gray-100 border border-l-0 border-gray-300 rounded-r-md p-1 hover:bg-gray-200"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                  </svg>
                </button>
              </div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
              {{ file.view_count }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
              <button 
                @click="unshareFile(file.share_id)" 
                class="text-red-600 hover:text-red-900"
              >
                Unshare
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { formatFileSize } from '../services/fileService';

// Interfaces
interface SharedFile {
  id: number;
  share_id: string;
  filename: string;
  size: number;
  shared_at: string;
  view_count: number;
}

// State
const loading = ref(true);
const error = ref<string | null>(null);
const sharedFiles = ref<SharedFile[]>([]);

// Get all shared files
async function fetchSharedFiles() {
  const API_URL = import.meta.env.VITE_API_URL || '/api';
  
  try {
    loading.value = true;
    error.value = null;
    
    const response = await fetch(`${API_URL}/shares`);
    
    if (!response.ok) {
      throw new Error('Failed to fetch shared files');
    }
    
    const data = await response.json();
    sharedFiles.value = data;
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load shared files';
    console.error('Error loading shared files:', err);
  } finally {
    loading.value = false;
  }
}

// Format date
function formatDate(dateString: string): string {
  const date = new Date(dateString);
  return new Intl.DateTimeFormat('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date);
}

// Get shareable link
function getShareLink(shareId: string): string {
  const baseUrl = window.location.origin;
  return `${baseUrl}/share/${shareId}`;
}

// Copy share link to clipboard
function copyShareLink(shareId: string): void {
  const shareLink = getShareLink(shareId);
  navigator.clipboard.writeText(shareLink)
    .then(() => {
      alert('Link copied to clipboard!');
    })
    .catch(err => {
      console.error('Failed to copy link:', err);
    });
}

// Unshare a file
async function unshareFile(shareId: string) {
  if (!confirm('Are you sure you want to unshare this file? This action cannot be undone.')) {
    return;
  }
  
  const API_URL = import.meta.env.VITE_API_URL || '/api';
  
  try {
    const response = await fetch(`${API_URL}/shares/${shareId}`, {
      method: 'DELETE'
    });
    
    if (!response.ok) {
      throw new Error('Failed to unshare file');
    }
    
    // Remove from list
    sharedFiles.value = sharedFiles.value.filter(file => file.share_id !== shareId);
  } catch (err) {
    console.error('Error unsharing file:', err);
    alert('Failed to unshare file. Please try again.');
  }
}

// Fetch shared files on component mount
onMounted(fetchSharedFiles);
</script>
