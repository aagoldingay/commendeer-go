package main

import (
	"testing"

	"github.com/aagoldingay/commendeer-go/server/data"
	// 	sqlmock "github.com/DATA-DOG/go-sqlmock"
	// 	pb "github.com/aagoldingay/commendeer-go/pb"
	// 	utils "github.com/aagoldingay/commendeer-go/server/utilities"
)

func Test_HTMLQuestions(t *testing.T) {
	d, err := dbSetup()
	if err != nil {
		t.Errorf("error connecting to database : %v\n", err)
	}
	// db = d
	defer d.Close()
	// defer db.Close()

	err = d.Ping()
	if err != nil {
		t.Errorf("error connecting to database : %v\n", err)
	}
	q := data.GetQuestions(1, d)
	data := data.HTMLQuestionnaire(q, "helloworld")
	if data.AccessCode != "<input type=\"hidden\" name=\"accesscode\" value=\"helloworld\"/>" {
		t.Fail()
	}
}

// func Test_Server(t *testing.T) {
// 	// setup tests
// 	d, err := dbSetup()
// 	if err != nil {
// 		t.Errorf("error connecting to database : %v\n", err)
// 	}
// 	db = d
// 	defer d.Close()
// 	defer db.Close()

// 	err = db.Ping()
// 	if err != nil {
// 		t.Errorf("error connecting to database : %v\n", err)
// 	}

// 	t.Run("Test_CreateAccessCodes_Errors", func(t *testing.T) { // run before code altering UserInfo table
// 		CreateAccessCodesErrors(t)
// 	})
// }

// func CreateAccessCodesErrors(t *testing.T) {
// 	s := server{}
// 	reqs := []*pb.CreateCodeRequest{
// 		&pb.CreateCodeRequest{Authcode: "", QuestionnaireID: 0},                     // invalid questionnaire
// 		&pb.CreateCodeRequest{Authcode: "test", QuestionnaireID: 1},                 // invalid code
// 		&pb.CreateCodeRequest{Authcode: "helloworldhelloworld", QuestionnaireID: 1}, // no authorisation
// 	}
// 	expResp := []*pb.CreateCodeResponse{
// 		&pb.CreateCodeResponse{Error: pb.Error_BADREQUEST, ErrorDetails: "Invalid questionnaire"},
// 		&pb.CreateCodeResponse{Error: pb.Error_BADREQUEST, ErrorDetails: "Invalid code"},
// 		&pb.CreateCodeResponse{Error: pb.Error_FORBIDDEN, ErrorDetails: "No authorisation"},
// 	}
// 	for i := 0; i < len(reqs); i++ {
// 		resp, err := s.CreateAccessCodes(context.Background(), reqs[i])
// 		if err != nil {
// 			t.Errorf("problem with CreateAccessCodes server method - shouldn't return an error")
// 		}
// 		if resp.Error != expResp[i].Error {
// 			t.Errorf("[%v] : expected error not correct", i+1)
// 		}
// 		if resp.ErrorDetails != expResp[i].ErrorDetails {
// 			t.Errorf("[%v] : expected error details not correct", i+1)
// 		}
// 	}
// }

// func Test_CreateForm_Errors(t *testing.T) {

// }

// func Test_GetFeedbackForm_Errors(t *testing.T) {

// }

// func Test_LoginUser_Errors(t *testing.T) {
// 	d, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer d.Close()
// 	db = d

// 	utils.Setup(0)
// 	s := server{}
// 	u, p := "admin1", "4dm1n123"
// 	mr := sqlmock.NewRows([]string{"userid", "username", "pass", "administrator", "salt"}).
// 		AddRow("1", "admin1", "57d8da63dbcfd720673fd0622ac91549", "true", "zRvjFZ8Amq")
// 	mock.ExpectQuery("^SELECT (.+) FROM userinfo where Username = (.+)").
// 		WithArgs(u).WillReturnRows(mr)
// 	mock.ExpectExec(`^INSERT INTO authcodes \(UserID, Code, Administrator\) VALUES \(.+\)`).
// 		WithArgs(1, "mUNERA9rI2cvTK4UHomc", true).WillReturnResult(sqlmock.NewResult(1, 1))

// 	req := &pb.LoginRequest{Username: u, Password: p}
// 	resp, err := s.LoginUser(context.Background(), req)
// 	if err != nil {
// 		t.Errorf("problem with LoginUser server method - shouldn't return an error")
// 	}
// 	if resp.Error != pb.Error_OK {
// 		t.Errorf("expected OK - error : %v\n", resp.Error.String())
// 	}
// 	if resp.ErrorDetails != "" {
// 		t.Errorf("expected no error details : %v\n", resp.ErrorDetails)
// 	}
// 	if resp.Username != u {
// 		t.Errorf("incorrect user returned : %v\n", resp.Username)
// 	}
// }

// func Test_LogoutUser_Errors(t *testing.T) {

// }

// func Test_SubmitFeedback_Errors(t *testing.T) {

// }

// func Test_ViewResponses_Errors(t *testing.T) {

// }
