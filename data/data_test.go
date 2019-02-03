package data_test

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "commendeer"

	codedataCleanupQuery = "update accesscode set code = null where code is not null"
)

var db *sql.DB

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

func Test_DataPackage(t *testing.T) {
	// setup tests
	d, err := dbSetup()
	if err != nil {
		t.Errorf("error connecting to database : %v\n", err)
	}
	db = d
	defer d.Close()
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Errorf("error connecting to database : %v\n", err)
	}

	// userdata_test.go
	t.Run("Test_GetUserInfo_Success", func(t *testing.T) { // run before code altering UserInfo table
		GetUserInfo_Success(t)
	})
	t.Run("Test_GetUserInfo_UnknownUser", func(t *testing.T) { // run before code altering UserInfo table
		GetUserInfo_UnknownUser(t)
	})
	t.Run("Test_GetUserInfo_IncorrectPassword", func(t *testing.T) { // run before code altering UserInfo table
		GetUserInfo_IncorrectPassword(t)
	})

	// codedata_test.go
	t.Run("Test_InitialSendCodes", func(t *testing.T) { // run before get from AccessCode table
		InitialSendCodes(t)
	})
	t.Run("Test_SecondSendCodes", func(t *testing.T) {
		SecondSendCodes(t)
	})

	// teardown db (codedata_test.go)
	_, err = db.Exec(codedataCleanupQuery)
	if err != nil {
		t.Errorf("error on SendCodes cleanup - %v\n", err)
	}
}
