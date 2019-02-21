package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"

	pb "github.com/aagoldingay/commendeer-go/pb"
	"github.com/aagoldingay/commendeer-go/server/data"
	utils "github.com/aagoldingay/commendeer-go/server/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "commendeer"
)

var (
	db *sql.DB
)

type server struct{}

func (s *server) CreateAccessCodes(ctx context.Context, in *pb.CreateCodeRequest) (*pb.CreateCodeResponse, error) {
	adminOnly := true
	if in.QuestionnaireID < 1 {
		return &pb.CreateCodeResponse{Error: pb.Error_BADREQUEST, ErrorDetails: "Invalid questionnaire"}, nil
	}

	auth, err := data.CheckAuthorised(in.Authcode, adminOnly, db)
	if err != nil {
		if err.Error() == "invalid code" {
			return &pb.CreateCodeResponse{Error: pb.Error_BADREQUEST, ErrorDetails: "Invalid code"}, nil
		}
	}
	if !auth {
		return &pb.CreateCodeResponse{Error: pb.Error_FORBIDDEN, ErrorDetails: "No authorisation"}, nil
	}
	sent, err := data.SendCodes(int(in.QuestionnaireID), db)
	if !sent {
		return &pb.CreateCodeResponse{Error: pb.Error_NIL, ErrorDetails: "Codes did not send - confirm with administrator"}, nil
	}
	return &pb.CreateCodeResponse{Error: pb.Error_OK, ErrorDetails: ""}, nil
}

func (s *server) CreateForm(ctx context.Context, in *pb.CreateFormRequest) (*pb.CreateResponse, error) {
	adminOnly := true
	qCount := len(in.DateQuestions) + len(in.LongTextQuestions) + len(in.MultiChoiceQuestions) + len(in.LongTextQuestions) + len(in.RadioQuestions)
	if qCount == 0 {
		return &pb.CreateResponse{Error: pb.Error_NOTACCEPTABLE, ErrorDetails: "No questions provided"}, nil
	}
	if in.Title == "" {
		return &pb.CreateResponse{Error: pb.Error_NOTACCEPTABLE, ErrorDetails: "No title provided"}, nil
	}
	auth, err := data.CheckAuthorised(in.AuthCode, adminOnly, db)
	if err != nil {
		if err.Error() == "invalid code" {
			return &pb.CreateResponse{Error: pb.Error_BADREQUEST, ErrorDetails: "Invalid code"}, nil
		}
	}
	if !auth {
		return &pb.CreateResponse{Error: pb.Error_FORBIDDEN, ErrorDetails: "No authorisation"}, nil
	}
	nQuestions := utils.ProcessQuestions(in.ShortTextQuestions, in.LongTextQuestions, in.DateQuestions, in.RadioQuestions, in.MultiChoiceQuestions)
	err = data.CreateForm(in.Title, nQuestions, db)
	if err != nil {
		if err.Error() == "problem creating questionnaire" {
			return &pb.CreateResponse{Error: pb.Error_INTERNALERROR, ErrorDetails: "Problem creating questionnaire"}, nil
		}
	}
	return &pb.CreateResponse{Error: pb.Error_OK, ErrorDetails: ""}, nil
}

func (s *server) GetFeedbackForm(ctx context.Context, in *pb.GetFormRequest) (*pb.FormResponse, error) {
	if in.Email == "" {
		return &pb.FormResponse{Error: pb.Error_BADREQUEST, ErrorDetails: "Invalid code or email"}, nil
	}
	qid, err := data.GetAccessCode(in.Email, in.AccessCode, db)
	if err != nil {
		if strings.Contains(err.Error(), "code not of desired length: ") || err.Error() == "code or user not found" {
			return &pb.FormResponse{Error: pb.Error_BADREQUEST, ErrorDetails: "Invalid code or email"}, nil
		}
		if err.Error() == "problem on GetAccessCode" {
			return &pb.FormResponse{Error: pb.Error_BADREQUEST, ErrorDetails: "Unable to authenticate"}, nil
		}
		if err.Error() == "unbound code" {
			return &pb.FormResponse{Error: pb.Error_INTERNALERROR, ErrorDetails: "Unbound code"}, nil
		}
	}
	if qid < 1 {
		return &pb.FormResponse{Error: pb.Error_FORBIDDEN, ErrorDetails: "Code already used"}, nil
	}

	questionnaire := data.GetQuestions(qid, db)
	if len(questionnaire.Questions) < 1 {
		return &pb.FormResponse{Error: pb.Error_NIL, ErrorDetails: "Problem encountered"}, nil
	}

	return &pb.FormResponse{Error: pb.Error_OK, ErrorDetails: "", Title: questionnaire.Title, Questions: questionnaire.Questions}, nil
}

