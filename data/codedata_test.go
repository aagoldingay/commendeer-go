package data_test

import (
	"testing"

	"github.com/aagoldingay/commendeer-go/data"
)

// ATTENTION: tests are called and run in data_test.go

func InitialSendCodes(t *testing.T) {
	r := data.SendCodes(db)
	if r != 1 { // for success
		t.Errorf("update query status not 1, actual %v\n", r)
	}
}

func SecondSendCodes(t *testing.T) {
	r := data.SendCodes(db) // should return 0 as codes already populated above
	if r != 0 {
		t.Errorf("second SendCodes test was not 0 as expected. actual %v\n", r)
	}
}
