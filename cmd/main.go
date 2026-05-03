package main

import (
	"log"
	"net/http"

	"github.com/blueberry-adii/CineMe/internal/booking"
	"github.com/blueberry-adii/CineMe/internal/redis"
	"github.com/blueberry-adii/CineMe/internal/utils"
)

func main() {

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./static")))
	mux.HandleFunc("GET /api/movies", listMovies)

	redisStore := booking.NewRedisStore(redis.NewRedisClient("localhost:6379"))
	service := booking.NewService(redisStore)
	handler := booking.NewHandler(service)

	mux.HandleFunc("GET /api/health", getHealth)

	mux.HandleFunc("GET /api/movies/{movieId}/bookings", handler.ListBookings)
	mux.HandleFunc("POST /api/movies/{movieId}/seats/{seatId}/hold", handler.HoldSeat)

	mux.HandleFunc("PUT /api/sessions/{sessionID}/confirm", handler.ConfirmSession)
	mux.HandleFunc("DELETE /api/sessions/{sessionID}", handler.ReleaseSession)

	log.Println("Server listening on localhost")
	if err := http.ListenAndServe(":80", mux); err != nil {
		panic(err)
	}
}

func listMovies(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, booking.Movies)
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "ok",
	}
	utils.WriteJSON(w, http.StatusOK, response)
}
