package Services

import (
	"context"
	"d_uber_golang/internal/Database/MongoDB"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"sync"
	"time"
)

type UserRequirements struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	From      string `json:"from"`
	To        string `json:"to"`

	//Requester Person  `json:"requester"`
	TimeStamp time.Time `json:"time_stamp"`
}

var (
	mu   sync.Mutex
	user UserRequirements
)

func HandlerCreateTravelIntent(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	defer r.Body.Close()

	user.TimeStamp = time.Now()

	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"from":       user.From,
		"to":         user.To,
		"time_stamp": user.TimeStamp.Format(time.RFC3339), // Formating to JSON
	}

	err := MongoDB.InsertOneRequest(response)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

type DriverRequirements struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Registration string `json:"registration"`
}

func HandlerAcceptRequester(w http.ResponseWriter, r *http.Request) {

	var driver DriverRequirements
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&driver); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	defer r.Body.Close()

	if MongoDB.GetCollection("requests").FindOne(context.Background(), bson.M{"email": user.Email}).Err() != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"first_name":   driver.FirstName,
		"last_name":    driver.LastName,
		"email":        driver.Email,
		"registration": driver.Registration,
	}

	err := MongoDB.InsertOneResponse(response)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusAccepted)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

}
