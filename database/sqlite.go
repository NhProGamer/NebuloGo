package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type User struct {
	Username     string
	PasswordHash string
}

var UsersCache map[string]User

func InitSqliteDB() {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}

	err = loadUsersToCache(db)
	if err != nil {
		log.Fatal(err)
	}

}

func loadUsersToCache(db *sql.DB) error {

	rows, err := db.Query("SELECT username, password_hash FROM users")
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

func insertUserAndUpdateCache(db *sql.DB, username, passwordHash string) error {

	insertQuery := `INSERT INTO users (username, password_hash) VALUES (?, ?)`
	_, err := db.Exec(insertQuery, username, UsersCache)
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
