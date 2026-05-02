package booking

import (
	"net/http"

	"github.com/blueberry-adii/CineMe/internal/utils"
)

type Handler struct {
}

type IHandler interface {
	GetMovies(w http.ResponseWriter, r *http.Request)
	GetMovieById(w http.ResponseWriter, r *http.Request)
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) GetMovies(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, 200, "Ok")
}

func (h *Handler) GetMovieById(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, 200, "Ok")
}
