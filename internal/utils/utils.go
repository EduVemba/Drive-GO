package utils

import (
	"crypto/rand"
	"d_uber_golang/internal/Database/PostgreSQL"
	"database/sql"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
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

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	//psql := `SELECT 1 from requester where password = $1`
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

func ValidUser(email, password string) bool {
	var hashedPassword string

	query := "SELECT password FROM requester WHERE email = $1"
	err := PostgreSQL.Db.QueryRow(query, email).Scan(&hashedPassword)
	if err != nil {

		return false
	}
	return CheckPasswordHash(password, hashedPassword)
}

func ValidDriver(email, password string) bool {
	var hashedPassword string

	query := "SELECT password FROM driver WHERE email = $1"
	err := PostgreSQL.Db.QueryRow(query, email).Scan(&hashedPassword)
	if err != nil {
		return false
	}
	return CheckPasswordHash(password, hashedPassword)
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
