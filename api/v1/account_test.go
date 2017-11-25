package v1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/TinyKitten/TimelineServer/models"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	validator "gopkg.in/go-playground/validator.v9"
)

func TestAccountCreate(t *testing.T) {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	reqParams := SignupReq{
		ID:       "PikkaPikka1Nensei",
		Email:    "elemschoo@example.com",
		Password: "password",
	}
	j, err := json.Marshal(reqParams)
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(echo.POST, "/1.0/account/create.json", strings.NewReader(string(j)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, th.AccountCreate(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		resp := models.LoginSuccessResponse{}
		respb := []byte(rec.Body.String())
		err := json.Unmarshal(respb, &resp)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, reqParams.ID, resp.UserID)
	}

}
