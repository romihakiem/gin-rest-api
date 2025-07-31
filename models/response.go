package models

type Response struct {
	Status  int `json:"status"`
	Message any `json:"message"`
	Data    any `json:"data"`
}
