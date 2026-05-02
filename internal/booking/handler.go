package booking

import "net/http"

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) ListMovies(w http.ResponseWriter, r *http.Request) {

}
