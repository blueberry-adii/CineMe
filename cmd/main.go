package main

import (
	"log"
	"net/http"

	"github.com/blueberry-adii/CineMe/internal/booking"
	"github.com/blueberry-adii/CineMe/internal/utils"
)

func main() {

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./static")))
	mux.HandleFunc("GET /api/movies", listMovies)

	// redisClient := redis.NewRedisClient("localhost:6379")
	handler := booking.NewHandler()

	mux.HandleFunc("GET /api/health", getHealth)

	mux.HandleFunc("GET /api/movies/{movieId}/seats", handler.ListSeats)
	mux.HandleFunc("GET /api/movies/{movieId}/seats/{seatId}/hold", handler.HoldSeat)

	log.Println("Server listening on localhost")
	if err := http.ListenAndServe(":80", mux); err != nil {
		panic(err)
	}
}

type movieResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Rows        int    `json:"rows"`
	SeatsPerRow int    `json:"seats_per_row"`
}

var movies = []movieResponse{
	{ID: "inception", Title: "Inception", Rows: 5, SeatsPerRow: 8},
	{ID: "interstellar", Title: "Interstellar", Rows: 4, SeatsPerRow: 6},
}

func listMovies(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, movies)
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "ok",
	}
	utils.WriteJSON(w, http.StatusOK, response)
}
