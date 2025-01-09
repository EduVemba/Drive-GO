package utils

import (
	"crypto/rand"
	"d_uber_golang/internal/Database/PostgreSQL"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
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

func ContainsToken(token string) bool {
	PostgreSQL.Db.QueryRow("SELECT * FROM requester WHERE sessiontoken = $1", token)
	if token == "" {
		return false
	}
	return true
}
