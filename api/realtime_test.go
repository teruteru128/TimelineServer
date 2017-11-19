package api

import (
	"encoding/json"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"github.com/TinyKitten/TimelineServer/models"
)

func TestCheckFollow_Followers(t *testing.T) {
	follower := models.NewUser("follower1", "password", "follower1@example.com", false)
	err := th.db.Create("users", follower)
	if err != nil {
		t.Errorf(err.Error())
	}
	followee := models.NewUser("followee1", "password", "followee1@example.com", false)
	err = th.db.Create("users", followee)
	if err != nil {
		t.Errorf(err.Error())
	}
	u := models.NewUser("hagehoge", "password", "hagehoge@example.com", false)
	u.Followers = append(u.Followers, follower.ID)
	u.Following = append(u.Following, followee.ID)
	err = th.db.Create("users", u)
	if err != nil {
		t.Errorf(err.Error())
	}

	post := models.NewPost(u.UserID, "hoge")

	j, followers := th.checkFollow(u.ID.Hex(), *post)
	if j == nil {
		t.Error("not matched")
	}
	if followers == nil {
		t.Error("not matched")
	}
	if followers[0] != follower.ID {
		t.Error("not matched")
	}
}
func TestCheckFollow_Own(t *testing.T) {
	follow := models.NewUser("Luigi", "password", "Luigi@example.com", false)
	err := th.db.Create("users", follow)
	if err != nil {
		t.Errorf(err.Error())
	}
	u := models.NewUser("Mario", "password", "Mario@example.com", false)
	u.Following = append(u.Following, follow.ID)
	err = th.db.Create("users", u)
	if err != nil {
		t.Errorf(err.Error())
	}

	if err != nil {
		t.Errorf(err.Error())
	}

	post := models.NewPost(u.UserID, "hoge")

	j, followers := th.checkFollow(u.ID.Hex(), *post)
	if j == nil {
		t.Error("not matched")
	}
	if followers != nil {
		t.Error("not matched")
	}
}

func TestCheckFollow_Followings(t *testing.T) {
	warioID := bson.NewObjectId()
	waluigiID := bson.NewObjectId()
	wario := models.NewUser("Wario", "password", "Wario@example.com", false)
	wario.ID = warioID
	wario.Followers = append(wario.Followers, waluigiID)
	err := th.db.Create("users", wario)
	if err != nil {
		t.Errorf(err.Error())
	}
	u := models.NewUser("Waluigi", "password", "waluigi@example.com", false)
	u.ID = waluigiID
	u.Following = append(u.Following, warioID)
	err = th.db.Create("users", u)
	if err != nil {
		t.Errorf(err.Error())
	}

	post := models.NewPost(wario.UserID, "Hi waluigi")

	j, followers := th.checkFollow(u.ID.Hex(), *post)
	if j == nil {
		t.Error("not matched")
	}
	var uresp StreamPostResp
	err = json.Unmarshal(*j, &uresp)
	if err != nil {
		t.Error(err)
	}
	for f := range followers {
		if string(f) != u.ID.Hex() {
			t.Log(f)
			t.Errorf("not matched")
		}
	}
	if followers != nil {
		t.Error("not matched")
	}
}
