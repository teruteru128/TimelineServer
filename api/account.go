package api

import (
	"net/http"
	"strconv"

	"gopkg.in/mgo.v2/bson"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	"github.com/TinyKitten/TimelineServer/config"
	"github.com/TinyKitten/TimelineServer/utils"

	"github.com/TinyKitten/TimelineServer/token"
	"go.uber.org/zap"

	"github.com/TinyKitten/TimelineServer/models"
	"github.com/labstack/echo"
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
	AccountSettingsResponse struct {
		DisplayName string `json:"screen_name"`
	}
	AccountSettingsRequest struct {
		Name        string `json:"name"`
		URL         string `json:"url"`
		Location    string `json:"location"`
		Description string `json:"description"`
	}
	AccountImageRequest struct {
		Image string `json:"image"`
	}
)

func (h *handler) signupHandler(c echo.Context) error {
	reqUser := new(SignupReq)
	if err := c.Bind(reqUser); err != nil {
		h.logger.Debug("API Error", zap.String("Error", err.Error()))
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrParamsRequired}
	}
	if err := c.Validate(reqUser); err != nil {
		h.logger.Debug("API Error", zap.String("Error", err.Error()))
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrBadFormat}
	}
	hashed, err := utils.HashPassword(reqUser.Password)
	if err != nil {
		h.logger.Debug("API Error", zap.String("Error", err.Error()))
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: ErrUnknown}
	}
	u := models.NewUser(reqUser.ID, hashed, reqUser.Email, false)
	err = h.db.Insert("users", u)
	if err != nil {
		h.logger.Debug("API Error", zap.String("Error", err.Error()))
		return handleMgoError(err)
	}

	token, err := token.CreateToken(u.ID, false)
	if err != nil {
		h.logger.Debug("API Error", zap.String("Error", err.Error()))
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: ErrLoginFailed}
	}
	resp := models.UserToLoginSucessResponse(*u, token)

	return c.JSON(http.StatusCreated, resp)
}

func (h *handler) loginHandler(c echo.Context) error {
	reqUser := new(LoginReq)
	if err := c.Bind(reqUser); err != nil {
		h.logger.Debug("API Error", zap.String("Error", err.Error()))
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrParamsRequired}
	}

	u, err := h.db.FindUser(reqUser.ID)
	if err != nil {
		h.logger.Debug("API Error", zap.String("Error", err.Error()))
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
	resp := models.UserToLoginSucessResponse(*u, token)
	return c.JSON(http.StatusOK, resp)
}

func (h *handler) getUserHandler(c echo.Context) error {
	// Jwtチェック
	config := config.GetAPIConfig()
	tokenStr := c.QueryParam("token")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Jwt), nil
	})
	if err != nil {
		h.logger.Debug("API Error", zap.String("Error", err.Error()))
		return &echo.HTTPError{Code: http.StatusForbidden, Message: ErrInvalidJwt}
	}
	if !token.Valid {
		h.logger.Debug("API Error", zap.String("Error", err.Error()))
		return &echo.HTTPError{Code: http.StatusForbidden, Message: ErrInvalidJwt}
	}

	screenName := c.QueryParam("screen_name")
	userId := c.QueryParam("user_id")

	if screenName != "" {
		user, err := h.db.FindUser(screenName)
		if err != nil {
			return handleMgoError(err)
		}
		resp := models.UserToUserResponse(*user)
		return c.JSON(http.StatusOK, resp)
	}

	if userId != "" {
		user, err := h.db.FindUserByOID(bson.ObjectId(userId))
		if err != nil {
			return handleMgoError(err)
		}
		resp := models.UserToUserResponse(*user)
		return c.JSON(http.StatusOK, resp)
	}

	h.logger.Debug("API Error", zap.String("Error", ErrParamsRequired))
	return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrParamsRequired}
}

func (h *handler) userDeleteHandler(c echo.Context) error {
	id := c.Param("id")

	err := h.db.DeleteUser(id)
	if err != nil {
		h.logger.Debug("API Error", zap.String("Error", err.Error()))
		return handleMgoError(err)
	}
	return c.JSON(http.StatusNoContent, &messageResponse{Message: RespDeleted})
}

func (h *handler) getSettingsHandler(c echo.Context) error {
	// Jwtチェック
	config := config.GetAPIConfig()
	tokenStr := c.QueryParam("token")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Jwt), nil
	})
	if err != nil {
		h.logger.Debug("API Error", zap.String("Error", err.Error()))
		return &echo.HTTPError{Code: http.StatusForbidden, Message: ErrInvalidJwt}
	}
	if !token.Valid {
		h.logger.Debug("API Error", zap.String("Error", err.Error()))
		return &echo.HTTPError{Code: http.StatusForbidden, Message: ErrInvalidJwt}
	}
	claims := token.Claims.(jwt.MapClaims)
	id := claims["id"].(string)

	user, err := h.db.FindUserByOID(bson.ObjectIdHex(id))
	if err != nil {
		h.logger.Debug("API Error", zap.String("Error", err.Error()))
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: ErrUnknown}
	}

	return c.JSON(http.StatusOK, &AccountSettingsResponse{
		DisplayName: user.DisplayName,
	})
}

