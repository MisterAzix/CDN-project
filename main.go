package main

import (
	"hetic-cdn-project/app"
	"log"
	"net/http"
	"github.com/rs/cors"
)

const PORT = ":8080"
const HOST = "http://localhost" + PORT

func main() {
	app.LoadEnv()
	app.ConnectDB()
	app.InitS3Client()
	router := app.NewRouter()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := corsHandler.Handler(router)
	
	log.Println("Server is running on", HOST)
	err := http.ListenAndServe(PORT, handler)
	if err != nil {
		log.Fatal("Error starting server!", err)
		return
	}
}
