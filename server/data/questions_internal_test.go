package data

import (
	"testing"

	"github.com/aagoldingay/commendeer-go/pb"
)

func Test_orderQuestionsToArray(t *testing.T) {
	questions := map[int]Question{
		1: Question{1, 1, 4, "question 4", nil},
		2: Question{2, 3, 2, "question 2", nil},
		3: Question{3, 3, 5, "question 5", nil},
		4: Question{4, 4, 1, "question 1", nil},
		5: Question{5, 2, 3, "question 3", nil},
	}
	correctOrder := []*pb.Question{
		&pb.Question{Type: 4, Order: 1, Title: "question 1", Options: nil},
		&pb.Question{Type: 3, Order: 2, Title: "question 2", Options: nil},
		&pb.Question{Type: 2, Order: 3, Title: "question 3", Options: nil},
		&pb.Question{Type: 1, Order: 4, Title: "question 4", Options: nil},
		&pb.Question{Type: 3, Order: 5, Title: "question 5", Options: nil},
	}

	oq := orderQuestionsToArray(questions)
	for i := 0; i < len(oq); i++ {
		if oq[i].Order != correctOrder[i].Order {
			t.Errorf("expected order: %v, actual: %v", correctOrder, oq)
		}
	}
}
