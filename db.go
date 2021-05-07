package main

import (
	"database/sql"
	"io"
	"log"
	"os"

	"github.com/gchaincl/dotsql"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

type DatabaseUtility struct {
	dot       *dotsql.DotSql
	dbContext *sql.DB
}

const dbName = "sqlite-database.db"
const queries = `
-- name: create-stream-key-table
CREATE TABLE IF NOT EXISTS stream_keys (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	stream_key TEXT NOT NULL UNIQUE,
	is_valid INTEGER NOT NULL
);

--name: create-users-table
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	is_admin INTEGER NOT NULL
);

-- name: get-streams
SELECT id, is_valid, stream_key FROM stream_keys;

-- name: create-new-stream
INSERT INTO stream_keys (stream_key, is_valid) VALUES (?, 0);

-- name: delete-stream
DELETE FROM stream_keys WHERE id = ?;

-- name: toggle-stream
UPDATE stream_keys SET is_valid = ? WHERE stream_key = ?;

--name: create-new-user
INSERT INTO users (username, password, is_admin) VALUES (?, ?, ?);

--name: validate-user
SELECT id, is_admin, password FROM users WHERE username = ? LIMIT 1;
`

func DbInitialize() *DatabaseUtility {
	db := new(DatabaseUtility)
	db.createDb()
	sqlContext, openError := sql.Open("sqlite3", dbName)
	if openError == nil {
		db.dbContext = sqlContext
	}
	dot, loadError := dotsql.LoadFromString(queries)
	if loadError == nil {
		db.dot = dot
	}
	db.seedDb()
	initialAdminPassword := os.Getenv("ADMIN_PASSWORD")
	if initialAdminPassword == "" {
		log.Println("Admin password not set, generating...")
		initialAdminPassword = generatePassword(32)
	}
	db.CreateNewUser("admin", initialAdminPassword, 1)
	file, fileError := os.Create("initial-admin-password")
	if fileError != nil {
		log.Fatalf("Cannot open file! %s", fileError)
	}
	defer file.Close()
	_, writeError := io.WriteString(file, initialAdminPassword)
	if writeError != nil {
		log.Fatalf("Cannot write to file! %s", fileError)
	}
	file.Sync()
	log.Printf("Inital login is username: admin, password: %s", initialAdminPassword)
	return db
}

func (db *DatabaseUtility) CloseDb() {
	db.dbContext.Close()
}

func (db *DatabaseUtility) createDb() {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		// path/to/whatever does not exist
		log.Printf("Creating %s...", dbName)
		file, err := os.Create(dbName) // Create SQLite file
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
	}
	log.Printf("%s created", dbName)
}

func (db *DatabaseUtility) seedDb() {
	_, streamKeyTableErr := db.dot.Exec(db.dbContext, "create-stream-key-table")
	if streamKeyTableErr != nil {
		panic(streamKeyTableErr)
	}
	_, err := db.dot.Exec(db.dbContext, "create-users-table")
	if err != nil {
		panic(err)
	}
}

func (db *DatabaseUtility) ToggleStream(status int, streamKey string) (sql.Result, error) {
	return db.dot.Exec(db.dbContext, "toggle-stream", status, streamKey)
}

func (db *DatabaseUtility) GetStreams() (*sql.Rows, error) {
	return db.dot.Query(db.dbContext, "get-streams")
}

func (db *DatabaseUtility) CreateNewStream(guid string) (sql.Result, error) {
	return db.dot.Exec(db.dbContext, "create-new-stream", guid)
}

func (db *DatabaseUtility) DeleteStream(id string) (sql.Result, error) {
	return db.dot.Exec(db.dbContext, "delete-stream", id)
}

func (db *DatabaseUtility) Login(username, password string) (*sql.Rows, error) {
	return db.dot.Query(db.dbContext, "validate-user", username, password)
}

func (db *DatabaseUtility) CreateNewUser(username, password string, isAdmin int) (sql.Result, error) {
	hashedPwd, _ := HashPassword(password)
	return db.dot.Exec(db.dbContext, "create-new-user", username, hashedPwd, isAdmin)
}
