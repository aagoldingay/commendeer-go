package main

import (
	"html/template"
	"log"
	"net/http"
)

// PageData models page information for templating
type PageData struct {
	PageTitle string
	Body      template.HTML
}

func adminLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// load dashboard
	} else {
		// load login
	}
}

func feedbackHandler(w http.ResponseWriter, r *http.Request) {
	// parse and save feedback to db
}

func fromCreatorHandler(w http.ResponseWriter, r *http.Request) {
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

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", handler)
	http.HandleFunc("/thanks", feedbackHandler)

	http.HandleFunc("/admin", adminLoginHandler)
	http.HandleFunc("/formCreator", fromCreatorHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
