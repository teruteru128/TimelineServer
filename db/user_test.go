package db

import (
	"testing"

	"github.com/TinyKitten/TimelineServer/models"
	"gopkg.in/mgo.v2/bson"
)

func TestFindUserByOID(t *testing.T) {
	dummy := models.NewUser("hello2", "password", "hello2@example.com")
	err := ins.Create("users", dummy)
	if err != nil {
		t.Errorf(err.Error())
	}

	_, err = ins.FindUserByOID(dummy.ID.Hex())
	if err != nil {
		t.Errorf(err.Error())
	}

	_, err = ins.FindUserByOID(bson.NewObjectId().Hex())
	if err == nil {
		t.Errorf("not registered")
	}
}
func TestFindUser(t *testing.T) {
	dummy := models.NewUser("hello", "password", "hello@example.com")
	err := ins.Create("users", dummy)
	if err != nil {
		t.Errorf(err.Error())
	}
	exist, err := ins.FindUser("hello")
	if err != nil {
		t.Errorf(err.Error())
	}
	if exist == nil {
		t.Errorf("User not found")
	}

	_, err = ins.FindUser("ugly_betty")
	if err == nil {
		t.Errorf("not registered")
	}
}

func TestDeleteUser(t *testing.T) {
	id := "waste"
	dummy := models.NewUser(id, "password", "garbage@example.com")
	err := ins.Create("users", dummy)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = ins.DeleteUser(id)
	if err != nil {
		t.Errorf(err.Error())
	}
	_, err = ins.FindUser(id)
	if err == nil {
		t.Errorf("AHHHH!!!! Gopher has found Zombie!!!")
	}
}

func TestSuspendUser(t *testing.T) {
	id := "banned"
	ban := models.NewUser(id, "password", "banned@example.com")
	err := ins.Create("users", ban)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = ins.SuspendUser(id, true)
	if err != nil {
		t.Errorf(err.Error())
	}
	ban, err = ins.FindUser(id)
	if err != nil {
		t.Errorf(err.Error())
	}
	if ban.Suspended != true {
		t.Errorf("%s still alive!!", ban.UserID)
	}

	err = ins.SuspendUser(id, false)
	if err != nil {
		t.Errorf(err.Error())
	}
	ban, err = ins.FindUser(id)
	if err != nil {
		t.Errorf(err.Error())
	}
	if ban.Suspended != false {
		t.Errorf("Ohh... %s is dead...", ban.UserID)
	}
}
