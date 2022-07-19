package models

// User schema of the user table
type ToDo struct {
	ID     int64  `json:"id"`
	Task   string `json:"task"`
	Status bool   `json:"status"`
}
