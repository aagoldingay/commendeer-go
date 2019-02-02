package data_test

import (
	"testing"

	"github.com/aagoldingay/commendeer-go/data"
)

const (
	query = "update accesscode set code = null where code is not null"
)

// dbSetup and connection constants for this package are located in userdata_test.go

func Test_SendCodes(t *testing.T) {
	// setup
	db, err := dbSetup()
	if err != nil {
		t.Errorf("error connecting to database : %v\n", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Errorf("error connecting to database : %v\n", err)
	}

	// tests
	r := data.SendCodes(db)
	if r != 1 { // for success
		t.Errorf("update query status not 1, actual %v\n", r)
	}
	r = data.SendCodes(db) // should return 0 as codes already populated above
	if r != 0 {
		t.Errorf("second SendCodes test was not 0 as expected. actual %v\n", r)
	}

	// cleanup db
	_, err = db.Exec(query)
	if err != nil {
		t.Errorf("error on SendCodes cleanup - %v\n", err)
	}
}
