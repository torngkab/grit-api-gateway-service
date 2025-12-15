package getbalance

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{}
)

type GetBalanceRequest struct {
	AccountId string `json:"account_id" validate:"required"`
}

type GetBalanceResponse struct {
	Balance         float64 `json:"balance"`
	LatestUpdatedAt string  `json:"latest_updated_at"`
}

type GetBalanceResponseError struct {
	Error string `json:"error"`
}

type Handler struct {
	svc GetBalanceService
}

func NewHandler(svc GetBalanceService) *Handler {
	return &Handler{svc: svc}
}

// GetBalance godoc
//
//	@Summary		Get Balance
//	@Description	Get the balance of an account
//	@Tags			accounts
//	@Accept			json
//	@Produce		plain
//	@Success		200	{object}	GetBalanceResponse	"Success"
//	@Failure		400	{object}	GetBalanceResponseError	"Bad Request"
//	@Failure		500	{object}	GetBalanceResponseError	"Internal Server Error"
//	@Router			/api/v1/accounts/:id/balance [get]
//	@Param			request	body		GetBalanceRequest	true	"Get balance request"
func (h *Handler) GetBalance(c echo.Context) error {
	var request GetBalanceRequest

	// get account id from path
	accountId := c.Param("id")
	log.Println("accountId", accountId)
	if accountId == "" {
		return c.JSON(http.StatusBadRequest, GetBalanceResponseError{
			Error: "AccountId is required",
		})
	}

	// set account id to request
	request.AccountId = accountId

	getBalanceResponse, err := h.svc.GetBalance(c.Request().Context(), request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GetBalanceResponseError{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, GetBalanceResponse{
		Balance:         getBalanceResponse.Balance,
		LatestUpdatedAt: getBalanceResponse.LatestUpdatedAt,
	})
}

func (h *Handler) GetBalanceWebSocket(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		var request GetBalanceRequest

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

		getBalanceResponse, err := h.svc.GetBalance(c.Request().Context(), request)
		if err != nil {
			c.Logger().Error(err)
		}

		// Write
		responseBytes, marshalErr := json.Marshal(getBalanceResponse)
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
