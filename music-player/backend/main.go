package main

import (
	"backend/database"
	"backend/handlers"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// initialize db
	if err := database.InitDB(); err != nil {
		log.Fatal("Database connection failed: ", err)
	}
	// properly close the db connection when the program ends
	defer database.CloseDB()

	r := gin.Default() // set up gin

	// add listener for GET requests on root directory saying that the backend is running
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Backend is running!"})
	})

	// add listener for POST request to upload songs
	r.POST("/upload", handlers.UploadSong)

	// view all stored songs
	r.GET("/songs", func(c *gin.Context) {
		songs, err := database.GetAllSongs()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "There was an error while retrieving stored songs!"})
			log.Fatal(err)
			return
		}
		c.JSON(http.StatusOK, songs)
	})

	// reset db
	r.GET("/songs/clear", func(c *gin.Context) {
		if err := database.ResetSongsTable(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "There was an error while clearing the database!"})
			log.Fatal(err)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cleared stored song metadata."})
	})

	r.Run(":8080")
}
