package app

import (
	"context"
	"fmt"
    "encoding/json"

	// "image"
	"log"
	"net/http"
	"os"

	// "os/exec"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/v2/bson"
	
	// _ "golang.org/x/image/bmp"
	// _ "golang.org/x/image/tiff"
	// _ "golang.org/x/image/webp"
)

const (
	UPLOAD_DIR = "datas"
)

// ----------------------------------------------------- //
// ------------------- Interfaces ---------------------- //
// ----------------------------------------------------- //

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
	ParentID      bson.ObjectID          `bson:"parent_id,omitempty"`
}

type Folder struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	Name      string        `bson:"name"`
	UploaderID bson.ObjectID `bson:"uploader_id"`
	ParentID   bson.ObjectID `bson:"parent_id,omitempty"`
	CreatedAt time.Time     `bson:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"`
}

type FolderWithFiles struct {
	Folder  Folder        `bson:"folder"`
	Files   []FileMetadata `bson:"files"`
	Subfolders []FolderWithFiles `bson:"subfolders"`
}

type AccessControl struct {
	Public      bool     `bson:"public"`
	Permissions []string `bson:"permissions"`
}


func buildDirectoryPath(uploaderID, parentID bson.ObjectID) (string, error) {
	db := client.Database("file_manager")
	collection := db.Collection("folders")

	var pathParts []string
	currentID := parentID

	// Traverse the hierarchy of parent folders
	for currentID != uploaderID {
		var folder Folder
		err := collection.FindOne(context.TODO(), bson.M{"_id": currentID}).Decode(&folder)
		if err != nil {
			return "", err
		}
		pathParts = append([]string{folder.ID.Hex()}, pathParts...)
		currentID = folder.ParentID
	}

	// Add the uploader ID as the base directory
	pathParts = append([]string{uploaderID.Hex()}, pathParts...)

	// Construct the full directory path
	dirPath := filepath.Join(UPLOAD_DIR, filepath.Join(pathParts...))
	return dirPath, nil
}
// ----------------------------------------------------- //
// ------------- Métadonnées Fichiers ----------------- //
// ----------------------------------------------------- //

func generateTag(fileExtension string) string {
	switch strings.ToLower(fileExtension) {
	case ".pdf":
		return "pdf"
	case ".doc", ".docx":
		return "doc"
	case ".xls", ".xlsx", ".csv":
		return "sheet"
	case ".ppt", ".pptx":
		return "slide"
	case ".mp3", ".wav", ".flac":
		return "audio"
	case ".mp4", ".avi", ".mov":
		return "video"
	case ".jpeg", ".jpg", ".png", ".bmp", ".tiff", ".webp":
		return "image"
	default:
		return "other"
	}
}

// func getImageDimensions(filePath string) (int, int, error) {
// 	file, err := os.Open(filePath)
// 	m, _, err := image.Decode(reader)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		bounds := m.Bounds()
// 		w := bounds.Dx()
// 		h := bounds.Dy()
// 	return img.w, img.h, nil
// }

// func getVideoDuration(filePath string) (time.Duration, error) {
// 	cmd := exec.Command("ffmpeg", "-i", filePath)
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		return 0, err
// 	}

// 	durationStr := "0s"
// 	for _, line := range strings.Split(string(output), "\n") {
// 		if strings.Contains(line, "Duration:") {
// 			parts := strings.Split(line, ",")
// 			for _, part := range parts {
// 				if strings.HasPrefix(part, "Duration:") {
// 					durationStr = strings.TrimSpace(strings.Split(part, "Duration:")[1])
// 					break
// 				}
// 			}
// 			break
// 		}
// 	}

// 	duration, err := time.ParseDuration(durationStr)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return duration, nil
// }


// ---------------------------------------------------------- //
// ------------- Upload Fichiers / Dossiers ----------------- //
// ---------------------------------------------------------- //

