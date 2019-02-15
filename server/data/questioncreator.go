package data

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	utils "github.com/aagoldingay/commendeer-go/server/utilities"
)

const (
	comm                  = ", "
	getQuestionTypesQuery = "SELECT * FROM QuestionType" // returned: int, string

	newQuestionnaireQuery = "WITH new_questionnaire as (INSERT INTO Questionnaire (Title) VALUES ('%v') returning questionnaireID), "
	newQuestionQuery      = "new_questions as (INSERT INTO Question (questionTypeID, questionorder, title, questionnaireID) VALUES %v returning questionID, title) "
	newOptionQuery        = "INSERT INTO multichoicequestionoption (OptionDescription, questionID) VALUES %v;"

	questionValues = "(%v, %v, '%v', (select questionnaireID from new_questionnaire))"         // (typeid, order, title)
	optionValues   = "('%v', (select q.questionID from new_questions q where q.title = '%v'))" // (option-title ... = question-title))
)

// QuestionType struct contains the ID and description of a questiontype stored in the database
// for use elsewhere in the system where necessary
type QuestionType struct {
	ID   int
	Desc string
}

// CreateForm adds questions to the database
func CreateForm(title string, questions []utils.QuestionInfo, db *sql.DB) error {
	query := ""

	// compile questionnaire data
	query += fmt.Sprintf(newQuestionnaireQuery, title)

	q, o := "", ""
	// compile questions (q), options (o) queries
	for i := 0; i < len(questions); i++ {
		q += fmt.Sprintf(questionValues, questions[i].QuestionType, questions[i].Order, questions[i].Title)

		if questions[i].Options != nil {
			for j := 0; j < len(questions[i].Options); j++ {
				o += fmt.Sprintf(optionValues, questions[i].Options[j].Title, questions[i].Title) + comm
			}
		}
		if i < len(questions)-1 { // final question does not need a comma
			q += comm
		}
	}
	o = strings.TrimRight(o, comm)

	// compile whole query
	query += fmt.Sprintf(newQuestionQuery, q)
	query += fmt.Sprintf(newOptionQuery, o)

	_, err := db.Exec(query)
	if err != nil {
		fmt.Printf("%v: error on CreateForm executing query - %v", time.Now(), err)
		return errors.New("problem creating questionnaire")
	}
	return nil
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
