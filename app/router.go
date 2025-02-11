package app

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Go!")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server is healthy")
}

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", handler)
	r.HandleFunc("/health", healthHandler)
    r.HandleFunc("/upload", uploadFileHandler).Methods("POST")

	return r
}
