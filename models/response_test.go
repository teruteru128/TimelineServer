package models

import "testing"

func TestUserToUserResponse(t *testing.T) {
	u := NewUser("hage", "pass", "mail@example.com", true)
	resp := UserToUserResponse(*u)
	if u.ID.Hex() != resp.ID {
		t.Fatalf("id not matched: %s", resp.ID)
	}
}

func TestUsersToUserResponseArray(t *testing.T) {
	u0 := NewUser("hage", "pass", "mail@example.com", true)
	u1 := NewUser("hage", "pass", "mail@example.com", true)
	u2 := NewUser("hage", "pass", "mail@example.com", true)
	uArr := []User{*u0, *u1, *u2}
	respArr := UsersToUserResponseArray(uArr)
	if len(respArr) != 3 {
		t.Fatalf("Response array length not matched: %d", len(respArr))
	}
}
