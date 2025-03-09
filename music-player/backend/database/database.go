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

	// check if any errors occured during iteration
	if err := rows.Err(); err != nil {
		return nil, err
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

	if err := row.Scan(&song.ID, &song.Filename, &song.Title, &song.Artist, &song.Album, &song.UploadedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // no song in db with that id
		}
		return nil, err
	}

	return &song, nil
}

func GetPlaylistById(id int) (*models.Playlist, []int, error) {
	var playlist models.Playlist

	row := db.QueryRow("SELECT id, name, created_at FROM playlists WHERE id = $1", id)

	if err := row.Scan(&playlist.ID, &playlist.Name, &playlist.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil // no playlist with that id
		}
		return nil, nil, err
	}

	songIDs, err := getSongsForPlaylist(id)
	if err != nil {
		return nil, nil, err
	}

	return &playlist, songIDs, nil
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

func DeletePlaylistById(id int) error {
	// create transaction with db
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// delete from playlists_songs first since this table relies on the IDs from playlists
	// things could break if an ID is deleted from playlists but remains in playlists_songs
	_, err = tx.Exec("DELETE FROM playlists_songs WHERE playlist_id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM playlists WHERE id = $1", id)
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

func CreatePlaylist(name string) error {
	_, err := db.Exec(`
		INSERT INTO playlists (name)
		VALUES ($1)
	`, name)

	return err
}

func GetAllPlaylists() ([]models.Playlist, [][]int, error) {
	rows, err := db.Query("SELECT id, name, created_at FROM playlists")
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var playlists []models.Playlist
	var songIDs [][]int

	for rows.Next() {
		var playlist models.Playlist
		if err := rows.Scan(&playlist.ID, &playlist.Name, &playlist.CreatedAt); err != nil {
			return nil, nil, err
		}
		playlists = append(playlists, playlist)

		playlistSongIDs, err := getSongsForPlaylist(playlist.ID)
		if err != nil {
			return nil, nil, err
		}

		songIDs = append(songIDs, playlistSongIDs)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return playlists, songIDs, nil
}

func getSongsForPlaylist(playlistID int) ([]int, error) {
	rows, err := db.Query("SELECT song_id FROM playlists_songs WHERE playlist_id = $1", playlistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	songIDs := []int{}
	for rows.Next() {
		var songID int
		if err := rows.Scan(&songID); err != nil {
			return nil, err
		}
		songIDs = append(songIDs, songID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return songIDs, nil
}

func AddSongToPlaylist(playlistID, songID int) error {
	// database already prevents duplicates with its key system
	_, err := db.Exec(`
		INSERT INTO playlists_songs (playlist_id, song_id)
		VALUES ($1, $2) ON CONFLICT DO NOTHING
	`, playlistID, songID)

	return err
}

func DeleteSongFromPlaylist(playlistID, songID int) error {
	// database already prevents duplicates with its key system
	_, err := db.Exec(`
		DELETE FROM playlists_songs
		WHERE playlist_id = $1 AND song_id = $2
	`, playlistID, songID)

	return err
}
