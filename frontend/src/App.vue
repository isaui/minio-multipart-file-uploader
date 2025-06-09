<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100">
    <!-- Navigation Header -->
    <header class="bg-white shadow-sm">
      <div class="container mx-auto px-4 py-4 flex justify-between items-center">
        <div class="flex items-center space-x-2">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-primary" viewBox="0 0 20 20" fill="currentColor">
            <path d="M5.5 13a3.5 3.5 0 01-.369-6.98 4 4 0 117.753-1.977A4.5 4.5 0 1113.5 13H9V9.413l1.293 1.293a1 1 0 001.414-1.414l-3-3a1 1 0 00-1.414 0l-3 3a1 1 0 001.414 1.414L7 9.414V13H5.5z" />
            <path d="M9 13h2v5a1 1 0 11-2 0v-5z" />
          </svg>
          <h1 class="text-xl font-bold text-gray-800">FileShare</h1>
        </div>
        
        <nav class="space-x-6">
          <RouterLink to="/" class="text-gray-600 hover:text-primary transition-colors">
            Files
          </RouterLink>
          <RouterLink to="/upload" class="text-gray-600 hover:text-primary transition-colors">
            Upload
          </RouterLink>
          <RouterLink to="/shared" class="text-gray-600 hover:text-primary transition-colors">
            Shared
          </RouterLink>
        </nav>
      </div>
    </header>

    <!-- Main Content -->
    <main class="container mx-auto px-4 py-8">
      <RouterView />
    </main>
    
    <!-- Footer -->
    <footer class="bg-white border-t mt-auto py-4">
      <div class="container mx-auto px-4 text-center text-sm text-gray-500">
        &copy; 2025 FileShare - Modern file sharing platform
      </div>
    </footer>
  </div>
</template>

