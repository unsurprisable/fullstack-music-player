package database

import (
	"backend/models"
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const schemaPath = "../db/schema.sql"

var db *sql.DB

func InitDB() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSL_MODE")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", dbUser, dbPassword, dbName, dbSSLMode)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	// make sure the connection worked
	err = db.Ping()
	if err != nil {
		return err
	}

	// find schema.sql file
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		return err
	}

	// set up database tables from schema.sql
	_, err = db.Exec(string(schema))
	if err != nil {
		return err
	}

	return nil
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}

func InsertSongMetadata(title, artist, album, filename string) error {
	// insert metadata under correct columns
	_, err := db.Exec(`
		INSERT INTO songs (filename, title, artist, album)
		VALUES ($1, $2, $3, $4)
	`, filename, title, artist, album)

	return err
}

func GetAllSongs() ([]models.Song, error) {
	// retrieve list of song metadata stored in db
	rows, err := db.Query("SELECT id, filename, title, artist, album, uploaded_at FROM songs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []models.Song
	for rows.Next() {
		var song models.Song
		if err := rows.Scan(&song.ID, &song.Filename, &song.Title, &song.Artist, &song.Album, &song.UploadedAt); err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}

	return songs, nil
}

func ResetSongsTable() error {
	_, err := db.Exec("TRUNCATE TABLE playlists_songs, songs, playlists RESTART IDENTITY CASCADE")
	if err != nil {
		return err
	}

	return nil
}

func GetSongById(id int) (*models.Song, error) {
	var song models.Song

	row := db.QueryRow("SELECT id, filename, title, artist, album, uploaded_at FROM songs WHERE id = $1", id)

	err := row.Scan(&song.ID, &song.Filename, &song.Title, &song.Artist, &song.Album, &song.UploadedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // no song in db with that id
		}
		return nil, err
	}

	return &song, nil
}

func DeleteSongById(id int) error {
	// create transaction with db
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// delete from playlists_songs first since this table relies on the IDs from songs
	// things could break if an ID is deleted from songs but remains in playlists_songs
	_, err = tx.Exec("DELETE FROM playlists_songs WHERE song_id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM songs WHERE id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// finalize the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