func createUserFolder(uploaderID bson.ObjectID) error {
	uploaderIDString := uploaderID.Hex()
	dirPath := filepath.Join(UPLOAD_DIR, uploaderIDString)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func createFolderHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting folder creation process")
	name := r.FormValue("name")
	uploaderIDStr := r.FormValue("uploader_id")
	parentIDStr := r.FormValue("parent_id")

	uploaderID, err := bson.ObjectIDFromHex(uploaderIDStr)
	if err != nil {
		http.Error(w, "Invalid uploader ID", http.StatusBadRequest)
		return
	}

	var parentID bson.ObjectID
	if parentIDStr == "" {
		parentID = uploaderID
	} else {
		parentID, err = bson.ObjectIDFromHex(parentIDStr)
		if err != nil {
			http.Error(w, "Invalid parent ID", http.StatusBadRequest)
			return
		}
	}

	db := client.Database("file_manager")
	collection := db.Collection("folders")

	folder := Folder{
		ID:         bson.NewObjectID(),
		Name:       name,
		UploaderID: uploaderID,
		ParentID:   parentID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	_, err = collection.InsertOne(context.TODO(), folder)
	if err != nil {
		http.Error(w, "Failed to create folder in database", http.StatusInternalServerError)
		return
	}

	// Build the full directory path considering the hierarchy of parent folders
	directoryPath, err := buildDirectoryPath(uploaderID, folder.ID)
	if err != nil {
		http.Error(w, "Failed to build directory path", http.StatusInternalServerError)
		return
	}

	if err := os.MkdirAll(directoryPath, os.ModePerm); err != nil {
		http.Error(w, "Failed to create directory", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Folder created successfully: %s\n", name)
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

	parentIDString := r.FormValue("parent_id")
	uploaderIDString := r.FormValue("uploader_id")

	parentID, err := bson.ObjectIDFromHex(parentIDString)
	if err != nil {
		http.Error(w, "Invalid parent ID", http.StatusBadRequest)
		return
	}

	uploaderID, err := bson.ObjectIDFromHex(uploaderIDString)
	if err != nil {
		http.Error(w, "Invalid uploader ID", http.StatusBadRequest)
		return
	}

	// Construct the directory path using uploader ID and parent ID
	dirPath, err := buildDirectoryPath(uploaderID, parentID)
	if err != nil {
		http.Error(w, "Failed to build directory path", http.StatusInternalServerError)
		return
	}

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

	if _, err := io.Copy(outFile, file); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	if err := outFile.Sync(); err != nil {
		http.Error(w, "Failed to sync file", http.StatusInternalServerError)
		return
	}

	fileInfo, err := outFile.Stat()
	if err != nil {
		http.Error(w, "Failed to get file info", http.StatusInternalServerError)
		return
	}

	fileExtension := filepath.Ext(handler.Filename)
	tag := generateTag(fileExtension)

	db := client.Database("file_manager")
	collection := db.Collection("files")
	_, err = collection.InsertOne(context.TODO(), FileMetadata{
		FileName:   handler.Filename,
		FileSize:   fileInfo.Size(),
		FileType:   handler.Header.Get("Content-Type"),
		FilePath:   filePath,
		UploadedAt: time.Now(),
		UpdatedAt:  time.Now(),
		UploaderID: uploaderID,
		Metadata: map[string]interface{}{
			"width":    0,
			"height":   0,
			"duration": 0,
			"tags":     []string{tag},
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



// ---------------------------------------------------------- //
// ------------- Get Fichiers / Dossiers ----------------- //
// ---------------------------------------------------------- //

func fetchFoldersHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting folder fetching process")
	uploaderIDStr := r.URL.Query().Get("uploader_id")
	uploaderID, err := bson.ObjectIDFromHex(uploaderIDStr)
	if err != nil {
		http.Error(w, "Invalid uploader ID", http.StatusBadRequest)
		return
	}

	db := client.Database("file_manager")
	collection := db.Collection("folders")
	// Fetch folders where the parent ID is the uploader ID
	filter := bson.M{"uploader_id": uploaderID, "parent_id": uploaderID}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		http.Error(w, "Failed to fetch folders", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var folders []Folder
	if err = cursor.All(context.TODO(), &folders); err != nil {
		http.Error(w, "Failed to decode folders", http.StatusInternalServerError)
		return
	}

	var foldersWithFiles []FolderWithFiles
	for _, folder := range folders {
		files, err := listFilesForFolder(folder.ID)
		if err != nil {
			http.Error(w, "Failed to list files", http.StatusInternalServerError)
			return
		}

		subfolders, err := fetchSubfolders(folder.ID)
		if err != nil {
			http.Error(w, "Failed to fetch subfolders", http.StatusInternalServerError)
			return
		}

		foldersWithFiles = append(foldersWithFiles, FolderWithFiles{
			Folder:      folder,
			Files:       files,
			Subfolders:  subfolders,
		})
	}

	log.Println("Successfully fetched folders")
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(foldersWithFiles)
	if err != nil {
		log.Printf("Error marshalling folders: %v", err)
		http.Error(w, "Failed to encode folders", http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(jsonData); err != nil {
		log.Printf("Error writing JSON data: %v", err)
		http.Error(w, "Failed to write JSON data", http.StatusInternalServerError)
	}
}





func listFilesForFolder(parentID bson.ObjectID) ([]FileMetadata, error) {
	log.Printf("Listing files for folder: %s", parentID)
	db := client.Database("file_manager")
	collection := db.Collection("files")

	filter := bson.M{"parent_id": parentID}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var files []FileMetadata
	if err = cursor.All(context.TODO(), &files); err != nil {
		return nil, err
	}

	log.Printf("Found %d files for folder: %s", len(files), parentID)
	return files, nil
}

func fetchSubfolders(parentID bson.ObjectID) ([]FolderWithFiles, error) {
	log.Printf("Fetching subfolders for folder: %s", parentID)
	db := client.Database("file_manager")
	collection := db.Collection("folders")

	filter := bson.M{"parent_id": parentID}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var subfolders []Folder
	if err = cursor.All(context.TODO(), &subfolders); err != nil {
		return nil, err
	}

	var subfoldersWithFiles []FolderWithFiles

	for _, subfolder := range subfolders {
		files, err := listFilesForFolder(subfolder.ID)
		if err != nil {
			return nil, err
		}

		nestedSubfolders, err := fetchSubfolders(subfolder.ID)
		if err != nil {
			return nil, err
		}

		subfoldersWithFiles = append(subfoldersWithFiles, FolderWithFiles{
			Folder:      subfolder,
			Files:       files,
			Subfolders:  nestedSubfolders,
		})
	}

	log.Printf("Found %d subfolders for folder: %s", len(subfoldersWithFiles), parentID)
	return subfoldersWithFiles, nil
}

func serveFileHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting file serving process")
	fileID := mux.Vars(r)["id"]

	db := client.Database("file_manager")
	collection := db.Collection("files")

	objID, err := bson.ObjectIDFromHex(fileID)
	if err != nil {
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	var fileMetadata FileMetadata
	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&fileMetadata)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	filePath := fileMetadata.FilePath
	http.ServeFile(w, r, filePath)
}


// ---------------------------------------------------------- //
// ------------- Delete Fichiers / Dossiers ----------------- //
// ---------------------------------------------------------- //
func deleteFolderHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting folder deletion process")
	uploaderIDStr := r.FormValue("uploader_id")
	folderIDStr := r.FormValue("id")

	uploaderID, err := bson.ObjectIDFromHex(uploaderIDStr)
	if err != nil {
		http.Error(w, "Invalid uploader ID", http.StatusBadRequest)
		return
	}

	folderID, err := bson.ObjectIDFromHex(folderIDStr)

	db := client.Database("file_manager")
	collection := db.Collection("folders")

	if err != nil {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	// ------------- Suppression des fichiers serveur ------------ //

	// First, delete all files within the folder
	err = deleteFilesInFolder(folderIDStr)
	if err != nil {
		http.Error(w, "Failed to delete files in folder", http.StatusInternalServerError)
		return
	}
	var dirPath string
	dirPath, err = buildDirectoryPath(uploaderID, folderID)
	if err != nil {
		http.Error(w, "Failed to build directory path", http.StatusInternalServerError)
		return
	}
	if err := os.RemoveAll(dirPath); err != nil {
		http.Error(w, "Failed to delete folder directory", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Folder and its files deleted successfully"))

	// ------------- Suppression dans la base ------------------- //

	// Then, delete the folder itself from the database
	_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": folderID})
	if err != nil {
		http.Error(w, "Failed to delete folder", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Folder and its files deleted successfully"))
}

func deleteFilesInFolder(folderID string) error {
	db := client.Database("file_manager")
	collection := db.Collection("files")

	filter := bson.M{"parent_id": folderID}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return err
	}
	defer cursor.Close(context.TODO())

	var files []FileMetadata
	if err = cursor.All(context.TODO(), &files); err != nil {
		return err
	}

	for _, file := range files {
		err = deleteFile(file.ID.Hex())
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteFile(fileID string) error {
	db := client.Database("file_manager")
	collection := db.Collection("files")

	objID, err := bson.ObjectIDFromHex(fileID)
	if err != nil {
		return err
	}

	var fileMetadata FileMetadata
	err = collection.FindOneAndDelete(context.TODO(), bson.M{"_id": objID}).Decode(&fileMetadata)
	if err != nil {
		return err
	}

	err = os.Remove(fileMetadata.FilePath)
	if err != nil {
		return err
	}

	return nil
}

func deleteFileHandler(w http.ResponseWriter, r *http.Request) {
	fileID := mux.Vars(r)["id"]
	err := deleteFile(fileID)
	if err != nil {
		http.Error(w, "Failed to delete file", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("File deleted successfully"))
}
