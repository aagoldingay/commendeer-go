package data_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/aagoldingay/commendeer-go/server/data"
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

func GetAccessCode_Success(t *testing.T) {
	notUsed, err := data.GetAccessCode("e@email.com", "helloworld", db)
	if !notUsed || err != nil {
		t.Errorf("Test_GetAccessCode_Success did not succeed\n")
	}
}

func getAccessCode_Success_Setup(db *sql.DB) {
	iq := "INSERT INTO AccessCode (Email, SystemUsername, Code, Used) VALUES ('e@email.com', 'fake1', 'helloworld', FALSE);"
	_, err := db.Exec(iq)
	if err != nil {
		fmt.Printf("Test_GetAccessCode_Success problem on setup")
	}
}

func getAccessCode_Success_TD(db *sql.DB) {
	dq := "delete from AccessCode where codeid = 51"
	_, err := db.Exec(dq)
	if err != nil {
		fmt.Printf("Test_GetAccessCode_Success problem on cleanup")
	}
}

func GetAccessCode_UsedCode(t *testing.T) {
	notUsed, err := data.GetAccessCode("e@email.com", "helloworld", db)
	if notUsed || err != nil {
		t.Errorf("Test_GetAccessCode_UsedCode did not fail\n")
	}
}

func getAccessCode_UsedCode_Setup(db *sql.DB) {
	iq := "INSERT INTO AccessCode (Email, SystemUsername, Code, Used) VALUES ('e@email.com', 'fake1', 'helloworld', TRUE);"
	_, err := db.Exec(iq)
	if err != nil {
		fmt.Printf("Test_GetAccessCode_UsedCode problem on setup")
	}
}

func getAccessCode_UsedCode_TD(db *sql.DB) {
	dq := "delete from AccessCode where codeid = 52"
	_, err := db.Exec(dq)
	if err != nil {
		fmt.Printf("Test_GetAccessCode_UsedCode problem on cleanup")
	}
}

func GetAccessCode_InvalidCode(t *testing.T) {
	used, err := data.GetAccessCode("fakeemail@this.com", "invalidcode", db)
	if used || err.Error() != "code not of desired length: 10" {
		t.Errorf("GetAccessCode_InvalidCode did not fail\n")
	}
}

func GetAccessCode_UnknownEmail(t *testing.T) {
	used, err := data.GetAccessCode("fakeemail@this.com", "frjoghriug", db)
	if used || err.Error() != "code or user not found" {
		t.Errorf("GetAccessCode_UnknownEmail did not fail\n")
	}
}
