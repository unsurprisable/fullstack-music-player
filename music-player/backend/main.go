package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

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

	r.Run(":8080")
}
