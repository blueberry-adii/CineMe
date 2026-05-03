package booking

import (
	"net/http"

	"github.com/blueberry-adii/CineMe/internal/utils"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc}
}

func (h *Handler) ListSeats(w http.ResponseWriter, r *http.Request) {
	movieId := r.PathValue("movieId")
	if Seats[movieId] == nil {
		utils.WriteJSON(w, 404, "Movie Not Listed")
		return
	}
	utils.WriteJSON(w, 200, Seats[movieId])
}

func (h *Handler) HoldSeat(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		UserID string `json:"user_id"`
	}
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteJSON(w, 400, "Invalid Request")
		return
	}

	booking := Booking{
		MovieID: r.PathValue("movieId"),
		SeatID:  r.PathValue("seatId"),
		UserID:  payload.UserID,
	}

	if err := h.svc.Book(booking); err != nil {
		utils.WriteJSON(w, 400, err.Error())
		return
	}

	utils.WriteJSON(w, 200, booking)
}

func (h *Handler) ConfirmSession(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, 200, "Ok")
}

func (h *Handler) ReleaseSession(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, 200, "Ok")
}
