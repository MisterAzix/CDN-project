package main

import (
	"github.com/joho/godotenv"
	"hetic-cdn-project/app"
	"log"
	"net/http"
)

const PORT = ":8080"
const HOST = "http://localhost" + PORT

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app.ConnectDB()
	app.InitS3Client()
	router := app.NewRouter()

	log.Println("Server is running on", HOST)
	err = http.ListenAndServe(PORT, router)
	if err != nil {
		log.Fatal("Error starting server!", err)
		return
	}
}
