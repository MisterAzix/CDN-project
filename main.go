package main

import (
	"hetic-cdn-project/app"
	"log"
	"net/http"
)

const PORT = ":8080"
const HOST = "http://localhost" + PORT

func main() {
	app.LoadEnv()
	app.ConnectDB()
	app.InitS3Client()
	app.InitRedisClient()
	router := app.NewRouter()
	
	log.Println("Server is running on", HOST)
	err := http.ListenAndServe(PORT, router)
	if err != nil {
		log.Fatal("Error starting server!", err)
		return
	}
}
