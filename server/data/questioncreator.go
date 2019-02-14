package data

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	utils "github.com/aagoldingay/commendeer-go/server/utilities"
)

const (
	getQuestionTypesQuery                 = "SELECT * FROM QuestionType" // returned: int, string
	getMultiQuestionsByQuestionnaireQuery = "SELECT QuestionID FROM Question WHERE (QuestionnaireID = $1) AND (QuestionTypeID = 1 OR QuestionTypeID = 2);"
	addQuestionnaireQuery                 = "INSERT INTO Questionnaire (Title) VALUES ($1);"
	getQuestionnaireQuery                 = "SELECT QuestionnaireID FROM Questionnaire WHERE Title = $1"
	addQuestionQuery                      = "INSERT INTO Question (QuestionTypeID, QuestionOrder, Title, QuestionnaireID) VALUES ('%v', '%v', '%v', '%v'); "
	addQuestionOptionQuery                = "INSERT INTO MultiChoiceQuestionOption (QuestionID, OptionDescription) VALUES ('%v', '%v'); "
)

// QuestionType struct contains the ID and description of a questiontype stored in the database
// for use elsewhere in the system where necessary
type QuestionType struct {
	ID   int
	Desc string
}

// CreateForm adds questions to the database
func CreateForm(title string, questions []utils.QuestionInfo, db *sql.DB) error {
	// add questionnaire
	_, err := db.Exec(addQuestionnaireQuery, title)
	if err != nil {
		fmt.Printf("%v: error on CreateForm create questionnaire - %v\n", time.Now(), err)
		return errors.New("problem creating questionnaire")
	}

	// retrieve questionnaire
	var questionnaireID int
	rows, err := db.Query(getQuestionnaireQuery, title)
	rows.Next()
	rows.Scan(&questionnaireID)
	if questionnaireID < 1 {
		fmt.Printf("%v: error on CreateForm retrieving questionnaire\n", time.Now())
		return errors.New("problem creating questionnaire")
	}

	qQuery := ""
	multiQuestionIndexes := make(map[int]int) // [ index ] questionID
	mQKeys := []int{}                         // stores map keys (index from input)

	// setup questions query
	for i := 0; i < len(questions); i++ {
		q := fmt.Sprintf(addQuestionQuery, questions[i].QuestionType, questions[i].Order, questions[i].Title, questionnaireID)
		qQuery += q
		if questions[i].Options != nil {
			multiQuestionIndexes[i] = 0
			mQKeys = append(mQKeys, i)
		}
	}

	_, err = db.Exec(qQuery)
	if err != nil {
		fmt.Printf("%v: error on CreateForm add questions - %v\n", time.Now(), err)
		return errors.New("problem creating questionnaire")
	}

	// setup questionoptions query
	// / get multi choice question IDs from table
	rows, err = db.Query(getMultiQuestionsByQuestionnaireQuery, questionnaireID)
	if err != nil {
		fmt.Printf("%v: error on CreateForm retrieving questions - %v\n", time.Now(), err)
		return errors.New("problem creating questionnaire")
	}
	i := 0
	for rows.Next() {
		var id int
		rows.Scan(&id)
		multiQuestionIndexes[mQKeys[i]] = id
		i++
	}

	// / configure query
	oQuery := ""
	for i := 0; i < len(mQKeys); i++ {
		id := multiQuestionIndexes[mQKeys[i]]
		os := questions[mQKeys[i]].Options
		for j := 0; j < len(os); i++ {
			q := fmt.Sprintf(addQuestionOptionQuery, id, os[j])
			oQuery += q
		}
	}

	_, err = db.Exec(oQuery)
	if err != nil {
		fmt.Printf("%v: error on CreateForm adding multi question options - %v\n", time.Now(), err)
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
