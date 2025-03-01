package handlers

import (
	"backend/database"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bogem/id3v2/v2"
	"github.com/gin-gonic/gin"
)

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

	const uploadDir = "./uploads"

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
	if err := database.InsertSongMetadata(title, artist, album, dst); err != nil {
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
