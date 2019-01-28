package main

import "testing"

func Test_dbSetup(t *testing.T) {
	_, err := dbSetup()
	if err != nil {
		t.Error(err)
	}
}
