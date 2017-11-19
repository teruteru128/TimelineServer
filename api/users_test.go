package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/TinyKitten/TimelineServer/config"
	"github.com/TinyKitten/TimelineServer/models"
	"github.com/TinyKitten/TimelineServer/token"
	"github.com/TinyKitten/TimelineServer/utils"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/stretchr/testify/assert"
	validator "gopkg.in/go-playground/validator.v9"
)

func TestSignupHandler(t *testing.T) {
	e := echo.New()
	postReq := SignupReq{
		ID:       "hoge",
		Email:    "hoge@example.com",
		Password: "hogePass",
	}
	j, err := json.Marshal(postReq)
	if err != nil {
		t.Error(err)
	}
	req := httptest.NewRequest(echo.POST, "/v1/signup/", strings.NewReader(string(j)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	e.Validator = &customValidator{validator: validator.New()}

	if assert.NoError(t, th.signupHandler(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		expectedResp := messageResponse{Message: RespCreated}
		json, err := json.Marshal(expectedResp)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, string(json), rec.Body.String())
	}
}
func TestSignupHandlerNoParams(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/v1/signup/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	e.Validator = &customValidator{validator: validator.New()}

	err := th.signupHandler(c)

	if err == nil {
		t.Errorf(ErrParamsRequired)
	}
}

func TestSignupHandlerBadFormat(t *testing.T) {
	e := echo.New()
	postReq := SignupReq{
		ID:       "hoge",
		Email:    "hogeexample.com",
		Password: "hogePass",
	}
	j, err := json.Marshal(postReq)
	if err != nil {
		t.Error(err)
	}
	req := httptest.NewRequest(echo.POST, "/v1/signup/", strings.NewReader(string(j)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	e.Validator = &customValidator{validator: validator.New()}

	err = th.signupHandler(c)

	if err == nil {
		t.Errorf(ErrBadFormat)
	}
}

func TestLoginHandler(t *testing.T) {
	e := echo.New()
	hashed, err := utils.HashPassword("password")
	if err != nil {
		t.Error(err)
	}
	u := models.NewUser("haa", hashed, "haa@example.com", false)
	err = th.db.Create("users", u)
	if err != nil {
		t.Errorf(err.Error())
	}
	postReq := LoginReq{
		ID:       "haa",
		Password: "password",
	}
	j, err := json.Marshal(postReq)
	if err != nil {
		t.Error(err)
	}
	req := httptest.NewRequest(echo.POST, "/v1/login/", strings.NewReader(string(j)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	e.Validator = &customValidator{validator: validator.New()}

	if assert.NoError(t, th.loginHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestLoginHandlerNoParams(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/v1/login/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	e.Validator = &customValidator{validator: validator.New()}

	err := th.loginHandler(c)

	if err == nil {
		t.Errorf(ErrParamsRequired)
	}
}

func TestLoginHandlerNotExist(t *testing.T) {
	e := echo.New()
	postReq := LoginReq{
		ID:       "hogehage",
		Password: "password",
	}
	j, err := json.Marshal(postReq)
	if err != nil {
		t.Error(err)
	}
	req := httptest.NewRequest(echo.POST, "/v1/login/", strings.NewReader(string(j)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	e.Validator = &customValidator{validator: validator.New()}

	err = th.loginHandler(c)
	if err == nil {
		t.Errorf(ErrLoginFailed)
	}
}

func TestLoginHandlerSuspended(t *testing.T) {
	e := echo.New()
	hashed, err := utils.HashPassword("password")
	if err != nil {
		t.Error(err)
	}
	u := models.NewUser("maguro", hashed, "maguro@example.com", false)
	u.Suspended = true
	err = th.db.Create("users", u)
	if err != nil {
		t.Errorf(err.Error())
	}
	postReq := LoginReq{
		ID:       "maguro",
		Password: "password",
	}
	j, err := json.Marshal(postReq)
	if err != nil {
		t.Error(err)
	}
	req := httptest.NewRequest(echo.POST, "/v1/login/", strings.NewReader(string(j)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	e.Validator = &customValidator{validator: validator.New()}

	err = th.loginHandler(c)
	if err == nil {
		t.Errorf(ErrSuspended)
	}
}

func TestGetUserHandler(t *testing.T) {
	e := echo.New()
	u := models.NewUser("nagura", "password", "nagura@example.com", false)
	err := th.db.Create("users", u)
	if err != nil {
		t.Errorf(err.Error())
	}

	req := httptest.NewRequest(echo.GET, "/v1/users/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(u.UserID)
	e.Validator = &customValidator{validator: validator.New()}

	if assert.NoError(t, th.getUserHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestUserDeleteHandler(t *testing.T) {
	e := echo.New()
	u := models.NewUser("nabra", "password", "nabra@example.com", false)
	err := th.db.Create("users", u)
	if err != nil {
		t.Errorf(err.Error())
	}

	req := httptest.NewRequest(echo.DELETE, "/v1/users/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	token, err := token.CreateToken(u.ID, false)
	if err != nil {
		t.Errorf(err.Error())
	}
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(u.UserID)
	e.Validator = &customValidator{validator: validator.New()}

	err = middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(config.MockJwtToken),
	})(th.userDeleteHandler)(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}

	u, err = th.db.FindUser(u.UserID)
	if err == nil {
		t.Log(u.ID.Hex())
		t.Errorf("Not deleted")
	}
}
