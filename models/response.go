package models

type Response struct {
	Success bool               `json:"success"`
	Rates   map[string]float32 `json:"rates"`
}
