package Authentication

import (
	"d_uber_golang/internal/models"
	"d_uber_golang/internal/routes"
	"d_uber_golang/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var Users = map[string]models.Person{}

var Drivers = map[string]models.DriverUser{}

/*
* Normal User Register function
 */
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Login should require a Post method", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
	}

	firstName := input.FirstName
	lastName := input.LastName
	email := input.Email
	password := input.Password

	if firstName == "" || lastName == "" || email == "" || password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	if !utils.IsEmailValid(email) || len(password) < 8 || len(firstName) < 3 || len(lastName) < 3 {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}

	if _, ok := Users[email]; ok {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	hashedPassword, _ := utils.HashPassword(password)
	Users[email] = models.Person{
		Password: hashedPassword,
	}
	// Insert of the user in the Postgres HERE
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User Created"))
}

/*
* Driver user Register function
 */
func RegisterDriver(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Register should require a Post method", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		Email        string `json:"email"`
		Password     string `json:"password"`
		Registration string `json:"registration"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
	}

	firstName := input.FirstName
	lastName := input.LastName
	email := input.Email
	password := input.Password
	registration := input.Registration

	if firstName == "" || lastName == "" || email == "" || password == "" || registration == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	if !utils.IsEmailValid(email) || len(password) < 8 ||
		len(firstName) < 3 ||
		len(lastName) < 3 ||
		!utils.IsValidRegistration(registration) {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}

	if _, ok := Drivers[email]; ok {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	hashedPassword, _ := utils.HashPassword(password)
	Drivers[email] = models.DriverUser{
		Password: hashedPassword,
	}
	// Insert of the user in the Postgres HERE
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Driver Created"))
}

/*
* Normal User Login function
 */
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Login should require a Post method", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
	}

	email := input.Email
	password := input.Password

	user, ok := Users[email]
	if !ok || !utils.CheckPasswordHash(password, user.Password) {
		http.Error(w, "User does not exist", http.StatusNotFound)
		return
	}
	sessionToken := utils.GenerateToken(32)
	csrfToken := utils.GenerateToken(32)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     routes.USER,
		Secure:   true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
	})

	user.SessionToken = sessionToken
	user.CSRFToken = csrfToken
	Users[email] = user

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Welcome Back!"))
}

/*
* Driver user Login function
 */
func LoginDriver(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Login should require a Post method", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
	}

	email := input.Email
	password := input.Password

	drivers, ok := Drivers[email]
	if !ok || !utils.CheckPasswordHash(password, drivers.Password) {
		http.Error(w, "Driver does not exist", http.StatusNotFound)
		return
	}

	sessionToken := utils.GenerateToken(32)
	csrfToken := utils.GenerateToken(32)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     routes.DRIVER,
		Secure:   true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
	})

	drivers.SessionToken = sessionToken
	drivers.CSRFToken = csrfToken
	Drivers[email] = drivers

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Welcome Back!"))
}

// This Portected func is used in the Services that Require cookies
func Protected(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Request Method.", http.StatusMethodNotAllowed)
		return
	}
	if err := Authorize(r); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	fname := r.FormValue("firstname")
	lname := r.FormValue("lastname")
	fmt.Fprintf(w, "CSRF validation successful! Welcome, %s %s", fname, lname)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if err := Authorize(r); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: false,
	})
	email := r.FormValue("email")
	user, _ := Users[email]
	user.SessionToken = ""
	user.CSRFToken = ""
	Users[email] = user

	w.WriteHeader(http.StatusNoContent)

}
