package main

import (
    "log"
    "net/http"
    "github.com/joho/godotenv"
    "hetic-cdn-project/app"
)

const PORT = ":8080"
const HOST = "http://localhost" + PORT

func main() {
    // Charger les variables d'environnement à partir du fichier .env
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    app.ConnectDB()
    router := app.NewRouter()

    log.Println("Server is running on", HOST)
    err = http.ListenAndServe(PORT, router)
    if err != nil {
        log.Fatal("Error starting server!", err)
        return
    }
}