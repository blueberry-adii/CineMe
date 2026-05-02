package main

import (
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	// redisClient := redis.NewRedisClient("localhost:6379")

	mux.Handle("/", http.FileServer(http.Dir("./static")))
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API Healthy"))
	})

	if err := http.ListenAndServe(":80", mux); err != nil {
		panic(err)
	}
}
