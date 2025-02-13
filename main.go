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

<<<<<<< HEAD
	log.Println("Server is running on", HOST)
	err := http.ListenAndServe(PORT, corsHandler)
=======
	log.Println("Server is running on", HOST)
	err := http.ListenAndServe(PORT, corsHandler)
=======
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := corsHandler.Handler(router)
	
	log.Println("Server is running on", HOST)
	err := http.ListenAndServe(PORT, handler)
>>>>>>> 4ae975125751ec90717dba71fd37f4727adfe13d
	if err != nil {
		log.Fatal("Error starting server!", err)
		return
	}
}
