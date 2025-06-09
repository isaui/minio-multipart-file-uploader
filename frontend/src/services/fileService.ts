// File service for handling API calls related to files

// Define types for our data
export interface FileItem {
  id: number;
  filename: string;
  filepath: string;
  size: number;
  uploaded_at: string;
  content_type?: string;
}

// API base URL - dapat diinjeksi melalui environment variable di Docker Compose
// Default ke relative URL '/api' untuk production atau localhost untuk development
const API_URL = import.meta.env.VITE_API_URL || '/api';

/**
 * Get all files
 * @returns Promise with array of files
 */
export async function getAllFiles(): Promise<FileItem[]> {
  try {
    const response = await fetch(`${API_URL}/files`);
    
    if (!response.ok) {
      throw new Error(`Error: ${response.status}`);
    }
    
    return await response.json();
  } catch (error) {
    console.error('Error fetching files:', error);
    throw error;
  }
}

/**
 * Get a file by ID
 * @param id File ID
 * @returns Promise with file data
 */
export async function getFileById(id: number): Promise<FileItem> {
  try {
    const response = await fetch(`${API_URL}/files/${id}`);
    
    if (!response.ok) {
      throw new Error(`Error: ${response.status}`);
    }
    
    return await response.json();
  } catch (error) {
    console.error(`Error fetching file ${id}:`, error);
    throw error;
  }
}

/**
 * Delete a file
 * @param id File ID to delete
 * @returns Promise with deletion status
 */
export async function deleteFile(id: number): Promise<any> {
  try {
    const response = await fetch(`${API_URL}/files/${id}`, {
      method: 'DELETE',
    });
    
    if (!response.ok) {
      throw new Error(`Error: ${response.status}`);
    }
    
    return await response.json();
  } catch (error) {
    console.error(`Error deleting file ${id}:`, error);
    throw error;
  }
}

/**
 * Generate file type icon based on filename
 * @param filename The name of the file
 * @returns Icon class name
 */
export function getFileIcon(filename: string): string {
  const extension = filename.split('.').pop()?.toLowerCase() || '';
  
  switch (extension) {
    case 'pdf':
      return 'file-pdf';
    case 'doc':
    case 'docx':
      return 'file-word';
    case 'xls':
    case 'xlsx':
      return 'file-excel';
    case 'ppt':
    case 'pptx':
      return 'file-powerpoint';
    case 'jpg':
    case 'jpeg':
    case 'png':
    case 'gif':
    case 'webp':
      return 'file-image';
    case 'mp4':
    case 'mov':
    case 'avi':
      return 'file-video';
    case 'mp3':
    case 'wav':
    case 'ogg':
      return 'file-audio';
    case 'zip':
    case 'rar':
    case '7z':
      return 'file-archive';
    default:
      return 'file';
  }
}

/**
 * Format file size to human readable format
 * @param bytes Size in bytes
 * @returns Formatted string (e.g., "2.5 MB")
 */
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes';
  
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
  
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}
