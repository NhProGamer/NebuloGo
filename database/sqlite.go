package sqlite

import (
	"NebuloGo/utils"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

type User struct {
	Username     string
	PasswordHash string
}

var database *sql.DB

var UsersCache map[string]User

func InitSqliteDB() {
	databaseFile := "./database.db"

	exist, err := utils.DoesExistFile(databaseFile)
	if err != nil {
		log.Fatal(err)
	}
	if !exist {
		file, err := os.Create(databaseFile)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
	}

	db, err := sql.Open("sqlite3", databaseFile)
	database = db
	if err != nil {
		log.Fatal(err)
	}
	if !exist {
		err := generateDbStructure()
		if err != nil {
			log.Fatal(err)
		}
	}

	err = loadUsersToCache()
	if err != nil {
		log.Fatal(err)
	}

}

func generateDbStructure() error {
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        password_hash TEXT NOT NULL
    );`

	_, err := database.Exec(createTableQuery)
	if err != nil {
		return err
	}

	return nil
}

func loadUsersToCache() error {

	rows, err := database.Query("SELECT username, password_hash FROM users")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var username, passwordHash string
		if err := rows.Scan(&username, &passwordHash); err != nil {
			return err
		}
		UsersCache[username] = User{
			Username:     username,
			PasswordHash: passwordHash,
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

func verifyUserInCache(username, passwordHash string) bool {
	user, exists := UsersCache[username]
	if !exists {
		return false
	}

	return passwordHash == user.PasswordHash
}

func insertUserAndUpdateCache(username, passwordHash string) error {

	insertQuery := `INSERT INTO users (username, password_hash) VALUES (?, ?)`
	_, err := database.Exec(insertQuery, username, UsersCache)
	if err != nil {
		return err
	}

	// Mise à jour du cache
	UsersCache[username] = User{
		Username:     username,
		PasswordHash: passwordHash,
	}

	return nil
}
