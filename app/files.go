package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	UPLOAD_DIR = "datas"
)

type FileMetadata struct {
	ID            bson.ObjectID          `bson:"_id,omitempty"`
	FileName      string                 `bson:"file_name"`
	FileSize      int64                  `bson:"file_size"`
	FileType      string                 `bson:"file_type"`
	FilePath      string                 `bson:"file_path"`
	UploadedAt    time.Time              `bson:"uploaded_at"`
	UpdatedAt     time.Time              `bson:"updated_at"`
	UploaderID    bson.ObjectID          `bson:"uploader_id"`
	Metadata      map[string]interface{} `bson:"metadata"`
	Status        string                 `bson:"status"`
	AccessControl AccessControl          `bson:"access_control"`
	ParentID      string                 `bson:"parent_id,omitempty"`
}

type AccessControl struct {
	Public      bool     `bson:"public"`
	Permissions []string `bson:"permissions"`
}

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting file upload process")
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "File too big", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	parentID := r.FormValue("parent_id")
	dirPath := filepath.Join(UPLOAD_DIR, parentID)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		http.Error(w, "Failed to create directory", http.StatusInternalServerError)
		return
	}

	filePath := filepath.Join(dirPath, handler.Filename)
	outFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	if _, err := outFile.ReadFrom(file); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	fileInfo, err := outFile.Stat()
	if err != nil {
		http.Error(w, "Failed to get file info", http.StatusInternalServerError)
		return
	}

	db := client.Database("file_manager")
	collection := db.Collection("files")
	_, err = collection.InsertOne(context.TODO(), FileMetadata{
		FileName:   handler.Filename,
		FileSize:   fileInfo.Size(),
		FileType:   handler.Header.Get("Content-Type"),
		FilePath:   filePath,
		UploadedAt: time.Now(),
		UpdatedAt:  time.Now(),
		UploaderID: bson.NewObjectID(),
		Metadata: map[string]interface{}{
			"width":    1920,
			"height":   1080,
			"duration": nil,
			"tags":     []string{"example", "image", "cdn"},
		},
		Status: "active",
		AccessControl: AccessControl{
			Public:      true,
			Permissions: []string{"read", "download"},
		},
		ParentID: parentID,
	})
	if err != nil {
		http.Error(w, "Failed to save file metadata", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %s\n", handler.Filename)
}

func deleteFileHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting file deletion process")
	fileID := mux.Vars(r)["id"]

	db := client.Database("file_manager")
	collection := db.Collection("files")

	var fileMetadata FileMetadata
	objID, err := bson.ObjectIDFromHex(fileID)
	if err != nil {
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	err = collection.FindOneAndDelete(context.TODO(), bson.M{"_id": objID}).Decode(&fileMetadata)
	if err != nil {
		http.Error(w, "Failed to delete file metadata", http.StatusInternalServerError)
		return
	}

	err = os.Remove(fileMetadata.FilePath)
	if err != nil {
		http.Error(w, "Failed to delete file", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("File deleted successfully"))
}

func listFilesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting file listing process")
	parentID := r.URL.Query().Get("parent_id")

	db := client.Database("file_manager")
	collection := db.Collection("files")

	filter := bson.M{}
	if parentID != "" {
		filter["parent_id"] = parentID
	}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		http.Error(w, "Failed to list files", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var files []FileMetadata
	if err = cursor.All(context.TODO(), &files); err != nil {
		http.Error(w, "Failed to decode files", http.StatusInternalServerError)
		return
	}

	for _, file := range files {
		fmt.Fprintf(w, "ID: %s, Name: %s, Path: %s, ParentID: %s\n", file.ID.Hex(), file.FileName, file.FilePath, file.ParentID)
	}
}
