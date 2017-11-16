package db

import (
	"testing"

	"github.com/TinyKitten/TimelineServer/models"
	"gopkg.in/mgo.v2/bson"
)

func TestFindUserByOID(t *testing.T) {
	dummy := models.NewUser("hello2", "password", "hello2@example.com", false)
	err := ins.Create("users", dummy)
	if err != nil {
		t.Errorf(err.Error())
	}

	_, err = ins.FindUserByOID(dummy.ID)
	if err != nil {
		t.Errorf(err.Error())
	}

	_, err = ins.FindUserByOID(bson.NewObjectId())
	if err == nil {
		t.Errorf("not registered")
	}
}
func TestFindUser(t *testing.T) {
	dummy := models.NewUser("hello", "password", "hello@example.com", false)
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
	dummy := models.NewUser(id, "password", "garbage@example.com", false)
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

func TestFindUserByOIDArray(t *testing.T) {
	arr1 := models.NewUser("arr1", "password", "arr1@example.com", false)
	arr2 := models.NewUser("arr2", "password", "arr2@example.com", false)
	arr3 := models.NewUser("arr3", "password", "arr3@example.com", false)
	arr4 := models.NewUser("arr4", "password", "arr4@example.com", false)
	arr5 := models.NewUser("arr5", "password", "arr5@example.com", false)
	arr := []models.User{*arr1, *arr2, *arr3, *arr4, *arr5}
	err := createUserFromArray("users", arr)
	if err != nil {
		t.Error(err)
	}
	oids := []bson.ObjectId{arr1.ID, arr2.ID, arr3.ID, arr4.ID, arr5.ID}
	dbArr, err := ins.FindUserByOIDArray(oids)
	if err != nil {
		t.Error(err)
	}

	if len(arr) != len(dbArr) {
		t.Error("length not matched")
	}
}

func TestSuspendUser(t *testing.T) {
	id := "banned"
	ban := models.NewUser(id, "password", "banned@example.com", false)
	err := ins.Create("users", ban)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = ins.SuspendUser(ban.ID, true)
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

	err = ins.SuspendUser(ban.ID, false)
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

func TestFollowUser(t *testing.T) {
	followID := "follow1"
	followerID := "follower1"
	follow := models.NewUser(followID, "password", "follow@example.com", false)
	follower := models.NewUser(followerID, "password", "follower@example.com", false)
	err := ins.Create("users", follow)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = ins.Create("users", follower)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = ins.FollowUser(follow.ID, follower.ID)
	if err != nil {
		t.Errorf(err.Error())
	}

	follow, err = ins.FindUserByOID(follow.ID)
	if err != nil {
		t.Errorf(err.Error())
	}
	follower, err = ins.FindUserByOID(follower.ID)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(follow.Following) == 0 {
		t.Fatal("not followed")
	}
	if len(follower.Followers) == 0 {
		t.Fatal("not followed")
	}
}

func TestUnfollowUser(t *testing.T) {
	followID := "follow2"
	followerID := "follower2"
	follow := models.NewUser(followID, "password", "follow2@example.com", false)
	follower := models.NewUser(followerID, "password", "follower2@example.com", false)
	err := ins.Create("users", follow)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = ins.Create("users", follower)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = ins.UnfollowUser(follow.ID, follower.ID)
	if err != nil {
		t.Errorf(err.Error())
	}

	follow, err = ins.FindUserByOID(follow.ID)
	if err != nil {
		t.Errorf(err.Error())
	}
	follower, err = ins.FindUserByOID(follower.ID)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(follow.Following) != 0 {
		t.Fatal("not unfollowed")
	}
	if len(follower.Followers) != 0 {
		t.Fatal("not unfollowed")
	}
}

func TestSetOfficial(t *testing.T) {
	u := models.NewUser("tabunerai", "password", "tabunerai@example.com", false)
	ins.Create("users", u)
	ins.SetOfficial(u.ID, true)
	u, err := ins.FindUserByOID(u.ID)
	if err != nil {
		t.Error(err)
	}
	if !u.Official {
		t.Fatal("failed")
	}

	ins.SetOfficial(u.ID, false)
	u, err = ins.FindUserByOID(u.ID)
	if err != nil {
		t.Error(err)
	}
	if u.Official {
		t.Fatal("failed")
	}
}

func createUserFromArray(key string, arr []models.User) error {
	for _, item := range arr {
		err := ins.Create(key, item)
		if err != nil {
			return err
		}
	}
	return nil
}
