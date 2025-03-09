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

type Playlist struct {
	ID        int
	Name      string
	CreatedAt time.Time
}

type CreatePlaylistRequest struct {
	Name string `json:"name"`
}
