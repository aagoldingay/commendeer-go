package data_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aagoldingay/commendeer-go/server/data"

	utils "github.com/aagoldingay/commendeer-go/server/utilities"
)

func Test_CreateForm(t *testing.T) {
	d, err := dbSetup()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer d.Close()
	err = d.Ping()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	processedQs := []utils.QuestionInfo{}
	err := data.CreateForm("Questionnaire 1", processedQs, d)
}
