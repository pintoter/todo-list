package entity

import "time"

const (
	StatusDone    = "done"
	StatusNotDone = "not_done"
)

type Note struct {
	ID          int       `json:"id,omitempty"`
	UserId      int       `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Date        time.Time `json:"date"`
	Status      string    `json:"status"`
}
