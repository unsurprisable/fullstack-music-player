package main

import (
	"backend/database"
	"backend/handlers"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// initialize db
	if err := database.InitDB(); err != nil {
		log.Fatal("Database connection failed: ", err)
	}
	log.Println("Database connection established")
	// properly close the db connection when the program ends
	defer database.CloseDB()

	r := gin.Default() // set up gin

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	// add listener for GET requests on root directory saying that the backend is running
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Backend is running!"})
	})

	// add listener for POST request to upload songs
	r.POST("/upload", handlers.UploadSong)

	// get all songs or by id/filename
	r.GET("/songs", handlers.GetAllSongs)
	r.GET("/songs/:id", handlers.GetSongByID)
	r.GET("/songs/file/:filename", handlers.ServeSongFile)

	// delete all songs or by id
	r.GET("/songs/clear", handlers.ResetStoredData)
	r.DELETE("/songs/:id", handlers.DeleteSongByID)

	// TODO
	r.GET("/playlists")
	r.GET("/playlists/:id")

	// TODO
	r.POST("/playlists")
	r.POST("/playlists/:id/songs/:song_id")

	// TODO
	r.DELETE("/playlists/:id")
	r.DELETE("/playlists/:id/songs/:song_id")

	r.Run(":8080")
}
