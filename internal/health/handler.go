package health

import (
	"net/http"

	"github.com/blueberry-adii/CineMe/internal/utils"
)

func GetHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "ok",
	}
	utils.WriteJSON(w, http.StatusOK, response)
}
