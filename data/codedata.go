package data

import (
	"database/sql"
	"fmt"
	"time"

	utils "github.com/aagoldingay/commendeer-go/utilities"
)

const (
	codeLen = 10
	query   = "SELECT CodeID FROM AccessCode WHERE Code IS NULL;"
	update  = "UPDATE AccessCode SET Code = '%v' WHERE CodeID = %v; "
)

// SendCodes updates AccessCode table to find any entries without a code
// generates codes with utilities pkg, then updates the table with the generated, unique codes
func SendCodes(db *sql.DB) int {
	codeIDs := []int{}

	// get count of codes to create
	rows, err := db.Query(query)
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
