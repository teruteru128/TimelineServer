package api

import (
	"net/http"

	"github.com/TinyKitten/TimelineServer/utils"

	"github.com/TinyKitten/TimelineServer/token"
	"go.uber.org/zap"

	"github.com/TinyKitten/TimelineServer/models"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	mgo "gopkg.in/mgo.v2"
)

type (
	SignupReq struct {
		ID       string `json:"id" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}
	LoginReq struct {
		ID       string `json:"id" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
)

func (h *handler) signupHandler(c echo.Context) error {
	reqUser := new(SignupReq)
	if err := c.Bind(reqUser); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrParamsRequired}
	}
	if err := c.Validate(reqUser); err != nil {
		log.Error(err.Error())
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrBadFormat}
	}
	hashed, err := utils.HashPassword(reqUser.Password)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: ErrUnknown}
	}
	u := models.NewUser(reqUser.ID, hashed, reqUser.Email, false)
	err = h.db.Create("users", u)
	if err != nil {
		return handleMgoError(err)
	}

	return c.JSON(http.StatusCreated, &messageResponse{Message: RespCreated})
}

func (h *handler) loginHandler(c echo.Context) error {
	reqUser := new(LoginReq)
	if err := c.Bind(reqUser); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrParamsRequired}
	}

	u, err := h.db.FindUser(reqUser.ID)
	if err != nil {
		if err == mgo.ErrNotFound {
			return &echo.HTTPError{Code: http.StatusUnauthorized, Message: ErrLoginFailed}
		}
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: ErrLoginFailed}
	}
	if matched := utils.CheckPasswordHash(reqUser.Password, u.Password); !matched {
		return &echo.HTTPError{Code: http.StatusUnauthorized, Message: ErrLoginFailed}
	}

	// 凍結
	if u.Suspended {
		// TODO: どこかで凍結情報をキャッシュする
		return &echo.HTTPError{Code: http.StatusForbidden, Message: ErrSuspended}
	}

	token, err := token.CreateToken(u.ID, false)
	if err != nil {
		h.logger.Error("Failed to create jwt token", zap.String("Reason", err.Error()))
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: ErrLoginFailed}
	}
	resp := models.LoginSuccessResponse{
		ID:           u.ID.Hex(),
		UserID:       u.UserID,
		CreatedDate:  u.CreatedDate,
		UpdatedDate:  u.UpdatedDate,
		SessionToken: token,
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *handler) getUserHandler(c echo.Context) error {
	id := c.Param("id")

	u, err := h.db.FindUser(id)
	if err != nil {
		return handleMgoError(err)
	}
	resp := models.UserToUserResponse(*u)
	return c.JSON(http.StatusOK, resp)
}

func (h *handler) userDeleteHandler(c echo.Context) error {
	id := c.Param("id")

	err := h.db.DeleteUser(id)
	if err != nil {
		return handleMgoError(err)
	}
	return c.JSON(http.StatusNoContent, &messageResponse{Message: RespDeleted})
}
