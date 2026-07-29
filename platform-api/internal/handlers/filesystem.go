package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/marcosalbano/platform-api/internal/system"
)

func FilesystemHandler(w http.ResponseWriter, r *http.Request) {
	filesystems, err := system.GetFilesystems()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(filesystems); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
