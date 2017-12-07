package utils

import "testing"

func TestHashPassword(t *testing.T) {
	password := "password"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if matched := CheckPasswordHash(password, hash); !matched {
		t.Fatalf("Passwords not match")
	}
}
