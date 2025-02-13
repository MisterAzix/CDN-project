package main

import (
	"hetic-cdn-project/app"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
)

const PORT = ":8080"
const HOST = "http://localhost" + PORT

func main() {
	app.LoadEnv()
	app.ConnectDB()
	app.InitS3Client()
	app.InitAuth()
	router := app.NewRouter()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{os.Getenv("FRONTEND_URL")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(router)

	log.Println("Server is running on", HOST)
	err := http.ListenAndServe(PORT, corsHandler)
	if err != nil {
		log.Fatal("Error starting server!", err)
		return
	}
}
