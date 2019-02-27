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
	comm                     = ", "
	getQuestionTypesQuery    = "SELECT * FROM QuestionType"                                                                                                // returned: int, string
	getQuestionnaireQuery    = "SELECT Title FROM Questionnaire WHERE QuestionnaireID = $1"                                                                // returned: string
	getQuestionsQuery        = "SELECT QuestionID, QuestionTypeID, QuestionOrder, Title FROM Question WHERE QuestionnaireID = $1"                          // returned: int,  int, int, string
	getQuestionOptionsQuery  = "SELECT multichoicequestionoptionid, QuestionID, OptionDescription FROM MultiChoiceQuestionOption WHERE QuestionID IN (%v)" // returned: int, string
	getQuestionsReturnIDType = "SELECT QuestionID, QuestionTypeID FROM Question WHERE QuestionnaireID = $1"                                                // returns int, int

	newQuestionnaireQuery = "WITH new_questionnaire as (INSERT INTO Questionnaire (Title) VALUES ('%v') returning questionnaireID), "
	newQuestionQuery      = "new_questions as (INSERT INTO Question (questionTypeID, questionorder, title, questionnaireID) VALUES %v returning questionID, title) "
	newOptionQuery        = "INSERT INTO multichoicequestionoption (OptionDescription, questionID) VALUES %v;"

	newAnswerQuery       = "INSERT INTO Question_Result (QuestionID, CodeID, Answer) VALUES %v; "
	newAnswerOptionQuery = "INSERT INTO MultiChoiceQuestionOption_Result (QuestionID, MultiChoiceQuestionOptionID, CodeID) VALUES %v;"

	questionValues = "(%v, %v, '%v', (select questionnaireID from new_questionnaire))"         // (typeid, order, title)
	optionValues   = "('%v', (select q.questionID from new_questions q where q.title = '%v'))" // (option-title ... = question-title))

	// int, int, int, string, int, string, int, string
	//getResponsesQuery = "select q.questionorder, q.questionid, q.questiontypeid, q.title, qr.question_resultid as answer_id, qr.answer, mr.multichoicequestionoptionid as option_id, mo.optiondescription from question q left join question_result qr on qr.questionid = q.questionid left join multichoicequestionoption_result mr on mr.questionid = q.questionid left join multichoicequestionoption mo on mo.multichoicequestionoptionid = mr.multichoicequestionoptionid where questionnaireid = $1 order by q.questionorder"
	getResponsesQuery = "select q.questionorder, q.questiontypeid, q.title, (case when qr.answer is null then '' else qr.answer end) as answer, (case when mr.multichoicequestionoptionid is null then 0 else mr.multichoicequestionoptionid end) as option_id, (case when mo.optiondescription is null then '' else mo.optiondescription end) as optiondescription from question q left join question_result qr on qr.questionid = q.questionid left join multichoicequestionoption_result mr on mr.questionid = q.questionid left join multichoicequestionoption mo on mo.multichoicequestionoptionid = mr.multichoicequestionoptionid where questionnaireid = $1 order by q.questionorder"
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

// Qu used for HTML templating
type Qu struct {
	TypeClass, QuestionID, Title, Answer string
}

// HTMLQ models a questionnaire containing template elements required for questionnaire_template.html
type HTMLQ struct {
	Title, AccessCode string
	Questions         []Qu
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
			id, oid int
			title   string
		)
		rows.Scan(&oid, &id, &title)
		q := questions[id]
		q.options = append(q.options, &pb.QuestionOption{Id: int32(oid), Title: title})
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

// GetResponses queries the database for results relating to the supplied questionnaire
func GetResponses(qid int, db *sql.DB) ([]*pb.QuestionResponse, error) {
	if qid < 1 {
		return nil, errors.New("invalid questionnaire")
	}
	rows, err := db.Query(getResponsesQuery, qid)
	if err != nil {
		fmt.Printf("%v: error on GetResponses query - %v", time.Now(), err)
		return nil, errors.New("problem encountered")
	}
	defer rows.Close()

	var (
		order, typeid, optid       int
		title, answer, optionTitle string
	)
	currOrder := 0
	var currQuestion *pb.QuestionResponse
	qs := []*pb.QuestionResponse{}
	qTextAs := []string{}
	qops := make(map[int]string) // [optid]title
	optVal := make(map[int]int)  // [optid]total

	// int, int, int, string, int, string, int, string
	for rows.Next() {
		rows.Scan(&order, &typeid, &title, &answer, &optid, &optionTitle)
		if order < 1 {
			fmt.Printf("%v: error on GetResponses read rows - %v", time.Now(), err)
			return nil, errors.New("problem encountered")
		}
		if currOrder != 0 && currOrder != order {
			// commit previous to map
			if int(currQuestion.Type) < 3 {
				currQuestion.TextAnswers = nil
				a := []*pb.MultiChoiceAnswers{}
				for k, v := range qops {
					a = append(a, &pb.MultiChoiceAnswers{Id: int32(k), Title: v, Values: int32(optVal[k])})
				}
				currQuestion.OptionAnswers = a

			} else {
				currQuestion.OptionAnswers = nil
				currQuestion.TextAnswers = qTextAs
			}
			qs = append(qs, currQuestion)

			// reset data structures
			qTextAs, qops, optVal = []string{}, make(map[int]string), make(map[int]int)
			currQuestion = &pb.QuestionResponse{Type: int32(typeid), Title: title}
			currOrder++
		}
		if currOrder == 0 {
			currQuestion = &pb.QuestionResponse{Type: int32(typeid), Title: title}
			currOrder++
		}

		if int32(currQuestion.Type) < 3 {
			if _, ok := qops[optid]; !ok {
				qops[optid] = optionTitle
			}
			if _, ok := optVal[optid]; !ok {
				optVal[optid] = 0
			}
			optVal[optid]++
		} else {
			qTextAs = append(qTextAs, answer)
		}
	}
	// last commit to qs, as rows.Next() returns false after final result
	if int(currQuestion.Type) < 3 {
		currQuestion.TextAnswers = nil
		a := []*pb.MultiChoiceAnswers{}
		for k, v := range qops {
			a = append(a, &pb.MultiChoiceAnswers{Id: int32(k), Title: v, Values: int32(optVal[k])})
		}
		currQuestion.OptionAnswers = a

	} else {
		currQuestion.OptionAnswers = nil
		currQuestion.TextAnswers = qTextAs
	}
	qs = append(qs, currQuestion)
	return qs, nil
}

