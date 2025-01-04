package Services

import (
	"context"
	"d_uber_golang/internal/Database/MongoDB"
	"d_uber_golang/internal/Database/PostgreSQL"
	"database/sql"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strings"
	"time"
)

type UserRequirements struct {
	SessionToken string `json:"session_token"`
	Email        string `json:"email"`
	From         string `json:"from"`
	To           string `json:"to"`

	//Requester Person  `json:"requester"`
	TimeStamp time.Time `json:"time_stamp"`
}

var (
	//	mu   sync.Mutex
	user UserRequirements
)

func HandlerCreateTravelIntent(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user UserRequirements

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	/*
		if user.From == "" || user.To == "" || user.Email == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}
	*/

	user.SessionToken = r.Header.Get("Authorization")
	if user.SessionToken == "" {
		http.Error(w, "Unauthorized - Missing session token", http.StatusUnauthorized)
		return
	}

	//	user.SessionToken = user.SessionToken[len("Bearer "):]

	tokenString := strings.TrimPrefix(user.SessionToken, "Bearer ")

	var firstName, lastName, email string
	err := PostgreSQL.Db.QueryRow(`SELECT first_name, last_name, email FROM requester WHERE sessiontoken = $1`, tokenString).
		Scan(&firstName, &lastName, &email)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Unauthorized - Invalid session token", http.StatusUnauthorized)
		} else {
			http.Error(w, "Database Error", http.StatusBadRequest)
		}
		return
	}

	response := map[string]interface{}{
		"first_name": firstName,
		"last_name":  lastName,
		"email":      email,
		"from":       user.From,
		"to":         user.To,
		"time_stamp": user.TimeStamp.Format(time.RFC3339), // Formating to JSON
	}

	err = MongoDB.InsertOneRequest(response)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
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

	//userEmail := `SELECT `

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
