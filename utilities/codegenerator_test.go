package utilities_test // allows testing of exported functions, similar to integration, but stores in the same folder

import (
	"fmt"
	"testing"

	utils "github.com/aagoldingay/commendeer-go/utilities"
)

var expectedCodes = []string{"A9rI2", "cvTK4", "UHomc", "jcEQv", "mUNER"}

func Test_GenerateCodes(t *testing.T) {
	utils.Setup(0)
	codes := utils.GenerateCodes(5, 5) // 5 codes, length of 5
	fmt.Printf("codes : %v", codes)
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
