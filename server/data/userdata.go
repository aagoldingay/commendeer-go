package data

import (
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"time"

	utils "github.com/aagoldingay/commendeer-go/server/utilities"
)

// UserReturn is a returned object with information required for a user session
type UserReturn struct {
	Username, Code string
}

type userData struct {
	id                   int
	username, pass, salt string
	admin                bool
}

const (
	getUserQuery       = "SELECT * FROM userinfo where Username = $1"
	addAuthRowQuery    = "INSERT INTO authcodes (UserID, Code, Administrator) VALUES ($1, $2, $3);"
	getAuthRowQuery    = "SELECT userid FROM authcodes WHERE code = $1"
	deleteAuthRowQuery = "DELETE FROM authcodes WHERE userid = $1"
)

// Login checks the database for an occurrence of a user by username, then compares hashed passwords
// Authorises if valid
func Login(u, p string, db *sql.DB) (UserReturn, error) {
	var (
		id                   int
		username, pass, salt string
		admin                bool
	)

	rows, err := db.Query(getUserQuery, u)
	if err != nil {
		fmt.Printf("%v: error on GetUserInfo query - %v\n", time.Now(), err)
		return UserReturn{}, errors.New("error on authorisation")
	}
	defer rows.Close()

	for rows.Next() { // per row returned
		rows.Scan(&id, &username, &pass, &admin, &salt)
	}
	ud := userData{id, username, pass, salt, admin}

	if ud.id < 1 || !comparePasswords(ud.pass, p, ud.salt) { // no user returned || incorrect password
		return UserReturn{}, errors.New("incorrect username or password")
	}

	code := authorise(ud.id, ud.admin, db)
	if code == "" {
		return UserReturn{}, errors.New("error on authorisation")
	}
	return UserReturn{ud.username, code}, nil
}

// Logout removes rows from database related to a user based on a valid code being given
func Logout(code string, db *sql.DB) error {
	// search db for code
	if len(code) != 20 {
		return errors.New("invalid code")
	}
	rows, err := db.Query(getAuthRowQuery, code)
	if err != nil {
		fmt.Printf("%v: error on getAuthRow query - %v\n", time.Now(), err)
		return errors.New("error on logout")
	}

	var (
		userid int
	)
	for rows.Next() {
		rows.Scan(&userid)
	}
	if userid < 1 {
		return errors.New("unknown code")
	}

	// remove any codes for that user
	_, err = db.Exec(deleteAuthRowQuery, userid)
	if err != nil {
		fmt.Printf("%v: error on deleteAuthRows query - %v\n", time.Now(), err)
		return errors.New("error on logout")
	}
	return nil
}

func authorise(id int, admin bool, db *sql.DB) string {
	code := utils.GenerateCode(20)

	_, err := db.Exec(addAuthRowQuery, id, code, admin)
	if err != nil {
		fmt.Printf("%v: error on Authorisation query - %v\n", time.Now(), err)
		return ""
	}

	return code
}

func comparePasswords(dbPass, userPass, salt string) bool {
	uHashPass := fmt.Sprintf("%x", md5.Sum([]byte(userPass+salt)))
	return uHashPass == dbPass
}
