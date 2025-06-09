package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"time"

	"backend/config"
	"backend/models"

	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
)

// GetAllFiles retrieves all files from the database
func GetAllFiles(w http.ResponseWriter, _ *http.Request) {
	var files []models.File
	
	result := config.DB.Order("uploaded_at desc").Find(&files)
	if result.Error != nil {
		http.Error(w, "Error retrieving files", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

// GetFileByID retrieves a specific file by ID
func GetFileByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	
	var file models.File
	result := config.DB.First(&file, id)
	
	if result.Error != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(file)
}

// UploadFile handles file uploads to MinIO and records metadata in the database
func UploadFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(8 << 20) 
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	
	if handler.Size > 300 << 20 { // 300MB
		http.Error(w, "File is too large (max 300MB)", http.StatusBadRequest)
		return
	}
	
	bucketName := "uploads"
	objectName := fmt.Sprintf("%d-%s", time.Now().Unix(), handler.Filename)
	contentType := handler.Header.Get("Content-Type")
	
	_, err = config.MinioClient.PutObject(
		context.Background(),
		bucketName,
		objectName,
		file,
		handler.Size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error uploading file to MinIO: %v", err), http.StatusInternalServerError)
		return
	}
	
	fileRecord := models.File{
		Filename:   handler.Filename,
		Filepath:   objectName,
		Size:       handler.Size,
		UploadedAt: time.Now(),
	}
	
	result := config.DB.Create(&fileRecord)
	if result.Error != nil {
		http.Error(w, fmt.Sprintf("Error saving file metadata: %v", result.Error), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "File uploaded successfully",
		"file":    fileRecord,
	})
}

// DeleteFile removes a file from both storage and database
func DeleteFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	
	var file models.File
	result := config.DB.First(&file, id)
	
	if result.Error != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	
	err := config.MinioClient.RemoveObject(
		context.Background(),
		"uploads",
		file.Filepath,
		minio.RemoveObjectOptions{},
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error removing file from storage: %v", err), http.StatusInternalServerError)
		return
	}
	
	config.DB.Delete(&file)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "File deleted successfully",
	})
}

// MultipartUploadInfo tracks in-progress multipart uploads
type MultipartUploadInfo struct {
	UploadID    string
	FileName    string
	Bucket      string
	ObjectName  string
	ContentType string
	Parts       []minio.CompletePart
}

// In-memory store for multipart uploads
var multipartUploads = make(map[string]*MultipartUploadInfo)

