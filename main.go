package main

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gchaincl/dotsql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

const dbName = "sqlite-database.db"
const queries = `
-- name: create-stream-key-table
CREATE TABLE IF NOT EXISTS stream_keys (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	stream_key TEXT NOT NULL UNIQUE,
	is_valid INTEGER NOT NULL
);

-- name: get-streams
SELECT id, is_valid, stream_key FROM stream_keys;

-- name: create-new-stream
INSERT INTO stream_keys (stream_key, is_valid) VALUES (?, 0);

-- name: delete-stream
DELETE FROM stream_keys WHERE id = ?;

-- name: toggle-stream
UPDATE stream_keys (is_valid) VALUES (?) WHERE stream_key = ?;
`
// Stream obj
type Stream struct {
	Id, IsValid int
	StreamKey string
}

func main() {
	createDb()
	sqlContext, _ := sql.Open("sqlite3", dbName)
	defer sqlContext.Close()
	dot, _ := dotsql.LoadFromString(queries)
	seedDb(sqlContext, dot)
	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/verify-stream", func(c *gin.Context) {
		key := c.DefaultPostForm("name", "")
		if key == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"responseMessage": "Something went wrong getting the stream key!",
			})
			return
		}
		result, _ := dot.Exec(sqlContext, "toggle-stream", 1, key)
		rowsAffected, _ := result.RowsAffected()
		isValidKey := rowsAffected == 1
		switch {
		case isValidKey:
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"responseMessage": "Stream key is good!",
			})
		case !isValidKey:
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"responseMessage": "No stream key here!",
			})
		}
	})
	r.GET("/streams", func(c *gin.Context) {
		streams := make([]Stream, 0)
		rows, err := dot.Query(sqlContext, "get-streams")
		defer rows.Close()
		
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				c.JSON(http.StatusNotFound, gin.H{"success": false, "responseMessage": err.Error()})
			default:
				c.JSON(http.StatusBadRequest, gin.H{"success": false, "responseMessage": err.Error()})
			}
			return
		}
		for rows.Next() {
			var (
				id, isValid int
				streamKey string
			)
			err := rows.Scan(&id, &isValid, &streamKey)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "responseMessage": err.Error()})
				return
			}
			stream := Stream{id, isValid, streamKey}
			streams = append(streams, stream)
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "streams": streams})
	})
	r.POST("/create-key", func(c *gin.Context) {
		guid := generateGUID()
		_, err := dot.Exec(sqlContext, "create-new-stream", guid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"responseMessage": fmt.Sprintf("Failed to insert new record, %s", err.Error()),
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"responseMessage": "Insert new stream ok!",
		})
	})
	r.DELETE("/streams/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"responseMessage": "Something went wrong getting the stream id!",
			})
			return
		}
		_, err := dot.Exec(sqlContext, "delete-stream", id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"responseMessage": fmt.Sprintf("Failed to delete existing record, %s", err.Error()),
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"responseMessage": "Deleted stream!",
		})
	})
	r.POST("/stream-over", func(c *gin.Context) {
		key := c.DefaultPostForm("name", "")
		if key == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"responseMessage": "Something went wrong getting the stream key!",
			})
			return
		}
		result, _ := dot.Exec(sqlContext, "toggle-stream", 0, key)
		rowsAffected, _ := result.RowsAffected()
		isValidKey := rowsAffected == 1
		switch {
		case isValidKey:
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"responseMessage": "Stream key is good!",
			})
		case !isValidKey:
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"responseMessage": "No stream key here!",
			})
		}
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func createDb() {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		// path/to/whatever does not exist
		log.Println("Creating %s...", dbName)
		file, err := os.Create(dbName) // Create SQLite file
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
	}
	log.Println("%s created", dbName)
}

func seedDb(context *sql.DB, dot *dotsql.DotSql) {
	_, err := dot.Exec(context, "create-stream-key-table")
	if err != nil {
		panic(err)
	}
}

func generateGUID() (s string) {
	b := make([]byte, 10)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x",
		b[0:4], b[4:6], b[6:])
	return uuid
}