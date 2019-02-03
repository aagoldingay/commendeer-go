package data_test

import (
	"testing"

	"github.com/aagoldingay/commendeer-go/data"
	_ "github.com/lib/pq"
)

// ATTENTION: tests are called and run in data_test.go

func GetUserInfo_Success(t *testing.T) {
	u, p := "admin1", "4dm1n123"
	d := data.GetUserInfo(u, p, db)
	if d.ID != 1 {
		t.Errorf("user %v not correct : expected (%v,%v), actual (%v,%v)\n", u, u, 1, d.Username, d.ID)
	}
}

func GetUserInfo_UnknownUser(t *testing.T) {
	u, p := "admin5", "4dm1n123"
	d := data.GetUserInfo(u, p, db)
	if d.ID != 0 {
		t.Errorf("user returned successfully : expected (%v,%v), actual (%v,%v)\n", u, 0, d.Username, d.ID)
	}
}

func GetUserInfo_IncorrectPassword(t *testing.T) {
	u, p := "admin1", "randompassword"
	d := data.GetUserInfo(u, p, db)
	if d.ID != 0 {
		t.Errorf("user returned successfully : expected (%v,%v), actual (%v,%v)\n", u, 0, d.Username, d.ID)
	}
}
