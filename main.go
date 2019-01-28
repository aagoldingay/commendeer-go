package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "commendeer"
)

// PageData models page information for templating
type PageData struct {
	PageTitle string
	Body      template.HTML
}

var db *sql.DB

func adminLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if r.FormValue("action") == "Send Codes" {
			fmt.Println("sent codes")
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
	if r.Method == "POST" {
		// save config
	} else {
		// configure selection
		// load existing settings
	}
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

func resultsHandler(w http.ResponseWriter, r *http.Request) {

}

func dbSetup() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	d, err := sql.Open("postgres", psqlInfo) // configures db info for connection via code
	if err != nil {
		return nil, err // if there's an error, we don't want to continue listening to the port!
	}
	defer d.Close()

	err = d.Ping() // open connection to the database
	if err != nil {
		return nil, err
	}
	return d, nil
}

func main() {
	d, err := dbSetup()
	if err != nil {
		panic(err)
	}
	db = d

	// UI access setup
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", handler)
	http.HandleFunc("/thanks", feedbackHandler)

	http.HandleFunc("/admin", adminLoginHandler)
	http.HandleFunc("/formCreator", formCreatorHandler)
	http.HandleFunc("/results", resultsHandler)

	// listen to port 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}
