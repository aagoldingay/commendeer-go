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

	authcodeCleanupQuery      = "delete from authcodes where code is not null"
	codedataCleanupQuery      = "update accesscode set code = null where code is not null"
	questionnaireCleanupQuery = "delete from questionnaire where questionnaireid != 1"
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
	t.Run("Test_Login_Success", func(t *testing.T) { // run before code altering UserInfo table
		Login_Success(t)
	})
	t.Run("Test_Login_UnknownUser", func(t *testing.T) { // run before code altering UserInfo table
		Login_UnknownUser(t)
	})
	t.Run("Test_Login_IncorrectPassword", func(t *testing.T) { // run before code altering UserInfo table
		Login_IncorrectPassword(t)
	})
	t.Run("Test_Logout_NoError", func(t *testing.T) {
		Logout_NoError(t)
	})
	t.Run("Test_Logout_IncorrectCode", func(t *testing.T) {
		Logout_IncorrectCode(t)
	})
	t.Run("Test_Logout_InvalidCode", func(t *testing.T) {
		Logout_InvalidCode(t)
	})
	// teardown db (userdata_test.go)
	_, err = db.Exec(authcodeCleanupQuery)
	if err != nil {
		t.Errorf("error on AuthCode cleanup - %v\n", err)
	}

	// questioncreator_test.go
	t.Run("Test_CreateForm", func(t *testing.T) {
		CreateForm(t)
	})

	// codedata_test.go
	t.Run("Test_InitialSendCodes", func(t *testing.T) { // run before get from AccessCode table
		InitialSendCodes(t)
	})
	t.Run("Test_SecondSendCodes", func(t *testing.T) {
		SecondSendCodes(t)
	})
	t.Run("Test_GetAccessCode_InvalidCode", func(t *testing.T) {
		GetAccessCode_InvalidCode(t)
	})
	t.Run("Test_GetAccessCode_UnknownEmail", func(t *testing.T) {
		GetAccessCode_UnknownEmail(t)
	})
	t.Run("Test_GetAccessCode_UsedCode", func(t *testing.T) {
		getAccessCode_UsedCode_Setup(db)
		GetAccessCode_UsedCode(t)
		getAccessCode_UsedCode_TD(db)
	})
	t.Run("Test_GetAccessCode_Success", func(t *testing.T) {
		getAccessCode_Success_Setup(db)
		GetAccessCode_Success(t)
		getAccessCode_Success_TD(db)
	})

	// teardown db (codedata_test.go)
	_, err = db.Exec(codedataCleanupQuery)
	if err != nil {
		t.Errorf("error on SendCodes cleanup - %v\n", err)
	}

	// teardown db (questioncreator_test.go)
	_, err = db.Exec(questionnaireCleanupQuery)
	if err != nil {
		t.Errorf("error on QuestionCreator cleanup - %v\n", err)
	}
}
