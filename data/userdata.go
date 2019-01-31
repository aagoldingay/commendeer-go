package data

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"time"
)

// UserReturn is a returned object with information required for a user session
type UserReturn struct {
	ID       int
	Username string
	Admin    bool
}

type userData struct {
	id                   int
	username, pass, salt string
	admin                bool
}

const getUserQuery = "SELECT * FROM userinfo where Username = $1"

// GetUserInfo checks the database for an occurrence of a user by username, then compares hashed passwords
func GetUserInfo(u, p string, db *sql.DB) UserReturn {
	var (
		id                   int
		username, pass, salt string
		admin                bool
	)

	rows, err := db.Query(getUserQuery, u)
	if err != nil {
		fmt.Printf("%v: error on GetUserInfo query - %v\n", time.Now(), err)
	}
	defer rows.Close()

	for rows.Next() { // per row returned
		rows.Scan(&id, &username, &pass, &admin, &salt)
	}
	ud := userData{id, username, pass, salt, admin}

	if ud.id < 1 || !comparePasswords(ud.pass, p, ud.salt) { // no user returned || incorrect password
		return UserReturn{}
	}

	return UserReturn{ud.id, ud.username, ud.admin}
}

func comparePasswords(dbPass, userPass, salt string) bool {
	uHashPass := fmt.Sprintf("%x", md5.Sum([]byte(userPass+salt)))
	return uHashPass == dbPass
}
