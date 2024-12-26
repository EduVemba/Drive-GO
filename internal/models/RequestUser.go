package models

import "time"

type Person struct {
	firstName string    `json:"first_name"`
	lastName  string    `json:"last_name"`
	dateBirth time.Time `json:"date_birth"`
	timeStamp time.Time `json:"time_stamp"`
	address   string    `json:"address"`
	local     string    `json:"local"`
}

func CreateRequester(firstname, lastname, local string) Person {
	concretePerson := Person{firstName: firstname, lastName: lastname, timeStamp: time.Now()}
	return concretePerson
}

func (p *Person) getName() string {
	return p.firstName + " " + p.lastName
}

func (p *Person) getAddress() string {
	return p.address
}
