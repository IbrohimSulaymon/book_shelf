package server

type CreateNewUserRequest struct {
	Name   string `json:"name"`
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

type CreateNewBookRequest struct {
	ISBN string `json:"isbn"`
}

type EditBookRequest struct {
	Status int `json:"status"`
}
