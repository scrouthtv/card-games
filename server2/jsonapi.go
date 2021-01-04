package main

type CreateResponse struct {
	Success bool   `json:"success"`
	Err     string `json:"error"`
	Id      byte   `json:"id"`
}
