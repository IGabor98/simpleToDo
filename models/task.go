package models

type Task struct {
	Id     uint64 `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	Status bool   `json:"status"`
}
