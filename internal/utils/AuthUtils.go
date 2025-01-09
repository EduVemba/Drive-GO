package utils

import (
	"d_uber_golang/internal/Database/PostgreSQL"
	"database/sql"
	"log"
	"regexp"
)

func IsValidRegistration(r string) bool {
	registrationRegex := regexp.MustCompile("^\\d{2}-[A-Z]{2}-[A-Z]\\d$")
	return registrationRegex.MatchString(r)
}

func IsEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

func EmailExists(emailT string) bool {
	psqlUser := `SELECT email from requester where email = $1 `
	psqlDriver := `SELECT email from driver where email = $1 `

	var email string

	err := PostgreSQL.Db.QueryRow(psqlUser, emailT).Scan(&email)
	if err != nil {
		if err != sql.ErrNoRows {
			return false
		}

	} else {
		return true
	}

	err = PostgreSQL.Db.QueryRow(psqlDriver, emailT).Scan(&email)
	if err != nil {
		if err != sql.ErrNoRows {
			return false
		}
	} else {
		return true
	}

	return false
}

func PasswordExists(passwordT string) bool {
	psqlUser := `SELECT password from requester`
	psqlDriver := `SELECT password from driver`

	var hashedPassword string

	rows, err := PostgreSQL.Db.Query(psqlUser)
	if err != nil {
		log.Println("Error querying requester table:", err)
		return false
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&hashedPassword)
		if err != nil {
			log.Println("Error scanning hashed password:", err)
			continue
		}

		if CheckPasswordHash(passwordT, hashedPassword) {
			return true
		}
	}
	rows, err = PostgreSQL.Db.Query(psqlDriver)
	if err != nil {
		log.Println("Error querying driver table:", err)
		return false
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&hashedPassword)
		if err != nil {
			log.Println("Error scanning hashed password:", err)
			continue
		}

		if CheckPasswordHash(passwordT, hashedPassword) {
			return true
		}
	}
	return false
}
