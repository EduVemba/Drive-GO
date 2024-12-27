package Controllers

import (
	"d_uber_golang/internal/Services"
	"net/http"
)

func CreateTravelIntent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post required", http.StatusBadRequest)
		return
	}
	Services.HandlerCreateTravelIntent(w, r)
}

func GetTravelIntent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Get required", http.StatusBadRequest)
	}
	Services.HandlerAcceptRequester(w, r)
}
