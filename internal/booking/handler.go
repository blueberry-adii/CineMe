package booking

import (
	"net/http"
	"time"

	"github.com/blueberry-adii/CineMe/internal/utils"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc}
}

func (h *Handler) ListBookings(w http.ResponseWriter, r *http.Request) {
	movieId := r.PathValue("movieId")
	bookings := h.svc.ListBookings(movieId)
	seats := make([]seatInfo, 0, len(bookings))

	for _, b := range bookings {
		if b.MovieID == movieId {
			seats = append(seats, seatInfo{
				SeatId:    b.SeatID,
				UserId:    b.UserID,
				Booked:    true,
				ExpiresAt: b.ExpiresAt,
			})
		}
	}
	utils.WriteJSON(w, 200, seats)
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

	session, err := h.svc.Book(booking)

	if err != nil {
		utils.WriteJSON(w, 400, err.Error())
		return
	}

	type holdResponse struct {
		SessionId string `json:"session_id"`
		MovieID   string `json:"movie_id"`
		SeatId    string `json:"seat_id"`
		ExpiresAt string `json:"expires_at"`
	}

	utils.WriteJSON(w, 200, holdResponse{
		SessionId: session.ID,
		MovieID:   session.MovieID,
		SeatId:    session.SeatID,
		ExpiresAt: session.ExpiresAt.Format(time.RFC3339),
	})
}

func (h *Handler) ConfirmSession(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, 200, "Ok")
}

func (h *Handler) ReleaseSession(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, 200, "Ok")
}

type seatInfo struct {
	SeatId    string    `json:"seat_id"`
	UserId    string    `json:"user_id"`
	Booked    bool      `json:"booked"`
	ExpiresAt time.Time `json:"expires_at"`
}
