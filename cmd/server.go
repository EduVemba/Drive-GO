package main

import (
	"d_uber_golang/internal/Controllers"
	"d_uber_golang/internal/Database/MongoDB"
	"d_uber_golang/internal/routes"
	"log"
	"net/http"
)

func init() {
	MongoDB.Connect()
}

func main() {

	http.HandleFunc(routes.CREATE_REQUEST_INTENT, Controllers.CreateTravelIntent)
	http.HandleFunc(routes.ACCPET_REQUEST_INTENT, Controllers.GetTravelIntent)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