// SubmitResponse takes questionnaire response from the client and adds it do the database
func SubmitResponse(f *pb.PostFeedbackRequest, db *sql.DB) error {
	if f.QuestionnaireID < 1 {
		return errors.New("invalid questionnaire")
	}
	if len(f.AccessCode) != CodeLen {
		return errors.New("code changed")
	}
	codeID, used := GetAccessCodeID(f.AccessCode, int(f.QuestionnaireID), db)
	if codeID < 1 {
		return errors.New("no code found")
	}
	if used {
		return errors.New("feedback has already been submitted")
	}

	rows, err := db.Query(getQuestionsReturnIDType, f.QuestionnaireID)
	if err != nil {
		fmt.Printf("%v: SubmitResponse problem getting questions by id - %v", time.Now(), err)
		return errors.New("problem executing")
	}
	defer rows.Close()

	questions := make(map[int]int) // [id]type
	var (
		id, qType int
	)
	for rows.Next() {
		rows.Scan(&id, &qType)
		questions[id] = qType
	}

	if len(questions) != len(f.Questions) {
		return errors.New("incorrect number of questions")
	}

	r := ""
	or := ""
	for i := 0; i < len(f.Questions); i++ {
		if t, ok := questions[int(f.Questions[i].Id)]; ok && (t == 1 || t == 2) { // if type is multi choice
			for j := 0; j < len(f.Questions[i].SelectedOptions); j++ { // (QuestionID, MultiChoiceQuestionOptionID, CodeID)
				if or != "" {
					or += comm
				}
				or += fmt.Sprintf("(%v, %v, %v)", f.Questions[i].Id, f.Questions[i].SelectedOptions[j].Id, codeID)
			}
			continue
		}
		if _, ok := questions[int(f.Questions[i].Id)]; ok { // catch if any other type
			if r != "" {
				r += comm
			}
			r += fmt.Sprintf("(%v, %v, '%v')", f.Questions[i].Id, codeID, f.Questions[i].Answer)
		}
	}

	// exec query
	q := fmt.Sprintf(newAnswerQuery, r)
	q += fmt.Sprintf(newAnswerOptionQuery, or)

	_, err = db.Exec(q)
	if err != nil {
		fmt.Printf("%v: error on SubmitResponse insert answers - %v\n", time.Now(), err)
		return errors.New("problem executing")
	}

	if !UpdateCode(codeID, db) {
		fmt.Printf("%v: error updating code - %v\n", time.Now(), codeID)
		return errors.New("problem executing")
	}
	return nil
}

// HTMLQuestionnaire formats a questionnaire from data pkg into HTML
func HTMLQuestionnaire(q Questionnaire, accessCode string) HTMLQ {
	data := HTMLQ{Title: q.Title, AccessCode: fmt.Sprintf("<input type=\"hidden\" name=\"accesscode\" value=\"%v\"/>", accessCode)}
	questions := []Qu{}

	for i := 0; i < len(q.Questions); i++ {
		n := Qu{QuestionID: fmt.Sprint(q.Questions[i].Id), TypeClass: fmt.Sprint(q.Questions[i].Type), Title: q.Questions[i].Title}
		switch q.Questions[i].Type {
		case 1:
		case 2:
			div := "<div class=\"options\">%v</div>"
			opts := ""
			t := "radio"
			if q.Questions[i].Type == 2 {
				t = "checkbox"
			}
			for j := 0; j < len(q.Questions[i].Options); j++ {
				opts += fmt.Sprintf("<input type=\"%v\" name=\"option_%v\" value=\"%v\">%v<br>", t, q.Questions[i].Options[j].Id, q.Questions[i].Options[j].Id, q.Questions[i].Options[j].Title)
			}
			n.Answer = fmt.Sprintf(div, opts)
			break
		case 3: // type="text"
			n.Answer = fmt.Sprint("<input type=\"text\" name=\"answer\"><br>")
			break
		case 4:
			n.Answer = fmt.Sprint("<textarea name=\"answer\"></textarea>")
			break
		case 5: // type="date"
			n.Answer = fmt.Sprint("<input type=\"date\" name=\"answer\"><br>")
			break
		}
		questions = append(questions, n)
	}
	data.Questions = questions
	return data
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
