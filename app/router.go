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
	r.HandleFunc("/register", registerHandler).Methods("POST")
    r.HandleFunc("/login", loginHandler).Methods("POST")
	return r
}
