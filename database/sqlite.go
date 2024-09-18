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
	UsersCache = make(map[string]User)
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
	createTableUsers := `
    CREATE TABLE "users" (
	"userId"	INTEGER NOT NULL UNIQUE,
	"username"	VARCHAR(25) NOT NULL UNIQUE,
	"passwordHash"	TEXT NOT NULL,
	"firstName"	VARCHAR(25),
	"lastName"	VARCHAR(25),
	PRIMARY KEY("userId" AUTOINCREMENT)
)`

	_, err := database.Exec(createTableUsers)
	if err != nil {
		return err
	}

	return nil
}

func loadUsersToCache() error {

	rows, err := database.Query("SELECT username, passwordHash FROM users")
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

func VerifyUserInCache(username, passwordHash string) bool {
	user, exists := UsersCache[username]
	if !exists {
		return false
	}
	return passwordHash == user.PasswordHash
}

func InsertUserAndUpdateCache(username, passwordHash string) error {

	insertQuery := `INSERT INTO users (username, passwordHash) VALUES (?, ?)`
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
