package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/marcosalbano/platform-api/internal/system"
)

func NetworkHandler(w http.ResponseWriter, r *http.Request) {

	network, err := system.GetNetwork()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(network); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
