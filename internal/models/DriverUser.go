package models

type DriverUser struct {
	Person       `json:"request_user_model"`
	registration string `json:"registration"`
}

func AcceptRequester(conreteRequester Person) {}
