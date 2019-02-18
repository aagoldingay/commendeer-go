package data_test

import (
	"testing"

	"github.com/aagoldingay/commendeer-go/pb"

	"github.com/aagoldingay/commendeer-go/server/data"

	utils "github.com/aagoldingay/commendeer-go/server/utilities"
)

func CreateForm(t *testing.T) {
	processedQs := []utils.QuestionInfo{
		utils.QuestionInfo{QuestionType: 3, Order: 1, Title: "Question 1"},
		utils.QuestionInfo{QuestionType: 3, Order: 3, Title: "Question 3"},
		utils.QuestionInfo{QuestionType: 1, Order: 2, Title: "Question 2", Options: []utils.AnswerOption{
			utils.AnswerOption{Title: "Option 1"},
			utils.AnswerOption{Title: "Option 2"},
		}},
		utils.QuestionInfo{QuestionType: 2, Order: 4, Title: "Question 4", Options: []utils.AnswerOption{
			utils.AnswerOption{Title: "Option 1"},
			utils.AnswerOption{Title: "Option 2"},
			utils.AnswerOption{Title: "Option 3"},
		}},
		utils.QuestionInfo{QuestionType: 4, Order: 5, Title: "Question 5"},
	}
	err := data.CreateForm("Questionnaire 1", processedQs, db)
	if err != nil {
		t.Errorf("error on CreateForm: %v", err)
	}
}

func GetQuestions(t *testing.T) {
	correct := data.Questionnaire{
		Title: "Questionnaire 1",
		Questions: []*pb.Question{
			&pb.Question{Id: 1, Type: 3, Order: 1, Title: "Example Question"},
			&pb.Question{Id: 2, Type: 2, Order: 2, Title: "Example Multi Choice Question", Options: []*pb.QuestionOption{
				&pb.QuestionOption{Title: "Option 1"},
				&pb.QuestionOption{Title: "Option 2"},
			}},
		},
	}
	q := data.GetQuestions(1, db)

	if len(q.Questions) != len(correct.Questions) {
		t.Errorf("Incorrect number of questions returned - expected: %v, actual: %v\n", len(correct.Questions), len(q.Questions))
	}
	if q.Title != correct.Title {
		t.Errorf("Incorrect title - expected: %v, actual:  %v\n", correct.Title, q.Title)
	}
	if q.Questions[0].Id != correct.Questions[0].Id {
		t.Errorf("Incorrect id in first position - expected: %v, actual: %v\n", correct.Questions[0].Id, q.Questions[0].Id)
	}
	if q.Questions[0].Order != correct.Questions[0].Order {
		t.Errorf("First question order did not match - expected: %v, actual: %vzn", correct.Questions[0].Order, q.Questions[0].Order)
	}
	if len(q.Questions[1].Options) != len(correct.Questions[1].Options) {
		t.Errorf("Incorrect number of question options - expected: %v, actual: %v\n", len(correct.Questions[1].Options), len(q.Questions[1].Options))
	}
}
