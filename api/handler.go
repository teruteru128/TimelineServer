package api

import (
	"net/http"

	"github.com/TinyKitten/TimelineServer/db"
	"github.com/labstack/echo"
	"go.uber.org/zap"
	validator "gopkg.in/go-playground/validator.v9"
	mgo "gopkg.in/mgo.v2"
)

type (
	handler struct {
		db     *db.MongoInstance
		logger *zap.Logger
	}
	customValidator struct {
		validator *validator.Validate
	}
	messageResponse struct {
		Message string `json:"message"`
	}
)

func (cv *customValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

const (
	ErrParamsRequired = "parameters required"
	ErrBadFormat      = "bad format"
	ErrLoginFailed    = "invalid credential"
	ErrUnknown        = "unknown error"
	ErrNotFound       = "not found"
	ErrSuspended      = "account suspended"
	RespCreated       = "created"
	RespDeleted       = "deleted"
	RespSuspended     = "suspended"
	ErrDuplicated     = "resource duplicated"
	ErrTooLong        = "post text too long"
	RespFollowed      = "followed"
	RespUnfollowed    = "unfollowed"
)

func handleMgoError(err error) *echo.HTTPError {
	if mgo.IsDup(err) {
		return &echo.HTTPError{Code: http.StatusConflict, Message: ErrDuplicated}
	}
	switch err {
	case mgo.ErrNotFound:
		return &echo.HTTPError{Code: http.StatusNotFound, Message: ErrNotFound}
	default:
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: ErrUnknown}
	}
}
