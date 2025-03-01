package models

import "time"

type Song struct {
	ID         int
	Filename   string
	Title      string
	Artist     string
	Album      string
	UploadedAt time.Time
}
