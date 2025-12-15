package transfer

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type TransferRequest struct {
	FromAccountId string  `json:"from_account_id" validate:"required"`
	ToAccountId   string  `json:"to_account_id" validate:"required"`
	Amount        float64 `json:"amount" validate:"required,min=0"`
	Note          *string `json:"note,omitempty"`
}

type TransferResponse struct {
	TransactionId string `json:"transaction_id"`
}

type TransferResponseError struct {
	Error string `json:"error"`
}

type Handler struct {
	svc TransferService
}

func NewHandler(svc TransferService) *Handler {
	return &Handler{svc: svc}
}

// Transfer godoc
//
//	@Summary		Transfer
//	@Description	Transfer money
//	@Tags			transfer
//	@Accept			json
//	@Produce		plain
//	@Success		200	{object}	TransferResponse	"Success"
//	@Failure		400	{object}	TransferResponseError	"Bad Request"
//	@Failure		500	{object}	TransferResponseError	"Internal Server Error"
//	@Router			/api/v1/transfer [post]
//	@Param			request	body		TransferRequest	true	"Transfer request"
func (h *Handler) Transfer(c echo.Context) error {
	var request TransferRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, TransferResponseError{
			Error: err.Error(),
		})
	}

	if err := c.Validate(&request); err != nil {
		return c.JSON(http.StatusBadRequest, TransferResponseError{
			Error: err.Error(),
		})
	}

	transferResponse, err := h.svc.Transfer(c.Request().Context(), request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, TransferResponseError{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, TransferResponse{
		TransactionId: transferResponse,
	})
}
