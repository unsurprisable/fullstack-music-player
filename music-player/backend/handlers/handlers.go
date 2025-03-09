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

func HandleUploadedSong(c *gin.Context) {
	// get the file
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files attached!"})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files attached!"})
		return
	}

	// keep track of uploads since i cant just return instantly now
	successfulUploads := 0
	failedUploads := 0
	results := []map[string]string{}

	for _, file := range files {
		// this seems very fragile? there must be a way to check filetype directly
		if !strings.HasSuffix(file.Filename, ".mp3") {
			results = append(results, map[string]string{
				"filename": file.Filename,
				"status":   "failed",
				"reason":   "Only MP3 files are allowed!",
			})
			failedUploads++
			continue
		}

		// save destination with inline string & filename
		dst := fmt.Sprintf("%s/%s", uploadDir, file.Filename)

		// check if uploads directory even exists, if not make it
		_, err = os.Stat(uploadDir)
		if os.IsNotExist(err) {
			err := os.MkdirAll(uploadDir, os.ModePerm)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store the files"})
				log.Fatal("Failed to create uploads directory: ", err)
			}
		}

		// attempt to save the file locally in /backend/uploads
		if err := c.SaveUploadedFile(file, dst); err != nil {
			log.Println(err)
			results = append(results, map[string]string{
				"filename": file.Filename,
				"status":   "failed",
				"reason":   "Failed to store the file",
			})
			failedUploads++
			continue
		}

		title, artist, album := getMetadata(dst)
		// so it has some form of identifier
		if title == "" {
			title = strings.TrimSuffix(file.Filename, ".mp3")
		}

		if err := database.InsertSongMetadata(title, artist, album, file.Filename); err != nil {
			results = append(results, map[string]string{
				"filename": file.Filename,
				"status":   "failed",
				"reason":   "Failed to save metadata to database",
			})
			failedUploads++
			continue
		}

		// yay this file uploaded successfully
		results = append(results, map[string]string{
			"filename": file.Filename,
			"status":   "success",
			"title":    title,
			"artist":   artist,
			"album":    album,
		})
		successfulUploads++
		continue
	}

	// send the results
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%d files uploaded successfully, %d failed", successfulUploads, failedUploads),
		"results": results,
	})
}

// parse incoming mp3 and retrieve some metadata
func getMetadata(filePath string) (string, string, string) {
	tag, err := id3v2.Open(filePath, id3v2.Options{Parse: true})
	if err != nil {
		log.Println(err)
		return "", "", ""
	}
	defer tag.Close()

	title := strings.TrimSpace(tag.Title())
	artist := strings.TrimSpace(tag.Artist())
	album := strings.TrimSpace(tag.Album())

	return title, artist, album
}

func DeleteSongByID(c *gin.Context) {
	rawId := c.Param("id")

	id, err := strconv.Atoi(rawId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID."})
		return
	}

	if err = database.DeleteSongById(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete song."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song deleted!"})
}

func GetAllSongs(c *gin.Context) {
	songs, err := database.GetAllSongs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve stored songs."})
		log.Println(err)
		return
	}

	formattedSongs := make([]gin.H, len(songs))

	for i, song := range songs {
		formattedSongs[i] = exportSong(&song)
	}

	c.JSON(http.StatusOK, formattedSongs)
}

func GetAllPlaylists(c *gin.Context) {
	playlists, songIDs, err := database.GetAllPlaylists()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve stored playlists."})
		log.Println()
		return
	}

	formattedPlaylists := make([]gin.H, len(playlists))

	for i, playlist := range playlists {
		formattedPlaylists[i] = exportPlaylist(&playlist, songIDs[i])
	}

	c.JSON(http.StatusOK, formattedPlaylists)
}

func ResetStoredData(c *gin.Context) {
	// reset db table
	if err := database.ResetSongsTable(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "There was an error while clearing the database!"})
		log.Println(err)
		return
	}
	// delete /uploads directory
	if err := os.RemoveAll(uploadDir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "There was an error while deleting stored files!"})
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cleared stored songs and playlists."})
}

func GetSongByID(c *gin.Context) {
	rawId := c.Param("id")

	id, err := strconv.Atoi(rawId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID."})
		return
	}

	song, err := database.GetSongById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't fetch from database!"})
		log.Println(err)
		return
	}
	if song == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "That song does not exist."})
		return
	}

	c.JSON(http.StatusOK, exportSong(song))
}

func GetPlaylistByID(c *gin.Context) {
	rawId := c.Param("id")

	id, err := strconv.Atoi(rawId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID."})
		return
	}

	playlist, songIDs, err := database.GetPlaylistById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't fetch from database!"})
		log.Println(err)
		return
	}
	if playlist == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "That playlist does not exist."})
		return
	}

	c.JSON(http.StatusOK, exportPlaylist(playlist, songIDs))
}

func ServeSongFile(c *gin.Context) {
	fileName := c.Param("filename")

	filePath := fmt.Sprintf("%s/%s", uploadDir, fileName)

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File not found."})
		log.Println(err)
		return
	}

	c.Header("Content-Disposition", "inline")
	c.Header("Content-Type", "audio/mpeg")

	c.File(filePath)
}

func CreatePlaylist(c *gin.Context) {
	var request models.CreatePlaylistRequest

	// map the JSON request into the CreatePlaylistRequest struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format!"})
		return
	}

	if request.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Playlist name cannot be empty!"})
		return
	}

	if err := database.CreatePlaylist(request.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create playlist."})
		log.Println(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Playlist successfully uploaded!"})
}

func DeletePlaylistByID(c *gin.Context) {
	rawId := c.Param("id")

	id, err := strconv.Atoi(rawId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID."})
		return
	}

	if err = database.DeletePlaylistById(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete song."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Playlist deleted!"})
}

func AddSongToPlaylist(c *gin.Context) {
	rawPlaylistId := c.Param("id")
	rawSongId := c.Param("song_id")

	playlistId, err := strconv.Atoi(rawPlaylistId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID!"})
		return
	}
	songId, err := strconv.Atoi(rawSongId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID!"})
		return
	}

	if err := database.AddSongToPlaylist(playlistId, songId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to add song to the playlist."})
		log.Println(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Song successfully added to playlist!"})
}

func DeleteSongFromPlaylist(c *gin.Context) {
	rawPlaylistId := c.Param("id")
	rawSongId := c.Param("song_id")

	playlistId, err := strconv.Atoi(rawPlaylistId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID!"})
		return
	}
	songId, err := strconv.Atoi(rawSongId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID!"})
		return
	}

	if err := database.DeleteSongFromPlaylist(playlistId, songId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove song to the playlist."})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song successfully removed from playlist!"})
}
