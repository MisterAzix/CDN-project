package main

import (
	"log"
	"net/http"

	"hetic-cdn-project/app"
)

const PORT = ":8080"
const HOST = "http://localhost" + PORT

func main() {
	app.NewRouter()
	app.ConnectDB()

	log.Println("Server is running on", HOST)
	err := http.ListenAndServe(PORT, app.NewRouter())
	if err != nil {
		log.Fatal("Error starting server!", err)
		return
	}
}
