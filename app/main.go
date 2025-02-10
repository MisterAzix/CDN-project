package main

import (
    "fmt"
    "net/http"
)

func main() {
	NewRouter()
	fmt.Println("Serveur en cours sur http://localhost:8080")
    http.ListenAndServe(":8080", NewRouter())
}