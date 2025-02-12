package main

import (
	"hetic-cdn-project/app"
	"log"
	"net/http"
	"golang.org/x/time/rate"
	"time"
)

const PORT = ":8080"
const HOST = "http://localhost" + PORT

func main() {
	app.LoadEnv()
	app.ConnectDB()
	app.InitS3Client()
	router := app.NewRouter()

	// 5 requests per second with a burst of 10
    limiter := app.NewRateLimiter(rate.Every(time.Second/5), 10)
    limitedRouter := app.RateLimitMiddleware(limiter, router)

	log.Println("Server is running on", HOST)
	err := http.ListenAndServe(PORT, limitedRouter)
	if err != nil {
		log.Fatal("Error starting server!", err)
		return
	}
}
