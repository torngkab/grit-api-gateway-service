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
	Transactions []TransactionModel `json:"transactions"`
	Total        int32              `json:"total"`
	Page         int32              `json:"page"`
	Limit        int32              `json:"limit"`
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

	// convert protobuf transactions to TransactionModel
	transactions := make([]TransactionModel, len(getTransactionsResponse.GetTransactions()))
	for i, pbTx := range getTransactionsResponse.GetTransactions() {
		transactions[i] = TransactionModel{
			TransactionId: pbTx.GetTransactionId(),
			AccountId:     pbTx.GetAccountId(),
			Type:          pbTx.GetType(),
			Amount:        float64(pbTx.GetAmount()),
			Note:          pbTx.GetNote(),
			CreatedAt:     pbTx.GetCreatedAt(),
		}
	}

	return c.JSON(http.StatusOK, GetTransactionsResponse{
		Transactions: transactions,
		Total:        getTransactionsResponse.Total,
		Page:         getTransactionsResponse.Page,
		Limit:        getTransactionsResponse.Limit,
	})
}
