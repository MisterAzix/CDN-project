package main

import (
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"hetic-cdn-project/app"
	"log"
)

const PORT = ":8080"
const HOST = "http://localhost" + PORT

func main() {
	app.LoadEnv()
	app.ConnectDB()
	app.InitS3Client()
	app.InitAuth()
	router := app.NewRouter()

	// Metrics endpoint for Prometheus
	router.Handle("/metrics", promhttp.Handler())

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