// InitiateMultipartUpload starts a new multi-part upload process
func InitiateMultipartUpload(w http.ResponseWriter, r *http.Request) {
	var request struct {
		FileName    string `json:"fileName"`
		ContentType string `json:"contentType"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	bucketName := "uploads"
	objectName := fmt.Sprintf("%d-%s", time.Now().Unix(), request.FileName)

	// Use MinioCore from config for multipart operations
	uploadID, err := config.MinioCore.NewMultipartUpload(
		context.Background(),
		bucketName,
		objectName,
		minio.PutObjectOptions{ContentType: request.ContentType},
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error initiating multipart upload: %v", err), http.StatusInternalServerError)
		return
	}

	// Store upload info
	multipartUploads[uploadID] = &MultipartUploadInfo{
		UploadID:    uploadID,
		FileName:    request.FileName,
		Bucket:      bucketName,
		ObjectName:  objectName,
		ContentType: request.ContentType,
		Parts:       []minio.CompletePart{},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"uploadId":   uploadID,
		"objectName": objectName,
	})
}

// UploadPart handles one part of a multi-part upload
func UploadPart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uploadID := vars["uploadId"]

	uploadInfo, exists := multipartUploads[uploadID]
	if !exists {
		http.Error(w, "Upload not found or expired", http.StatusNotFound)
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	partNumberStr := r.FormValue("partNumber")
	partNumber, err := strconv.Atoi(partNumberStr)
	if err != nil || partNumber < 1 {
		http.Error(w, "Invalid part number", http.StatusBadRequest)
		return
	}

	filePart, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file part", http.StatusBadRequest)
		return
	}
	defer filePart.Close()

	// Use MinioCore from config for multipart operations
	partInfo, err := config.MinioCore.PutObjectPart(
		context.Background(),
		uploadInfo.Bucket,
		uploadInfo.ObjectName,
		uploadID,
		partNumber,
		filePart,
		handler.Size,
		minio.PutObjectPartOptions{},
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error uploading part: %v", err), http.StatusInternalServerError)
		return
	}

	uploadInfo.Parts = append(uploadInfo.Parts, minio.CompletePart{
		PartNumber: partNumber,
		ETag:       partInfo.ETag,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"partNumber": partNumber,
		"etag":       partInfo.ETag,
	})
}

// CompleteMultipartUpload finalizes a multi-part upload
func CompleteMultipartUpload(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uploadID := vars["uploadId"]

	uploadInfo, exists := multipartUploads[uploadID]
	if !exists {
		http.Error(w, "Upload not found or expired", http.StatusNotFound)
		return
	}

	sort.Slice(uploadInfo.Parts, func(i, j int) bool {
		return uploadInfo.Parts[i].PartNumber < uploadInfo.Parts[j].PartNumber
	})

	// Use MinioCore from config for multipart operations
	_, err := config.MinioCore.CompleteMultipartUpload(
		context.Background(),
		uploadInfo.Bucket,
		uploadInfo.ObjectName,
		uploadID,
		uploadInfo.Parts,
		minio.PutObjectOptions{},
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error completing multipart upload: %v", err), http.StatusInternalServerError)
		return
	}

	// Get object info for file size
	objInfo, err := config.MinioClient.StatObject(
		context.Background(), 
		uploadInfo.Bucket, 
		uploadInfo.ObjectName, 
		minio.StatObjectOptions{},
	)
	
	var totalSize int64 = 0
	if err == nil {
		totalSize = objInfo.Size
	}

	fileRecord := models.File{
		Filename:   uploadInfo.FileName,
		Filepath:   uploadInfo.ObjectName,
		Size:       totalSize,
		UploadedAt: time.Now(),
	}

	result := config.DB.Create(&fileRecord)
	if result.Error != nil {
		http.Error(w, fmt.Sprintf("Error saving file metadata: %v", result.Error), http.StatusInternalServerError)
		return
	}

	delete(multipartUploads, uploadID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Multipart upload completed successfully",
		"file":    fileRecord,
	})
}

// AbortMultipartUpload cancels an in-progress multi-part upload
func AbortMultipartUpload(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uploadID := vars["uploadId"]

	uploadInfo, exists := multipartUploads[uploadID]
	if !exists {
		http.Error(w, "Upload not found or already completed", http.StatusNotFound)
		return
	}

	// Use MinioCore from config for multipart operations
	err := config.MinioCore.AbortMultipartUpload(
		context.Background(),
		uploadInfo.Bucket,
		uploadInfo.ObjectName,
		uploadID,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error aborting multipart upload: %v", err), http.StatusInternalServerError)
		return
	}

	delete(multipartUploads, uploadID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Upload aborted successfully",
	})
}

// DownloadFile retrieves a file from MinIO storage and sends it to the client
func DownloadFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	
	var file models.File
	result := config.DB.First(&file, id)
	if result.Error != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	
	bucketName := "uploads"
	object, err := config.MinioClient.GetObject(
		context.Background(),
		bucketName,
		file.Filepath,
		minio.GetObjectOptions{},
	)
	if err != nil {
		http.Error(w, "Error retrieving file from storage", http.StatusInternalServerError)
		return
	}
	defer object.Close()
	
	objInfo, err := object.Stat()
	if err != nil {
		http.Error(w, "Error retrieving file information", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.Filename))
	w.Header().Set("Content-Type", objInfo.ContentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", objInfo.Size))
	
	if _, err := io.Copy(w, object); err != nil {
		http.Error(w, "Error streaming file", http.StatusInternalServerError)
		return
	}
}