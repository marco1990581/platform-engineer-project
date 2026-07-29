package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/marcosalbano/platform-api/internal/system"
)

func HostnameHandler(w http.ResponseWriter, r *http.Request) {

	hostname, err := system.GetHostname()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hostname)
}
