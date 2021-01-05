package main

// CreateResponse is the response to the api?create request
type CreateResponse struct {
	Success bool   `json:"success"`
	Err     string `json:"error"`
	ID      byte   `json:"id"`
}
