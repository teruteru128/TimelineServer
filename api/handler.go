package api

import (
	"github.com/TinyKitten/Timeline/db"
	"go.uber.org/zap"
	validator "gopkg.in/go-playground/validator.v9"
)

type (
	handler struct {
		db     *db.MongoInstance
		logger *zap.Logger
	}
	customValidator struct {
		validator *validator.Validate
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
	RespCreated       = "created"
)
