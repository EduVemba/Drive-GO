package models

import "time"

type Person struct {
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	SessionToken string    `json:"session_token"`
	CSRFToken    string    `json:"csrf_token"`
	DateBirth    time.Time `json:"date_birth"`
}

/*
func CreateRequester(firstname, lastname, email string, time time.Time) Person {
	concretePerson := Person{FirstName: firstname, LastName: lastname, Email: email, TimeStamp: time}
	return concretePerson
}

*/

func (p *Person) GetName() string {
	return p.FirstName + " " + p.LastName
}

/*
func (p *Person) GetAddress() string {
	return p.Address
}

*/

func (p *Person) GetEmail() string {
	return p.Email
}
