package Authentication

import (
	"d_uber_golang/internal/Database/PostgreSQL"
	"d_uber_golang/internal/models"
	"database/sql"
	"errors"
	"net/http"
)

var AuthError = errors.New("Unauthorized")

func Authorize(r *http.Request) error {
	//TODO: For some reason the Authorize function still giving errors
	st, err := r.Cookie("session_token")
	if err != nil || st.Value == "" {
		return AuthError
	}

	var user models.Person
	err = PostgreSQL.Db.QueryRow("SELECT email, sessiontoken, csrftoken FROM requester WHERE sessiontoken = $1", st.Value).
		Scan(&user.Email, &user.SessionToken, &user.CSRFToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return AuthError
		}
		return err
	}

	if st.Value != user.SessionToken {
		return AuthError
	}

	csrf := r.Header.Get("X-Csrf-Token")
	if csrf != user.CSRFToken || csrf == "" {
		return AuthError
	}
	return nil
}
