package data

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aagoldingay/commendeer-go/pb"
	utils "github.com/aagoldingay/commendeer-go/server/utilities"
)

const (
	comm                    = ", "
	getQuestionTypesQuery   = "SELECT * FROM QuestionType"                                                                       // returned: int, string
	getQuestionnaireQuery   = "SELECT Title FROM Questionnaire WHERE QuestionnaireID = $1"                                       // returned: string
	getQuestionsQuery       = "SELECT QuestionID, QuestionTypeID, QuestionOrder, Title FROM Question WHERE QuestionnaireID = $1" // returned: int,  int, int, string
	getQuestionOptionsQuery = "SELECT QuestionID, OptionDescription FROM MultiChoiceQuestionOption WHERE QuestionID IN (%v)"     // returned: int, string

	newQuestionnaireQuery = "WITH new_questionnaire as (INSERT INTO Questionnaire (Title) VALUES ('%v') returning questionnaireID), "
	newQuestionQuery      = "new_questions as (INSERT INTO Question (questionTypeID, questionorder, title, questionnaireID) VALUES %v returning questionID, title) "
	newOptionQuery        = "INSERT INTO multichoicequestionoption (OptionDescription, questionID) VALUES %v;"

	questionValues = "(%v, %v, '%v', (select questionnaireID from new_questionnaire))"         // (typeid, order, title)
	optionValues   = "('%v', (select q.questionID from new_questions q where q.title = '%v'))" // (option-title ... = question-title))
)

// Question models a question from the database
type Question struct {
	id, qType, order int // private id field used for assigning options to relevant questions
	title            string
	options          []*pb.QuestionOption
}

// Questionnaire contains information related to a single questionnaire
type Questionnaire struct {
	Title     string
	Questions []*pb.Question
}

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

// GetQuestions queries database for questionnaire information in a generic format
func GetQuestions(qid int, db *sql.DB) Questionnaire {
	if qid < 1 {
		return Questionnaire{}
	}

	// get questionnaire
	row := db.QueryRow(getQuestionnaireQuery, qid)
	var qTitle string
	err := row.Scan(&qTitle)
	if err != nil {
		if err == sql.ErrNoRows {
			return Questionnaire{}
		}
		fmt.Printf("%v: error on get Questionnaire by ID - %v\n", time.Now(), err)
		return Questionnaire{}
	}
	questionnaire := Questionnaire{Title: qTitle}

	// get questions
	qrows, err := db.Query(getQuestionsQuery, qid)
	if err != nil {
		fmt.Printf("%v: error on get questions by questionnaire - %v\n", time.Now(), err)
		return Questionnaire{}
	}
	defer qrows.Close()

	questions := make(map[int]Question) // [id]Question
	multiQs := []int{}                  // id's
	for qrows.Next() {
		var (
			id, qType, order int
			title            string
		)
		err = qrows.Scan(&id, &qType, &order, &title)
		if err != nil {
			fmt.Printf("%v, error on GetQuestions read question rows - %v\n", time.Now(), err)
			return Questionnaire{}
		}
		questions[id] = Question{id, qType, order, title, nil}
		if qType == 1 || qType == 2 {
			multiQs = append(multiQs, id)
		}
	}
	ids := utils.IntArrayToString(multiQs, ", ")
	rows, err := db.Query(fmt.Sprintf(getQuestionOptionsQuery, ids))
	if err != nil {
		fmt.Printf("%v: error on get question options query - %v\n", time.Now(), err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id    int
			title string
		)
		rows.Scan(&id, &title)
		q := questions[id]
		q.options = append(q.options, &pb.QuestionOption{Id: int32(id), Title: title})
		questions[id] = q
	}
	questionnaire.Questions = orderQuestionsToArray(questions)
	return questionnaire
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

func orderQuestionsToArray(q map[int]Question) []*pb.Question {
	o := []*pb.Question{}
	curr := 0
	for curr < len(q) {
		for _, v := range q {
			if v.order == curr+1 {
				e := &pb.Question{Id: int32(v.id), Type: int32(v.qType), Order: int32(v.order), Title: v.title, Options: v.options}
				o = append(o, e)
				curr++
			}
		}
	}
	return o
}
