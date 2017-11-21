package models

import (
	"testing"
)

func TestNewPost(t *testing.T) {
	uid := "hage"
	text := "kono hage"

	genPost := NewPost(uid, text)

	if genPost.UserID != uid {
		t.Fatalf("UserID not matched: %s", genPost.UserID)
	}

	if genPost.Text != text {
		t.Fatalf("Text not matched: %s", genPost.Text)
	}
}
