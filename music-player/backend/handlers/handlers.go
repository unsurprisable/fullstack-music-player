package handlers

import (
	"backend/database"
	"backend/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bogem/id3v2/v2"
	"github.com/gin-gonic/gin"
)

const uploadDir = "./uploads"

func UploadSong(c *gin.Context) {
	// get the file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file attached!"})
		return
	}

	// this seems very fragile? there must be a way to check filetype directly
	if !strings.HasSuffix(file.Filename, ".mp3") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only MP3 files are allowed!"})
		return
	}

	// save destination with inline string & filename
	dst := fmt.Sprintf("%s/%s", uploadDir, file.Filename)
	// check if uploads directory even exists, if not make it
	_, err = os.Stat(uploadDir)
	if os.IsNotExist(err) {
		err := os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			log.Fatal("Failed to create uploads directory: ", err)
		}
	}

	// attempt to save the file locally in /backend/uploads
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the file!"})
		return
	}

	// yay it worked
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("File '%s' uploaded successfully!", file.Filename),
	})

	title, artist, album := getMetadata(dst)
	if err := database.InsertSongMetadata(title, artist, album, file.Filename); err != nil {
		log.Fatal(err)
	}
}

// parse incoming mp3 and retrieve some metadata
func getMetadata(filePath string) (string, string, string) {
	tag, err := id3v2.Open(filePath, id3v2.Options{Parse: true})
	if err != nil {
		log.Fatal(err)
	}
	defer tag.Close()

	title := tag.Title()
	artist := tag.Artist()
	album := tag.Album()

	return title, artist, album
}

func GetAllSongs(c *gin.Context) {
	songs, err := database.GetAllSongs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "There was an error while retrieving stored songs!"})
		log.Println(err)
		return
	}

	formattedSongs := make([]gin.H, len(songs))

	for i, song := range songs {
		formattedSongs[i] = exportSong(&song)
	}

	c.JSON(http.StatusOK, formattedSongs)
}

func ResetStoredData(c *gin.Context) {
	// reset db table
	if err := database.ResetSongsTable(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "There was an error while clearing the database!"})
		log.Println(err)
		return
	}
	// delete /uploads directory
	if err := os.RemoveAll(uploadDir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "There was an error while deleting stored files!"})
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cleared stored song metadata."})
}

func GetSongById(c *gin.Context) {
	rawId := c.Param("id")

	id, err := strconv.Atoi(rawId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID."})
		return
	}

	song, err := database.GetSongById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Couldn't fetch from database!"})
		log.Println(err)
		return
	}
	if song == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "That song does not exist."})
		return
	}

	c.JSON(http.StatusOK, exportSong(song))
}

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

func ServeSongFile(c *gin.Context) {
	fileName := c.Param("filename")

	filePath := fmt.Sprintf("%s/%s", uploadDir, fileName)

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "File not found."})
		log.Println(err)
		return
	}

	c.Header("Content-Disposition", "inline")
	c.Header("Content-Type", "audio/mpeg")

	c.File(filePath)
}
