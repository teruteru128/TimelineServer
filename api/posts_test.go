package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/TinyKitten/TimelineServer/config"
	"github.com/TinyKitten/TimelineServer/db"
	"github.com/TinyKitten/TimelineServer/logger"
	"github.com/TinyKitten/TimelineServer/models"
	"github.com/TinyKitten/TimelineServer/token"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	validator "gopkg.in/go-playground/validator.v9"
	dockertest "gopkg.in/ory-am/dockertest.v3"
)

var th *handler

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("mongo", "latest", []string{})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err := pool.Retry(func() error {
		conf := config.DBConfig{
			Server:   fmt.Sprintf("localhost:%s", resource.GetPort("27017/tcp")),
			Database: "testing",
		}
		ins := &db.MongoInstance{Conf: conf}
		logger := logger.GetLogger()

		th = &handler{
			db:     ins,
			logger: logger,
		}

		return ins.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	code := m.Run()
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	os.Exit(code)
}
func TestGetPublicPostsHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/v1/posts/public/", nil)
	u := models.NewUser("id", "password", "mail@example.com")
	err := th.db.Create("users", u)
	if err != nil {
		t.Errorf(err.Error())
	}
	token, err := token.CreateToken(u.ID)
	if err != nil {
		t.Errorf(err.Error())
	}
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	exec := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(config.MockJwtToken),
	})(th.getPublicPostsHandler)(c)

	if assert.NoError(t, exec) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "null", rec.Body.String())
	}
}

func TestSuspendedGetPublicPostsHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/v1/posts/public/", nil)
	u := models.NewUser("id2", "password", "mail2@example.com")
	u.Suspended = true
	err := th.db.Create("users", u)
	if err != nil {
		t.Errorf(err.Error())
	}
	token, err := token.CreateToken(u.ID)
	if err != nil {
		t.Errorf(err.Error())
	}
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err = middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(config.MockJwtToken),
	})(th.getPublicPostsHandler)(c)

	if err.Error() != "code=403, message=account suspended" {
		t.Errorf("Error code not matched: %s", err.Error())
	}
}

func TestPostHandler(t *testing.T) {
	e := echo.New()
	postReq := `{"text": "post message"}`
	req := httptest.NewRequest(echo.POST, "/v1/posts/", strings.NewReader(postReq))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	u := models.NewUser("id3", "password", "mail3@example.com")
	err := th.db.Create("users", u)
	if err != nil {
		t.Errorf(err.Error())
	}
	token, err := token.CreateToken(u.ID)
	if err != nil {
		t.Errorf(err.Error())
	}
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	e.Validator = &customValidator{validator: validator.New()}
	err = middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(config.MockJwtToken),
	})(th.postHandler)(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestSuspendedPostHandler(t *testing.T) {
	e := echo.New()
	postReq := `{"text": "post message"}`
	req := httptest.NewRequest(echo.POST, "/v1/posts/", strings.NewReader(postReq))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	u := models.NewUser("id4", "password", "mail4@example.com")
	u.Suspended = true
	err := th.db.Create("users", u)
	if err != nil {
		t.Errorf(err.Error())
	}
	token, err := token.CreateToken(u.ID)
	if err != nil {
		t.Errorf(err.Error())
	}
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	e.Validator = &customValidator{validator: validator.New()}
	err = middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(config.MockJwtToken),
	})(th.postHandler)(c)

	if err.Error() != "code=403, message=account suspended" {
		t.Errorf("Error code not matched: %s", err.Error())
	}
}
func TestEmptyPostHandler(t *testing.T) {
	e := echo.New()
	postReq := `{}`
	req := httptest.NewRequest(echo.POST, "/v1/posts/", strings.NewReader(postReq))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	u := models.NewUser("id5", "password", "mail5@example.com")
	err := th.db.Create("users", u)
	if err != nil {
		t.Errorf(err.Error())
	}
	token, err := token.CreateToken(u.ID)
	if err != nil {
		t.Errorf(err.Error())
	}
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	e.Validator = &customValidator{validator: validator.New()}
	err = middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(config.MockJwtToken),
	})(th.postHandler)(c)

	if err.Error() != "code=400, message=bad format" {
		t.Errorf("Error code not matched: %s", err.Error())
	}
}

func TestBindPostHandler(t *testing.T) {
	e := echo.New()
	postReq := `{"text": "post message"}`
	req := httptest.NewRequest(echo.POST, "/v1/posts/", strings.NewReader(postReq))
	u := models.NewUser("id6", "password", "mail6@example.com")
	err := th.db.Create("users", u)
	if err != nil {
		t.Errorf(err.Error())
	}
	token, err := token.CreateToken(u.ID)
	if err != nil {
		t.Errorf(err.Error())
	}
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	e.Validator = &customValidator{validator: validator.New()}
	err = middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(config.MockJwtToken),
	})(th.postHandler)(c)

	if err.Error() != "code=400, message=parameters required" {
		t.Errorf("Error code not matched: %s", err.Error())
	}
}
