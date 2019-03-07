package main

import (
	"context"
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	pb "github.com/aagoldingay/commendeer-go/pb"
	utils "github.com/aagoldingay/commendeer-go/server/utilities"
)

func Test_CreateAccessCodes_InvalidQuestionError(t *testing.T) {
	req := &pb.CreateCodeRequest{QuestionnaireID: 0}

	s := server{}
	resp, err := s.CreateAccessCodes(context.Background(), req)

	if err != nil {
		t.Errorf("problem with CreateAccessCodes server method - shouldn't return an error")
	}
	if resp.Error != pb.Error_BADREQUEST {
		t.Errorf("expected BAD REQUEST - error : %v\n", resp.Error.String())
	}
	if resp.ErrorDetails != "Invalid questionnaire" {
		t.Errorf("expected ErrorDetails - %v\n", resp.ErrorDetails)
	}
}

func Test_CreateAccessCodes_InvalidCodeError(t *testing.T) {
	s := server{}

	req := &pb.CreateCodeRequest{Authcode: "helloworld", QuestionnaireID: 1}
	resp, err := s.CreateAccessCodes(context.Background(), req)
	if err != nil {
		t.Errorf("problem with CreateAccessCodes server method - shouldn't return an error")
	}
	if resp.Error != pb.Error_BADREQUEST {
		t.Errorf("expected BAD REQUEST - error : %v\n", resp.Error.String())
	}
	if resp.ErrorDetails != "Invalid code" {
		t.Errorf("expected Errordetails - %v\n", resp.ErrorDetails)
	}
}

func Test_CreateAccessCodes_NoAuthError(t *testing.T) {
	// mock := setupSQLMock(t)
	d, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer d.Close()
	db = d

	utils.Setup(0)
	s := server{}

	req := &pb.CreateCodeRequest{QuestionnaireID: 1, Authcode: "helloworldhelloworld"}
	mr := sqlmock.NewRows([]string{"userid", "administrator"}).
		AddRow("1", false)
	mock.ExpectQuery("SELECT (.+) FROM authcodes WHERE code = (.+)").
		WithArgs(req.Authcode).WillReturnRows(mr)

	resp, err := s.CreateAccessCodes(context.Background(), req)
	if err != nil {
		t.Errorf("problem with CreateAccessCodes server method - shouldn't return an error")
	}
	if resp.Error != pb.Error_FORBIDDEN {
		t.Errorf("expected FORBIDDEN - error : %v\n", resp.Error.String())
	}
	if resp.ErrorDetails != "No authorisation" {
		t.Errorf("expected error details : %v\n", resp.ErrorDetails)
	}
}

func Test_CreateAccessCodes_NotSentError(t *testing.T) {
	d, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer d.Close()
	db = d

	utils.Setup(0)
	s := server{}

	req := &pb.CreateCodeRequest{QuestionnaireID: 1, Authcode: "helloworldhelloworld"}
	mr := sqlmock.NewRows([]string{"userid", "administrator"}).
		AddRow("1", true)
	mock.ExpectQuery("SELECT (.+) FROM authcodes WHERE code = (.+)").
		WithArgs(req.Authcode).WillReturnRows(mr)

	mock.ExpectQuery("SELECT (.+) FROM AccessCode WHERE Code IS NULL AND QuestionnaireID = (.+);").
		WithArgs(1).WillReturnError(sql.ErrNoRows)

	resp, err := s.CreateAccessCodes(context.Background(), req)
	if err != nil {
		t.Errorf("problem with CreateAccessCodes server method - shouldn't return an error")
	}
	if resp.Error != pb.Error_NIL {
		t.Errorf("Expected NIL error - %v\n", resp.Error.String())
	}
	if resp.ErrorDetails != "Codes did not send - confirm with administrator" {
		t.Errorf("expected error details - %v\n", resp.ErrorDetails)
	}
}

