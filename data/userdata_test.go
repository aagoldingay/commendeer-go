package data_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/aagoldingay/commendeer-go/data"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "commendeer"
)

func dbSetup() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	d, err := sql.Open("postgres", psqlInfo) // configures db info for connection via code
	if err != nil {
		return nil, err // if there's an error, we don't want to continue listening to the port!
	}

	err = d.Ping() // open connection to the database
	if err != nil {
		return nil, err
	}
	return d, nil
}

func Test_GetUserInfo_Success(t *testing.T) {
	db, err := dbSetup()
	if err != nil {
		t.Errorf("error connecting to database : %v\n", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Errorf("error connecting to database : %v\n", err)
	}

	u, p := "admin1", "4dm1n123"
	d := data.GetUserInfo(u, p, db)
	if d.ID != 1 {
		t.Errorf("user %v not correct : expected (%v,%v), actual (%v,%v)\n", u, u, 1, d.Username, d.ID)
	}
}

func Test_GetUserInfo_UnknownUser(t *testing.T) {
	db, err := dbSetup()
	if err != nil {
		t.Errorf("error connecting to database : %v\n", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Errorf("error connecting to database : %v\n", err)
	}

	u, p := "admin5", "4dm1n123"
	d := data.GetUserInfo(u, p, db)
	if d.ID != 0 {
		t.Errorf("user returned successfully : expected (%v,%v), actual (%v,%v)\n", u, 0, d.Username, d.ID)
	}
}

func Test_GetUserInfo_IncorrectPassword(t *testing.T) {
	db, err := dbSetup()
	if err != nil {
		t.Errorf("error connecting to database : %v\n", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Errorf("error connecting to database : %v\n", err)
	}

	u, p := "admin1", "randompassword"
	d := data.GetUserInfo(u, p, db)
	if d.ID != 0 {
		t.Errorf("user returned successfully : expected (%v,%v), actual (%v,%v)\n", u, 0, d.Username, d.ID)
	}
}