func (h *handler) setSettingsHandler(c echo.Context) error {
	jwtUser := c.Get("user").(*jwt.Token)
	claims := jwtUser.Claims.(jwt.MapClaims)
	idStr := claims["id"].(string)
	id := bson.ObjectId(bson.ObjectIdHex(idStr))

	req := new(AccountSettingsRequest)
	if err := c.Bind(req); err != nil {
		h.logger.Debug("API Error", zap.String("Error", err.Error()))
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrParamsRequired}
	}

	if req.Name != "" {
		err := h.db.UpdateUser(id, "displayName", req.Name)
		if err != nil {
			h.logger.Debug("API Error", zap.String("Error", err.Error()))
			return &echo.HTTPError{Code: http.StatusInternalServerError, Message: ErrUnknown}
		}
	}
	if req.URL != "" {
		err := h.db.UpdateUser(id, "url", req.URL)
		if err != nil {
			h.logger.Debug("API Error", zap.String("Error", err.Error()))
			return &echo.HTTPError{Code: http.StatusInternalServerError, Message: ErrUnknown}
		}
	}
	if req.Location != "" {
		err := h.db.UpdateUser(id, "location", req.Location)
		if err != nil {
			h.logger.Debug("API Error", zap.String("Error", err.Error()))
			return &echo.HTTPError{Code: http.StatusInternalServerError, Message: ErrUnknown}
		}
	}
	if req.Description != "" {
		err := h.db.UpdateUser(id, "description", req.Description)
		if err != nil {
			h.logger.Debug("API Error", zap.String("Error", err.Error()))
			return &echo.HTTPError{Code: http.StatusInternalServerError, Message: ErrUnknown}
		}
	}

	if req.Name == "" && req.URL == "" && req.Location == "" && req.Description == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrParamsRequired}
	}

	newUser, err := h.db.FindUserByOID(id)
	if err != nil {
		h.logger.Debug("API Error", zap.String("Error", err.Error()))
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: ErrUnknown}
	}

	resp := models.UserToUserResponse(*newUser)

	return c.JSON(http.StatusOK, resp)
}

func (h *handler) updateProfileImageHandler(c echo.Context) error {
	jwtUser := c.Get("user").(*jwt.Token)
	claims := jwtUser.Claims.(jwt.MapClaims)
	idStr := claims["id"].(string)
	id := bson.ObjectId(bson.ObjectIdHex(idStr))

	req := new(AccountImageRequest)
	if err := c.Bind(req); err != nil {
		h.logger.Debug("API Error", zap.String("Error", err.Error()))
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrParamsRequired}
	}

	ext := utils.DetectFileExtension(req.Image)
	if ext == "" {
		return &echo.HTTPError{Code: http.StatusUnsupportedMediaType, Message: ErrMediaNotSupported}
	}

	dat, err := utils.DecodeImage(req.Image)
	if err != nil {
		h.logger.Debug("API Error", zap.String("Error", err.Error()))
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: ErrUnknown}
	}

	path := config.GetUploadImagePath()
	uuid, err := uuid.NewUUID()
	if err != nil {
		h.logger.Debug("API Error", zap.String("Error", err.Error()))
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: ErrUnknown}
	}

	filePath := path + uuid.String() + ext
	err = utils.SaveFile(dat, filePath, 1)
	if err != nil && err != utils.ErrFileHuge {
		h.logger.Debug("API Error", zap.String("Error", err.Error()))
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: ErrUnknown}
	}
	if err == utils.ErrFileHuge {
		h.logger.Debug("API Error", zap.String("Error", err.Error()))
		return &echo.HTTPError{Code: http.StatusRequestEntityTooLarge, Message: ErrTooLargeImage}
	}

	cfg := config.GetAPIConfig()
	portStr := strconv.Itoa(cfg.Port)
	avatarUrl := "https://" + cfg.Endpoint + ":" + portStr + "/" + cfg.Version + "/" + filePath
	err = h.db.UpdateUser(id, "avatarUrl", avatarUrl)
	if err != nil {
		h.logger.Debug("API Error", zap.String("Error", err.Error()))
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: ErrUnknown}
	}

	u, err := h.db.FindUserByOID(id)
	if err != nil {
		h.logger.Debug("API Error", zap.String("Error", err.Error()))
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: ErrUnknown}
	}

	return c.JSON(http.StatusCreated, u)
}
