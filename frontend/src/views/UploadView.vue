<template>
  <div class="container mx-auto px-4 py-6">
    <h1 class="text-2xl font-bold text-gray-800 mb-6">Upload Files</h1>
    
    <!-- Upload Zone -->
    <div 
      class="w-full border-2 border-dashed rounded-lg p-6 mb-6"
      :class="{ 
        'border-indigo-300 bg-indigo-50': isDragging,
        'border-gray-300 bg-gray-50': !isDragging
      }"
      @dragenter.prevent="isDragging = true"
      @dragleave.prevent="isDragging = false"
      @dragover.prevent="isDragging = true"
      @drop.prevent="onFileDrop"
    >
      <div class="flex flex-col items-center">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 text-indigo-500 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
        </svg>
        <p class="text-lg mb-2">Drag & drop files here</p>
        <p class="text-gray-500 text-sm mb-4">or</p>
        <label class="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-md cursor-pointer">
          Browse Files
          <input type="file" class="hidden" multiple @change="onFileSelect" />
        </label>
      </div>
    </div>
    
    <!-- File upload progress -->
    <div v-if="uploadQueue.length > 0" class="bg-white rounded-lg shadow-sm p-4 mb-6">
      <h3 class="text-lg font-medium mb-4">Upload Progress</h3>
      <div v-for="file in uploadQueue" :key="file.name" class="mb-4">
        <div class="flex items-center justify-between mb-1">
          <div class="flex items-center">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            <span class="text-sm truncate max-w-md">{{ file.name }}</span>
          </div>
          <div class="text-xs text-gray-500">{{ formatFileSize(file.size) }}</div>
        </div>
        
        <div class="w-full bg-gray-200 rounded-full h-2.5">
          <div 
            class="h-2.5 rounded-full transition-all" 
            :class="{
              'bg-indigo-600': file.status !== 'error',
              'bg-red-600': file.status === 'error'
            }"
            :style="{ width: file.status === 'error' ? '100%' : `${file.progress}%` }" 
          ></div>
        </div>
        
        <div class="flex justify-between mt-1">
          <span class="text-xs text-gray-500">
            {{ file.status === 'uploading' ? `${file.progress}%` : file.status }}
          </span>
          <button 
            v-if="file.status === 'uploading'" 
            @click="cancelUpload(file)"
            class="text-xs text-red-600 hover:text-red-800"
          >
            Cancel
          </button>
          <button 
            v-if="file.status === 'error'" 
            @click="retryUpload(file)"
            class="text-xs text-indigo-600 hover:text-indigo-800"
          >
            Retry
          </button>
        </div>
      </div>
    </div>
    
    <!-- Recently uploaded files -->
    <div v-if="recentlyUploaded.length > 0" class="bg-white rounded-lg shadow-sm p-4">
      <h3 class="text-lg font-medium mb-4">Recently Uploaded</h3>
      <ul class="divide-y">
        <li v-for="file in recentlyUploaded" :key="file.id" class="py-3 flex justify-between items-center">
          <div class="flex items-center">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-3 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            <span>{{ file.filename }}</span>
          </div>
          <div class="flex">
            <router-link :to="`/share/${file.id}`" class="text-indigo-600 hover:text-indigo-800 text-sm mr-4">
              Share
            </router-link>
            <router-link :to="'/'" class="text-gray-600 hover:text-gray-800 text-sm">
              View All Files
            </router-link>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { formatFileSize } from '../services/fileService';

interface UploadFile {
  id?: string;
  name: string;
  size: number;
  file: File;
  progress: number;
  status: 'uploading' | 'completed' | 'error';
  uploadId?: string;
  partNumber?: number;
}

interface UploadedFile {
  id: number;
  filename: string;
  size: number;
  uploaded_at: string;
}

const isDragging = ref(false);
const uploadQueue = ref<UploadFile[]>([]);
const recentlyUploaded = ref<UploadedFile[]>([]);

// Handle file drop
const onFileDrop = (event: DragEvent) => {
  isDragging.value = false;
  const files = event.dataTransfer?.files;
  if (files) {
    addFilesToQueue(files);
  }
};

// Handle file select from browse button
const onFileSelect = (event: Event) => {
  const input = event.target as HTMLInputElement;
  const files = input.files;
  if (files) {
    addFilesToQueue(files);
    input.value = ''; // Reset input so same file can be selected again
  }
};

