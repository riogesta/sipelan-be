package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sipelan/common"
	"time"
)

// UploadFile handles file uploads for transaction attachments.
func UploadFile(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form, 10 << 20 specifies a maximum upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Error(http.StatusBadRequest, "Error retrieving the file: "+err.Error()))
		return
	}
	defer file.Close()

	// Create a unique filename
	ext := filepath.Ext(handler.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join("uploads", filename)

	// Create the file on the server
	dst, err := os.Create(filePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(common.Error(http.StatusInternalServerError, "Error creating the file: "+err.Error()))
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the destination file
	if _, err := io.Copy(dst, file); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(common.Error(http.StatusInternalServerError, "Error saving the file: "+err.Error()))
		return
	}

	// Return the file path/URL
	fileURL := "/uploads/" + filename
	json.NewEncoder(w).Encode(common.Success("File uploaded successfully", map[string]string{
		"url": fileURL,
	}))
}