func Test_CreateAccessCodes_NoError(t *testing.T) {
	d, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer d.Close()
	db = d

	utils.Setup(0)
	s := server{}

	req := &pb.CreateCodeRequest{QuestionnaireID: 1, Authcode: "helloworldhelloworld"}
	mr := sqlmock.NewRows([]string{"userid", "administrator"}).
		AddRow("1", true)
	mock.ExpectQuery("SELECT (.+) FROM authcodes WHERE code = (.+)").
		WithArgs(req.Authcode).WillReturnRows(mr)

	mr2 := sqlmock.NewRows([]string{"codeid"}).AddRow("1")
	mock.ExpectQuery("SELECT (.+) FROM AccessCode WHERE Code IS NULL AND QuestionnaireID = (.+);").
		WithArgs(1).WillReturnRows(mr2)
	mock.ExpectExec("UPDATE AccessCode SET Code = (.+)").WillReturnResult(sqlmock.NewResult(1, 1))

	resp, err := s.CreateAccessCodes(context.Background(), req)
	if err != nil {
		t.Errorf("problem with CreateAccessCodes server method - shouldn't return an error")
	}
	if resp.Error != pb.Error_OK {
		t.Errorf("Expected OK error - %v\n", resp.Error.String())
	}
	if resp.ErrorDetails != "" {
		t.Errorf("expected no error details - %v\n", resp.ErrorDetails)
	}
}

// func Test_CreateForm_Errors(t *testing.T) {

// }

// func Test_GetFeedbackForm_Errors(t *testing.T) {

// }

func Test_LoginUser_NoErrors(t *testing.T) {
	// mock := setupSQLMock(t)
	d, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer d.Close()
	db = d

	utils.Setup(0)
	s := server{}

	u, p := "admin1", "4dm1n123"
	mr := sqlmock.NewRows([]string{"userid", "username", "pass", "administrator", "salt"}).
		AddRow("1", "admin1", "57d8da63dbcfd720673fd0622ac91549", "true", "zRvjFZ8Amq")
	mock.ExpectQuery("^SELECT (.+) FROM userinfo where Username = (.+)").
		WithArgs(u).WillReturnRows(mr)
	mock.ExpectExec(`^INSERT INTO authcodes \(UserID, Code, Administrator\) VALUES \(.+\)`).
		WithArgs(1, "mUNERA9rI2cvTK4UHomc", true).WillReturnResult(sqlmock.NewResult(1, 1))

	req := &pb.LoginRequest{Username: u, Password: p}
	resp, err := s.LoginUser(context.Background(), req)
	if err != nil {
		t.Errorf("problem with LoginUser server method - shouldn't return an error")
	}
	if resp.Error != pb.Error_OK {
		t.Errorf("expected OK - error : %v\n", resp.Error.String())
	}
	if resp.ErrorDetails != "" {
		t.Errorf("expected no error details : %v\n", resp.ErrorDetails)
	}
	if resp.Username != u {
		t.Errorf("incorrect user returned : %v\n", resp.Username)
	}
}

func Test_LoginUser_InformationError(t *testing.T) {
	d, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer d.Close()
	db = d

	utils.Setup(0)
	s := server{}

	u, p := "user", "falsepass"
	mr := sqlmock.NewRows([]string{}).AddRow()
	mock.ExpectQuery("^SELECT (.+) FROM userinfo where Username = (.+)").
		WithArgs(u).WillReturnRows(mr)

	req := &pb.LoginRequest{Username: u, Password: p}
	resp, err := s.LoginUser(context.Background(), req)

	if err != nil {
		t.Errorf("problem with LoginUser server method - shouldn't return an error")
	}
	if resp.Error != pb.Error_BADREQUEST {
		t.Errorf("expected BAD REQUEST - error : %v\n", resp.Error.String())
	}
	if resp.ErrorDetails != "Incorrect username or password" {
		t.Errorf("expected error details : %v\n", resp.ErrorDetails)
	}
}

func Test_LoginUser_AuthorisationError(t *testing.T) {
	mock := setupSQLMock(t)

	utils.Setup(0)
	s := server{}

	u, p := "admin1", "4dm1n123"
	mock.ExpectQuery("^SELECT (.+) FROM userinfo where Username = (.+)").
		WithArgs(u).WillReturnError(sql.ErrNoRows)

	req := &pb.LoginRequest{Username: u, Password: p}
	resp, err := s.LoginUser(context.Background(), req)

	if err != nil {
		t.Errorf("problem with LoginUser server method - shouldn't return an error")
	}
	if resp.Error != pb.Error_INTERNALERROR {
		t.Errorf("expected INTERNAL ERROR - error : %v\n", resp.Error.String())
	}
	if resp.ErrorDetails != "Problem logging in" {
		t.Errorf("expected error details : %v\n", resp.ErrorDetails)
	}
}

// func Test_LogoutUser_Errors(t *testing.T) {

// }

// func Test_SubmitFeedback_Errors(t *testing.T) {

// }

// func Test_ViewResponses_Errors(t *testing.T) {

// }

func setupSQLMock(t *testing.T) sqlmock.Sqlmock {
	d, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer d.Close()
	db = d
	return mock
}
