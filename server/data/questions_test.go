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

func SubmitResponse(t *testing.T) {
	qid := 1
	a := &pb.PostFeedbackRequest{
		AccessCode: "helloworld", QuestionnaireID: int32(qid), Questions: []*pb.AnsweredQuestion{
			&pb.AnsweredQuestion{Id: 1, Type: 3, Answer: "hello", SelectedOptions: nil},
			&pb.AnsweredQuestion{Id: 2, Type: 2, Answer: "", SelectedOptions: []*pb.SelectedOption{
				&pb.SelectedOption{Id: 1},
				&pb.SelectedOption{Id: 2},
			}},
		}}

	err := data.SubmitResponse(a, db)
	if err != nil {
		t.Errorf("SubmitResponse errored: %v\n", err)
	}

	// refer to GetResponses test for default values before totalling expected response including additional above
	q, err := data.GetResponses(qid, db)
	if len(q) != 2 {
		t.Errorf("questionnaire did not return as expected")
	}
	if q[0].Type != 3 && len(q[0].TextAnswers) != 3 {
		t.Errorf("text question did not return as expected")
	}
	if q[1].Type != 2 && q[1].OptionAnswers[0].Values != 2 && q[1].OptionAnswers[1].Values != 3 {
		t.Errorf("multi choice question did not return as expected")
	}

	err = data.SubmitResponse(a, db)
	if err.Error() != "feedback has already been submitted" {
		t.Errorf("second submit did not error")
	}
}

func GetResponses(t *testing.T) {
	qid := 1
	q, err := data.GetResponses(qid, db)
	if err != nil {
		t.Errorf("Test_GetResponses errored; shouldn't have - %v\n", err)
	}
	if len(q) != 2 {
		t.Errorf("incorrect number of questions returned. expected: %v, actual: %v\n", 2, len(q))
	}
	if q[0].Type != 3 && q[1].Type != 2 {
		t.Errorf("an incorrect question type: (1 = t:%v, wanted:%v), (2 = t:%v, wanted:%v)\n", q[0].Type, 3, q[1].Type, 2)
	}
	if q[0].Title != "Example Question" && q[1].Title != "Example Multi Choice Question" {
		t.Errorf("an incorrect question title returned (1 = %v), (2 = %v)\n", q[0].Title, q[1].Title)
	}
	if len(q[0].TextAnswers) != 2 {
		t.Errorf("question 1 should have 2 answers\n")
	}
	if len(q[1].OptionAnswers) != 2 {
		t.Errorf("question 2 should have 2 options\n")
	}
	if q[1].OptionAnswers[0].Values != 1 && q[1].OptionAnswers[0].Values != 2 {
		t.Errorf("options do not have correct values\n")
	}
}
