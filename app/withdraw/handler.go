package withdraw

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type WithdrawRequest struct {
	AccountId string  `json:"account_id" validate:"required"`
	Amount    float64 `json:"amount" validate:"required,min=0"`
	Note      *string `json:"note,omitempty"`
}

type WithdrawResponse struct {
	TransactionId string `json:"transaction_id"`
}

type WithdrawResponseError struct {
	Error string `json:"error"`
}

type Handler struct {
	svc WithdrawService
}

func NewHandler(svc WithdrawService) *Handler {
	return &Handler{svc: svc}
}

// Withdraw godoc
//
//	@Summary		Withdraw
//	@Description	Withdraw money
//	@Tags			withdraw
//	@Accept			json
//	@Produce		plain
//	@Success		200	{object}	WithdrawResponse	"Success"
//	@Failure		400	{object}	WithdrawResponseError	"Bad Request"
//	@Failure		500	{object}	WithdrawResponseError	"Internal Server Error"
//	@Router			/api/v1/withdraw [post]
//	@Param			request	body		WithdrawRequest	true	"Withdraw request"
func (h *Handler) Withdraw(c echo.Context) error {
	var request WithdrawRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, WithdrawResponseError{
			Error: err.Error(),
		})
	}

	if err := c.Validate(&request); err != nil {
		return c.JSON(http.StatusBadRequest, WithdrawResponseError{
			Error: err.Error(),
		})
	}

	withdrawResponse, err := h.svc.Withdraw(c.Request().Context(), request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, WithdrawResponseError{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, WithdrawResponse{
		TransactionId: withdrawResponse,
	})
}
