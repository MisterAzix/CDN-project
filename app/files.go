package app

import (
	"context"
	"fmt"
	"log" // Add log package
	"net/http"
	"os"
	"path/filepath"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	UPLOAD_DIR = "datas" // Répertoire de téléchargement des fichiers
)

type FileMetadata struct {
	ID            primitive.ObjectID     `bson:"_id,omitempty"`       // ID du fichier
	FileName      string                 `bson:"file_name"`           // Nom du fichier
	FileSize      int64                  `bson:"file_size"`           // Taille du fichier
	FileType      string                 `bson:"file_type"`           // Type du fichier
	FilePath      string                 `bson:"file_path"`           // Chemin du fichier
	UploadedAt    time.Time              `bson:"uploaded_at"`         // Date de téléchargement
	UpdatedAt     time.Time              `bson:"updated_at"`          // Date de mise à jour
	UploaderID    string                 `bson:"uploader_id"`         // ID de l'uploader
	Metadata      map[string]interface{} `bson:"metadata"`            // Métadonnées
	Status        string                 `bson:"status"`              // Statut du fichier
	AccessControl AccessControl          `bson:"access_control"`      // Contrôle d'accès
	ParentID      string                 `bson:"parent_id,omitempty"` // ID du dossier parent
}

type AccessControl struct {
	Public      bool     `bson:"public"`      // Public ou non
	Permissions []string `bson:"permissions"` // Permissions
}

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting file upload process")
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Println("Error parsing multipart form:", err)
		http.Error(w, "File too big", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println("Error getting file from form:", err)
		http.Error(w, "Failed to get file", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	log.Println("File received:", handler.Filename)

	parentID := r.FormValue("parent_id")
	dirPath := filepath.Join(UPLOAD_DIR, parentID)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		log.Println("Error creating directory:", err)
		http.Error(w, "Failed to create directory", http.StatusInternalServerError)
		return
	}
	log.Println("Directory created:", dirPath)

	filePath := filepath.Join(dirPath, handler.Filename)
	outFile, err := os.Create(filePath)
	if err != nil {
		log.Println("Error creating file:", err)
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()
	log.Println("File created:", filePath)

	if _, err := outFile.ReadFrom(file); err != nil {
		log.Println("Error saving file:", err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	log.Println("File saved successfully")

	fileInfo, err := outFile.Stat()
	if err != nil {
		log.Println("Error getting file info:", err)
		http.Error(w, "Failed to get file info", http.StatusInternalServerError)
		return
	}

	db := client.Database("file_manager")
	collection := db.Collection("files")
	result, err := collection.InsertOne(context.TODO(), FileMetadata{
		FileName:   handler.Filename,
		FileSize:   fileInfo.Size(),
		FileType:   handler.Header.Get("Content-Type"),
		FilePath:   filePath,
		UploadedAt: time.Now(),
		UpdatedAt:  time.Now(),
		UploaderID: "user_identifier", // Remplacez par l'ID réel de l'uploader
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
		log.Println("Error saving file metadata:", err)
		http.Error(w, "Failed to save file metadata", http.StatusInternalServerError)
		return
	}
	log.Println("File metadata saved successfully with ID:", result.InsertedID)

	fmt.Fprintf(w, "File uploaded successfully: %s\n", handler.Filename)
}

