package Authentication

import (
	"errors"
	"net/http"
)

var AuthError = errors.New("Unauthorized")

func Authorize(r *http.Request) error {
	email := r.FormValue("email")
	user, ok := users[email]
	if !ok {
		return AuthError
	}

	st, err := r.Cookie("session_token")
	if err != nil || st.Value == "" || st.Value != user.SessionToken {
		return AuthError
	}

	csrf := r.Header.Get("X-Csrf-Token")
	if csrf != user.CSRFToken || csrf == "" {
		return AuthError
	}
	return nil
}