func (s *server) LoginUser(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	usr, err := data.Login(in.Username, in.Password, db)
	if err != nil {
		if err.Error() == "incorrect username or password" {
			return &pb.LoginResponse{Error: pb.Error_BADREQUEST, ErrorDetails: "Incorrect username or password"}, nil
		}
		if err.Error() == "error on authorisation" {
			return &pb.LoginResponse{Error: pb.Error_INTERNALERROR, ErrorDetails: "Problem logging in"}, nil
		}
	}
	return &pb.LoginResponse{Username: usr.Username, Authcode: usr.Code, Error: pb.Error_OK, ErrorDetails: ""}, nil
}

func (s *server) LogoutUser(ctx context.Context, in *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	err := data.Logout(in.Authcode, db)
	if err != nil {
		if err.Error() == "invalid code" {
			return &pb.LogoutResponse{Error: pb.Error_BADREQUEST, ErrorDetails: "Authorisation invalid"}, nil
		}
		if err.Error() == "error on logout" {
			return &pb.LogoutResponse{Error: pb.Error_INTERNALERROR, ErrorDetails: "Problem logging out"}, nil
		}
		if err.Error() == "unknown code" {
			return &pb.LogoutResponse{Error: pb.Error_BADREQUEST, ErrorDetails: "Problem identifying user"}, nil
		}
	}
	return &pb.LogoutResponse{Error: pb.Error_OK, ErrorDetails: ""}, nil
}

func (s *server) SubmitFeedback(ctx context.Context, in *pb.PostFeedbackRequest) (*pb.PostFeedbackResponse, error) {
	err := data.SubmitResponse(in, db)
	if err != nil {
		if err.Error() == "code changed" || err.Error() == "invalid questionnaire" || err.Error() == "no code found" {
			return &pb.PostFeedbackResponse{Error: pb.Error_BADREQUEST, ErrorDetails: "Session no longer matches"}, nil
		}
		if err.Error() == "feedback has already been submitted" {
			return &pb.PostFeedbackResponse{Error: pb.Error_FORBIDDEN, ErrorDetails: "Feedback already submitted by this user"}, nil
		}
		if err.Error() == "problem executing" {
			return &pb.PostFeedbackResponse{Error: pb.Error_INTERNALERROR, ErrorDetails: "Problem executing"}, nil
		}
		if err.Error() == "incorrect number of questions" {
			return &pb.PostFeedbackResponse{Error: pb.Error_BADREQUEST, ErrorDetails: "Incorrect number of questions answered"}, nil
		}
	}
	return &pb.PostFeedbackResponse{Error: pb.Error_OK, ErrorDetails: ""}, nil
}

func (s *server) ViewResponses(ctx context.Context, in *pb.ViewRequest) (*pb.ViewResponse, error) {
	adminOnly := false
	auth, err := data.CheckAuthorised(in.AuthCode, adminOnly, db)
	if err != nil {
		if err.Error() == "invalid code" {
			return &pb.ViewResponse{Error: pb.Error_BADREQUEST, ErrorDetails: "Invalid code"}, nil
		}
	}
	if !auth {
		return &pb.ViewResponse{Error: pb.Error_FORBIDDEN, ErrorDetails: "No authorisation"}, nil
	}
	d, err := data.GetResponses(int(in.QuestionnaireID), db)
	if err != nil {
		switch err.Error() {
		case "invalid questionnaire":
			return &pb.ViewResponse{Error: pb.Error_BADREQUEST, ErrorDetails: "Invalid questionnaire"}, nil
		case "problem encountered":
			return &pb.ViewResponse{Error: pb.Error_INTERNALERROR, ErrorDetails: "Problem encountered"}, nil
		}
	}
	return &pb.ViewResponse{Error: pb.Error_OK, ErrorDetails: "Problem encountered", Questions: d}, nil
}

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

func main() {
	d, err := dbSetup()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer d.Close()
	db = d
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		os.Exit(1)
	}
	s := grpc.NewServer()
	pb.RegisterCommendeerServer(s, &server{})

	go func() {
		// Register reflection service on gRPC server.
		reflection.Register(s)
		if err := s.Serve(lis); err != nil {
			fmt.Printf("failed to serve: %v", err)
			os.Exit(1)
		}
	}()
	<-stop
	fmt.Printf("shutting down service")

	s.Stop()
}
