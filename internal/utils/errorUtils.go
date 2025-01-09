package utils

import "net/http"

func RegisterHandling(firstName, lastName, email, password string) {
	var w http.ResponseWriter

	if firstName == "" || lastName == "" || email == "" || password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	if !IsEmailValid(email) || len(password) < 8 || len(firstName) < 3 || len(lastName) < 3 {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}
	if PasswordExists(password) {
		http.Error(w, "Account Password already exists", http.StatusBadRequest)
		return
	}

	if EmailExists(email) {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}
}
