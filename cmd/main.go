package main

import (
	"log"
	"net/http"

	"github.com/blueberry-adii/CineMe/internal/booking"
	"github.com/blueberry-adii/CineMe/internal/redis"
)

func main() {
	/*
	* Mux which serves on http requests
	 */
	mux := http.NewServeMux()

	/*
	* "/" serving frontend static files
	 */
	mux.Handle("/", http.FileServer(http.Dir("./static")))

	/*
	* Dependency Injection - HTTP Handler Depends on Service
	* - Service depends on Redis Store
	* - Redis Store depends on Redis Client
	 */
	redisStore := booking.NewRedisStore(redis.NewRedisClient("localhost:6379"))
	service := booking.NewService(redisStore)
	handler := booking.NewHandler(service)

	/*
	* returns list of all movies
	 */
	mux.HandleFunc("GET /api/movies", handler.ListMovies)

	/*
	* Basic Health Endpoints
	 */
	mux.HandleFunc("GET /api/health", handler.GetHealth)

	/*
	* Get All Bookings
	 */
	mux.HandleFunc("GET /api/movies/{movieId}/bookings", handler.ListBookings)

	/*
	* Hold a Seat
	 */
	mux.HandleFunc("POST /api/movies/{movieId}/seats/{seatId}/hold", handler.HoldSeat)

	/*
	* Confirm Seat
	 */
	mux.HandleFunc("PUT /api/sessions/{sessionID}/confirm", handler.ConfirmSession)

	/*
	* Release Seat from Hold
	 */
	mux.HandleFunc("DELETE /api/sessions/{sessionID}", handler.ReleaseSession)

	/*
	* Server listens to port 80 (HTTP)
	 */
	log.Println("Server listening on localhost")
	if err := http.ListenAndServe(":80", mux); err != nil {
		panic(err)
	}
}
