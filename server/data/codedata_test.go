package data_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/aagoldingay/commendeer-go/server/data"
)

// ATTENTION: tests are called and run in data_test.go

func InitialSendCodes(t *testing.T) {
	sent, err := data.SendCodes(1, db) // should return false as codes already populated above
	if err != nil {
		if err.Error() == "no codes to create" {
			t.Errorf("InitialSendCodes should have created codes\n")
		}
		if err.Error() == "problem encountered while creating codes" {
			t.Errorf("InitialSendCodes error during queries\n")
		}
	}
	if !sent { // for success
		t.Errorf("InitialSendCodes should have sent, didn't\n")
	}
}

func SecondSendCodes(t *testing.T) {
	sent, err := data.SendCodes(1, db) // should return false as codes already populated above
	if err != nil {
		if err.Error() == "problem encountered while creating codes" {
			t.Errorf("InitialSendCodes error during queries\n")
		}
	}
	if sent {
		t.Errorf("second SendCodes test sent. shouldnt have\n")
	}
}

func GetAccessCode_Success(t *testing.T) {
	qid, err := data.GetAccessCode("e@email.com", "helloworld", db)
	if err != nil {
		t.Errorf("Test_GetAccessCode_Success errored, shouldnt have : %v\n", err)
	}
	if qid < 1 {
		t.Errorf("Test_GetAccessCode_Success did not succeed\n")
	}
}

func getAccessCode_Success_Setup(db *sql.DB) {
	iq := "INSERT INTO AccessCode (Email, SystemUsername, Code, Used, QuestionnaireID) VALUES ('e@email.com', 'fake1', 'helloworld', FALSE, 1);"
	_, err := db.Exec(iq)
	if err != nil {
		fmt.Printf("Test_GetAccessCode_Success problem on setup\n")
	}
}

func getAccessCode_Success_TD(db *sql.DB) {
	dq := "delete from AccessCode where email = 'e@email.com'"
	_, err := db.Exec(dq)
	if err != nil {
		fmt.Printf("Test_GetAccessCode_Success problem on cleanup\n")
	}
}

func GetAccessCode_UsedCode(t *testing.T) {
	qid, err := data.GetAccessCode("e@email.com", "helloworld", db)
	if qid > 0 || err != nil {
		t.Errorf("Test_GetAccessCode_UsedCode did not fail\n")
	}
}

func getAccessCode_UsedCode_Setup(db *sql.DB) {
	iq := "INSERT INTO AccessCode (Email, SystemUsername, Code, Used, QuestionnaireID) VALUES ('e@email.com', 'fake1', 'helloworld', TRUE, 1);"
	_, err := db.Exec(iq)
	if err != nil {
		fmt.Printf("Test_GetAccessCode_UsedCode problem on setup\n")
	}
}

func getAccessCode_UsedCode_TD(db *sql.DB) {
	dq := "delete from AccessCode where email = 'e@email.com'"
	_, err := db.Exec(dq)
	if err != nil {
		fmt.Printf("Test_GetAccessCode_UsedCode problem on cleanup\n")
	}
}

func GetAccessCode_InvalidCode(t *testing.T) {
	qid, err := data.GetAccessCode("fakeemail@this.com", "invalidcode", db)
	if qid > 0 || err.Error() != "code not of desired length: 10" {
		t.Errorf("GetAccessCode_InvalidCode did not fail\n")
	}
}

func GetAccessCode_UnknownEmail(t *testing.T) {
	qid, err := data.GetAccessCode("fakeemail@this.com", "frjoghriug", db)
	if qid > 0 || err.Error() != "code or user not found" {
		t.Errorf("GetAccessCode_UnknownEmail did not fail\n")
	}
}
