<template>
  <div class="container mx-auto px-4 py-6">
    <!-- Header section -->
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-gray-800">My Files</h1>
      <router-link 
        to="/upload" 
        class="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-md flex items-center"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
        </svg>
        Upload New File
      </router-link>
    </div>

    <!-- Search and filter bar -->
    <div class="bg-white p-4 rounded-lg shadow-sm mb-6 flex flex-wrap gap-4">
      <div class="flex-1 min-w-[200px]">
        <input 
          v-model="searchTerm"
          type="text"
          placeholder="Search files..."
          class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
        />
      </div>
      <div class="w-48">
        <select 
          v-model="sortBy"
          class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
        >
          <option value="name_asc">Name (A-Z)</option>
          <option value="name_desc">Name (Z-A)</option>
          <option value="date_asc">Date (Oldest)</option>
          <option value="date_desc">Date (Newest)</option>
          <option value="size_asc">Size (Smallest)</option>
          <option value="size_desc">Size (Largest)</option>
        </select>
      </div>
      <div class="flex gap-2">
        <button 
          @click="viewMode = 'grid'"
          :class="['p-2 rounded-md', viewMode === 'grid' ? 'bg-indigo-100 text-indigo-600' : 'text-gray-500']"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
          </svg>
        </button>
        <button 
          @click="viewMode = 'list'"
          :class="['p-2 rounded-md', viewMode === 'list' ? 'bg-indigo-100 text-indigo-600' : 'text-gray-500']"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
          </svg>
        </button>
      </div>
    </div>

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
          <p class="text-sm text-red-700">
            {{ error }}
          </p>
        </div>
      </div>
    </div>

    <!-- File list component -->
    <FileList
      v-else
      :files="filteredFiles"
      :view-mode="viewMode"
      :search-term="searchTerm"
      @download="downloadFile"
      @share="shareFile"
      @delete="confirmDelete"
    />

    <!-- Delete confirmation modal -->
    <DeleteConfirmationModal
      :show="showDeleteModal"
      @confirm="deleteFile"
      @cancel="cancelDelete"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import FileList from '../components/FileList.vue';
import DeleteConfirmationModal from '../components/DeleteConfirmationModal.vue';
import { getAllFiles, deleteFile as deleteFileAPI } from '../services/fileService';
import type { FileItem } from '../services/fileService';

// Router
const router = useRouter();

// State
const files = ref<FileItem[]>([]);
const loading = ref(true);
const error = ref<string | null>(null);
const searchTerm = ref('');
const sortBy = ref('date_desc');
const viewMode = ref<'grid' | 'list'>('grid');
const showDeleteModal = ref(false);
const fileToDelete = ref<number | null>(null);

// Fetch files
async function fetchFiles() {
  loading.value = true;
  error.value = null;
  
  try {
    files.value = await getAllFiles();
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load files';
    console.error('Error loading files:', err);
  } finally {
    loading.value = false;
  }
}

// Computed properties
const filteredFiles = computed(() => {
  let result = [...files.value];
  
  // Apply search filter
  if (searchTerm.value) {
    const search = searchTerm.value.toLowerCase();
    result = result.filter(file => 
      file.filename.toLowerCase().includes(search)
    );
  }
  
  // Apply sorting
  switch (sortBy.value) {
    case 'name_asc':
      return result.sort((a, b) => a.filename.localeCompare(b.filename));
    case 'name_desc':
      return result.sort((a, b) => b.filename.localeCompare(a.filename));
    case 'date_asc':
      return result.sort((a, b) => new Date(a.uploaded_at).getTime() - new Date(b.uploaded_at).getTime());
    case 'date_desc':
      return result.sort((a, b) => new Date(b.uploaded_at).getTime() - new Date(a.uploaded_at).getTime());
    case 'size_asc':
      return result.sort((a, b) => a.size - b.size);
    case 'size_desc':
      return result.sort((a, b) => b.size - a.size);
    default:
      return result;
  }
});

// Methods
function downloadFile(id: number) {
  const file = files.value.find(f => f.id === id);
  if (!file) return;
  
  // Create a temporary anchor to trigger download
  const anchor = document.createElement('a');
  anchor.href = `${import.meta.env.VITE_API_URL || '/api'}/files/${id}/download`;
  anchor.download = file.filename;
  anchor.target = '_blank';
  document.body.appendChild(anchor);
  anchor.click();
  document.body.removeChild(anchor);
}

function shareFile(id: number) {
  router.push(`/share/${id}`);
}

function confirmDelete(id: number) {
  fileToDelete.value = id;
  showDeleteModal.value = true;
}

function cancelDelete() {
  showDeleteModal.value = false;
  fileToDelete.value = null;
}

async function deleteFile() {
  if (fileToDelete.value === null) return;
  
  try {
    await deleteFileAPI(fileToDelete.value);
    files.value = files.value.filter(file => file.id !== fileToDelete.value);
    showDeleteModal.value = false;
    fileToDelete.value = null;
  } catch (err) {
    console.error('Failed to delete file:', err);
    // You could show an error toast here
  }
}

// Fetch files on component mount
onMounted(fetchFiles);
</script>
