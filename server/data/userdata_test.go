package data_test

import (
	"testing"

	"github.com/aagoldingay/commendeer-go/server/data"
	_ "github.com/lib/pq"
)

// ATTENTION: tests are called and run in data_test.go

func Login_Success(t *testing.T) {
	u, p := "admin1", "4dm1n123"
	d, err := data.Login(u, p, db)
	if err != nil {
		t.Errorf("Login_Success errored - should not have : %v", err)
	}
	if d.Username != u {
		t.Errorf("username not correct : expected %v, actual %v\n", u, d.Username)
	}
	if len(d.Code) != 20 {
		t.Errorf("code length not correct : %v", len(d.Code))
	}
}

func Login_UnknownUser(t *testing.T) {
	u, p := "admin5", "4dm1n123"
	_, err := data.Login(u, p, db)
	if err.Error() != "incorrect username or password" {
		t.Errorf("incorrect error returned : %v\n", err)
	}
}

func Login_IncorrectPassword(t *testing.T) {
	u, p := "admin1", "randompassword"
	_, err := data.Login(u, p, db)
	if err.Error() != "incorrect username or password" {
		t.Errorf("incorrect error returned : %v\n", err)
	}
}

func Logout_NoError(t *testing.T) {
	u, p := "admin1", "4dm1n123"
	d, err := data.Login(u, p, db)
	if err != nil {
		t.Errorf("login errored - should not have : %v", err)
	}

	err = data.Logout(d.Code, db)
	if err != nil {
		t.Errorf("Logout_NoError errored - should not have : %v", err)
	}
}

func Logout_IncorrectCode(t *testing.T) {
	err := data.Logout("randomfakecodehere12", db)
	if err.Error() != "unknown code" {
		t.Errorf("Logout_IncorrectCode incorrect error: %v", err)
	}
}

func Logout_InvalidCode(t *testing.T) {
	err := data.Logout("code", db)
	if err.Error() != "invalid code" {
		t.Errorf("Logout_InvalidCode incorrect error: %v", err)
	}
}
