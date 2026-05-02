package booking

import (
	"net/http"

	"github.com/blueberry-adii/CineMe/internal/utils"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) ListSeats(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, 200, "Ok")
}

func (h *Handler) HoldSeat(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, 200, "Ok")
}

func (h *Handler) ConfirmSession(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, 200, "Ok")
}

func (h *Handler) ReleaseSession(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, 200, "Ok")
}
