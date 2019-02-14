package utilities_test

import (
	"testing"

	pb "github.com/aagoldingay/commendeer-go/pb"
	utils "github.com/aagoldingay/commendeer-go/server/utilities"
)

var typeSequence = [...]int{3, 3, 4, 4, 5, 1, 2}

func Test_ProcessQuestions(t *testing.T) {
	sTQs := []*pb.ShortTextQuestion{
		&pb.ShortTextQuestion{Order: 1, Title: "Question 1"},
		&pb.ShortTextQuestion{Order: 3, Title: "Question 3"},
	}
	lTQs := []*pb.LongTextQuestion{
		&pb.LongTextQuestion{Order: 2, Title: "Question 2"},
		&pb.LongTextQuestion{Order: 5, Title: "Question 5"},
	}
	dQs := []*pb.DateQuestion{
		&pb.DateQuestion{Order: 7, Title: "Question 7"},
	}
	rQs := []*pb.RadioQuestion{
		&pb.RadioQuestion{Order: 4, Title: "Question 4", Options: []*pb.AnswerOption{
			&pb.AnswerOption{Title: "Option 1"},
			&pb.AnswerOption{Title: "Option 2"},
		}},
	}
	mQs := []*pb.MultiChoiceQuestion{
		&pb.MultiChoiceQuestion{Order: 6, Title: "Question 6", Options: []*pb.AnswerOption{
			&pb.AnswerOption{Title: "Option 1"},
			&pb.AnswerOption{Title: "Option 2"},
			&pb.AnswerOption{Title: "Option 3"},
		}},
	}
	processedQs := utils.ProcessQuestions(sTQs, lTQs, dQs, rQs, mQs)
	if len(processedQs) != 7 {
		t.Errorf("expected 7 questions, actual: %v\n", len(processedQs))
	}
	for i := 0; i < len(typeSequence); i++ {
		if processedQs[i].QuestionType != typeSequence[i] {
			t.Errorf("question %v: type %v did not match %v\n", i+1, processedQs[i].QuestionType, typeSequence[i])
		}
	}
	if len(processedQs[5].Options) != 2 {
		t.Errorf("radio question did not contain 2 options\n")
	}
	if len(processedQs[6].Options) != 3 {
		t.Errorf("multi question did not contain 3 options\n")
	}
}
