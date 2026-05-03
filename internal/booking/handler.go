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
				SessionId: b.ID,
				Status:    b.Status,
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

	utils.WriteJSON(w, 200, seatInfo{
		SessionId: session.ID,
		SeatId:    session.SeatID,
		ExpiresAt: session.ExpiresAt,
		UserId:    session.UserID,
		Booked:    true,
		Status:    session.Status,
	})
}

func (h *Handler) ConfirmSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("sessionID")
	var payload struct {
		UserID string `json:"user_id"`
	}
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteJSON(w, 400, "Invalid Request")
		return
	}

	session, err := h.svc.Confirm(r.Context(), sessionID, payload.UserID)
	if err != nil {
		utils.WriteJSON(w, 400, err.Error())
		return
	}

	utils.WriteJSON(w, 200, session)
}

func (h *Handler) ReleaseSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("sessionID")
	var payload struct {
		UserID string `json:"user_id"`
	}
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteJSON(w, 400, "Invalid Request")
		return
	}

	if err := h.svc.Release(r.Context(), sessionID, payload.UserID); err != nil {
		utils.WriteJSON(w, 400, err.Error())
		return
	}

	utils.WriteJSON(w, 200, "Released")
}

type seatInfo struct {
	SeatId    string    `json:"seat_id"`
	UserId    string    `json:"user_id"`
	Booked    bool      `json:"booked"`
	ExpiresAt time.Time `json:"expires_at"`
	SessionId string    `json:"session_id"`
	Status    string    `json:"status"`
}
