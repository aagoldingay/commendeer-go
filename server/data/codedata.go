package data

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"time"

	utils "github.com/aagoldingay/commendeer-go/server/utilities"
)

const (
	// CodeLen enables external verification of desired length of a code
	CodeLen          = 10
	getCodeQuery     = "SELECT CodeID, Used, QuestionnaireID FROM AccessCode WHERE Email = $1 AND Code = $2"
	sendCodeQuery    = "SELECT CodeID FROM AccessCode WHERE Code IS NULL AND QuestionnaireID = $1;"
	getCodeDataQuery = "SELECT systemusername, Email, Code FROM AccessCode WHERE Email IS NOT NULL AND Code IS NOT NULL;"
	codeDataFile     = "testerdata.csv"
	update           = "UPDATE AccessCode SET Code = '%v' WHERE CodeID = %v AND QuestionnaireID = %v; "
	setUsed          = "UPDATE AccessCode SET Used = TRUE WHERE CodeID = $1 AND USED = FALSE;"
)

// GenerateCodeCSV will return data from the database containing username, email and code per registered beta user
func GenerateCodeCSV(db *sql.DB) string {
	// TODO - AMEND TO ADD QUESTIONNAIREID
	rows, err := db.Query(getCodeDataQuery)
	if err != nil {
		fmt.Printf("%v: error on GenerateCodeCSV query - %v\n", time.Now(), err)
		return ""
	}
	defer rows.Close()

	var (
		name, email, code string
	)
	data := [][]string{}
	data = append(data, []string{"Username", "Email", "Code"}) // column headers
	for rows.Next() {
		rows.Scan(&name, &email, &code)
		data = append(data, []string{name, email, code}) //append row to data slice
	}

	if _, err := os.Stat(codeDataFile); !os.IsNotExist(err) {
		err = os.Remove(codeDataFile)
		if err != nil {
			fmt.Printf("%v: error on previous code deletion - %v", time.Now(), err)
		}
	}
	file, err := os.Create(codeDataFile)
	if err != nil {
		fmt.Printf("%v: error on file creation (GenerateCodeCSV) - %v", time.Now(), err)
		return ""
	}
	defer file.Close()

	w := csv.NewWriter(file)

	w.WriteAll(data)

	w.Flush()
	if err := w.Error(); err != nil {
		fmt.Printf("%v: error on GenerateCodeCSV - %v", time.Now(), err)
		return ""
	}
	return codeDataFile
}

// GetAccessCode takes an email and code, then searches the database for a relevant entry
// errors return based on incorrect code length or 'no code or user found'
// no error when code already used
func GetAccessCode(email, code string, db *sql.DB) (int, error) {
	if len(code) != CodeLen {
		return 0, fmt.Errorf("code not of desired length: %v", CodeLen)
	}
	var (
		id, questionnaireID int
		codeUsed            bool
	)
	row := db.QueryRow(getCodeQuery, email, code)
	err := row.Scan(&id, &codeUsed, &questionnaireID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("code or user not found")
		}
		fmt.Printf("%v: error on GetAccessCode query - %v\n", time.Now(), err)
		return 0, errors.New("problem on GetAccessCode")
	}

	if id < 1 {
		return 0, errors.New("code or user not found")
	}

	if questionnaireID < 1 {
		return 0, errors.New("unbound code")
	}

	if codeUsed {
		return 0, nil // no error if code has been used before
	}
	return questionnaireID, nil // code not used before
}

// GetAccessCodeID searches the database for an id relating to a supplied code in the AccessCode table
func GetAccessCodeID(code string, qid int, db *sql.DB) int {
	if len(code) != CodeLen {
		return 0
	}
	var id int
	row := db.QueryRow("SELECT 1 FROM AccessCode WHERE Code = $1 AND QuestionnaireID = $2", code, qid)
	err := row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0
		}
		fmt.Printf("%v: error on GetAccessCodeID query - %v\n", time.Now(), err)
	}
	if id < 1 {
		return 0
	}
	return id
}

// SendCodes updates AccessCode table to find any entries without a code
// generates codes with utilities pkg, then updates the table with the generated, unique codes
func SendCodes(qid int, db *sql.DB) (bool, error) {
	codeIDs := []int{}

	// get count of codes to create
	rows, err := db.Query(sendCodeQuery, qid)
	if err == sql.ErrNoRows {
		fmt.Printf("%v: no codes to create\n", time.Now())
		return false, errors.New("no codes to create")
	}
	if err != nil {
		fmt.Printf("%v: error on SendCodes get count - %v\n", time.Now(), err)
		return false, errors.New("problem encountered while creating codes")
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			fmt.Printf("%v: error on SendCodes read get rows - %v\n", time.Now(), err)
			return false, errors.New("problem encountered while creating codes")
		}
		codeIDs = append(codeIDs, id) // maintain slice containing codeIDs to update
	}
	if len(codeIDs) == 0 {
		fmt.Printf("%v: no codes to create\n", time.Now())
		return false, errors.New("no codes to create")
	}
	utils.Setup(-1)

	// generate codes
	codes := utils.GenerateCodes(len(codeIDs), CodeLen) // quantity = total codes;

	// create query to insert codes
	var fullQuery string
	for i := 0; i < len(codes); i++ {
		q := fmt.Sprintf(update, codes[i], codeIDs[i], qid)
		fullQuery += q
	}

	// run query
	_, err = db.Exec(fullQuery)
	if err != nil {
		fmt.Printf("%v: error on SendCodes update query - %v\n", time.Now(), err)
		return false, errors.New("problem encountered while creating codes")
	}

	return true, nil
}

// UpdateCode sets a code as used upon completing a questionnaire
func UpdateCode(id int, db *sql.DB) bool {
	_, err := db.Exec(setUsed, id)
	if err != nil {
		fmt.Printf("%v, error on UpdateCode query - %v\n", time.Now(), err)
		return false
	}
	return true
}
