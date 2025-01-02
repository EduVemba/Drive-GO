package PostgreSQL

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var Db *sql.DB

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error opening the archive .env: %v", err)
	}
}

func Connect() {

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		log.Fatalf("The DB_PASSWORD environment variable is not set")
	}

	host := "localhost"
	port := 5432
	user := "postgres"
	dbname := "Drive-gO"

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	Db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error opening the connection with the Database: %v", err)
	}

	err = Db.Ping()
	if err != nil {
		log.Fatalf("Error to verify the connection with the Database: %v", err)
	}

	log.Println("Connection with PostgreSQL succeeded.")
}

func Close() {
	if Db != nil {
		err := Db.Close()
		if err != nil {
			log.Printf("Error closing the connection with  PostgreSQL: %v", err)
		} else {
			log.Println("Connection with PostgreSQL closed.")
		}
	}
}
