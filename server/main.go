package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"os/signal"

	pb "github.com/aagoldingay/commendeer-go/pb"
	"github.com/aagoldingay/commendeer-go/server/data"
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
	// sessionKey = []byte{35, 250, 103, 131, 245, 255, 194, 76, 198, 188, 157, 217, 82, 104, 157, 5}
	// store      *sessions.CookieStore
)

type server struct{}

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

/*
func adminLoginHandler(w http.ResponseWriter, r *http.Request) {
	// https://gowebexamples.com/sessions/

	if r.Method == "POST" {
		if r.FormValue("loginrequest") == "true" { // attempted log in
			session, _ := store.Get(r, "cookie-name")
			// take username and password from the submitted form

		}
		if r.FormValue("action") == "Send Codes" { // admin attempts to generate codes
			session, _ := store.Get(r, "cookie-name")
			if admin, ok := session.Values["admin"].(bool); !ok || !admin {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			numSent := data.SendCodes(db)
			if numSent > 0 {
				fmt.Println("sent codes")
			}
		}
		if r.FormValue("action") == "Download User Data" { // admin attempts to download user data
			session, _ := store.Get(r, "cookie-name")
			if admin, ok := session.Values["admin"].(bool); !ok || !admin {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			file := data.GenerateCodeCSV(db)
			http.ServeFile(w, r, file)
		}
		tmpl := template.Must(template.ParseFiles("tmpl/dashboard.html"))

		data := PageData{
			PageTitle: "Aston",
		}

		tmpl.Execute(w, data) // loads feedback form
	} else {
		tmpl := template.Must(template.ParseFiles("tmpl/admin.html"))

		data := PageData{
			PageTitle: "Aston",
		}

		tmpl.Execute(w, data)
	}
}

func feedbackHandler(w http.ResponseWriter, r *http.Request) {
	// parse and save feedback to db
	r.ParseForm()
	tmpl := template.Must(template.ParseFiles("tmpl/thanks.html"))

	data := PageData{
		PageTitle: "Aston",
	}

	tmpl.Execute(w, data)
}

func formCreatorHandler(w http.ResponseWriter, r *http.Request) {
	// configure feedback form, save config to db
	session, _ := store.Get(r, "cookie-name")
	if admin, ok := session.Values["admin"].(bool); !ok || !admin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if r.Method == "POST" {
		// save config
	} else {
		// configure selection
		// load existing settings
	}
	tmpl := template.Must(template.ParseFiles("tmpl/formcreator.html"))

	qts, err := data.GetQuestionTypesHTML(db)
	if err != nil {
		http.Error(w, "Problem loading content", http.StatusExpectationFailed)
		return
	}

	data := PageData{
		PageTitle:     "Aston",
		QuestionTypes: template.HTML(qts),
		Body:          template.HTML("<h1>Testing</h1>"), // populate via db
	}

	tmpl.Execute(w, data) // loads feedback form
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" { // on code submit
		r.ParseForm()
		// code := html.EscapeString(r.FormValue("code"))

		// check db for code

		// fmt.Printf("code : %v\n", code)

		if true { // if code is present
			// (db) code now been used
		}

		tmpl := template.Must(template.ParseFiles("tmpl/feedback.html"))

		data := PageData{
			PageTitle: "Aston",
			Body:      template.HTML("<h1>Testing</h1>"), // populate via db
		}

		tmpl.Execute(w, data) // loads feedback form

	} else { // initial website load
		tmpl := template.Must(template.ParseFiles("tmpl/index.html"))

		data := PageData{
			PageTitle: "Aston",
		}

		tmpl.Execute(w, data)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// revoke permissions
	session.Values["authenticated"] = false
	session.Values["admin"] = false
	session.Save(r, w)

	tmpl := template.Must(template.ParseFiles("tmpl/admin.html"))

	data := PageData{
		PageTitle: "Aston",
	}

	tmpl.Execute(w, data)
}

func resultsHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	if admin, ok := session.Values["authenticated"].(bool); !ok || !admin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
}
*/
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
