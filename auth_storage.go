package main

import (
	"crypto/sha256"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

func initDb() {
	db, err := sql.Open("sqlite3", "./auth_user.db")
	defer func(db *sql.DB) {
		checkErr(db.Close())
	}(db)

	if !checkErr(err) {
		panic(err)
	}

	_, _ = db.Exec("CREATE TABLE `userinfo` (" +
		"`uid` INTEGER PRIMARY KEY AUTOINCREMENT," +
		"`username` VARCHAR(64) NULL UNIQUE," +
		"`password` VARCHAR(64) NULL," +
		"`project` VARCHAR(64) NULL UNIQUE," +
		"`created` DATE NULL" +
		");")

	defer addUser("admin", sha256.Sum256([]byte(cfg.Admin)), "admin")

}

func addUser(username string, password [32]byte, project string) bool {
	db, err := sql.Open("sqlite3", "./auth_user.db")
	defer func(db *sql.DB) {
		checkErr(db.Close())
	}(db)
	if !checkErr(err) {
		return false
	}

	// insert
	stmt, err := db.Prepare("INSERT INTO userinfo(username, password, project, created) values(?,?,?,?)")
	if !checkErr(err) {
		return false
	}
	created := time.Now()

	_, err = stmt.Exec(username, string(password[:]), project, created)
	return checkErr(err)
}

func getUser(username string) (User, bool) {

	db, err := sql.Open("sqlite3", "./auth_user.db")
	defer func(db *sql.DB) {
		checkErr(db.Close())
	}(db)
	checkErr(err)

	row := db.QueryRow("SELECT * FROM userinfo where username=?", username)

	user := User{}
	if row == nil {
		return user, false
	}

	err = row.Scan(&user.id, &user.user, &user.pass, &user.project, &user.created)
	if !checkErr(err) {
		return User{}, false
	}

	return user, true
}
