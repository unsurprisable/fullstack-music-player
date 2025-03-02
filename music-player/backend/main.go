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

	// view all stored songs
	r.GET("/songs", handlers.GetAllSongs)

	// reset db
	r.GET("/songs/clear", handlers.ResetStoredData)

	// get a single song by id
	r.GET("/songs/:id", handlers.GetSongById)

	r.GET("/songs/file/:filename", handlers.ServeSongFile)

	r.Run(":8080")
}
