package handlers

import (
	"backend/models"
	"fmt"

	"github.com/gin-gonic/gin"
)

func exportSong(song *models.Song) gin.H {
	return gin.H{
		"id":         song.ID,
		"title":      song.Title,
		"artist":     song.Artist,
		"album":      song.Album,
		"uploadedAt": song.UploadedAt,
		"fileURL":    fmt.Sprintf("http://localhost:8080/songs/file/%s", song.Filename),
	}
}

func exportPlaylist(playlist *models.Playlist, songIDs []int) gin.H {

	return gin.H{
		"id":        playlist.ID,
		"name":      playlist.Name,
		"createdAt": playlist.CreatedAt,
		"songIDs":   songIDs,
	}
}
