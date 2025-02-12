package app

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var logger = GetLogger()


func handler(w http.ResponseWriter, r *http.Request) {
	logger.WithFields(logrus.Fields{
        "method": r.Method,
        "path":   r.URL.Path,
		"level" : "info",
    }).Info("200")
	fmt.Fprintf(w, "Hello, Go!")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server is healthy")
}

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", handler)
	r.HandleFunc("/health", healthHandler)
    r.HandleFunc("/file/upload", uploadFileHandler).Methods("POST")
	r.HandleFunc("/folder/upload", createFolderHandler).Methods("POST")
	r.HandleFunc("/file/delete/{id}", deleteFileHandler).Methods("DELETE")
	r.HandleFunc("/folder/delete/{id}", deleteFolderHandler).Methods("DELETE")
    r.HandleFunc("/fetch-folders", fetchFoldersHandler).Methods("GET")
	return r
}
