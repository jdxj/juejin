package db

import "testing"

func TestMySQL(t *testing.T) {
	err := MySQL.Ping()
	if err != nil {
		t.Fatalf("%s", err)
	}
}
