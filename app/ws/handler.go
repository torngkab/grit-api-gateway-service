package ws

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"github.com/torngkab/grit-api-gateway-service/app/accounts/getbalance"
)

var (
	upgrader = websocket.Upgrader{}
)

type Handler struct {
	svc getbalance.GetBalanceService
}

func NewHandler(svc getbalance.GetBalanceService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) ServeWebSocket(c echo.Context) error {

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		var request getbalance.GetBalanceRequest

		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}

		fmt.Printf("%s\n", msg)

		err = json.Unmarshal(msg, &request)
		if err != nil {
			c.Logger().Error(err)
		}

		getBalanceResponse, err := h.svc.GetBalance(c.Request().Context(), getbalance.GetBalanceRequest{
			AccountId: request.AccountId,
		})
		if err != nil {
			c.Logger().Error(err)
		}

		// Write
		responseBytes, marshalErr := json.Marshal(getbalance.GetBalanceResponse{
			Balance:         getBalanceResponse.Balance,
			LatestUpdatedAt: getBalanceResponse.LatestUpdatedAt,
		})
		if marshalErr != nil {
			c.Logger().Error(marshalErr)
			continue
		}

		err = ws.WriteMessage(websocket.TextMessage, responseBytes)
		if err != nil {
			c.Logger().Error(err)
		}
	}
}
