package data

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	utils "github.com/aagoldingay/commendeer-go/utilities"
)

const (
	codeLen       = 10
	getCodeQuery  = "SELECT CodeID, Used FROM AccessCode WHERE Email = '%v' AND Code = '%v'"
	sendCodeQuery = "SELECT CodeID FROM AccessCode WHERE Code IS NULL;"
	update        = "UPDATE AccessCode SET Code = '%v' WHERE CodeID = %v; "
)

// GetAccessCode takes an email and code, then searches the database for a relevant entry
// errors return based on incorrect code length or 'no code or user found'
// no error when code already used
func GetAccessCode(email, code string, db *sql.DB) (bool, error) {
	if len(code) != codeLen {
		return false, fmt.Errorf("code not of desired length: %v", codeLen)
	}
	var (
		id       int
		codeUsed bool
	)
	rows, err := db.Query(fmt.Sprintf(getCodeQuery, email, code))
	if err != nil {
		fmt.Printf("%v: error on GetAccessCode query - %v\n", time.Now(), err)
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&id, &codeUsed)
	}

	if id < 1 {
		return false, errors.New("code or user not found")
	}

	if codeUsed {
		return false, nil // no error if code has been used before
	}
	return true, nil // code not used before
}

// SendCodes updates AccessCode table to find any entries without a code
// generates codes with utilities pkg, then updates the table with the generated, unique codes
func SendCodes(db *sql.DB) int {
	codeIDs := []int{}

	// get count of codes to create
	rows, err := db.Query(sendCodeQuery)
	if err == sql.ErrNoRows {
		fmt.Printf("%v: no codes to create\n", time.Now())
		return 0
	}
	if err != nil {
		fmt.Printf("%v: error on SendCodes get count - %v\n", time.Now(), err)
		return 0
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			fmt.Printf("%v: error on SendCodes read get rows - %v\n", time.Now(), err)
		}
		codeIDs = append(codeIDs, id) // maintain slice containing codeIDs to update
	}
	if len(codeIDs) == 0 {
		fmt.Printf("%v: no codes to create\n", time.Now())
		return 0
	}
	utils.Setup(-1)

	// generate codes
	codes := utils.GenerateCodes(len(codeIDs), codeLen) // quantity = total codes;

	// create query to insert codes
	var fullQuery string
	for i := 0; i < len(codes); i++ {
		q := fmt.Sprintf(update, codes[i], codeIDs[i])
		fullQuery += q
	}

	// run query
	res, err := db.Exec(fullQuery)
	if err != nil {
		fmt.Printf("%v: error on SendCodes update query - %v\n", time.Now(), err)
		return 0
	}

	count, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("%v: error on SendCodes rows affected - %v\n", time.Now(), err)
		return int(count)
	}

	return int(count)
}
