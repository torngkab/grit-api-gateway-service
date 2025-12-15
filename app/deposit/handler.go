package deposit

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type DepositRequest struct {
	AccountId string  `json:"account_id" validate:"required"`
	Amount    float64 `json:"amount" validate:"required,min=0"`
	Note      *string `json:"note,omitempty"`
}

type DepositResponse struct {
	TransactionId string `json:"transaction_id"`
}

type DepositResponseError struct {
	Error string `json:"error"`
}

type Handler struct {
	svc DepositService
}

func NewHandler(svc DepositService) *Handler {
	return &Handler{svc: svc}
}

// Deposit godoc
//
//	@Summary		Deposit
//	@Description	Deposit money
//	@Tags			deposit
//	@Accept			json
//	@Produce		plain
//	@Success		200	{object}	DepositResponse	"Success"
//	@Failure		400	{object}	DepositResponseError	"Bad Request"
//	@Failure		500	{object}	DepositResponseError	"Internal Server Error"
//	@Router			/api/v1/deposit [post]
//	@Param			request	body		DepositRequest	true	"Deposit request"
func (h *Handler) Deposit(c echo.Context) error {
	var request DepositRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, DepositResponseError{
			Error: err.Error(),
		})
	}

	if err := c.Validate(&request); err != nil {
		return c.JSON(http.StatusBadRequest, DepositResponseError{
			Error: err.Error(),
		})
	}

	depositResponse, err := h.svc.Deposit(c.Request().Context(), request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, DepositResponseError{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, DepositResponse{
		TransactionId: depositResponse,
	})
}
