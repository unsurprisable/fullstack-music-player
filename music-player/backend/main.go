package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

func uploadSong(c *gin.Context) {
	// get the file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file attached!"})
		return
	}

	uploadDir := "./uploads"

	// save destination with inline string & filename
	dst := fmt.Sprintf("%s/%s", uploadDir, file.Filename)
	// check if uploads directory even exists, if not make it
	_, err = os.Stat(uploadDir)
	if os.IsNotExist(err) {
		err := os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			log.Fatal("Failed to create uplaods directory: ", err)
		}
	}

	// attempt to save the file locally in /backend/uploads
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the file!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("File %s uploaded successfully!", file.Filename),
	})
}

func main() {
	var err error // DANGER DANGER DANGER - USERNAME AND PASSWORD EXPOSED ON NEXT LINE:
	db, err = sql.Open("postgres", "user=postgres, password=peekaboosnakesql dbname=music sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // ensure the db connection is closed when the program ends

	r := gin.Default() // set up gin

	// add listener for GET requests on root directory saying that the backend is running
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Backend is running!"})
	})

	r.POST("/upload", uploadSong)

	r.Run(":8080")
}
