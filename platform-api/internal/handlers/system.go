package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/marcosalbano/platform-api/internal/system"
)

func SystemHandler(w http.ResponseWriter, r *http.Request) {

	systemInfo, err := system.GetSystemInfo()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(systemInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
