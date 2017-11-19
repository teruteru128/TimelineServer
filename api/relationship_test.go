package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/TinyKitten/TimelineServer/config"
	"github.com/TinyKitten/TimelineServer/models"
	"github.com/TinyKitten/TimelineServer/token"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/stretchr/testify/assert"
)

func TestFollowingListHandler(t *testing.T) {
	e := echo.New()

	u := models.NewUser("yaju", "password", "yaju@example.com", false)
	following1 := models.NewUser("tnok", "password", "tnok@example.com", false)
	follower1 := models.NewUser("kbtit", "password", "kbtit@example.com", false)
	u.Followers = append(u.Followers, follower1.ID)
	u.Following = append(u.Following, following1.ID)

	err := th.db.Create("users", u)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = th.db.Create("users", follower1)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = th.db.Create("users", following1)
	if err != nil {
		t.Errorf(err.Error())
	}

	token, err := token.CreateToken(u.ID, false)
	if err != nil {
		t.Errorf(err.Error())
	}

	q := make(url.Values)
	q.Set("token", token)
	req := httptest.NewRequest(echo.GET, "/v1/following/:id?"+q.Encode(), nil)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("yaju")

	if assert.NoError(t, th.followingListHandler(c)) {
		assert.NotEqual(t, "{\"users\":null}", rec.Body.String())

		following1Resp := models.UserToUserResponse(*following1)

		expect := &usersResponse{
			Users: []models.UserResponse{
				following1Resp,
			},
		}
		var actual usersResponse
		json.Unmarshal([]byte(rec.Body.String()), &actual)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expect.Users[0].DisplayName, actual.Users[0].DisplayName)
	}
}

func TestFollowerListHandler(t *testing.T) {
	e := echo.New()

	u := models.NewUser("yjsnpi2", "password", "yjsnpi2@example.com", false)
	following1 := models.NewUser("mur2", "password", "mur2@example.com", false)
	follower1 := models.NewUser("imp2", "password", "imp2@example.com", false)
	u.Followers = append(u.Followers, follower1.ID)
	u.Following = append(u.Following, following1.ID)

	err := th.db.Create("users", u)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = th.db.Create("users", follower1)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = th.db.Create("users", following1)
	if err != nil {
		t.Errorf(err.Error())
	}

	token, err := token.CreateToken(u.ID, false)
	if err != nil {
		t.Errorf(err.Error())
	}

	q := make(url.Values)
	q.Set("token", token)

	req := httptest.NewRequest(echo.GET, "/v1/following?"+q.Encode(), nil)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues("yjsnpi2")

	if err != nil {
		t.Error(err)
	}

	if assert.NoError(t, th.followerListHandler(c)) {
		assert.NotEqual(t, "{\"users\":null}", rec.Body.String())

		follower1Resp := models.UserToUserResponse(*follower1)

		expect := &usersResponse{
			Users: []models.UserResponse{
				follower1Resp,
			},
		}
		var actual usersResponse
		json.Unmarshal([]byte(rec.Body.String()), &actual)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expect.Users[0].DisplayName, actual.Users[0].DisplayName)
	}
}

func TestFollowingEmptyListHandler(t *testing.T) {
	e := echo.New()

	u := models.NewUser("yaju2", "password", "yaju2@example.com", false)

	err := th.db.Create("users", u)
	if err != nil {
		t.Errorf(err.Error())
	}

	token, err := token.CreateToken(u.ID, false)
	if err != nil {
		t.Errorf(err.Error())
	}

	q := make(url.Values)
	q.Set("token", token)

	req := httptest.NewRequest(echo.GET, "/v1/following/:id?"+q.Encode(), nil)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues("yaju2")

	if err != nil {
		t.Error(err)
	}

	if assert.NoError(t, th.followingListHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "{\"users\":null}", rec.Body.String())
	}
}

func TestFollowerEmptyListHandler(t *testing.T) {
	e := echo.New()

	u := models.NewUser("yjsnpi", "password", "yjsnpi@example.com", false)

	err := th.db.Create("users", u)
	if err != nil {
		t.Errorf(err.Error())
	}

	token, err := token.CreateToken(u.ID, false)
	if err != nil {
		t.Errorf(err.Error())
	}
	q := make(url.Values)
	q.Set("token", token)
	req := httptest.NewRequest(echo.GET, "/v1/follower/:id?"+q.Encode(), nil)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues("yjsnpi")

	if err != nil {
		t.Error(err)
	}

	if assert.NoError(t, th.followerListHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "{\"users\":null}", rec.Body.String())
	}
}

func TestFollowHandler(t *testing.T) {
	e := echo.New()

	u := models.NewUser("fromkitten", "password", "fromkitten@example.com", false)
	followUser := models.NewUser("tokotten", "password", "tokotten@example.com", false)

	err := th.db.Create("users", u)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = th.db.Create("users", followUser)
	if err != nil {
		t.Errorf(err.Error())
	}

	req := httptest.NewRequest(echo.PUT, "/v1/follow/:id", nil)

	token, err := token.CreateToken(u.ID, false)
	if err != nil {
		t.Errorf(err.Error())
	}
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues("tokotten")

	err = middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(config.MockJwtToken),
	})(th.followHandler)(c)

	if err != nil {
		t.Error(err)
	}

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "{\"message\":\"followed\"}", rec.Body.String())

		u, err = th.db.FindUserByOID(u.ID)
		followUser, err = th.db.FindUserByOID(followUser.ID)
		if err != nil {
			t.Errorf(err.Error())
		}
		if len(u.Following) == 0 {
			t.Fatal("not followed")
		} else {
			if u.Following[0] != followUser.ID {
				t.Fatalf("followed user not matched: %s %s", u.ID.Hex(), followUser.ID.Hex())
			}
		}
		if len(followUser.Followers) == 0 {
			t.Fatal("not followed: %d")
		} else {
			if followUser.Followers[0] != u.ID {
				t.Fatalf("followed user not matched: %s", u.ID.Hex())
			}
		}
	}
}

func TestUnfollowHandler(t *testing.T) {
	e := echo.New()

	u := models.NewUser("fromkitten2", "password", "fromkitten2@example.com", false)
	unfollowUser := models.NewUser("tokotten2", "password", "tokotten2@example.com", false)

	err := th.db.Create("users", u)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = th.db.Create("users", unfollowUser)
	if err != nil {
		t.Errorf(err.Error())
	}

	req := httptest.NewRequest(echo.PUT, "/v1/unfollow/:id", nil)

	token, err := token.CreateToken(u.ID, false)
	if err != nil {
		t.Errorf(err.Error())
	}
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues("tokotten2")

	err = middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(config.MockJwtToken),
	})(th.unfollowHandler)(c)

	if err != nil {
		t.Error(err)
	}

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "{\"message\":\"unfollowed\"}", rec.Body.String())

		u, err = th.db.FindUserByOID(u.ID)
		unfollowUser, err = th.db.FindUserByOID(unfollowUser.ID)
		if err != nil {
			t.Errorf(err.Error())
		}
		if len(u.Following) != 0 {
			t.Fatal("not unfollowed")
		}
		if len(unfollowUser.Followers) != 0 {
			t.Fatal("not unfollowed")
		}

	}
}
