package main

import (
	"net/http"

	"github.com/blueberry-adii/CineMe/internal/booking"
	"github.com/blueberry-adii/CineMe/internal/health"
)

func main() {

	mux := http.NewServeMux()
	// redisClient := redis.NewRedisClient("localhost:6379")

	var handler booking.IHandler = booking.NewHandler()

	mux.Handle("/", http.FileServer(http.Dir("./static")))
	mux.HandleFunc("GET /api/health", health.GetHealth)

	mux.HandleFunc("GET /api/movies", handler.GetMovies)
	mux.HandleFunc("GET /api/movies/{movieId}", handler.GetMovieById)

	if err := http.ListenAndServe(":80", mux); err != nil {
		panic(err)
	}
}
