package utilities_test // allows testing of exported functions, similar to integration, but stores in the same folder

import (
	"testing"

	utils "github.com/aagoldingay/commendeer-go/server/utilities"
)

var expectedCodes = []string{"A9rI2", "cvTK4", "UHomc", "jcEQv", "mUNER"}

func Test_GenerateCode_CorrectLength(t *testing.T) {
	l := 10
	c := utils.GenerateCode(l)
	if len(c) != l {
		t.Errorf("generatecode failed: expected %v, actual %v", l, len(c))
	}
}

func Test_GenerateCodes(t *testing.T) {
	utils.Setup(0)
	codes := utils.GenerateCodes(5, 5) // 5 codes, length of 5
	if len(codes) != 5 {
		t.Errorf("amount of codes generated is incorrect. expected : %v, actual : %v\n", 5, len(codes))
	}
	prs := 0
	for i := 0; i < len(expectedCodes); i++ {
		for j := 0; j < len(codes); j++ {
			if expectedCodes[i] == codes[j] {
				prs++
			}
		}
	}
	if prs != 5 {
		t.Errorf("unexpected codes; expected %v, actual : %v\n", expectedCodes, codes)
	}
}
