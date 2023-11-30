package entity

import "time"

type Note struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Date        time.Time `json:"date"`
	Status      string    `json:"status"`
}

const (
	StatusDone    = "done"
	StatusNotDone = "not_done"
)
