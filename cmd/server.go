package main

import (
	"d_uber_golang/internal/Authentication"
	"d_uber_golang/internal/Controllers"
	"d_uber_golang/internal/Database/MongoDB"
	"d_uber_golang/internal/Database/PostgreSQL"
	"d_uber_golang/internal/routes"
	"log"
	"net/http"
)

func main() {
	MongoDB.Connect()
	defer MongoDB.Close()

	PostgreSQL.Connect()
	defer PostgreSQL.Close()

	http.HandleFunc(routes.DRIVERREGISTER, Authentication.RegisterDriver)
	http.HandleFunc(routes.LOGINDRIVER, Authentication.LoginDriver)
	http.HandleFunc(routes.REGISTER, Authentication.Register)
	http.HandleFunc(routes.LOGIN, Authentication.Login)
	http.HandleFunc(routes.CREATE_REQUEST_INTENT, Authentication.Protected(Controllers.CreateTravelIntent))
	http.HandleFunc(routes.ACCPET_REQUEST_INTENT, Authentication.Protected(Controllers.GetTravelIntent))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
