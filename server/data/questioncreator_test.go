package data_test

import (
	"testing"

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
