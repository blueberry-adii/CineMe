package main

import (
	"net/http"

	"github.com/blueberry-adii/CineMe/internal/health"
)

func main() {

	mux := http.NewServeMux()
	// redisClient := redis.NewRedisClient("localhost:6379")

	mux.Handle("/", http.FileServer(http.Dir("./static")))
	mux.HandleFunc("GET /api/health", health.GetHealth)

	if err := http.ListenAndServe(":80", mux); err != nil {
		panic(err)
	}
}
