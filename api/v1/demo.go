package v1

import (
	"encoding/json"
	"time"

	"github.com/TinyKitten/TimelineServer/utils"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

func (h *APIHandler) SampleStreamHandler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			post := utils.GenerateSamplePostResponse()
			j, err := json.Marshal(post)
			if err != nil {
				return err
			}
			err = ws.WriteMessage(websocket.TextMessage, j)
			if err != nil {
				c.Logger().Error(err)
			}
		case <-quit:
			ticker.Stop()
		}
	}

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
	}

	return nil
}
