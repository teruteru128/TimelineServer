package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"github.com/TinyKitten/TimelineServer/config"
	"github.com/TinyKitten/TimelineServer/models"
	"github.com/TinyKitten/TimelineServer/token"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/stretchr/testify/assert"
)

func TestNormalGrantAccess(t *testing.T) {
	e := echo.New()
	dummyObjID := bson.NewObjectId()

	token, err := token.CreateToken(dummyObjID, false)
	if err != nil {
		t.Errorf(err.Error())
	}

	q := make(url.Values)
	q.Set("oid", dummyObjID.Hex())
	req := httptest.NewRequest(echo.GET, "/v1/suspend?"+q.Encode(), nil)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err = middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(config.MockJwtToken),
	})(th.userSuspendHandler)(c)

	if err == nil {
		t.Fatal("should reject")
	}

	q = make(url.Values)
	q.Set("oid", dummyObjID.Hex())
	q.Set("flag", "false")
	req = httptest.NewRequest(echo.GET, "/v1/offical?"+q.Encode(), nil)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	err = middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(config.MockJwtToken),
	})(th.setOfficalFlagHandler)(c)

	if err == nil {
		t.Fatal("should reject")
	}

}

func TestUserSuspendHandler(t *testing.T) {
	e := echo.New()

	u := models.NewUser("susp", "password", "susp@example.com", false)
	err := th.db.Create("users", u)
	if err != nil {
		t.Error(err)
	}

	token, err := token.CreateToken(u.ID, true)
	if err != nil {
		t.Errorf(err.Error())
	}

	q := make(url.Values)
	q.Set("oid", u.ID.Hex())
	req := httptest.NewRequest(echo.GET, "/v1/suspend?"+q.Encode(), nil)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err = middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(config.MockJwtToken),
	})(th.userSuspendHandler)(c)

	u, err = th.db.FindUserByOID(u.ID)
	if err != nil {
		t.Error(err)
	}

	if !u.Suspended {
		t.Fatal("not suspended")
	}

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "{\"message\":\"ok\"}", rec.Body.String())
	}

}

func TestSetOfficialFlagHandler(t *testing.T) {
	e := echo.New()

	u := models.NewUser("erai", "password", "erai@example.com", false)
	err := th.db.Create("users", u)
	if err != nil {
		t.Error(err)
	}

	token, err := token.CreateToken(u.ID, true)
	if err != nil {
		t.Errorf(err.Error())
	}

	q := make(url.Values)
	q.Set("oid", u.ID.Hex())
	q.Set("flag", "true")
	req := httptest.NewRequest(echo.GET, "/v1/official?"+q.Encode(), nil)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err = middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(config.MockJwtToken),
	})(th.setOfficalFlagHandler)(c)

	u, err = th.db.FindUserByOID(u.ID)
	if err != nil {
		t.Error(err)
	}

	if !u.Official {
		t.Fatal("not official")
	}

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "{\"message\":\"ok\"}", rec.Body.String())
	}

	q = make(url.Values)
	q.Set("oid", u.ID.Hex())
	q.Set("flag", "faa")
	req = httptest.NewRequest(echo.GET, "/v1/official?"+q.Encode(), nil)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	err = middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(config.MockJwtToken),
	})(th.setOfficalFlagHandler)(c)

	if err == nil {
		t.Fatal("No error")
	}
}
