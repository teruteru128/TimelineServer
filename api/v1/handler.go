package v1

import (
	"net/http"

	"github.com/TinyKitten/TimelineServer/config"
	"github.com/TinyKitten/TimelineServer/db"
	"github.com/TinyKitten/TimelineServer/logger"
	"github.com/labstack/echo"
	"go.uber.org/zap"
	validator "gopkg.in/go-playground/validator.v9"
	mgo "gopkg.in/mgo.v2"
)

type (
	APIHandler struct {
		db     *db.MongoInstance
		logger *zap.Logger
	}
	messageResponse struct {
		Message string `json:"message"`
	}
	CustomValidator struct {
		validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func NewHandler() APIHandler {
	logger := logger.GetLogger()
	conf := config.GetDBConfig()
	cacheConf := config.GetCacheConfig()
	mongoIns, err := db.NewMongoInstance(conf, cacheConf)
	if err != nil {
		logger.Panic("Failed to connect database.", zap.Skip())
	}
	return APIHandler{
		db:     mongoIns,
		logger: logger,
	}

}

const (
	ErrParamsRequired    = "parameters required"
	ErrBadFormat         = "bad format"
	ErrLoginFailed       = "login failed"
	ErrUnknown           = "unknown error"
	ErrNotFound          = "not found"
	ErrSuspended         = "account suspended"
	RespCreated          = "created"
	RespDeleted          = "deleted"
	RespSuspended        = "suspended"
	ErrDuplicated        = "resource duplicated"
	ErrTooLong           = "post text too long"
	RespFollowed         = "followed"
	RespUnfollowed       = "unfollowed"
	ErrAdminOnly         = "administration area"
	ErrInvalidJwt        = "invalid jwt token"
	ErrTooLargeImage     = "uploaded image is too large"
	ErrMediaNotSupported = "uploaded media type is not supported"
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
