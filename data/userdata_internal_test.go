package data

import (
	"testing"
)

func Test_comparePasswords(t *testing.T) {
	dbPass := "57d8da63dbcfd720673fd0622ac91549"
	salt := "zRvjFZ8Amq"
	correctP := "4dm1n123"
	wrongP := "helloworld"

	if !comparePasswords(dbPass, correctP, salt) {
		t.Errorf("passwords didn't match: should have\n")
	}
	if comparePasswords(dbPass, wrongP, salt) {
		t.Errorf("passwords matched: shouldn't have\n")
	}
}
