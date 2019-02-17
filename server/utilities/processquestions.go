package utilities

import (
	"fmt"
	"strings"

	pb "github.com/aagoldingay/commendeer-go/pb"
)

// AnswerOption are normalised options for multiple choice questions
type AnswerOption struct {
	Title string
}

//QuestionInfo are normalised question requests in a generic format
type QuestionInfo struct {
	QuestionType int
	Order        int32
	Title        string
	Options      []AnswerOption
}

// IntArrayToString converts array of ints to string with delimiter
func IntArrayToString(arr []int, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(arr), " ", delim, -1), "[]")
}

// ProcessQuestions normalises question requests into a generic format
// easily readable for table input
func ProcessQuestions(shortTexts []*pb.ShortTextQuestion, longTexts []*pb.LongTextQuestion,
	dates []*pb.DateQuestion, radioQs []*pb.RadioQuestion, multiQs []*pb.MultiChoiceQuestion) []QuestionInfo {
	questions := []QuestionInfo{}
	if len(shortTexts) > 0 { // type = 3
		for i := 0; i < len(shortTexts); i++ {
			questions = append(questions, QuestionInfo{3, shortTexts[i].GetOrder(), shortTexts[i].Title, nil})
		}
	}
	if len(longTexts) > 0 { // type = 4
		for i := 0; i < len(longTexts); i++ {
			questions = append(questions, QuestionInfo{4, longTexts[i].GetOrder(), longTexts[i].Title, nil})
		}
	}
	if len(dates) > 0 { // type = 5
		for i := 0; i < len(dates); i++ {
			questions = append(questions, QuestionInfo{5, dates[i].GetOrder(), dates[i].Title, nil})
		}
	}
	if len(radioQs) > 0 { // type = 1
		for i := 0; i < len(radioQs); i++ {
			options := []AnswerOption{}
			for j := 0; j < len(radioQs[i].Options); j++ {
				options = append(options, AnswerOption{radioQs[i].Options[j].Title})
			}
			questions = append(questions, QuestionInfo{1, radioQs[i].Order, radioQs[i].Title, options})
		}
	}
	if len(multiQs) > 0 { // type = 2
		for i := 0; i < len(multiQs); i++ {
			options := []AnswerOption{}
			for j := 0; j < len(multiQs[i].Options); j++ {
				options = append(options, AnswerOption{multiQs[i].Options[j].Title})
			}
			questions = append(questions, QuestionInfo{2, multiQs[i].Order, multiQs[i].Title, options})
		}
	}
	return questions
}
