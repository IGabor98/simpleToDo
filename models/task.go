package models

import "time"

type Task struct {
	ID        string    `json:"id"`
	Title     string    `json:"title" valid:"email"`
	Body      string    `json:"body"`
	Status    bool      `json:"status"`
	Timestamp time.Time `json:"time" bson:"time"`
}
