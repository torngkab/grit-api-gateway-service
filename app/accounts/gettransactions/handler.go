package gettransactions

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type GetTransactionsRequest struct {
	AccountIds []string
	Page       int32 `query:"page" json:"page" default:"1" validate:"required"`
	Limit      int32 `query:"limit" json:"limit" default:"10" validate:"required"`
}

type GetTransactionsResponse struct {
	Transactions []*TransactionModel `json:"transactions"`
	Total        int32               `json:"total"`
	Page         int32               `json:"page"`
	Limit        int32               `json:"limit"`
}

type GetTransactionsResponseError struct {
	Error string `json:"error"`
}

type Handler struct {
	svc GetTransactionsService
}

func NewHandler(svc GetTransactionsService) *Handler {
	return &Handler{svc: svc}
}

// GetTransactions godoc
//
//	@Summary		Get Transactions
//	@Description	Get the transactions of an account
//	@Tags			accounts
//	@Accept			json
//	@Produce		plain
//	@Success		200	{object}	GetTransactionsResponse	"Success"
//	@Failure		400	{object}	GetTransactionsResponseError	"Bad Request"
//	@Failure		500	{object}	GetTransactionsResponseError	"Internal Server Error"
//	@Router			/api/v1/transactions [get]
//	@Param			request	body		GetTransactionsRequest	true	"Get transactions request"
func (h *Handler) GetTransactions(c echo.Context) error {
	// declare request
	var request GetTransactionsRequest

	// bind request
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, GetTransactionsResponseError{
			Error: err.Error(),
		})
	}

	// validate request
	if err := c.Validate(&request); err != nil {
		return c.JSON(http.StatusBadRequest, GetTransactionsResponseError{
			Error: err.Error(),
		})
	}

	// get account ids by user id
	accountIds, err := h.svc.GetAccountIdsByUserId(c.Request().Context(), c.Get("userId").(string))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GetTransactionsResponseError{
			Error: err.Error(),
		})
	}

	// mapping account ids
	request.AccountIds = make([]string, len(accountIds))
	for i, account := range accountIds {
		request.AccountIds[i] = account.Id
	}

	// get transactions by account ids
	getTransactionsResponse, err := h.svc.GetTransactions(c.Request().Context(), request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GetTransactionsResponseError{
			Error: err.Error(),
		})
	}

	// convert transactions to transaction models
	transactionsPb := make([]*TransactionModel, len(getTransactionsResponse.Transactions))
	for i, transaction := range getTransactionsResponse.Transactions {
		transactionsPb[i] = &TransactionModel{
			TransactionId: transaction.TransactionId,
			AccountId:     transaction.AccountId,
			Type:          transaction.Type,
			Amount:        float64(transaction.Amount),
			Note:          transaction.Note,
			CreatedAt:     transaction.CreatedAt,
		}
	}

	return c.JSON(http.StatusOK, GetTransactionsResponse{
		Transactions: transactionsPb,
		Total:        getTransactionsResponse.Total,
		Page:         getTransactionsResponse.Page,
		Limit:        getTransactionsResponse.Limit,
	})
}
