package data

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

const (
	getQuestionTypesQuery = "SELECT * FROM QuestionType" // returned: int, string
)

// QuestionType struct contains the ID and description of a questiontype stored in the database
// for use elsewhere in the system where necessary
type QuestionType struct {
	ID   int
	Desc string
}

// GetQuestionTypesHTML returns HTML select input type containing returned types
func GetQuestionTypesHTML(db *sql.DB) (string, error) {
	types := GetQuestionTypes(db)
	if len(types) < 1 {
		return "", errors.New("no types returned")
	}
	html := "<select name=\"questiontypes\"><option value=\"0\">select one</option>"
	for _, t := range types {
		html += fmt.Sprintf("<option value=\"%v\">%v</option>", t.ID, t.Desc)
	}
	html += "</select>"

	return html, nil
}

// GetQuestionTypes queries the database for supported questiontypes
// returns array containing supported questiontypes
func GetQuestionTypes(db *sql.DB) []QuestionType {
	rows, err := db.Query(getQuestionTypesQuery)
	if err != nil {
		fmt.Printf("%v: error on GetQuestionTypes query - %v\n", time.Now(), err)
	}
	defer rows.Close()

	var (
		id   int
		desc string
	)
	types := []QuestionType{}

	// populate slice of QuestionTypes for returning
	for rows.Next() {
		rows.Scan(&id, &desc)
		types = append(types, QuestionType{id, desc})
	}
	if len(types) < 1 {
		fmt.Printf("%v: GetQuestionTypes - no data types returned", time.Now())
	}
	return types
}