// Add files to the upload queue
const addFilesToQueue = (fileList: FileList) => {
  for (let i = 0; i < fileList.length; i++) {
    const file = fileList[i];
    
    // Add file to the queue
    const uploadFile: UploadFile = {
      name: file.name,
      size: file.size,
      file,
      progress: 0,
      status: 'uploading'
    };
    
    uploadQueue.value.push(uploadFile);
    
    // Start upload process for this file
    startUpload(uploadFile);
  }
};

// Start multipart upload for a file
const startUpload = async (uploadFile: UploadFile) => {
  // Get API base URL from environment variable
  const API_URL = import.meta.env.VITE_API_URL || '/api';
  
  try {
    // Step 1: Initialize multipart upload
    const initResponse = await fetch(`${API_URL}/upload/init`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        filename: uploadFile.name,
        content_type: uploadFile.file.type,
      })
    });
    
    if (!initResponse.ok) {
      throw new Error('Failed to initialize upload');
    }
    
    const initData = await initResponse.json();
    const uploadId = initData.upload_id;
    uploadFile.uploadId = uploadId;
    
    // Step 2: Determine optimal chunk size based on file size
    // For smaller files (<5MB), use single part upload
    const CHUNK_SIZE = uploadFile.size < 5 * 1024 * 1024 
      ? uploadFile.size 
      : 5 * 1024 * 1024; // 5MB chunks for larger files
    
    // Step 3: Split file into chunks and upload each part
    const totalChunks = Math.ceil(uploadFile.size / CHUNK_SIZE);
    const parts = [];
    
    for (let partNumber = 1; partNumber <= totalChunks; partNumber++) {
      // Update current part number
      uploadFile.partNumber = partNumber;
      
      // Calculate chunk boundaries
      const start = (partNumber - 1) * CHUNK_SIZE;
      const end = Math.min(partNumber * CHUNK_SIZE, uploadFile.size);
      const chunk = uploadFile.file.slice(start, end);
      
      // Create form data for this part
      const formData = new FormData();
      formData.append('upload_id', uploadId);
      formData.append('part_number', partNumber.toString());
      formData.append('file_part', chunk);
      
      // Upload the part
      const partResponse = await fetch(`${API_URL}/upload/part`, {
        method: 'POST',
        body: formData,
      });
      
      if (!partResponse.ok) {
        throw new Error(`Failed to upload part ${partNumber}`);
      }
      
      const partData = await partResponse.json();
      parts.push({
        PartNumber: partNumber,
        ETag: partData.etag
      });
      
      // Update progress
      uploadFile.progress = Math.round((partNumber / totalChunks) * 100);
    }
    
    // Step 4: Complete the multipart upload
    const completeResponse = await fetch(`${API_URL}/upload/complete`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        upload_id: uploadId,
        parts: parts
      })
    });
    
    if (!completeResponse.ok) {
      throw new Error('Failed to complete upload');
    }
    
    const completeData = await completeResponse.json();
    
    // Mark as completed
    uploadFile.status = 'completed';
    uploadFile.progress = 100;
    
    // Add to recently uploaded files
    recentlyUploaded.value.push({
      id: completeData.id,
      filename: uploadFile.name,
      size: uploadFile.size,
      uploaded_at: new Date().toISOString()
    });
    
    // Remove from queue after a delay
    setTimeout(() => {
      uploadQueue.value = uploadQueue.value.filter(f => f.name !== uploadFile.name);
    }, 3000);
    
  } catch (error) {
    console.error('Upload error:', error);
    uploadFile.status = 'error';
    
    // Attempt to abort the multipart upload if it was initialized
    if (uploadFile.uploadId) {
      try {
        await fetch(`${API_URL}/upload/abort`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            upload_id: uploadFile.uploadId
          })
        });
      } catch (abortError) {
        console.error('Failed to abort upload:', abortError);
      }
    }
  }
};

// Cancel an ongoing upload
const cancelUpload = async (file: UploadFile) => {
  // Call abort API if upload was initialized
  if (file.uploadId) {
    const API_URL = import.meta.env.VITE_API_URL || '/api';
    
    try {
      await fetch(`${API_URL}/upload/abort`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          upload_id: file.uploadId
        })
      });
    } catch (error) {
      console.error('Failed to abort upload:', error);
    }
  }
  
  // Remove from queue
  uploadQueue.value = uploadQueue.value.filter(f => f !== file);
};

// Retry a failed upload
const retryUpload = (file: UploadFile) => {
  file.progress = 0;
  file.status = 'uploading';
  startUpload(file);
};
</script>
