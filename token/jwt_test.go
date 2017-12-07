package token

import (
	"testing"

	"gopkg.in/mgo.v2/bson"
)

func TestCreateToken(t *testing.T) {
	token, err := CreateToken(bson.NewObjectId(), false)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if token == "" {
		t.Fatalf("token is empty")
	}
}
