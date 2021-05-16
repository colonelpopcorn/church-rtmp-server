package main

import (
	"database/sql"
	"errors"
	"io"
	"log"
	"os"

	"github.com/gchaincl/dotsql"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

type IDatabaseUtility interface {
	CloseDb()
	CreateStream(guid string) (sql.Result, error)
	CreateUser(username, password string, isAdmin int) (sql.Result, error)
	DeleteStream(id string) (sql.Result, error)
	GetStreams() (*sql.Rows, error)
	GetUser(userId int64) (User, error)
	GetUsers() ([]User, error)
	Login(username, password string) (*sql.Rows, error)
	ToggleStream(status int, streamKey string) (sql.Result, error)
}

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

--name: get-user-by-id
SELECT id, is_admin, username FROM users WHERE id = ? LIMIT 1;

--name: get-users
SELECT id, is_admin, username FROM users ORDER BY is_admin DESC, id;
`

func DbInitialize() IDatabaseUtility {
	db := DatabaseUtility{}
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
	db.CreateUser("admin", initialAdminPassword, 1)
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

func (db DatabaseUtility) CloseDb() {
	db.dbContext.Close()
}

func (db DatabaseUtility) createDb() {
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

func (db DatabaseUtility) seedDb() {
	_, streamKeyTableErr := db.dot.Exec(db.dbContext, "create-stream-key-table")
	if streamKeyTableErr != nil {
		panic(streamKeyTableErr)
	}
	_, err := db.dot.Exec(db.dbContext, "create-users-table")
	if err != nil {
		panic(err)
	}
}

func (db DatabaseUtility) ToggleStream(status int, streamKey string) (sql.Result, error) {
	return db.dot.Exec(db.dbContext, "toggle-stream", status, streamKey)
}

func (db DatabaseUtility) GetStreams() (*sql.Rows, error) {
	return db.dot.Query(db.dbContext, "get-streams")
}

func (db DatabaseUtility) CreateStream(guid string) (sql.Result, error) {
	return db.dot.Exec(db.dbContext, "create-new-stream", guid)
}

func (db DatabaseUtility) DeleteStream(id string) (sql.Result, error) {
	return db.dot.Exec(db.dbContext, "delete-stream", id)
}

func (db DatabaseUtility) Login(username, password string) (*sql.Rows, error) {
	return db.dot.Query(db.dbContext, "validate-user", username, password)
}

func (db DatabaseUtility) CreateUser(username, password string, isAdmin int) (sql.Result, error) {
	hashedPwd, _ := HashPassword(password)
	return db.dot.Exec(db.dbContext, "create-new-user", username, hashedPwd, isAdmin)
}

func (db DatabaseUtility) GetUser(userId int64) (User, error) {
	var user User
	rows, err := db.dot.Query(db.dbContext, "get-user-by-id")
	if err != nil {
		return user, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&user.UserId, &user.IsAdmin, &user.Username)
		if err != nil {
			return user, nil
		}
	}
	return user, errors.New("failed to find user for this ID")
}

func (db DatabaseUtility) GetUsers() ([]User, error) {
	users := []User{}
	rows, err := db.dot.Query(db.dbContext, "get-users")
	if err != nil {
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		err := rows.Scan(&user.UserId, &user.IsAdmin, &user.Username)
		if err != nil {
			users = append(users, user)
		}
	}
	return users, nil
}
